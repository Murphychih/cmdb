package conf

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel

	// InvalidLevel is an invalid value for Level.
	//
	// Core implementations may panic if they see messages of this level.
	InvalidLevel = _maxLevel + 1
)

// 若level为空 则默认设置为 DebugLevel
func (l *log) setLevel() zap.AtomicLevel {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	if l.Level != "" {
		switch l.Level {
		case "debug":
			atomicLevel.SetLevel(zap.DebugLevel)
		case "info":
			atomicLevel.SetLevel(zap.InfoLevel)
		case "warn":
			atomicLevel.SetLevel(zap.WarnLevel)
		case "error":
			atomicLevel.SetLevel(zap.ErrorLevel)
		case "dpanic":
			atomicLevel.SetLevel(zap.DPanicLevel)
		case "panic":
			atomicLevel.SetLevel(zap.PanicLevel)
		case "fatal":
			atomicLevel.SetLevel(zap.FatalLevel)
		}
		return atomicLevel
	} else {
		atomicLevel.SetLevel(zap.DebugLevel)
		return atomicLevel
	}
}
