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

func NewFileLogger(file string, maxSize int64, opts ...Option) (*Logger, error) {
	handler, err := newFileHandler(file, maxSize)
	if err != nil {
		return nil, err
	}
	return New(handler, opts...), nil
}

func NewFileRotateLogger(dir, fileName string, maxAgeHours, hoursInterval int, opts ...Option) (*Logger, error) {
	handler, err := newFileRotateHandler(dir, fileName, maxAgeHours, hoursInterval)
	if err != nil {
		return nil, err
	}
	return New(handler, opts...), nil
}
