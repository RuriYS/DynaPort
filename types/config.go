package types

type Config struct {
	Server struct {
		Host 				string		`yaml:"host"`
		Port				uint16		`yaml:"port"`
		TTL					string		`yaml:"ttl"`
		AllowedIPs			[]string	`yaml:"allowed_ips"`
		AllowedPorts		[]uint16	`yaml:"allowed_ports"`
	}
	Client struct {
		Host				string 		`yaml:"host"`
		BroadcastInterval	string		`yaml:"broadcast_interval"`
		Timeout				string 		`yaml:"timeout"`
		Ports 				[]uint16	`yaml:"ports"`
		WhitelistMode		bool		`yaml:"whitelist_mode"`
	}
}
