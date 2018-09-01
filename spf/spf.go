package spf

import (
	"log"
	"net"
	"net/textproto"

	//"github.com/phalaaxx/milter"
	"github.com/porjo/milter"
	"github.com/porjo/spf"
)

type Milter struct {
	clientIP net.IP
	from     string
	Logger   *log.Logger

	milter.Milter
}

// Header parses message headers one by one
func (b *Milter) Header(name, value string, m *milter.Modifier) (milter.Response, error) {

	return milter.RespContinue, nil
}

// RcptTo is called to process filters on envelope TO address
//   supress with NoRcptTo
func (b *Milter) RcptTo(rcptTo string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil

}

// MailFrom is called on envelope from address
func (b *Milter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	b.from = from
	return milter.RespContinue, nil
}

func (b *Milter) Connect(host string, family string, port uint16, addr net.IP, mod *milter.Modifier) (milter.Response, error) {

	b.clientIP = addr
	return milter.RespContinue, nil
}

// Headers is called after the last of message headers
func (b *Milter) Headers(headers textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {

	result, err := spf.SPFTest(b.clientIP.String(), b.from)
	if err != nil {
		return nil, err
	}

	if result == spf.Fail {
		b.Logger.Printf("REJECT message, SPF Check: from IP %s with email %s\n", b.clientIP, b.from)
		return milter.RespReject, nil
	}
	return milter.RespContinue, nil
}

func (b *Milter) Body(m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}
