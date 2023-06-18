package service

import "blitzarx1/wisdom-fort/server/service/quotes"

type payloadRequestSolution struct {
	Solution solution `json:"solution"`
}

type payloadResponseSolution struct {
	Quote quotes.Quote `json:"quote"`
}

type payloadChallenge struct {
	Target difficulty `json:"target"`
}
