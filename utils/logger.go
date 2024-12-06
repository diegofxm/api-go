package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

type JSONLogger struct {
	Writer io.Writer
	logger.Config
}

type LogEntry struct {
	Timestamp   string      `json:"timestamp"`
	Level       string      `json:"level"`
	Message     string      `json:"message"`
	SQL         string      `json:"sql,omitempty"`
	Rows        int64       `json:"rows,omitempty"`
	Duration    string      `json:"duration,omitempty"`
	Error       string      `json:"error,omitempty"`
	Additional  interface{} `json:"additional,omitempty"`
}

func NewJSONLogger() logger.Interface {
	return &JSONLogger{
		Writer: os.Stdout,
		Config: logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		},
	}
}

func (l *JSONLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *JSONLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.log("info", msg, nil, data...)
	}
}

func (l *JSONLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.log("warn", msg, nil, data...)
	}
}

func (l *JSONLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.log("error", msg, nil, data...)
	}
}

func (l *JSONLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "trace",
		SQL:       sql,
		Rows:      rows,
		Duration:  elapsed.String(),
	}

	if err != nil {
		entry.Error = err.Error()
	}

	l.writeJSON(entry)
}

func (l *JSONLogger) log(level, msg string, err error, data ...interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   fmt.Sprintf(msg, data...),
	}

	if err != nil {
		entry.Error = err.Error()
	}

	if len(data) > 0 {
		entry.Additional = data[0]
	}

	l.writeJSON(entry)
}

func (l *JSONLogger) writeJSON(entry LogEntry) {
	json, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}
	
	l.Writer.Write(append(json, '\n'))
}
