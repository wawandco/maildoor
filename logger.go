package maildoor

import "log"

var (
	// default logger is a mute logger as we don't
	// want to spam the logs unless explicitly told by
	// the user.
	defaultLogger Logger = muteLogger(1)

	// BasicLogger is a simple logger that prints to stdout
	// using the `log` package.
	BasicLogger = stdOutLogger(2)
)

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

type muteLogger int

func (l muteLogger) Info(args ...interface{})                  {}
func (l muteLogger) Infof(format string, args ...interface{})  {}
func (l muteLogger) Error(args ...interface{})                 {}
func (l muteLogger) Errorf(format string, args ...interface{}) {}

type stdOutLogger int

func (l stdOutLogger) Info(args ...interface{}) {
	log.Print("level=info ", args)
}

func (l stdOutLogger) Infof(format string, args ...interface{}) {
	log.Printf("level=info "+format, args)
}

func (l stdOutLogger) Error(args ...interface{}) {
	log.Print("level=info ", args)
}

func (l stdOutLogger) Errorf(format string, args ...interface{}) {
	log.Printf("level=info "+format, args)
}
