package logger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync/atomic"
	"time"
	"unsafe"
)

// The csyslog function is necessary here because cgo does not appear
// to be able to call a variadic function directly and syslog has the
// same signature as printf.

// #include <stdlib.h>
// #ifdef _WIN32               /* for Windows builds, syslog functions do NOTHING */
// void csyslog(int p, const char *m) {}
// void openlog(const char *m, int i, int l) {}
// #define	LOG_ERR		3      /* error conditions */
// #define	LOG_WARNING	4      /* warning conditions */
// #define	LOG_INFO	6      /* informational */
// #define	LOG_DEBUG	7      /* debug-level messages */
// #define	LOG_PID		0x01   /* log the pid with each message */
// #define	LOG_NDELAY	0x08   /* don't delay open */
// #define	LOG_NOWAIT	0x10   /* don't wait for console forks: DEPRECATED */
// #define	LOG_USER	(1<<3) /* random user-level messages */
// #else
// #include <syslog.h>
// void csyslog(int p, const char *m) {
//     syslog(p, "%s", m);
// }
// #endif
import "C"

const (
	NumMessages   = 10 * 1024 // number of allowed log messages
	STDOUT_FORMAT = "2006-01-02T15:04:05.000 "
)

// logMessage contains a pending log message
type logMessage struct {
	bytes.Buffer
	level C.int
	time  time.Time
	le    logEntry
}

// logCaller stores where the logger public log method was called
type logCaller struct {
	File string
	Line int
}

func (lc logCaller) String() string {
	return fmt.Sprintf("%s:%d", lc.File, lc.Line)
}

// logEntry encapsulates all parameters to queueMsg
type logEntry struct {
	lvl  Level
	pre  string
	fmt  string
	fmtV []interface{}
	lc   logCaller
	tee  bool
}

var (
	ErrLogFullBuf           = errors.New("Log message queue is full")
	ErrFreeMessageOverflow  = errors.New("Too many free messages. Overflow of fixed	set.")
	ErrFreeMessageUnderflow = errors.New("Too few free messages. Underflow of fixed	set.")

	// the logName object for syslog to use
	logName       *C.char
	logNameString string

	// the message queue of pending or free messages
	// since only one can be full at a time, the total size will be about 10MB
	messages     chan *logMessage
	freeMessages chan *logMessage

	// mapping of our levels to syslog values
	levelSysLog = map[Level]C.int{
		Levels.Access: C.LOG_INFO,
		Levels.Off:    C.LOG_DEBUG,
		Levels.Panic:  C.LOG_ERR,
		Levels.Error:  C.LOG_ERR,
		Levels.Warn:   C.LOG_WARNING,
		Levels.Info:   C.LOG_INFO,
		Levels.Debug:  C.LOG_DEBUG,
	}

	// mirror of levelMap used to avoid making a new string with '[]' on every log
	// call
	levelMapFmt = map[Level][]byte{
		Levels.Access: []byte("[Access] "),
		Levels.Off:    []byte("[Off] "),
		Levels.Panic:  []byte("[Panic] "),
		Levels.Error:  []byte("[Error] "),
		Levels.Warn:   []byte("[Warn] "),
		Levels.Info:   []byte("[Info] "),
		Levels.Debug:  []byte("[Debug] "),
	}

	customSock net.Conn = nil

	logWriterFinished chan struct{}

	stdhdl io.Writer

	logTee chan string
)

// SetCustomSocket will switch over to writing log messages to the defined socket.
func SetCustomSocket(address, network string) (err error) {
	customSock, err = net.Dial(network, address)

	return err
}

func SetStdOut() {
	stdhdl = io.Writer(os.Stdout)
}

func SetStdErr() {
	stdhdl = io.Writer(os.Stderr)
}

func SetTee(tee chan string) {
	logTee = tee
}

// SetLogName sets the identifier used by syslog for this program
func SetLogName(p string) (err error) {

	logNameString = p
	if stdhdl != nil {
		return
	}

	if logName != nil {
		C.free(unsafe.Pointer(logName))
	}
	logName = C.CString(p)
	_, err = C.openlog(logName, C.LOG_NDELAY|C.LOG_NOWAIT|C.LOG_PID, C.LOG_USER)
	if err != nil {
		atomic.AddUint64(&errCount, 1)
	}

	return err
}

const reclaimThreshold int = 5120 // Seems to be around p99 on our runner-master log messages

// freeMsg releases the message back to be reused
func freeMsg(msg *logMessage) (err error) {
	if msg.Buffer.Cap() > reclaimThreshold {
		msg.Buffer = *bytes.NewBuffer(make([]byte, 0, reclaimThreshold))
	} else {
		msg.Reset()
	}
	select {
	case freeMessages <- msg: // no-op
	default:
		atomic.AddUint64(&errCount, 1)
		return ErrFreeMessageOverflow
	}

	return
}

