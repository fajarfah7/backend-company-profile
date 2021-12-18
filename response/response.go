package response

// Response will return formated data to the frontend
type Response struct {
	ActionStatus bool                `json:"action_status"`
	Status       int                 `json:"status"`
	Messages     []map[string]string `json:"messages"`
	Data         interface{}         `json:"data"`
}
