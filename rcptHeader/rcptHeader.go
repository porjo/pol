package rcptHeader

import (
	"log"
	"net/textproto"
	"strings"

	//"github.com/phalaaxx/milter"
	"github.com/porjo/milter"
)

type RcptMilter struct {
	milter.Milter

	Email  string
	Logger *log.Logger

	from       string
	rcptTo     string
	subject    string
	rcptInToCc bool
}

// Header parses message headers one by one
func (b *RcptMilter) Header(name, value string, m *milter.Modifier) (milter.Response, error) {
	// check if bogofilter has been run on the message already
	if name == "X-Pol" {
		// X-Polsity header is present, accept immediately
		return milter.RespAccept, nil
	}
	if name == "To" || name == "Cc" {
		if strings.Contains(strings.ToLower(value), b.Email) {
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
func (b *RcptMilter) RcptTo(rcptTo string, m *milter.Modifier) (milter.Response, error) {
	b.rcptTo = rcptTo
	return milter.RespContinue, nil

}

// MailFrom is called on envelope from address
func (b *RcptMilter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	// save from address for later reference
	b.from = from
	return milter.RespContinue, nil
}

// Headers is called after the last of message headers
func (b *RcptMilter) Headers(headers textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	if b.rcptTo == b.Email && !b.rcptInToCc {
		//TODO custom message here?
		b.Logger.Printf("REJECT message from %s to %s, subject '%s'\n", b.from, b.rcptTo, b.subject)
		return milter.RespReject, nil
	}
	b.Logger.Printf("OK message from %s to %s, subject '%s'\n", b.from, b.rcptTo, b.subject)
	return milter.RespContinue, nil
}

func (b *RcptMilter) Body(m *milter.Modifier) (milter.Response, error) {
	m.AddHeader("X-Pol", "1")
	return milter.RespAccept, nil
}
