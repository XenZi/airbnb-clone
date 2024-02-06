package client

import (
	"crypto/tls"
	"net/http"
)

func NewCustomClient(certPath, keyPath string) *http.Client {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		TLSClientConfig:     tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}
