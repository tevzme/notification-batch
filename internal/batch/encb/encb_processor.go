package encb

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"notification_batch/internal/api"
	"notification_batch/internal/config"
	"notification_batch/internal/logger"
	"notification_batch/internal/model"
	"notification_batch/internal/util"
)

// Define fixed positions and lengths for e-NCB file fields
const (
	encbUserTokenStart       = 0
	encbUserTokenLength      = 36
	encbTitleInboxTHStart    = 37
	encbTitleInboxTHLength   = 100
	encbMessageInboxTHStart  = 138
	encbMessageInboxTHLength = 200
	encbTitleInboxENStart    = 339
	encbTitleInboxENLength   = 100
	encbMessageInboxENStart  = 440
	encbMessageInboxENLength = 200
)

// ProcessENCBFile reads and processes each line of the e-NCB file.
func ProcessENCBFile(cfg *config.Config, filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %v", filePath, err)
	}
	defer file.Close()

	var results []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < encbUserTokenStart+encbUserTokenLength {
			logger.AppLogger.Sugar().Warnf("Skipping line due to insufficient length: '%s'", line)
			continue
		}

		userToken := strings.TrimSpace(util.SafeSubstring(line, encbUserTokenStart, encbUserTokenLength))
		titleInboxTH := strings.TrimSpace(util.SafeSubstring(line, encbTitleInboxTHStart, encbTitleInboxTHLength))
		messageInboxTH := strings.TrimSpace(util.SafeSubstring(line, encbMessageInboxTHStart, encbMessageInboxTHLength))
		titleInboxEN := strings.TrimSpace(util.SafeSubstring(line, encbTitleInboxENStart, encbTitleInboxENLength))
		messageInboxEN := strings.TrimSpace(util.SafeSubstring(line, encbMessageInboxENStart, encbMessageInboxENLength))

		notificationRequest := model.NotificationRequest{
			Usertoken:      userToken,
			Topiccode:      "test",
			TitleTH:        titleInboxTH,
			MessageTH:      messageInboxTH,
			TitleinboxTH:   titleInboxTH,
			MessageinboxTH: messageInboxTH,
			TitleEN:        titleInboxEN,
			MessageEN:      messageInboxEN,
			TitleinboxEN:   titleInboxEN,
			MessageinboxEN: messageInboxEN,
		}

		notificationClient := api.NewNotificationClient(cfg)
		notificationResponse, err := notificationClient.SendNotification(notificationRequest)
		if err != nil {
			logger.AppLogger.Sugar().Errorf("Failed to call Send Notification API (TH) for user token '%s': %v", userToken, err)
		} else {
			logger.AppLogger.Sugar().Infof("Notification (TH) sent for user token '%s', Response: %+v", userToken, notificationResponse)
			results = append(results, fmt.Sprintf("%s,%s,%s,%s,%+v", userToken, titleInboxTH, messageInboxTH, "TH", notificationResponse))
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filePath, err)
	}

	return results, nil
}
