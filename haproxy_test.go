package haproxystat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseHTTPNoMatch_ReturnsNilNoError(t *testing.T) {
	log := NewHaproxyLog("Not a haproxy syslog message")

	parsed, err := log.ParseHTTP()
	assert.Nil(t, parsed, "log returns nil when unable to parse message")
	assert.Nil(t, err, "No error return when cannot parse")
}

func Test_ParseHTTP_ReturnsHaproxyHttpLog(t *testing.T) {
	const haproxySyslogMessage = `192.168.9.185:56276 [29/May/2015:10:36:47.766] Service1~ Service1/host-1 2/0/0/10/12 200 423 - - ---- 282/36/0/0/0 0/0 {d7d9b784-4276-42bc-ae79-71e9e84d2b85} {d7d9b784-4276-42bc-ae79-71e9e84d2b85} "POST /path/to/app HTTP/1.1" ECDHE-RSA-AES128-GCM-SHA256/TLSv1.2`

	log := NewHaproxyLog(haproxySyslogMessage)
	parsed, err := log.ParseHTTP()
	assert.NotNil(t, parsed, "log returns HaproxyHTTPLog")
	assert.Nil(t, err, "No error return")

	assert.Equal(t, "POST", parsed.HTTPRequest.Method, "HTTPRequest method is POST")
	assert.Equal(t, "192.168.9.185", parsed.ClientIP, "ClientIP address matches log")
	assert.Equal(t, int64(200), parsed.HTTPStatusCode, "HTTP Status Code is 200")
}

func Test_ParseHTTPNoResponseHeaders_ReturnsHaproxyHttpLog(t *testing.T) {
	const haproxySyslogMessage = `192.168.9.185:56276 [29/May/2015:10:36:47.766] Service1~ Service1/host-1 2/0/0/10/12 200 423 - - ---- 282/36/0/0/0 0/0 {d7d9b784-4276-42bc-ae79-71e9e84d2b85} "POST /path/to/app HTTP/1.1" ECDHE-RSA-AES128-GCM-SHA256/TLSv1.2`

	log := NewHaproxyLog(haproxySyslogMessage)
	parsed, err := log.ParseHTTP()
	assert.NotNil(t, parsed, "log returns HaproxyHTTPLog")
	assert.Nil(t, err, "No error return")
}

func Test_ParseHTTPNoHeaders_ReturnsHaproxyHttpLog(t *testing.T) {
	const haproxySyslogMessage = `192.168.9.185:56276 [29/May/2015:10:36:47.766] Service1~ Service1/host-1 2/0/0/10/12 200 423 - - ---- 282/36/0/0/0 0/0 "POST /path/to/app HTTP/1.1" ECDHE-RSA-AES128-GCM-SHA256/TLSv1.2`

	log := NewHaproxyLog(haproxySyslogMessage)
	parsed, err := log.ParseHTTP()
	assert.NotNil(t, parsed, "log returns HaproxyHTTPLog")
	assert.Nil(t, err, "No error return")
}

func Test_ParseHTTPNoSSLInfo_ReturnsHaproxyHttpLog(t *testing.T) {
	const haproxySyslogMessage = `192.168.9.185:56276 [29/May/2015:10:36:47.766] Service1~ Service1/host-1 2/0/0/10/12 200 423 - - ---- 282/36/0/0/0 0/0 {d7d9b784-4276-42bc-ae79-71e9e84d2b85} {d7d9b784-4276-42bc-ae79-71e9e84d2b85} "POST /path/to/app HTTP/1.1"`

	log := NewHaproxyLog(haproxySyslogMessage)
	parsed, err := log.ParseHTTP()
	assert.NotNil(t, parsed, "log returns HaproxyHTTPLog")
	assert.Nil(t, err, "No error return")
}

func Test_ParseHTTPDecodeError_ReturnsError(t *testing.T) {
	const haproxySyslogMessage = `192.168.9.185:56276 [29/May/2015:10:36:47.766] Service1~ Service1/host-1 2/0/0/10/12 200 423 - - ---- 282/36/0/0/0 0/0 {d7d9b784-4276-42bc-ae79-71e9e84d2b85} {d7d9b784-4276-42bc-ae79-71e9e84d2b85} "POST /path/to/app" ECDHE-RSA-AES128-GCM-SHA256/TLSv1.2`
	log := NewHaproxyLog(haproxySyslogMessage)
	parsed, err := log.ParseHTTP()
	assert.Nil(t, parsed, "Invalid decode of object returns nil")
	assert.NotNil(t, err, "Invalid log line returns an error")
}
