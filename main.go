package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net"
	"net/textproto"
	"os"
	"strings"

	"github.com/phalaaxx/milter"
)

var Email = ""
var sysLog *syslog.Writer

type PolMilter struct {
	milter.Milter
	from       string
	rcptTo     string
	subject    string
	rcptInToCc bool
}

// Header parses message headers one by one
func (b *PolMilter) Header(name, value string, m *milter.Modifier) (milter.Response, error) {
	// check if bogofilter has been run on the message already
	if name == "X-Pol" {
		// X-Polsity header is present, accept immediately
		return milter.RespAccept, nil
	}
	if name == "To" || name == "Cc" {
		if strings.Contains(strings.ToLower(value), Email) {
			b.rcptInToCc = true
		}
	}
	if name == "Subject" {
		b.subject = value
	}
	return milter.RespContinue, nil
}

// RcptTo is called to process filters on envelope TO address
//   supress with NoRcptTo
func (b *PolMilter) RcptTo(rcptTo string, m *milter.Modifier) (milter.Response, error) {
	b.rcptTo = rcptTo
	return milter.RespContinue, nil

}

// MailFrom is called on envelope from address
func (b *PolMilter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	// save from address for later reference
	b.from = from
	return milter.RespContinue, nil
}

// Headers is called after the last of message headers
func (b *PolMilter) Headers(headers textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	if b.rcptTo == Email && !b.rcptInToCc {
		//TODO custom message here?
		sysLog.Info(fmt.Sprintf("REJECT message from %s to %s, subject '%s'\n", b.from, b.rcptTo, b.subject))
		return milter.RespReject, nil
	}
	sysLog.Info(fmt.Sprintf("OK message from %s to %s, subject '%s'\n", b.from, b.rcptTo, b.subject))
	return milter.RespContinue, nil
}

func (b *PolMilter) Body(m *milter.Modifier) (milter.Response, error) {
	m.AddHeader("X-Pol", "1")
	return milter.RespAccept, nil
}

func RunServer(socket net.Listener) {
	// declare milter init function
	init := func() (milter.Milter, uint32, uint32) {
		return &PolMilter{},
			milter.OptAddHeader,
			milter.OptNoConnect | milter.OptNoHelo | milter.OptNoBody
	}
	// start server
	if err := milter.RunServer(socket, init); err != nil {
		log.Fatal(err)
	}
}

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

	// run server
	go RunServer(socket)

	// sleep forever
	select {}
}
