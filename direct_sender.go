package libkflow

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/url"

	capnp "zombiezen.com/go/capnproto2"

	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/flow"
)

// messagePrefix is an 80-byte prefix for the message header when sending kflow to the Kentik API. This is a deprecated
// header, but the bytes must remain for backwards compatibility with the Kentik API.
var messagePrefix = [80]byte{}

// A DirectFlowSender transmits flows to Kentik.
type DirectFlowSender struct {
	client *api.Client

	url      string
	deviceID uint32
}

// NewDirectSender creates a new DirectFlowSender.
func NewDirectSender(client *api.Client, device *api.Device, url *url.URL) *DirectFlowSender {
	return &DirectFlowSender{
		client: client,

		url:      createURLString(*url, device.ClientID()),
		deviceID: uint32(device.ID),
	}
}

// SendFlows sends the flows to the Kentik API, returning the number of bytes sent as the payload. The device ID on
// the flows is set to the device ID of the sender, regardless of what it was previously set to. This is to ensure all
// data matches the expectations of the downstream URL/API.
//
// This will directly send the slice of flows without any additional downsampling or rate limiting. This does not
// contribute to the underlying Send call.
func (s *DirectFlowSender) SendFlows(flows []flow.Flow) (int64, error) {
	if len(flows) == 0 {
		return 0, nil
	}
	// ensure all flows have the device ID set; otherwise it may not be properly queried
	for i := range flows {
		flows[i].DeviceId = s.deviceID
	}

	// ensure the sample rate is matching the kentik api expectations
	flow.NormalizeSampleRate(flows, 0)

	// serialize the data
	_, segment, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return 0, fmt.Errorf("failed to create capn proto segment: %w", err)
	}
	message, err := flow.ToCapnProtoMessage(flows, segment)
	if err != nil {
		return 0, fmt.Errorf("failed to convert flows to capn proto: %w", err)
	}

	// write the data with additional gzip compression
	buf := &bytes.Buffer{}
	z := gzip.NewWriter(buf)
	_, err = z.Write(messagePrefix[:])
	if err != nil {
		return 0, fmt.Errorf("failed to write empty message header: %w", err)
	}
	err = capnp.NewPackedEncoder(z).Encode(message)
	if err != nil {
		return 0, fmt.Errorf("failed to encode packed capn proto message: %w", err)
	}
	err = z.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// send the compressed and packed message to the Kentik API
	payloadLength := int64(len(buf.Bytes()))
	err = s.client.SendFlow(s.url, buf)
	if err != nil {
		return 0, err
	}
	return payloadLength, nil
}

// createURLString creates the full URL string to use when sending data to the Kentik API.
func createURLString(u url.URL, clientID string) string {
	q := u.Query()
	q.Set("sid", "0")
	q.Set("sender_id", clientID)
	u.RawQuery = q.Encode()
	return u.String()
}
