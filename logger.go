package maildoor

import "log"

// default logger for the package
var defaultLogger Logger = logger(1)

// Logger interface defines the minimum set of methods
// that a logger should satisfy to be used by the library.
type Logger interface {
	// Log a message at the Info level.
	Info(args ...interface{})

	// Log a formatted message at the Info level.
	Infof(format string, args ...interface{})

	// Log a message at the Error level.
	Error(args ...interface{})

	// Log a formatted message at the Error level.
	Errorf(format string, args ...interface{})
}

type logger int

func (l logger) Info(args ...interface{}) {
	log.Print("level=info ", args)
}

func (l logger) Infof(format string, args ...interface{}) {
	log.Printf("level=info "+format, args)
}

func (l logger) Error(args ...interface{}) {
	log.Print("level=info ", args)
}

func (l logger) Errorf(format string, args ...interface{}) {
	log.Printf("level=info "+format, args)
}
