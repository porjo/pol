package pol

import (
	"net"
	"net/textproto"

	//"github.com/phalaaxx/milter"
	"github.com/porjo/milter"
)

type Milter struct {
	milter.Milter

	Milters []milter.Milter
}

// Header parses message headers one by one
func (b *Milter) Header(name, value string, mod *milter.Modifier) (milter.Response, error) {

	for _, m := range b.Milters {
		r, err := m.Header(name, value, mod)
		if err != nil {
			return nil, err
		}
		if r != milter.RespContinue {
			return r, nil
		}
	}

	return milter.RespContinue, nil
}

// RcptTo is called to process filters on envelope TO address
//   supress with NoRcptTo
func (b *Milter) RcptTo(rcptTo string, mod *milter.Modifier) (milter.Response, error) {
	for _, m := range b.Milters {
		r, err := m.RcptTo(rcptTo, mod)
		if err != nil {
			return nil, err
		}
		if r != milter.RespContinue {
			return r, nil
		}
	}

	return milter.RespContinue, nil
}

func (b *Milter) Connect(host string, family string, port uint16, addr net.IP, mod *milter.Modifier) (milter.Response, error) {

	for _, m := range b.Milters {
		r, err := m.Connect(host, family, port, addr, mod)
		if err != nil {
			return nil, err
		}
		if r != milter.RespContinue {
			return r, nil
		}
	}

	return milter.RespContinue, nil
}

// MailFrom is called on envelope from address
func (b *Milter) MailFrom(from string, mod *milter.Modifier) (milter.Response, error) {
	for _, m := range b.Milters {
		r, err := m.MailFrom(from, mod)
		if err != nil {
			return nil, err
		}
		if r != milter.RespContinue {
			return r, nil
		}
	}

	return milter.RespContinue, nil
}

// Headers is called after the last of message headers
func (b *Milter) Headers(headers textproto.MIMEHeader, mod *milter.Modifier) (milter.Response, error) {
	for _, m := range b.Milters {
		r, err := m.Headers(headers, mod)
		if err != nil {
			return nil, err
		}
		if r != milter.RespContinue {
			return r, nil
		}
	}

	return milter.RespContinue, nil
}

func (b *Milter) Body(mod *milter.Modifier) (milter.Response, error) {
	for _, m := range b.Milters {
		r, err := m.Body(mod)
		if err != nil {
			return nil, err
		}
		if r != milter.RespContinue {
			return r, nil
		}
	}

	return milter.RespContinue, nil
}
