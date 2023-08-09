package gologger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// severity of logs
type severity int

// Level describes the level of verbosity for info messages when using Verbose logging
type Level int

// Logger represents a logging object, multiple Loggers can be used
// simultaneously even if they are using the same writers.
type Logger struct {
	infoLog     *log.Logger
	warnLog     *log.Logger
	errorLog    *log.Logger
	fatalLog    *log.Logger
	closers     []io.Closer
	initialized bool
	level       Level
}

type Verbose struct {
	enabled bool
	logger  *Logger
}

const (
	sInfo severity = iota
	sWarn
	sError
	sFatal
)

// severity labels
const (
	labelInfo  = "[INFO]: "
	labelWarn  = "[WARN]: "
	labelErr   = "[ERROR]: "
	labelFatal = "[FATAL]: "
)

const (
	flags    = log.Ldate | log.Lmicroseconds | log.Lshortfile
	initText = "[ERROR]: Logging before logger initiated!\n"
)

var (
	logLock       sync.Mutex
	defaultLogger *Logger
)

// reset default logger for tests to reset environment
func init_logger() {
	defaultLogger = &Logger{
		infoLog:  log.New(os.Stderr, initText+labelInfo, flags),
		warnLog:  log.New(os.Stderr, initText+labelWarn, flags),
		errorLog: log.New(os.Stderr, initText+labelErr, flags),
		fatalLog: log.New(os.Stderr, initText+labelFatal, flags),
	}
}

func init() {
	init_logger()
}

/*
Init initializes logging and should be called in main().
Default log functions can be called before Init(),
but log output will only go to stderr (along with a warning).
The first call to Init populates the default logger and returns the
generated logger, subsequent calls to Init will only return the generated
logger. If the logFd passed in also satisfies io.Closer, logFd.Close will be called
when closing the logger.
*/
func Init(name string, verbose, systemLog bool, logFd io.Writer) *Logger {
	var il, wl, el io.Writer
	var syslogErr error
	if systemLog {
		il, wl, el, syslogErr = setup(name)
	}

	iLogs := []io.Writer{logFd}
	wLogs := []io.Writer{logFd}
	eLogs := []io.Writer{logFd}
	if il != nil {
		iLogs = append(iLogs, il)
	}
	if wl != nil {
		wLogs = append(wLogs, wl)
	}
	if el != nil {
		eLogs = append(eLogs, el)
	}

	eLogs = append(eLogs, os.Stderr)
	if verbose {
		iLogs = append(iLogs, os.Stdout)
		wLogs = append(wLogs, os.Stdout)
	}

	l := Logger{
		infoLog:  log.New(io.MultiWriter(iLogs...), labelInfo, flags),
		warnLog:  log.New(io.MultiWriter(wLogs...), labelWarn, flags),
		errorLog: log.New(io.MultiWriter(eLogs...), labelErr, flags),
		fatalLog: log.New(io.MultiWriter(eLogs...), labelFatal, flags),
	}
	for _, w := range []io.Writer{logFd, il, wl, el} {
		c, ok := w.(io.Closer)
		if ok && c != nil {
			l.closers = append(l.closers, c)
		}
	}
	l.initialized = true

	if syslogErr != nil {
		l.Error(syslogErr)
	}

	logLock.Lock()
	defer logLock.Unlock()
	if !defaultLogger.initialized {
		defaultLogger = &l
	}

	return &l

}

// Close closes the default logger.
func Close() {
	defaultLogger.Close()
}

func (l *Logger) output(s severity, depth int, txt string) {
	logLock.Lock()
	defer logLock.Unlock()
	switch s {
	case sInfo:
		l.infoLog.Output(3+depth, txt)
	case sWarn:
		l.warnLog.Output(3+depth, txt)
	case sError:
		l.errorLog.Output(3+depth, txt)
	case sFatal:
		l.fatalLog.Output(3+depth, txt)
	default:
		panic(fmt.Sprintln("[FATAL]: Unrecognized Severity:", s))
	}
}

/*
Close closes all log writers and will flush any cached logs.
Errors from closing the underlying log writers will be printed to stderr.
Once Close is called, all future calls to the logger will panic.
*/
func (l *Logger) Close() {
	logLock.Lock()
	defer logLock.Unlock()

	if !l.initialized {
		return
	}

	for _, c := range l.closers {
		if err := c.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]: Failed to close log %v: %v\n", c, err)
		}
	}
}

/**  INFO LOGS **/

