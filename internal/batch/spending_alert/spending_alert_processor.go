package spending_alert

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"notification_batch/internal/api"
	"notification_batch/internal/config"
	"notification_batch/internal/logger"
	"notification_batch/internal/model"
	"notification_batch/internal/util"
)

// Define fixed positions and lengths for Spending Alert file fields
const (
	cardNoStart        = 0
	cardNoLength       = 16
	userTokenStart     = 20
	userTokenLength    = 36
	originalDateStart  = 60
	originalDateLength = 10
	originalTimeStart  = 71
	originalTimeLength = 8
)

// ProcessSpendingAlertFile reads and processes each line of the Spending Alert file.
func ProcessSpendingAlertFile(cfg *config.Config, filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %v", filePath, err)
	}
	defer file.Close()

	var results []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < userTokenStart+userTokenLength {
			logger.AppLogger.Sugar().Warnf("Skipping line due to insufficient length: '%s'", line)
			continue
		}

		cardNo := strings.TrimSpace(util.SafeSubstring(line, cardNoStart, cardNoLength))
		userToken := strings.TrimSpace(util.SafeSubstring(line, userTokenStart, userTokenLength))
		originalDateStr := strings.TrimSpace(util.SafeSubstring(line, originalDateStart, originalDateLength))
		originalTimeStr := strings.TrimSpace(util.SafeSubstring(line, originalTimeStart, originalTimeLength))

		alertSettingClient := api.NewAlertSettingClient(cfg)
		alertSettingResponse, err := alertSettingClient.GetAlertSetting(userToken)
		if err != nil {
			logger.AppLogger.Sugar().Errorf("Failed to call Get Alert Setting API for user token '%s': %v", userToken, err)
			continue
		}

		if alertSettingResponse.SpendingAlertFlag && isLastLoginWithin90Days(alertSettingResponse.LastLogin) {
			notificationRequest := model.NotificationRequest{
				Usertoken:      userToken,
				Topiccode:      "test",
				TitleTH:        "แจ้งเตือนการใช้จ่าย",
				MessageTH:      fmt.Sprintf("คุณมีการใช้จ่ายผ่านบัตร %s เมื่อวันที่ %s เวลา %s", cardNo, originalDateStr, originalTimeStr),
				TitleEN:        "Spending Alert",
				MessageEN:      fmt.Sprintf("You have a spending transaction with card %s on %s at %s", cardNo, originalDateStr, originalTimeStr),
				TitleinboxTH:   "แจ้งเตือนการใช้จ่าย",
				MessageinboxTH: fmt.Sprintf("คุณมีการใช้จ่ายผ่านบัตร %s เมื่อวันที่ %s เวลา %s", cardNo, originalDateStr, originalTimeStr),
				TitleinboxEN:   "Spending Alert",
				MessageinboxEN: fmt.Sprintf("You have a spending transaction with card %s on %s at %s", cardNo, originalDateStr, originalTimeStr),
			}

			notificationClient := api.NewNotificationClient(cfg)
			notificationResponse, err := notificationClient.SendNotification(notificationRequest)
			if err != nil {
				logger.AppLogger.Sugar().Errorf("Failed to call Send Notification API for user token '%s': %v", userToken, err)
			} else {
				logger.AppLogger.Sugar().Infof("Notification sent for user token '%s', Response: %+v", userToken, notificationResponse)
				results = append(results, fmt.Sprintf("%s,%s,%s,%s,%s,%+v", cardNo, userToken, originalDateStr, originalTimeStr, userToken, notificationResponse))
			}
		} else {
			logger.AppLogger.Sugar().Infof("Spending Alert not triggered for user token '%s' (Flag: %t, LastLogin within 90 days: %t)", userToken, alertSettingResponse.SpendingAlertFlag, isLastLoginWithin90Days(alertSettingResponse.LastLogin))
			results = append(results, fmt.Sprintf("%s,%s,%s,%s,%s,Not Triggered", cardNo, userToken, originalDateStr, originalTimeStr, userToken))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filePath, err)
	}

	return results, nil
}

func isLastLoginWithin90Days(lastLogin string) bool {
	if lastLogin == "" {
		return false
	}
	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, lastLogin)
	if err != nil {
		logger.AppLogger.Sugar().Warnf("Failed to parse last login time '%s': %v", lastLogin, err)
		return false
	}
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	return parsedTime.After(ninetyDaysAgo)
}
