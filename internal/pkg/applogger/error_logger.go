package applogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ErrorLogger is the global logger instance for errors
var ErrorLogger *log.Logger

// dailyLogWriter is a custom io.Writer that manages daily file rotation.
type dailyLogWriter struct {
	logDir      string // Directory to store logs
	currentDate string // Current date (YYYY-MM-DD format)
	currentFile *os.File
	mutex       sync.Mutex
}

// newDailyLogWriter creates a new daily writer instance.
func newDailyLogWriter(logDir string) (*dailyLogWriter, error) {
	writer := &dailyLogWriter{
		logDir: logDir,
	}
	// Open the log file for today upon initial creation
	if err := writer.reopenFile(); err != nil {
		return nil, err
	}
	return writer, nil
}

// reopenFile (called within mutex) closes the old file and opens a new one.
func (w *dailyLogWriter) reopenFile() error {
	today := time.Now().Format("2006-01-02")

	if w.currentFile != nil {
		w.currentFile.Close()
	}

	logFileName := fmt.Sprintf("%s-error.log", today)
	logFilePath := filepath.Join(w.logDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", logFilePath, err)
	}

	w.currentFile = file
	w.currentDate = today

	// Use a more distinct log message for rotation/initial open
	log.Printf("🔄 Error logger now writing to: %s", logFilePath)
	return nil
}

// Write implements the io.Writer interface
func (w *dailyLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		if err := w.reopenFile(); err != nil {
			log.Printf("🚨 FATAL: Failed to rotate error log file: %v. Log message: %s", err, string(p))
			return os.Stdout.Write(p)
		}
	}

	return w.currentFile.Write(p)
}

// InitErrorLogger initializes the global error logger
func InitErrorLogger() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("🚨 CRITICAL: Failed to create log directory '%s': %v", logDir, err)
	}

	dailyWriter, err := newDailyLogWriter(logDir)
	if err != nil {
		log.Fatalf("🚨 CRITICAL: Failed to initialize daily log writer: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, dailyWriter)

	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	// Removed the final success log here, as reopenFile already logs the target file.
	// log.Println("✅ Daily error logger initialized successfully.")
}
