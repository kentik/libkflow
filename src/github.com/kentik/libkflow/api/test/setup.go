package test

import (
	"encoding/base64"
	"math/rand"
	"net/url"
	"time"

	"github.com/kentik/libkflow/api"
)

func NewClientServer() (*api.Client, *Server, *api.Device, error) {
	var (
		email  = randstr(8)
		token  = randstr(8)
		device = &api.Device{
			ID:          int(rand.Uint32()),
			Name:        randstr(8),
			MaxFlowRate: 10,
			CompanyID:   int(rand.Uint32()),
		}
	)

	server, err := NewServer("127.0.0.1", 0, false, true)
	if err != nil {
		return nil, nil, nil, err
	}
	go server.Serve(email, token, device)

	apiurl, _ := url.Parse(server.URL() + "/api/v5")

	client := api.NewClient(api.ClientConfig{
		Email:   email,
		Token:   token,
		Timeout: 1 * time.Second,
		API:     apiurl,
		Proxy:   nil,
	})

	return client, server, device, nil
}

func randstr(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
