package encb

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

// Run e-NCB Send Batch process.
func RunENCBSendBatch(cfg *config.Config) {
	logger.AppLogger.Info("Starting e-NCB Send Batch...")
	defer logger.AppLogger.Info("e-NCB Send Batch finished.")

	ftpClient, err := ftp.NewClient(ftp.Config{
		Host:     cfg.ENCB.FTP.Host,
		User:     cfg.ENCB.FTP.User,
		Password: cfg.ENCB.FTP.Password,
	})
	if err != nil {
		logger.AppLogger.Sugar().Errorf("Failed to create FTP client for e-NCB: %v", err)
		return
	}
	defer ftpClient.Close()

	remotePath := cfg.ENCB.FTP.RemotePathSend
	localDir := cfg.ENCB.FTP.LocalPath

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

		results, err := ProcessENCBFile(cfg, localFilePath)
		if err != nil {
			logger.AppLogger.Sugar().Errorf("Failed to process file '%s': %v", localFilePath, err)
			continue
		}

		if len(results) > 0 {
			resultFileName := fmt.Sprintf("%s_%s.txt", cfg.ENCB.ResultPrefix, time.Now().Format("20060102"))
			resultFilePath := filepath.Join(cfg.ENCB.FTP.LocalPath, resultFileName)
			err = util.WriteResultToFile(resultFilePath, results)
			if err != nil {
				logger.AppLogger.Sugar().Errorf("Failed to write result to file '%s': %v", resultFilePath, err)
				continue
			}
			logger.AppLogger.Sugar().Infof("Wrote result to file '%s'", resultFilePath)

			remoteResultPath := cfg.ENCB.FTP.RemotePathResult
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

// Run e-NCB Result Batch process.
func RunENCBResultBatch(cfg *config.Config) {
	logger.AppLogger.Info("Starting e-NCB Result Batch...")
	defer logger.AppLogger.Info("e-NCB Result Batch finished.")
	// Add logic for e-NCB Result Batch
}
