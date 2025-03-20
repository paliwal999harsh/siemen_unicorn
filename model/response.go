package model

type RequestAcceptedResponse struct {
	ReqId UnicornRequestId `json:"reqId"`
}

type ApiResponse struct {
	Msg string `json:"msg"`
}
