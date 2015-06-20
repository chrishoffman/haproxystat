# haproxystat [![Build Status](https://travis-ci.org/chrishoffman/haproxystat.png)](https://travis-ci.org/chrishoffman/haproxystat) [![GoDoc](https://godoc.org/github.com/chrishoffman/haproxystat?status.svg)](https://godoc.org/github.com/chrishoffman/haproxystat)

haproxystat is a go library that provides a syslog server that parses the [haproxy](http://www.haproxy.org/) syslog messages in realtime and easily streams them to external sources.

HAProxy provides incredibly detailed information from its logs that can be very useful in understanding the workloads of your proxy. The integrated stats page in HAProxy provides very useful overview data but does not always give the detail you want. By hooking into the syslog stream of detailed log data, you can easily have access to all the data your haproxy logs provide in realtime.

## Installation

Standard `go get`:

```
$ go get github.com/chrishoffman/haproxystat
```

## Example

haproxystat comes bundled with an example daemon that supports sending stat data to [statsd](https://github.com/etsy/statsd). You can use this daemon or build your own by creating your own handlers and attaching them to the configured syslog server.

To build the example daemon:

```
$ go build cmd/haproxystat
```

## Configuration

haproxystat acts as a syslog server and can be easily configured with your favorite syslog forwarder. In the case of rsyslog, add the following line to your rsyslog configuration:

```
local0.* @@localhost:10514
```
