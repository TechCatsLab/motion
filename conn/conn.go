package conn

type Conn interface {
	Dial() error
	OutPut(string) ([]byte, error)
	Close() error
}
