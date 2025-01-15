package response

type List struct {
	Pagination Pagination `json:"pagination"`
	Data       any        `json:"data"`
}

type Pagination struct {
	Total int64 `json:"total"`
}

type Created struct {
	ID int64 `json:"id"`
}

type Exception struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type ExceptionField struct {
	Field   string `json:"field"`
	Key     string `json:"key"`
	Message string `json:"message"`
}
