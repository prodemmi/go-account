package swagger

type SuccessJSONResult struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

type ErrorJSONResult struct {
	Status string `json:"status" example:"error"`
	Error  string `json:"error" example:"example error message."`
}
