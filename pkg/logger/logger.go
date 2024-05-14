package logger

import "github.com/sirupsen/logrus"

type LogLevel int

const (
	DebugLevel LogLevel = iota + 1
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	Infof(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Fatalf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
}

type AppLogger struct {
	Log   *logrus.Logger
	Level LogLevel
}

func NewLogger(level LogLevel) Logger {
	logger := logrus.New()
	return &AppLogger{Log: logger, Level: level}
}

func (a *AppLogger) setLogLevel() {
	switch a.Level {
	case DebugLevel:
		a.Log.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		a.Log.SetLevel(logrus.InfoLevel)
	case WarningLevel:
		a.Log.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		a.Log.SetLevel(logrus.ErrorLevel)
	case FatalLevel:
		a.Log.SetLevel(logrus.FatalLevel)
	}
}

func (a *AppLogger) Infof(message string, args ...interface{}) {
	a.setLogLevel()
	a.Log.Infof(message, args...)
}

func (a *AppLogger) Warningf(message string, args ...interface{}) {
	a.setLogLevel()
	a.Log.Warnf(message, args...)
}

func (a *AppLogger) Errorf(message string, args ...interface{}) {
	a.setLogLevel()
	a.Log.Errorf(message, args...)
}

func (a *AppLogger) Fatalf(message string, args ...interface{}) {
	a.setLogLevel()
	a.Log.Fatalf(message, args...)
}

func (a *AppLogger) Debugf(message string, args ...interface{}) {
	a.setLogLevel()
	a.Log.Debugf(message, args...)
}
