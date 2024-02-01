package config

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Logger struct {
	*log.Logger
}

func NewLogger(path string) *Logger {
	baseLogger := log.New()
	baseLogger.SetLevel(log.DebugLevel)
	baseLogger.SetFormatter(&log.JSONFormatter{})
	logFile := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	baseLogger.SetOutput(multiWriter)
	return &Logger{
		Logger: baseLogger,
	}
}

func (l Logger) Info(message string, fields map[string]interface{}) {
	l.WithFields(fields).Info(message)
}

func (l Logger) Warn(message string, fields map[string]interface{}) {
	l.WithFields(fields).Warn(message)
}

func (l Logger) Error(message string, fields map[string]interface{}) {
	l.WithFields(fields).Error(message)
}

func (l Logger) Debug(message string, fields map[string]interface{}) {
	l.WithFields(fields).Debug(message)
}

func (l Logger) Fatal(message string, fields map[string]interface{}) {
	l.WithFields(fields).Fatal(message)
}

func (l Logger) Panic(message string, fields map[string]interface{}) {
	l.WithFields(fields).Panic(message)
}

func (l Logger) Infof(message string, args ...interface{}) {
	l.Infoln(message)
}

func (l Logger) Warnf(message string, args ...interface{}) {
	l.Warnln(message)
}

func (l Logger) Errorf(message string, args ...interface{}) {
	l.Errorln(message)
}

func (l Logger) Debugf(message string, args ...interface{}) {
	l.Debugln(message)
}

func (l Logger) Fatalf(message string, args ...interface{}) {
	l.Fatalln(message)
}

func (l Logger) Panicf(message string, args ...interface{}) {
	l.Panicln(message)
}
