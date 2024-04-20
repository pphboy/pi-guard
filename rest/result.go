package rest

type SourceResult struct {
	Data any    `json:"data"`
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}
