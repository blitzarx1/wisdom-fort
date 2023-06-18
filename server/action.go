package server

type action uint

const (
	CHALLENGE action = iota
	SOLUTION
)

var actionStr = []string{
	"challenge",
	"solution",
}

func (a action) String() string {
	return actionStr[a]
}
