package api

type Action uint

const (
	CHALLENGE Action = iota
	SOLUTION
)

var actionStr = []string{
	"challenge",
	"solution",
}

func (a Action) String() string {
	return actionStr[a]
}
