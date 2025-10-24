package applogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ErrorLogger adalah instance logger global untuk error
var ErrorLogger *log.Logger

// InitErrorLogger menginisialisasi logger error untuk menulis ke file
// di direktori ./logs dengan nama file berstempel waktu.
func InitErrorLogger() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Buat nama file log berdasarkan waktu start aplikasi
	// Format: YYYY-MM-DD_HH-MM-SS_errors.log
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("%s_errors.log", timestamp)
	logFilePath := filepath.Join(logDir, logFileName)

	// Buka file log (atau buat jika tidak ada)
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Buat MultiWriter untuk menulis ke file DAN stdout
	mw := io.MultiWriter(os.Stdout, file)

	// Inisialisasi ErrorLogger global
	// Flag menyertakan: Tanggal, Waktu (UTC), dan File/Baris kode yang memanggil
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	log.Println("Error logger initialized. Writing errors to:", logFilePath)
}
