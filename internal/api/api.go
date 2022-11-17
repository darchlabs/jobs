package api

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

// handler's response
type HandlerRes struct {
	Payload    interface{}
	HttpStatus int
	Err        error
}
