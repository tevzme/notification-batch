package ftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"notification_batch/internal/logger"

	"github.com/jlaffaye/ftp"
)

// Config holds the FTP client configuration.
type Config struct {
	Host     string
	User     string
	Password string
}

// Client wraps the ftp.ServerConn.
type Client struct {
	conn   *ftp.ServerConn
	config Config
}

// NewClient creates a new FTP client connection.
func NewClient(cfg Config) (*Client, error) {
	conn, err := ftp.Dial(cfg.Host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FTP server '%s': %v", cfg.Host, err)
	}

	err = conn.Login(cfg.User, cfg.Password)
	if err != nil {
		conn.Quit()
		return nil, fmt.Errorf("failed to login to FTP server '%s' with user '%s': %v", cfg.Host, cfg.User, err)
	}

	logger.AppLogger.Sugar().Infof("Connected and logged in to FTP server '%s' as user '%s'", cfg.Host, cfg.User)

	return &Client{
		conn:   conn,
		config: cfg,
	}, nil
}

// Close closes the FTP connection.
func (c *Client) Close() {
	if c.conn != nil {
		err := c.conn.Quit()
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			logger.AppLogger.Sugar().Errorf("Failed to close FTP connection to '%s': %v", c.config.Host, err)
		} else {
			logger.AppLogger.Sugar().Infof("Closed FTP connection to '%s'", c.config.Host)
		}
		c.conn = nil
	}
}

// ListFiles lists files in the specified remote directory.
func (c *Client) ListFiles(remotePath string) ([]string, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("FTP connection is not established")
	}

	entries, err := c.conn.List(remotePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in '%s': %v", remotePath, err)
	}

	var files []string
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFile {
			files = append(files, entry.Name)
		}
	}

	logger.AppLogger.Sugar().Infof("Found %d files in '%s'", len(files), remotePath)
	return files, nil
}

// DownloadFile downloads a file from the remote path to the local directory.
func (c *Client) DownloadFile(remotePath, localDir string) (localFilePath string, err error) {
	if c.conn == nil {
		return "", fmt.Errorf("FTP connection is not established")
	}

	resp, err := c.conn.Retr(remotePath)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file '%s': %v", remotePath, err)
	}
	defer resp.Close()

	fileName := filepath.Base(remotePath)
	localFilePath = filepath.Join(localDir, fileName)

	outFile, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file '%s': %v", localFilePath, err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp)
	if err != nil {
		os.Remove(localFilePath)
		return "", fmt.Errorf("failed to copy data from '%s' to '%s': %v", remotePath, localFilePath, err)
	}

	logger.AppLogger.Sugar().Infof("Downloaded '%s' from FTP to '%s'", remotePath, localFilePath)
	return localFilePath, nil
}

// UploadFile uploads a local file to the remote path.
func (c *Client) UploadFile(localPath, remotePath string) error {
	if c.conn == nil {
		return fmt.Errorf("FTP connection is not established")
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file '%s': %v", localPath, err)
	}
	defer file.Close()

	err = c.conn.Stor(remotePath, file)
	if err != nil {
		return fmt.Errorf("failed to store file '%s' to '%s': %v", localPath, remotePath, err)
	}

	logger.AppLogger.Sugar().Infof("Uploaded '%s' to FTP as '%s'", localPath, remotePath)
	return nil
}
