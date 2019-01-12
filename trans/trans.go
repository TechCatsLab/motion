package trans

type Trans interface {
	Init() error
	SetFrom(string)
	SetDes(string)
	Run() error
	Close() error
}