// Info level logs.
// Arguments according to fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.output(sInfo, 0, fmt.Sprint(v...))
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth called with depth 0 is equivalent to Info.
func (l *Logger) InfoDepth(depth int, v ...interface{}) {
	l.output(sInfo, depth, fmt.Sprint(v...))
}

// Newline appended info level logs.
// Arguments according to fmt.Println.
func (l *Logger) Infoln(v ...interface{}) {
	l.output(sInfo, 0, fmt.Sprintln(v...))
}

// Formatted info level logs.
// Arguments according to fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(sInfo, 0, fmt.Sprintf(format, v...))
}

// Info level logs.
// Arguments according to fmt.Print.
func Info(v ...interface{}) {
	defaultLogger.output(sInfo, 0, fmt.Sprint(v...))
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth called with depth 0 is equivalent to Info.
func InfoDepth(depth int, v ...interface{}) {
	defaultLogger.output(sInfo, depth, fmt.Sprint(v...))
}

// Newline appended info level logs.
// Arguments according to fmt.Println.
func Infoln(v ...interface{}) {
	defaultLogger.output(sInfo, 0, fmt.Sprintln(v...))
}

// Formatted info level logs.
// Arguments according to fmt.Printf.
func Infof(format string, v ...interface{}) {
	defaultLogger.output(sInfo, 0, fmt.Sprintf(format, v...))
}

/*  WARNING LOGS  */

// Warning level logs.
// Arguments according to fmt.Print.
func (l *Logger) Warning(v ...interface{}) {
	l.output(sWarn, 0, fmt.Sprint(v...))
}

// WarningDepth acts as Warning but uses depth to determine which call frame to log.
// WarningDepth called with depth 0 is equivalent to Warning.
func (l *Logger) WarningDepth(depth int, v ...interface{}) {
	l.output(sWarn, depth, fmt.Sprint(v...))
}

// Newline appended warning level logs.
// Arguments according to fmt.Println.
func (l *Logger) Warningln(v ...interface{}) {
	l.output(sWarn, 0, fmt.Sprintln(v...))
}

// Formatted warning level logs.
// Arguments according to fmt.Printf.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.output(sWarn, 0, fmt.Sprintf(format, v...))
}

// Warning level logs.
// Arguments according to fmt.Print.
func Warning(v ...interface{}) {
	defaultLogger.output(sWarn, 0, fmt.Sprint(v...))
}

// WarningDepth acts as Warning but uses depth to determine which call frame to log.
// WarningDepth called with depth 0 is equivalent to Warning.
func WarningDepth(depth int, v ...interface{}) {
	defaultLogger.output(sWarn, depth, fmt.Sprint(v...))
}

// Newline appended warning level logs.
// Arguments according to fmt.Println.
func Warningln(v ...interface{}) {
	defaultLogger.output(sWarn, 0, fmt.Sprintln(v...))
}

// Formatted warning level logs.
// Arguments according to fmt.Printf.
func Warningf(format string, v ...interface{}) {
	defaultLogger.output(sWarn, 0, fmt.Sprintf(format, v...))
}

/*  ERROR LOGS  */

// Error level logs.
// Arguments according to fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.output(sError, 0, fmt.Sprint(v...))
}

// ErrorDepth acts as Error but uses depth to determine which call frame to log.
// ErrorDepth called with depth 0 is equivalent to Error.
func (l *Logger) ErrorDepth(depth int, v ...interface{}) {
	l.output(sError, depth, fmt.Sprint(v...))
}

// Newline appended error level logs.
// Arguments according to fmt.Println.
func (l *Logger) Errorln(v ...interface{}) {
	l.output(sError, 0, fmt.Sprintln(v...))
}

// Formatted error level logs.
// Arguments according to fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(sError, 0, fmt.Sprintf(format, v...))
}

// Error level logs.
// Arguments according to fmt.Print.
func Error(v ...interface{}) {
	defaultLogger.output(sError, 0, fmt.Sprint(v...))
}

// ErrorDepth acts as Error but uses depth to determine which call frame to log.
// ErrorDepth called with depth 0 is equivalent to Error.
func ErrorDepth(depth int, v ...interface{}) {
	defaultLogger.output(sError, depth, fmt.Sprint(v...))
}

// Newline appended error level logs.
// Arguments according to fmt.Println.
func Errorln(v ...interface{}) {
	defaultLogger.output(sError, 0, fmt.Sprintln(v...))
}

