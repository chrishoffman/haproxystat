package main

import (
	"fmt"
	"strings"

	"github.com/chrishoffman/haproxylog"
	"github.com/quipo/statsd"
)

type statsdHandler struct {
	client statsd.Statsd
}

func newStatsdHandler(address string, prefix string) *statsdHandler {
	client := statsd.NewStatsdClient(address, prefix)
	err := client.CreateSocket()
	if err != nil {
		panic(err)
	}

	return &statsdHandler{client}
}

func (s *statsdHandler) logHandler() func(*haproxy.Log) {
	return func(log *haproxy.Log) {
		switch log.GetFormat() {
		case haproxy.HTTP:
			s.sendHTTPStats(log)
		}
	}
}

func (s *statsdHandler) sendHTTPStats(log *haproxy.Log) {
	// Request stats
	requestStatPrefix := fmt.Sprintf("%s.%s.%s", cleanStatToken(log.FrontendName),
		cleanStatToken(log.BackendName), cleanAndLowerStatToken(log.ServerName))

	// HTTP Status Codes
	s.client.Incr(fmt.Sprintf("%s.http_status.%d", requestStatPrefix, log.HTTPStatusCode), 1)

	// Endpoint stats
	pathParts := strings.Split(log.HTTPRequest.URL.Path, "/")
	if len(pathParts) > 1 {
		basePath := pathParts[1]
		if basePath == "" {
			basePath = "_root_"
		}
		s.client.Timing(fmt.Sprintf("%s.endpoint.%s.%s", requestStatPrefix, cleanAndLowerStatToken(basePath),
			log.HTTPRequest.Method), log.Tt)
	}

	// SSL stats
	if log.SslVersion != "" {
		sslStat := fmt.Sprintf("%s.ssl.%s.%s", cleanStatToken(log.FrontendName),
			cleanStatToken(log.SslVersion), cleanStatToken(log.SslCipher))
		s.client.Incr(sslStat, 1)
	}
}

func cleanAndLowerStatToken(s string) string {
	return strings.ToLower(cleanStatToken(s))
}

func cleanStatToken(s string) string {
	r := map[string]string{
		".": "_",
		"~": "",
		"<": "-",
		">": "-",
	}
	for o, n := range r {
		s = strings.Replace(s, o, n, -1)
	}
	return s
}
