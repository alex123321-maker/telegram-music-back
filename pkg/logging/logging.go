package logging

import (
	"log"
	config "telegram-music/config/miniapp"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type simpleLogger struct{}

func (l *simpleLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}
func (l *simpleLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}
func (l *simpleLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}
func (l *simpleLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf("[FATAL] "+format, args...)
}

// NewLogger возвращает простой логгер или можно подключить logrus/zap
func NewLogger(cfg *config.Config) Logger {
	return &simpleLogger{}
}
