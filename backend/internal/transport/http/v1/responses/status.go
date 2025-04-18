package responses

type StatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
