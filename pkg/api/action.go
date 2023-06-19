package api

type Action uint

const (
	ActionChallenge Action = iota
	ActionSolution
)

var actionStr = []string{
	"challenge",
	"solution",
}

func (a Action) String() string {
	return actionStr[a]
}
