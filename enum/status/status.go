package status

type Status int

const (
	Die Status = iota
	Free
	Alive
	Busy
)
