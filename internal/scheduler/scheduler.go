package scheduler

import (
	"sync"
	"time"

	"notification_batch/internal/batch/encb"
	"notification_batch/internal/batch/spending_alert"
	"notification_batch/internal/config"
	"notification_batch/internal/logger"

	"github.com/go-co-op/gocron"
)

var (
	scheduler *gocron.Scheduler
	once      sync.Once
)

// InitScheduler initializes the scheduler and defines the batch jobs.
func InitScheduler(cfgMap map[string]*config.Config) {
	once.Do(func() {
		scheduler = gocron.NewScheduler(time.Local)
		setupBatchJobs(cfgMap)
	})
}

// StartScheduler starts the scheduler in a non-blocking way.
func StartScheduler() {
	if scheduler != nil {
		scheduler.StartAsync()
		logger.AppLogger.Info("Scheduler started.")
	} else {
		logger.AppLogger.Warn("Scheduler not initialized.")
	}
}

// StopScheduler stops the scheduler gracefully.
func StopScheduler() {
	if scheduler != nil {
		scheduler.Stop()
		logger.AppLogger.Info("Scheduler stopped.")
	} else {
		logger.AppLogger.Warn("Scheduler not initialized.")
	}
}

func setupBatchJobs(cfgMap map[string]*config.Config) {
	// Spending Alert Send
	if cfg, ok := cfgMap["spending_alert"]; ok {
		scheduler.Every(1).Day().At(cfg.SpendingAlert.Schedule.SendTime).Do(func() {
			logger.AppLogger.Info("Starting Spending Alert Send Batch (from scheduler)...")
			spending_alert.RunSpendingAlertSendBatch(cfg)
			logger.AppLogger.Info("Spending Alert Send Batch (from scheduler) finished.")
		})
	}

	// Spending Alert Result
	if cfg, ok := cfgMap["spending_alert"]; ok {
		scheduler.Every(1).Day().At(cfg.SpendingAlert.Schedule.ResultTime).Do(func() {
			logger.AppLogger.Info("Starting Spending Alert Result Batch (from scheduler)...")
			spending_alert.RunSpendingAlertResultBatch(cfg)
			logger.AppLogger.Info("Spending Alert Result Batch (from scheduler) finished.")
		})
	}

	// e-NCB Send
	if cfg, ok := cfgMap["encb"]; ok {
		scheduler.Every(1).Day().At(cfg.ENCB.Schedule.SendTime).Do(func() {
			logger.AppLogger.Info("Starting e-NCB Send Batch (from scheduler)...")
			encb.RunENCBSendBatch(cfg)
			logger.AppLogger.Info("e-NCB Send Batch (from scheduler) finished.")
		})
	}

	// e-NCB Result
	if cfg, ok := cfgMap["encb"]; ok {
		scheduler.Every(1).Day().At(cfg.ENCB.Schedule.ResultTime).Do(func() {
			logger.AppLogger.Info("Starting e-NCB Result Batch (from scheduler)...")
			encb.RunENCBResultBatch(cfg)
			logger.AppLogger.Info("e-NCB Result Batch (from scheduler) finished.")
		})
	}
}
