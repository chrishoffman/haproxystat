package main

import (
	"flag"
	"fmt"

	"github.com/chrishoffman/haproxystat"
)

// Config contains configuration data for the cli
type Config struct {
	StatsdHost  string
	StatsdPort  int
	StatPrefix  string
	BindAddress string
	Port        int
}

var config *Config

func init() {
	const (
		statsdHostDefault = "localhost"
		statsdHostDesc    = "The hostname for the StatsD server"

		statsdPortDefault = 8125
		statsdPortDesc    = "The port for the StatsD server"

		statPrefixDefault = "haproxy"
		statPrefixDesc    = "The prefix to use for stats. Use %HOST% for the current hostname."

		bindAddressDefault = "127.0.0.1"
		bindAddressDesc    = "The address to bind to for ingesting syslog messages"

		portDefault = 10514
		portDesc    = "The port to listen for syslog messages"
	)

	config = &Config{}

	flag.StringVar(&config.BindAddress, "bind-addr", bindAddressDefault, bindAddressDesc)
	flag.IntVar(&config.Port, "port", portDefault, portDesc)

	flag.StringVar(&config.StatsdHost, "statsd-host", statsdHostDefault, statsdHostDesc)
	flag.IntVar(&config.StatsdPort, "statsd-port", statsdPortDefault, statsdPortDesc)
	flag.StringVar(&config.StatPrefix, "stat-prefix", statPrefixDefault, statPrefixDesc)
}

func main() {
	flag.Parse()

	statsdAddress := fmt.Sprintf("%s:%d", config.StatsdHost, config.StatsdPort)
	statsdHandler := newStatsdHandler(statsdAddress, config.StatPrefix)

	s := haproxystat.NewServer()
	s.AddHandler(statsdHandler.logHandler())
	s.Start(config.BindAddress, config.Port)
}
