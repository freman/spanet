package flags

import (
	"errors"
	"net"
)

type IP struct {
	set     bool
	value   net.IP
	Default net.IP
}

func (i *IP) Set(x string) error {
	i.value = net.ParseIP(x)
	if i.value.String() != x {
		return errors.New("failed to parse IP")
	}
	i.set = true
	return nil
}

func (i *IP) String() string {
	return i.IP().String()
}

func (i *IP) IP() net.IP {
	if i.set {
		return i.value
	}
	return i.Default
}

func (i *IP) IsSet() bool {
	return i.set
}
