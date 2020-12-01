package main

type taskRequest struct {
	NextExecTime int64 `json:"exec_time"`
}

func (tr *taskRequest) Validate() error {
	return nil
}

type restErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
