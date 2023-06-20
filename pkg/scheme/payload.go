package scheme

type PayloadRequestSolution struct {
	Solution uint `json:"solution"`
}

type PayloadResponseSolution struct {
	Quote Quote `json:"quote"`
}

type PayloadResponseChallenge struct {
	Target uint8 `json:"target"`
}
