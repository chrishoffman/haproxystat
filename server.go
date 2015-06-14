package haproxystat

import (
	"fmt"

	"github.com/chrishoffman/haproxylog"
	"gopkg.in/mcuadros/go-syslog.v2"
)

// LogHandler reprents a handler that takes in a haproxy.Log message for processing
type LogHandler func(*haproxy.Log)

// Server is the container for the server instance
type Server struct {
	handlers []LogHandler
}

// NewServer creates a new intances of haproxystat.Server
func NewServer() *Server {
	return &Server{}
}

// AddHandler adds a new handler to the log processing pipeline
func (s *Server) AddHandler(handler LogHandler) {
	s.handlers = append(s.handlers, handler)
}

// Start binds the syslog server, attaches the handlers and
// waits for traffic to delegate to the handlers
func (s *Server) Start(bindAddress string, port int) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)

	listenAddress := fmt.Sprintf("%s:%d", bindAddress, port)
	listenErr := server.ListenTCP(listenAddress)
	if listenErr != nil {
		panic(listenErr)
	}

	bootErr := server.Boot()
	if bootErr != nil {
		panic(bootErr)
	}

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			s.logHandler(logParts["content"].(string))
		}
	}(channel)

	server.Wait()
}

func (s *Server) logHandler(rawLog string) {
	log, err := haproxy.NewLog(rawLog)
	if err != nil {
		return
	}

	for _, handler := range s.handlers {
		handler(log)
	}
}
