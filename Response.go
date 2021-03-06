package web

type Response struct {
	Code *int `json:"code,omitempty"`
	Success *bool `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

type Pager struct {
	Count int `json:"count,omitempty"`
	Data []interface{} `json:"data,omitempty"`
}

type PagerResponse struct {
	Response
	Pager
}

func (r *Response) SetCode(code int) *Response {
	r.Code = &code
	return r
}

func (r *Response) SetSuccess(success bool) *Response {
	r.Success = &success
	return r
}
