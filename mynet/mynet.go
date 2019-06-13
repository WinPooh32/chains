package mynet

type Connector interface {
	Connect(addr string) error
}
