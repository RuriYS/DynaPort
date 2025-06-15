package types

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

type Allocation struct {
	Protocol Protocol
	Port     uint16
}
