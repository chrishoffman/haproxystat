package main

import (
	"testing"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/chrishoffman/haproxylog"
	"github.com/stretchr/testify/assert"
)

func Test_BadReq_DoesNotPanic(t *testing.T) {
	haproxyRawLog := `192.168.9.185:56276 [29/May/2015:10:36:47.766] Service1~ Service1/host-1 2/0/0/10/12 200 423 - - ---- 282/36/0/0/0 0/0 {d7d9b784-4276-42bc-ae79-71e9e84d2b85} {d7d9b784-4276-42bc-ae79-71e9e84d2b85} "<BADREQ>" ECDHE-RSA-AES128-GCM-SHA256/TLSv1.2`
	log, _ := haproxy.NewLog(haproxyRawLog)

	statsdClient, _ := statsd.NewNoopClient()
	haproxyStatHandler := &statsdHandler{statsdClient}

	assert.NotPanics(t, func() { haproxyStatHandler.logHandler(log) }, "Bad request does not panic!")
}
