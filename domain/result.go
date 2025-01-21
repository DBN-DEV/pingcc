package domain

type PingResult struct {
	RttMicros uint32
	IsTimeout bool
	Tag       map[string]string
}

type TcpPingResult struct {
	RttMicros uint32
	IsTimeout bool
	Tag       map[string]string
}
