package patterns

type Connection struct {
	address    string
	port       string
	tpProtocol string
	speedLimit string
}

// ConnectionBuilder - network connection builder
type ConnectionBuilder interface {
	setAddress()
	setPort()
	setTpProtocol()
	setSpeedLimit()
}

// UDPBuilder - udp connection builder
type UDPBuilder struct {
	address    string
	port       string
	tpProtocol string
	speedLimit string
}

func (b *UDPBuilder) setAddress() {
	b.address = "address"
}

func (b *UDPBuilder) setPort() {
	b.port = "port"
}

func (b *UDPBuilder) setTpProtocol() {
	b.tpProtocol = "udp"
}

func (b *UDPBuilder) setSpeedLimit() {
	b.speedLimit = "limit"
}
