package rest

type RequestParam struct {
	QueryText string `json:"queryText"`
}

type RequestUriId struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type RequestUriNum struct {
	Num int `json:"num" uri:"num" binding:"required"`
}
