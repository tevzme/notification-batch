package model

// AlertSettingRequest defines the request structure for the Get Alert Setting API.
type AlertSettingRequest struct {
	RequestID string `json:"RequestID"`
	UserToken string `json:"UserToken"`
}

// AlertSettingResponse defines the response structure from the Get Alert Setting API.
type AlertSettingResponse struct {
	ResponseID        string `json:"ResponseID"`
	ResponseCode      string `json:"ResponseCode"`
	ResponseMessage   string `json:"ResponseMessage"`
	UserToken         string `json:"UserToken"`
	SpendingAlertFlag bool   `json:"spending_alert_flag"`
	LastLogin         string `json:"last_login"`
}
