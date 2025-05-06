package model

// NotificationRequest defines the request structure for the Send Notification API.
type NotificationRequest struct {
	Usertoken      string `json:"usertoken"`
	Topiccode      string `json:"topiccode"`
	TitleTH        string `json:"title_th,omitempty"`
	MessageTH      string `json:"message_th,omitempty"`
	TitleEN        string `json:"title_en,omitempty"`
	MessageEN      string `json:"message_en,omitempty"`
	TitleinboxTH   string `json:"titleinbox_th,omitempty"`
	MessageinboxTH string `json:"messageinbox_th,omitempty"`
	TitleinboxEN   string `json:"titleinbox_en,omitempty"`
	MessageinboxEN string `json:"messageinbox_en,omitempty"`
}

// NotificationResponse defines the response structure from the Send Notification API.
type NotificationResponse struct {
	ResponseID      string `json:"ResponseID"`
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
}
