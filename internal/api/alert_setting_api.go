package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"notification_batch/internal/config"
	"notification_batch/internal/logger"
	"notification_batch/internal/model"
	"notification_batch/internal/util"
)

// AlertSettingClient handles communication with the Alert Setting API.
type AlertSettingClient struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewAlertSettingClient creates a new AlertSettingClient.
func NewAlertSettingClient(cfg *config.Config) *AlertSettingClient {
	return &AlertSettingClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.APIEndpoints.Timeout * time.Second,
		},
	}
}

// GetAlertSetting retrieves alert settings for a given user token.
func (c *AlertSettingClient) GetAlertSetting(userToken string) (*model.AlertSettingResponse, error) {
	apiURL := c.cfg.APIEndpoints.GetAlertSetting
	requestID := util.GenerateRequestID()

	requestBody := map[string]string{
		"RequestID": requestID,
		"UserToken": userToken,
	}

	logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Calling Get Alert Setting API - Request: %+v, URL: %s", requestBody, apiURL))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Get Alert Setting API - Failed Response (Error: %v), URL: %s", err, apiURL))
		return nil, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Get Alert Setting API - Non-OK Status: %d, URL: %s, Error reading body: %v", resp.StatusCode, apiURL, err))
		} else {
			logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Get Alert Setting API - Non-OK Status: %d, URL: %s, Body: %s", resp.StatusCode, apiURL, string(errBodyBytes)))
		}
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	response := &model.AlertSettingResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Get Alert Setting API - Failed to decode response: %v, URL: %s", err, apiURL))
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Get Alert Setting API - Successful Response: %+v, URL: %s", response, apiURL))
	return response, nil
}
