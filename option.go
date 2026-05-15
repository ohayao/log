package log

type Option func(*Logger)

func WithMinLevel(minLevel int) Option {
	return func(l *Logger) {
		l.level = minLevel
	}
}

func WithShortName(enable bool) Option {
	return func(l *Logger) {
		l.shortName = enable
	}
}

func WithTimeStyle(style FLAG_TIME) Option {
	return func(l *Logger) {
		l.flagTime = style
	}
}

func WithColor(enable bool) Option {
	return func(l *Logger) {
		l.enableColor = enable
	}
}
