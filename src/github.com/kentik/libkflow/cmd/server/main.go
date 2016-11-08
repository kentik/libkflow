package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/api/server"
)

type Args struct {
	Host       string            `short:"h"          description:"listen on host"`
	Port       int               `short:"p"          description:"listen on port"`
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
		Email:      "test@example.com",
		Token:      "token",
		CompanyID:  1,
		DeviceID:   1,
		DeviceName: "dev1",
		Customs: map[string]uint64{
			"CUSTOM-STR": 1,
			"CUSTOM-U32": 2,
			"CUSTOM-F32": 3,
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

	s, err := server.New(args.Host, args.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s:%d", s.Host, s.Port)

	err = s.Serve(args.Email, args.Token, api.Device{
		ID:        args.DeviceID,
		Name:      args.DeviceName,
		CompanyID: args.CompanyID,
		Customs:   args.Customs,
	})

	if err != nil {
		log.Fatal(err)
	}
}
