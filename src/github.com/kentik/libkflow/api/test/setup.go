package test

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/kentik/libkflow/api"
)

func NewClientServer() (*api.Client, *Server, *api.Device, error) {
	var (
		email  = randstr(8)
		token  = randstr(8)
		device = api.Device{
			ID:          rand.Int(),
			Name:        randstr(8),
			MaxFlowRate: 10,
			CompanyID:   rand.Int(),
		}
	)

	client := api.NewClient(email, token, 1*time.Second)

	server, err := NewServer("127.0.0.1", 0, false)
	if err != nil {
		return nil, nil, nil, err
	}
	go server.Serve(email, token, device)

	return client, server, &device, nil
}

func randstr(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
