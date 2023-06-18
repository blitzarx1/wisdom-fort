package service

import "blitzarx1/wisdom-fort/server/service/quotes"

type PayloadSolutionRequest struct {
	Solution solution `json:"solution"`
}

type PayloadSolutionResponse struct {
	Quote quotes.Quote `json:"quote"`
}
