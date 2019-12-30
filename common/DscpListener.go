package common

import (
	"net"
)

type DscpListener struct {
	nl   net.Listener
	dscp int
}

func (l DscpListener) Accept() (c net.Conn, err error) {
	nc, err := l.nl.Accept()
	if err != nil {
		return nil, err
	}
	err = SetTos(nc, l.dscp*4)
	if FatalErrorln(err, "Error setting DSCP") != nil {
		return nil, err
	}
	return nc, nil
}

func (l DscpListener) Close() error {
	return l.nl.Close()
}

func (l DscpListener) Addr() net.Addr {
	return l.nl.Addr()
}
