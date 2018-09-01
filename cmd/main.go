package main

import (
	"flag"
	//"fmt"
	"log"
	"log/syslog"
	"net"
	"os"

	//"github.com/phalaaxx/milter"
	"github.com/porjo/milter"
	"github.com/porjo/pol"
	"github.com/porjo/pol/rcptHeader"
	"github.com/porjo/pol/spf"
)

var Email = ""
var sysLog *syslog.Writer

func main() {
	// parse commandline arguments
	var protocol, address string
	flag.StringVar(&protocol,
		"proto",
		"unix",
		"Protocol family (unix or tcp)")
	flag.StringVar(&address,
		"addr",
		"/tmp/pol.sock",
		"Bind to address or unix domain socket")
	flag.StringVar(&Email,
		"email",
		"",
		"Email address to check")
	flag.Parse()

	if Email == "" {
		log.Fatal("Email cannot be empty")
	}

	// make sure the specified protocol is either unix or tcp
	if protocol != "unix" && protocol != "tcp" {
		log.Fatal("invalid protocol name")
	}

	// make sure socket does not exist
	if protocol == "unix" {
		// ignore os.Remove errors
		os.Remove(address)
	}

	// bind to listening address
	socket, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	if protocol == "unix" {
		// set mode 0660 for unix domain sockets
		if err := os.Chmod(address, 0660); err != nil {
			log.Fatal(err)
		}
		// remove socket on exit
		defer os.Remove(address)
	}

	sysLog, err = syslog.New(syslog.LOG_INFO|syslog.LOG_MAIL, "pol")
	if err != nil {
		log.Fatal(err)
	}

	// declare milter init function
	init := func() (milter.Milter, uint32, uint32) {
		var milters []milter.Milter
		rcpt := &rcptHeader.Milter{
			Email:  Email,
			Logger: log.New(sysLog, "", log.LstdFlags),
		}

		spf_ := &spf.Milter{
			Logger: log.New(sysLog, "", log.LstdFlags),
		}
		milters = append(milters, spf_, rcpt)
		return &pol.Milter{Milters: milters},
			milter.OptAddHeader,
			milter.OptNoConnect | milter.OptNoHelo | milter.OptNoBody
	}
	// start server
	if err := milter.RunServer(socket, init); err != nil {
		log.Fatal(err)
	}
}
