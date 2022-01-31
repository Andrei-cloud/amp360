package amp360

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type BulkResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Failed  interface{} `json:"failed"`
	Updated interface{} `json:"updated"`
}
