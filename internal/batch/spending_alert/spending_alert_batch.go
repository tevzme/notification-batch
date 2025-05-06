package spending_alert

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"notification_batch/internal/config"
	"notification_batch/internal/ftp"
	"notification_batch/internal/logger"
	"notification_batch/internal/util"
)

// Run Spending Alert Send Batch process.
func RunSpendingAlertSendBatch(cfg *config.Config) {
	logger.AppLogger.Info("Starting Spending Alert Send Batch...")
	defer logger.AppLogger.Info("Spending Alert Send Batch finished.")

	ftpClient, err := ftp.NewClient(ftp.Config{
		Host:     cfg.SpendingAlert.FTP.Host,
		User:     cfg.SpendingAlert.FTP.User,
		Password: cfg.SpendingAlert.FTP.Password,
	})
	if err != nil {
		logger.AppLogger.Sugar().Errorf("Failed to create FTP client for Spending Alert: %v", err)
		return
	}
	defer ftpClient.Close()

	remotePath := cfg.SpendingAlert.FTP.RemotePathSend
	localDir := cfg.SpendingAlert.FTP.LocalPath

	if _, err := os.Stat(localDir); os.IsNotExist(err) {
		if err := os.MkdirAll(localDir, 0755); err != nil {
			logger.AppLogger.Sugar().Errorf("Failed to create local directory '%s': %v", localDir, err)
			return
		}
	}

	files, err := ftpClient.ListFiles(remotePath)
	if err != nil {
		logger.AppLogger.Sugar().Errorf("Failed to list files on FTP '%s': %v", remotePath, err)
		return
	}

	for _, file := range files {
		localFilePath, err := ftpClient.DownloadFile(filepath.Join(remotePath, file), localDir)
		if err != nil {
			logger.AppLogger.Sugar().Errorf("Failed to download file '%s': %v", file, err)
			continue
		}
		logger.AppLogger.Sugar().Infof("Downloaded file '%s' to '%s'", file, localFilePath)

		results, err := ProcessSpendingAlertFile(cfg, localFilePath)
		if err != nil {
			logger.AppLogger.Sugar().Errorf("Failed to process file '%s': %v", localFilePath, err)
			continue
		}

		if len(results) > 0 {
			resultFileName := fmt.Sprintf("%s_%s.txt", cfg.SpendingAlert.ResultPrefix, time.Now().Format("20060102"))
			resultFilePath := filepath.Join(cfg.SpendingAlert.FTP.LocalPath, resultFileName)
			err = util.WriteResultToFile(resultFilePath, results)
			if err != nil {
				logger.AppLogger.Sugar().Errorf("Failed to write result to file '%s': %v", resultFilePath, err)
				continue
			}
			logger.AppLogger.Sugar().Infof("Wrote result to file '%s'", resultFilePath)

			remoteResultPath := cfg.SpendingAlert.FTP.RemotePathResult
			err = ftpClient.UploadFile(resultFilePath, filepath.Join(remoteResultPath, resultFileName))
			if err != nil {
				logger.AppLogger.Sugar().Errorf("Failed to upload result file '%s' to '%s': %v", resultFilePath, remoteResultPath, err)
			} else {
				logger.AppLogger.Sugar().Infof("Uploaded result file '%s' to '%s'", resultFilePath, remoteResultPath)
				// Optionally delete the local file after successful upload
				os.Remove(localFilePath)
				os.Remove(resultFilePath)
			}
		} else {
			// Optionally delete the local file even if no results were processed
			os.Remove(localFilePath)
		}
	}
}

// Run Spending Alert Result Batch process.
func RunSpendingAlertResultBatch(cfg *config.Config) {
	logger.AppLogger.Info("Starting Spending Alert Result Batch...")
	defer logger.AppLogger.Info("Spending Alert Result Batch finished.")
	// Add logic for Spending Alert Result Batch
}
