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

// ErrorLogger adalah instance logger global untuk error
var ErrorLogger *log.Logger

// dailyLogWriter adalah io.Writer kustom yang mengelola rotasi file harian.
type dailyLogWriter struct {
	logDir      string // Direktori tempat menyimpan log
	currentDate string // Tanggal saat ini (format YYYY-MM-DD)
	currentFile *os.File
	mutex       sync.Mutex
}

// newDailyLogWriter membuat instance writer harian yang baru.
func newDailyLogWriter(logDir string) (*dailyLogWriter, error) {
	writer := &dailyLogWriter{
		logDir: logDir,
	}
	// Buka file log untuk hari ini saat pertama kali dibuat
	if err := writer.reopenFile(); err != nil {
		return nil, err
	}
	return writer, nil
}

// reopenFile (dipanggil di dalam mutex) untuk menutup file lama dan membuka file baru.
func (w *dailyLogWriter) reopenFile() error {
	// 1. Tentukan tanggal hari ini
	today := time.Now().Format("2006-01-02")

	// 2. Tutup file lama jika ada
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	// 3. Buat path file baru
	logFileName := fmt.Sprintf("%s-error.log", today)
	logFilePath := filepath.Join(w.logDir, logFileName)

	// 4. Buka file baru
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("gagal membuka file log %s: %w", logFilePath, err)
	}

	// 5. Perbarui state internal writer
	w.currentFile = file
	w.currentDate = today

	// Cetak ke stdout bahwa kita telah merotasi file
	log.Println("Error logger sekarang menulis ke:", logFilePath)
	return nil
}

// Write mengimplementasikan interface io.Writer
func (w *dailyLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 1. Cek apakah tanggal telah berubah
	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		// Tanggal telah berubah, lakukan rotasi file
		if err := w.reopenFile(); err != nil {
			// Gagal rotasi. Coba tulis ke stdout sebagai fallback
			log.Printf("FATAL: Gagal merotasi file log: %v. Pesan log: %s", err, string(p))
			return os.Stdout.Write(p)
		}
	}

	// 2. Tulis log ke file yang saat ini aktif
	return w.currentFile.Write(p)
}

// InitErrorLogger menginisialisasi logger error global
func InitErrorLogger() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Gagal membuat direktori log: %v", err)
	}

	// 1. Buat writer kustom kita
	dailyWriter, err := newDailyLogWriter(logDir)
	if err != nil {
		log.Fatalf("Gagal menginisialisasi daily log writer: %v", err)
	}

	// 2. Buat MultiWriter untuk menulis ke file (via dailyWriter) DAN stdout
	mw := io.MultiWriter(os.Stdout, dailyWriter)

	// 3. Inisialisasi ErrorLogger global
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	// Pesan "writing to..." sekarang ditangani oleh reopenFile()
	log.Println("Daily error logger berhasil diinisialisasi.")
}
