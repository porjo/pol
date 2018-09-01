package spf

import (
	//"fmt"
	"net/textproto"

	//"github.com/phalaaxx/milter"
	"github.com/porjo/milter"
)

type SPFMilter struct {
	milter.Milter
}

// Header parses message headers one by one
func (b *SPFMilter) Header(name, value string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

// RcptTo is called to process filters on envelope TO address
//   supress with NoRcptTo
func (b *SPFMilter) RcptTo(rcptTo string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil

}

// MailFrom is called on envelope from address
func (b *SPFMilter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

// Headers is called after the last of message headers
func (b *SPFMilter) Headers(headers textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (b *SPFMilter) Body(m *milter.Modifier) (milter.Response, error) {
	return milter.RespAccept, nil
}
