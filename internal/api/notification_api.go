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
)

// NotificationClient handles communication with the Notification API.
type NotificationClient struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewNotificationClient creates a new NotificationClient.
func NewNotificationClient(cfg *config.Config) *NotificationClient {
	return &NotificationClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.APIEndpoints.Timeout * time.Second,
		},
	}
}

// SendNotification sends a notification request.
func (c *NotificationClient) SendNotification(request model.NotificationRequest) (*model.NotificationResponse, error) {
	apiURL := c.cfg.APIEndpoints.SendNotification

	logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Calling Send Notification API - Request: %+v, URL: %s", request, apiURL))

	jsonData, err := json.Marshal(request)
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
		logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Send Notification API - Failed Response (Error: %v), URL: %s", err, apiURL))
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

	response := &model.NotificationResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Send Notification API - Failed to decode response: %v, URL: %s", err, apiURL))
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	logger.ApiLogger(c.cfg.LogPath, c.cfg.APILogPrefix, fmt.Sprintf("Send Notification API - Successful Response: %+v, URL: %s", response, apiURL))
	return response, nil
}
