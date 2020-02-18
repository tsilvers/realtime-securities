package markets

const (
	StatusError StatusType = iota
	StatusPre
	StatusOpen
	StatusClosed
)

type StatusType int