// queueMsg adds a message to the pending messages channel. It will drop the
// message and return an error if the channel is full.
func queueMsg(le *logEntry) (err error) {
	atomic.AddUint64(&logCount, 1)
	var msg *logMessage

	// get a message if possible
	select {
	case msg = <-freeMessages: // got a message-struct; proceed
	default:
		// no messages left, drop
		atomic.AddUint64(&dropCount, 1)
		return
	}

	msg.time = time.Now()
	msg.le = *le
	msg.level = levelSysLog[le.lvl]

	// formatting & teeing the log message is no longer done here, both are now in `logWriter`

	// queue the message
	select {
	case messages <- msg:
		// no-op
	default:
		// this should never happen since there is an exact number of messages
		atomic.AddUint64(&errCount, 1)
		return ErrLogFullBuf
	}

	return
}

// printTee writes the logMessage's buffer to logTee
func printTee(msg *logMessage) {
	select {
	case logTee <- string(msg.Bytes()):
	default:
		LogNoTee(Levels.Error, "[meta log]", "%s log tee is full", logTee)
		atomic.AddUint64(&errCount, 1)
	}
	return
}

// printStd prints msg to stdhdl
func printStd(msg *logMessage) (err error) {
	_, err = io.Copy(stdhdl, msg)
	return
}

// write function writes a message to syslog. This is a concrete, blocking event.
func write(msg *logMessage) (err error) {
	// terminate C string
	_ = msg.WriteByte(0) // *Buffer.WriteByte always returns nil

	start := (*C.char)(unsafe.Pointer(&msg.Bytes()[0]))
	_, err = C.csyslog(C.LOG_USER|msg.level, start)
	return
}

// writeCustomSocket writes a message to a pre-defined custom socket.
// This is a concrete, blocking event. Writes out using the syslog rfc5424 format.
func writeCustomSocket(msg *logMessage) (err error) {
	// terminate C string
	_ = msg.WriteByte(0) // *Buffer.WriteByte always returns nil

	_, err = customSock.Write(bytes.Join([][]byte{[]byte(fmt.Sprintf("<%d>", C.LOG_USER|msg.level)),
		msg.Bytes()}, []byte("")))
	return
}

// logWriter will write out messages to syslog. It may block if something breaks
// within the syslog call.
func logWriter() {
	var err error
	for msg := range messages {
		// format/build the log string
		if sendJSON {
			err = msg.asJSON()
		} else {
			err = msg.asString()
		}
		if err != nil {
			atomic.AddUint64(&errCount, 1)
			_ = freeMsg(msg)
			continue
		}

		// tee as needed
		if logTee != nil && msg.le.tee {
			printTee(msg)
		}

		// write the log string
		if stdhdl != nil {
			err = printStd(msg)
		} else if customSock == nil {
			err = write(msg)
		} else {
			err = writeCustomSocket(msg)
		}
		if err != nil {
			atomic.AddUint64(&errCount, 1)
		}
		_ = freeMsg(msg)
	}
	if customSock != nil {
		_ = customSock.Close()
	}
	close(logWriterFinished)
}

// Close shuts down the logger system. After Close is called, any additional
// logs will panic. Only call this if you are completely done.
func Close(ctx context.Context) error {
	close(messages)
	select {
	case <-logWriterFinished:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DrainContext blocks until it sees no pending messages or the context is canceled.
// Pending messages may never run out if another goroutine is constantly
// writing.
func DrainContext(ctx context.Context) error {
	for ctx.Err() == nil && (len(messages) > 0 || len(freeMessages) < cap(freeMessages)) {
		innerCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		<-innerCtx.Done() // Wait for 10ms and check the queues again
		cancel()
	}
	return ctx.Err()
}

// DrainWithTimeout is a helper function that uses the given timeout with
// DrainContext.
func DrainWithTimeout(d time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	_ = DrainContext(ctx)
	cancel()
}

// Drain is like DrainContext, but you didn't want to write context.Background().
// Outside of tests, you want to use DrainContext.
func Drain() {
	_ = DrainContext(context.Background())
}

func setup() {
	stdhdl = nil
	messages = make(chan *logMessage, NumMessages)
	freeMessages = make(chan *logMessage, NumMessages)
	msgArr := make([]logMessage, NumMessages)
	for i := range msgArr {
		if err := freeMsg(&msgArr[i]); err != nil {
			break
		}
	}

	logWriterFinished = make(chan struct{}, 1)
	go logWriter()
}

func init() {
	setup()
	setSendJSON()
}
