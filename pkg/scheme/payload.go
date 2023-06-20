package scheme

type PayloadRequestSolution struct {
	Solution uint64 `json:"solution"`
}

type PayloadResponseSolution struct {
	Quote Quote `json:"quote"`
}

type PayloadResponseChallenge struct {
	Target uint8 `json:"target"`
}