// Formatted error level logs.
// Arguments according to fmt.Printf.
func Errorf(format string, v ...interface{}) {
	defaultLogger.output(sError, 0, fmt.Sprintf(format, v...))
}

/*  FATAL LOGS  */

// Fatal level logs which terminates with os.Exit(1).
// Arguments according to fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	l.output(sFatal, 0, fmt.Sprint(v...))
	l.Close()
	os.Exit(1)
}

// FatalDepth acts as Fatal but uses depth to determine which call frame to log.
// FatalDepth called with depth 0 is equivalent to Fatal.
func (l *Logger) FatalDepth(depth int, v ...interface{}) {
	l.output(sFatal, depth, fmt.Sprint(v...))
	l.Close()
	os.Exit(1)
}

// Newline appended Fatal level logs which terminates with os.Exit(1).
// Arguments according to fmt.Println.
func (l *Logger) Fatalln(v ...interface{}) {
	l.output(sFatal, 0, fmt.Sprintln(v...))
	l.Close()
	os.Exit(1)
}

// Formatted Fatal level logs terminates with os.Exit(1).
// Arguments according to fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.output(sFatal, 0, fmt.Sprintf(format, v...))
	l.Close()
	os.Exit(1)
}

// Fatal level logs which terminates with os.Exit(1).
// Arguments according to fmt.Print.
func Fatal(v ...interface{}) {
	defaultLogger.output(sFatal, 0, fmt.Sprint(v...))
	defaultLogger.Close()
	os.Exit(1)
}

// FatalDepth acts as Fatal but uses depth to determine which call frame to log.
// FatalDepth called with depth 0 is equivalent to Fatal.
func FatalDepth(depth int, v ...interface{}) {
	defaultLogger.output(sFatal, depth, fmt.Sprint(v...))
	defaultLogger.Close()
	os.Exit(1)
}

// Newline appended Fatal level logs which terminates with os.Exit(1).
// Arguments according to fmt.Println.
func Fatalln(v ...interface{}) {
	defaultLogger.output(sFatal, 0, fmt.Sprintln(v...))
	defaultLogger.Close()
	os.Exit(1)
}

// Formatted Fatal level logs terminates with os.Exit(1).
// Arguments according to fmt.Printf.
func Fatalf(format string, v ...interface{}) {
	defaultLogger.output(sFatal, 0, fmt.Sprintf(format, v...))
	defaultLogger.Close()
	os.Exit(1)
}


// Set the logger verbosity level for verbose info logging.
func (l *Logger) SetLevel(lvl Level) {
	l.level = lvl
	l.output(sInfo, 0, fmt.Sprintf("Info verbosity set to %d", lvl))
}

// Sets the verbosity level for verbose info logging for the
// default logger.
func SetLevel(lvl Level) {
	defaultLogger.SetLevel(lvl)
}

/* VERBOSE LOGGING */

// Verbosity generates a log record depends on the setting of the Level; or none default.
// It uses the specified logger.
func (l *Logger) Verbosity(lvl Level) Verbose {
	return Verbose{
		enabled: l.level >= lvl,
		logger:  l,
	}
}

// Verbosity generates a log record, depends on the setting of the Level; or none
// by default using the default logger.
func Verbosity(lvl Level) Verbose {
	return defaultLogger.Verbosity(lvl)
}

// Info is equivalent to Info function, when verbosity(v) is enabled.
func (v Verbose) Info(args ...interface{}) {
	if v.enabled {
		v.logger.output(sInfo, 0, fmt.Sprint(args...))
	}
}

// Infoln is equivalent to Infoln function, when verbosity(v) is enabled.
// See the docs of Verbosity for usage.
func (v Verbose) Infoln(args ...interface{}) {
	if v.enabled {
		v.logger.output(sInfo, 0, fmt.Sprintln(args...))
	}
}

// Infof is equivalent to Infof function, when verbosity(v) is enabled.
// See the docs of Verbosity for usage.
func (v Verbose) Infof(format string, args ...interface{}) {
	if v.enabled {
		v.logger.output(sInfo, 0, fmt.Sprintf(format, args...))
	}
}

// Sets the output flags for the logger.
func SetFlags(flag int) {
	defaultLogger.infoLog.SetFlags(flag)
	defaultLogger.warnLog.SetFlags(flag)
	defaultLogger.errorLog.SetFlags(flag)
	defaultLogger.fatalLog.SetFlags(flag)
}

