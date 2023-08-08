# gologger #
GoLogger is a simple cross platform Go logging library for Linux(mainly), and
macOS, it can log to the Linux/macOS syslog, and an io.Writer.

## Usage ##

Set up the default gologger to log the system log (syslog) and a
file, include a flag to turn up verbosity:

```go
import (
  "flag"
  "os"

  "github.com/The-Debarghya/gologger"
)

const logPath = "/some/location/example.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
  flag.Parse()

  lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
  if err != nil {
    gologger.Fatalf("Failed to open log file: %v", err)
  }
  defer lf.Close()

  defer gologger.Init("LoggerExample", *verbose, true, lf).Close()

  gologger.Info("I'm about to do something!")
  if err := doSomething(); err != nil {
    gologger.Errorf("Error running doSomething: %v", err)
  }
}
```

The Init function returns a gologger so you can setup multiple instances if you
wish, only the first call to Init will set the default gologger:

```go
lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
if err != nil {
  gologger.Fatalf("Failed to open log file: %v", err)
}
defer lf.Close()

// Log to system log and a log file, Info logs don't write to stdout.
loggerOne := gologger.Init("LoggerExample", false, true, lf)
defer loggerOne.Close()
// Don't to system log or a log file, Info logs write to stdout..
loggerTwo := gologger.Init("LoggerExample", true, false, ioutil.Discard)
defer loggerTwo.Close()

loggerOne.Info("This will log to the log file and the system log")
loggerTwo.Info("This will only log to stdout")
gologger.Info("This is the same as using loggerOne")

```

## Custom Format ##

| Code                                 | Example                                                  |
|--------------------------------------|----------------------------------------------------------|
| `gologger.SetFlags(log.Ldate)`         | [ERROR]: 2018/11/11 Error running Foobar: message          |
| `gologger.SetFlags(log.Ltime)`         | [ERROR]: 09:42:45 Error running Foobar: message            |
| `gologger.SetFlags(log.Lmicroseconds)` | [ERROR]: 09:42:50.776015 Error running Foobar: message     |
| `gologger.SetFlags(log.Llongfile)`     | [ERROR]: /src/main.go:31: Error running Foobar: message    |
| `gologger.SetFlags(log.Lshortfile)`    | [ERROR]: main.go:31: Error running Foobar: message         |
| `gologger.SetFlags(log.LUTC)`          | [ERROR]: Error running Foobar: message                     |
| `gologger.SetFlags(log.LstdFlags)`     | [ERROR]: 2018/11/11 09:43:12 Error running Foobar: message |

```go
func main() {
    lf, err := os.OpenFile(logPath, …, 0660)
    defer gologger.Init("foo", *verbose, true, lf).Close()
    gologger.SetFlags(log.LstdFlags)
}
```
