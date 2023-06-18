package client

import "blitzarx1/wisdom-fort/server/service/quotes"

type payloadRequestSolution struct {
	Solution uint64 `json:"solution"`
}

type payloadResponseSolution struct {
	Quote quotes.Quote `json:"quote"`
}

type payloadChallenge struct {
	Target uint8 `json:"target"`
}
