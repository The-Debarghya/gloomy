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
	warningLog  *log.Logger
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
		infoLog:    log.New(os.Stderr, initText+labelInfo, flags),
		warningLog: log.New(os.Stderr, initText+labelWarn, flags),
		errorLog:   log.New(os.Stderr, initText+labelErr, flags),
		fatalLog:   log.New(os.Stderr, initText+labelFatal, flags),
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
		infoLog:    log.New(io.MultiWriter(iLogs...), labelInfo, flags),
		warningLog: log.New(io.MultiWriter(wLogs...), labelWarn, flags),
		errorLog:   log.New(io.MultiWriter(eLogs...), labelErr, flags),
		fatalLog:   log.New(io.MultiWriter(eLogs...), labelFatal, flags),
	}
	for _, w := range []io.Writer{logFd, il, wl, el} {
		c, ok := w.(io.Closer)
		if ok && c != nil {
			l.closers = append(l.closers, c)
		}
	}
	l.initialized = true

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
			l.warningLog.Output(3+depth, txt)
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


