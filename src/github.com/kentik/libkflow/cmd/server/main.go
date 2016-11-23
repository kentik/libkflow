package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/api/test"
)

type Args struct {
	Host       string            `short:"h"          description:"listen on host"`
	Port       int               `short:"p"          description:"listen on port"`
	TLS        bool              `long:"tls"         description:"require TLS   "`
	Email      string            `long:"email"       description:"API auth email"`
	Token      string            `long:"token"       description:"API auth token"`
	CompanyID  int               `long:"company-id"  description:"company ID    "`
	DeviceID   int               `long:"device-id"   description:"device ID     "`
	DeviceName string            `long:"device-name" description:"device name   "`
	Customs    map[string]uint64 `long:"custom"      description:"custom fields "`
}

func main() {
	args := Args{
		Host:       "127.0.0.1",
		Port:       8999,
		TLS:        false,
		Email:      "test@example.com",
		Token:      "token",
		CompanyID:  1,
		DeviceID:   1,
		DeviceName: "dev1",
		Customs: map[string]uint64{
			"RETRANSMITTED_IN_PKTS":  1,
			"RETRANSMITTED_OUT_PKTS": 2,
			"FRAGMENTS":              3,
			"CLIENT_NW_LATENCY_MS":   4,
			"SERVER_NW_LATENCY_MS":   5,
			"APPL_LATENCY_MS":        6,
			"OOORDER_IN_PKTS":        7,
			"OOORDER_OUT_PKTS":       8,
			"KFLOW_HTTP_URL":         9,
			"KFLOW_HTTP_STATUS":      10,
			"KFLOW_HTTP_UA":          11,
			"KFLOW_HTTP_REFERER":     12,
			"KFLOW_DNS_QUERY":        13,
			"KFLOW_DNS_QUERY_TYPE":   14,
			"KFLOW_DNS_RET_CODE":     15,
		},
	}

	parser := flags.NewParser(&args, flags.PassDoubleDash|flags.HelpFlag)
	if _, err := parser.Parse(); err != nil {
		switch err.(*flags.Error).Type {
		case flags.ErrHelp:
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		default:
			log.Fatal(err)
		}
	}

	s, err := test.NewServer(args.Host, args.Port, args.TLS)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s:%d", s.Host, s.Port)

	err = s.Serve(args.Email, args.Token, api.Device{
		ID:          args.DeviceID,
		Name:        args.DeviceName,
		MaxFlowRate: 4000,
		CompanyID:   args.CompanyID,
		Customs:     args.Customs,
	})

	if err != nil {
		log.Fatal(err)
	}
}
