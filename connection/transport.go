package connection

type Transport interface {

	// Status
	IsConnected() bool

	// Control
	Open() error
	Close() error

	// Send/receive
	Send(cmd string) error
	SendReceive(cmd string) (string, error)
}
