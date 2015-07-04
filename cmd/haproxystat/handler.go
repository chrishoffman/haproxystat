package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/chrishoffman/haproxylog"
)

type statsdHandler struct {
	client statsd.Statter
}

func newStatsdHandler(address string, prefix string) *statsdHandler {
	flushInterval := time.Second * 1

	client, err := statsd.NewBufferedClient(address, prefix, flushInterval, 0)
	if err != nil {
		panic(err)
	}

	return &statsdHandler{client}
}

func (s *statsdHandler) logHandler(log *haproxy.Log) {
	switch log.GetFormat() {
	case haproxy.HTTP:
		s.sendHTTPStats(log)
	}
}

func (s *statsdHandler) sendHTTPStats(log *haproxy.Log) {
	frontendName := cleanStatToken(log.BackendName)
	backendName := cleanStatToken(log.FrontendName)
	serverName := cleanAndLowerStatToken(log.ServerName)
	requestStatPrefix := fmt.Sprintf("%s.%s.%s", frontendName, backendName, serverName)

	// Request stats
	s.inc(fmt.Sprintf("%s.response_size", requestStatPrefix), log.BytesRead)
	s.inc(fmt.Sprintf("%s.hits", requestStatPrefix), 1)
	s.inc(fmt.Sprintf("%s.responses.%d", requestStatPrefix, log.HTTPStatusCode), 1)

	// Timing Stats
	s.timing(fmt.Sprintf("%s.response_time", requestStatPrefix), log.Tt)
	s.timing(fmt.Sprintf("%s.queue_time", requestStatPrefix), log.Tw)
	s.timing(fmt.Sprintf("%s.request_time", requestStatPrefix), log.Tq)

	// Misc Stats (stored in timing to average)
	s.timing(fmt.Sprintf("%s.retries", requestStatPrefix), log.Retries)
	s.timing(fmt.Sprintf("%s.queue", requestStatPrefix), log.ServerQueue)
	s.timing(fmt.Sprintf("%s.active_connections", requestStatPrefix), log.ActConn)
	s.timing(fmt.Sprintf("%s.backend_connections", requestStatPrefix), log.BeConn)
	s.timing(fmt.Sprintf("%s.frontend_connections", requestStatPrefix), log.FeConn)
	s.timing(fmt.Sprintf("%s.server_connections", requestStatPrefix), log.SrvConn)
	s.timing(fmt.Sprintf("%s.response_size", requestStatPrefix), log.BytesRead)

	// Backend Stats
	backendPrefix := fmt.Sprintf("backend.%s", backendName)
	s.timing(fmt.Sprintf("%s.connect_time", backendPrefix), log.Tc)
	s.timing(fmt.Sprintf("%s.response_time", backendPrefix), log.Tr)
	s.timing(fmt.Sprintf("%s.queue", backendPrefix), log.BackendQueue)

	// Endpoint stats
	pathParts := strings.Split(log.HTTPRequest.URL.Path, "/")
	if len(pathParts) > 1 {
		basePath := pathParts[1]
		if basePath == "" {
			basePath = "_root_"
		}
		s.timing(fmt.Sprintf("%s.endpoint.%s.%s.response_time", requestStatPrefix, cleanAndLowerStatToken(basePath),
			log.HTTPRequest.Method), log.Tt)
	}

	// SSL stats
	if log.SslVersion != "" {
		sslStat := fmt.Sprintf("%s.ssl.%s.%s", cleanStatToken(log.FrontendName),
			cleanStatToken(log.SslVersion), cleanStatToken(log.SslCipher))
		s.inc(sslStat, 1)
	}
}

func (s *statsdHandler) inc(stat string, count int64) {
	s.client.Inc(stat, count, 1.0)
}

func (s *statsdHandler) timing(stat string, delta int64) {
	if delta == -1 {
		return
	}
	s.client.Timing(stat, delta, 1.0)
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
