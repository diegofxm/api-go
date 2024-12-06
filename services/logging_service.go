package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type LogEntry struct {
	Timestamp  string      `json:"timestamp"`
	Level      string      `json:"level"`
	Message    string      `json:"message"`
	File       string      `json:"file,omitempty"`
	Line       int         `json:"line,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	StackTrace string      `json:"stack_trace,omitempty"`
}

type LoggingService struct {
	mu        sync.Mutex
	file      *os.File
	level     LogLevel
	maxSize   int64 // tamaño máximo del archivo en bytes
	filename  string
	showLine  bool
}

func NewLoggingService(filename string, level LogLevel, maxSize int64) (*LoggingService, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &LoggingService{
		file:     file,
		level:    level,
		maxSize:  maxSize,
		filename: filename,
		showLine: true,
	}, nil
}

func (l *LoggingService) log(level LogLevel, message string, data interface{}) error {
	if level < l.level {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Verificar el tamaño del archivo y rotarlo si es necesario
	if err := l.rotateIfNeeded(); err != nil {
		return err
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     l.getLevelString(level),
		Message:   message,
		Data:      data,
	}

	if l.showLine {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry.File = filepath.Base(file)
			entry.Line = line
		}
	}

	if level >= ERROR {
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		entry.StackTrace = string(buf[:n])
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	if _, err := l.file.Write(append(jsonData, '\n')); err != nil {
		return err
	}

	return nil
}

func (l *LoggingService) rotateIfNeeded() error {
	info, err := l.file.Stat()
	if err != nil {
		return err
	}

	if info.Size() < l.maxSize {
		return nil
	}

	// Cerrar el archivo actual
	if err := l.file.Close(); err != nil {
		return err
	}

	// Renombrar el archivo actual
	backupName := fmt.Sprintf("%s.%s", l.filename, time.Now().Format("2006-01-02-15-04-05"))
	if err := os.Rename(l.filename, backupName); err != nil {
		return err
	}

	// Crear un nuevo archivo
	file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.file = file
	return nil
}

func (l *LoggingService) getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l *LoggingService) Debug(message string, data ...interface{}) {
	l.log(DEBUG, message, data)
}

func (l *LoggingService) Info(message string, data ...interface{}) {
	l.log(INFO, message, data)
}

func (l *LoggingService) Warn(message string, data ...interface{}) {
	l.log(WARN, message, data)
}

func (l *LoggingService) Error(message string, data ...interface{}) {
	l.log(ERROR, message, data)
}

func (l *LoggingService) Fatal(message string, data ...interface{}) {
	l.log(FATAL, message, data)
	os.Exit(1)
}

func (l *LoggingService) SetLevel(level LogLevel) {
	l.level = level
}

func (l *LoggingService) SetShowLine(show bool) {
	l.showLine = show
}

func (l *LoggingService) Close() error {
	return l.file.Close()
}
