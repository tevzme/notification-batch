package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// AppLogger is the global logger instance for the application.
var AppLogger *zap.Logger

// InitLogger initializes the application logger.
func InitLogger(logPath string) error {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			return err
		}
	}

	appLogFile := filepath.Join(logPath, "app.log")

	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   appLogFile,
			MaxSize:    100, // megabytes
			MaxBackups: 5,
			MaxAge:     7, // days
			Compress:   true,
		}),
		zapcore.AddSync(os.Stdout),
	)

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zap.InfoLevel,
	)

	AppLogger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(AppLogger)

	AppLogger.Info("Logger initialized successfully", zap.String("log_path", logPath))
	return nil
}

// ApiLogger logs API requests and responses to a separate file with timestamp.
func ApiLogger(logPath, logPrefix, message string) {
	apiLogFile := filepath.Join(logPath, logPrefix+".log")
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	apiLogger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   apiLogFile,
			MaxSize:    100, // megabytes
			MaxBackups: 5,
			MaxAge:     7, // days
			Compress:   true,
		}),
		zap.InfoLevel,
	))
	defer apiLogger.Sync()
	apiLogger.Info(message)
}
