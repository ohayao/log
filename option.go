package log

type Option func(*Logger)

func WithMinLevel(minLevel int) Option {
	return func(l *Logger) {
		l.level = minLevel
	}
}

func WithShortName(ok bool) Option {
	return func(l *Logger) {
		l.shortName = ok
	}
}

func WithTimeStyle(style FLAG_TIME) Option {
	return func(l *Logger) {
		l.flagTime = style
	}
}

func WithColor(ok bool) Option {
	return func(l *Logger) {
		l.enableColor = ok
	}
}
