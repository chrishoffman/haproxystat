package haproxystat

import (
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

// HaproxyHTTPLog represents a HAProxy HTTP log line
type HaproxyHTTPLog struct {
	ClientIP                string
	ClientPort              int64
	AcceptDate              time.Time
	FrontendName            string
	BackendName             string
	ServerName              string
	Tq                      int64
	Tw                      int64
	Tc                      int64
	Tr                      int64
	Tt                      int64
	HTTPStatusCode          int64
	BytesRead               int64
	CapturedRequestCookie   string
	CapturedResponseCookie  string
	TerminationState        string
	ActConn                 int64
	FeConn                  int64
	BeConn                  int64
	SrvConn                 int64
	Retries                 int64
	ServerQueue             int64
	BackendQueue            int64
	CapturedRequestHeaders  []string
	CapturedResponseHeaders []string
	HTTPRequest             *HaproxyHTTPRequest
	SslCipher               string
	SslVersion              string
}

// HaproxyHTTPRequest is the HTTPRequest object in an HAProxy HTTP log line
type HaproxyHTTPRequest struct {
	Method  string
	URL     *url.URL
	Version string
}

// This represents the default HTTP format and adds additional optional SSL information. Both the default and the modified version work.
// log-format %ci:%cp\ [%t]\ %ft\ %b/%s\ %Tq/%Tw/%Tc/%Tr/%Tt\ %ST\ %B\ %CC\ \%CS\ %tsc\ %ac/%fc/%bc/%sc/%rc\ %sq/%bq\ %hr\ %hs\ %{+Q}r\ %sslc/%sslv
var haproxyHTTPLogRegexp = myRegexp{
	regexp.MustCompile(`(?P<ClientIp>(\d{1,3}\.){3}\d{1,3}):(?P<ClientPort>\d{1,5}) ` +
		`\[(?P<AcceptDate>\d{2}/\w{3}/\d{4}(:\d{2}){3}\.\d{3})\] ` +
		`(?P<FrontendName>\S+) (?P<BackendName>[\w-\.]+)/(?P<ServerName>\S+) ` +
		`(?P<Tq>(-1|\d+))/(?P<Tw>(-1|\d+))/(?P<Tc>(-1|\d+))/` +
		`(?P<Tr>(-1|\d+))/(?P<Tt>\+?\d+) ` +
		`(?P<HTTPStatusCode>\d{3}) (?P<BytesRead>\d+) ` +
		`(?P<CapturedRequestCookie>\S+) (?P<CapturedResponseCookie>\S+) ` +
		`(?P<TerminationState>[\w-]{4}) ` +
		`(?P<ActConn>\d+)/(?P<FeConn>\d+)/(?P<BeConn>\d+)/` +
		`(?P<SrvConn>\d+)/(?P<Retries>\d+) ` +
		`(?P<ServerQueue>\d+)/(?P<BackendQueue>\d+) ` +
		`(\{(?P<CapturedRequestHeaders>.*?)\} )?` +
		`(\{(?P<CapturedResponseHeaders>.*?)\} )?` +
		`"(?P<HTTPRequest>.+)"` +
		`( (?P<SslCipher>[\w-]+)/(?P<SslVersion>[\w\.]+))?`)}

type haproxyLog struct {
	rawLog string
}

func newHaproxyLog(log string) *haproxyLog {
	return &haproxyLog{log}
}

func (l *haproxyLog) ParseHTTP() (*HaproxyHTTPLog, error) {
	parsed := haproxyHTTPLogRegexp.FindStringSubmatchMap(l.rawLog)
	if len(parsed) == 0 {
		return nil, nil
	}

	var result HaproxyHTTPLog
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &result,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToHTTPRequestHook,
			stringToTimeHook,
			mapstructure.StringToSliceHookFunc("|"),
		),
	}

	decoder, _ := mapstructure.NewDecoder(config)
	err := decoder.Decode(parsed)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func stringToHTTPRequestHook(f reflect.Type, t reflect.Type, v interface{}) (interface{}, error) {
	if t == reflect.TypeOf(&HaproxyHTTPRequest{}) {
		// Split "POST /relative/path HTTP/1.1"
		parts := strings.Split(v.(string), " ")
		if len(parts) == 3 {
			u, _ := url.Parse(parts[1])
			v = &HaproxyHTTPRequest{Method: parts[0], URL: u, Version: parts[2]}
		}
	}
	return v, nil
}

func stringToTimeHook(f reflect.Type, t reflect.Type, v interface{}) (interface{}, error) {
	if t == reflect.TypeOf(time.Time{}) {
		const format = "02/Jan/2006:15:04:05.000"
		acceptDate, _ := time.Parse(format, v.(string))
		v = acceptDate
	}
	return v, nil
}
