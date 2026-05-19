package log

import (
	"os"
)

func New(handler IHandler, opts ...Option) *Logger {
	logger := &Logger{
		handler:     handler,
		enableColor: true,
		shortName:   false,
		flagTime:    FLAG_TIME_DATETIME,
		level:       LV_DEBUG,
		pool:        poolNew(),
	}
	for _, opt := range opts {
		opt(logger)
	}

	switch handler.(type) {
	case *fileHandler, *fileRotateHandler:
		logger.enableColor = false
	}
	return logger
}

func NewTerminalLogger(file *os.File, opts ...Option) *Logger {
	return New(newTerminalHandler(file), opts...)
}

func NewFileLogger(filePath string, maxSize int64, opts ...Option) (*Logger, error) {
	handler, err := newFileHandler(filePath, maxSize)
	if err != nil {
		return nil, err
	}
	return New(handler, opts...), nil
}

func NewFileRotateLogger(filePath string, hoursInterval, maxAgeHours int, opts ...Option) (*Logger, error) {
	handler, err := newFileRotateHandler(filePath, hoursInterval, maxAgeHours)
	if err != nil {
		return nil, err
	}
	return New(handler, opts...), nil
}
