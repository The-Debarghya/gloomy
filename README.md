# Gloomy 📜
Gloomy is a simple cross platform Go logging library for Linux(mainly), and
macOS, it can log to the Linux/macOS syslog, and an io.Writer.(similar to tee)

## Usage 💻

Set up the default gloomy to log the system log (syslog) and a
file, include a flag to turn up verbosity:

```bash
cd examples && go run main.go
```

The Init function returns a gloomy so you can setup multiple instances if you
wish, only the first call to Init will set the default gloomy:

```go
lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
if err != nil {
  gloomy.Fatalf("Failed to open log file: %v", err)
}
defer lf.Close()

// Log to system log and a log file, Info logs don't write to stdout.
loggerOne := gloomy.Init("GloomyExample", false, true, lf)
defer loggerOne.Close()
// Don't to system log or a log file, Info logs write to stdout..
loggerTwo := gloomy.Init("GloomyExample", true, false, ioutil.Discard)
defer loggerTwo.Close()

loggerOne.Info("This will log to the log file and the system log")
loggerTwo.Info("This will only log to stdout")
gloomy.Info("This is the same as using loggerOne")

```

## Custom Format 📑

| Code                                 | Example                                                  |
|--------------------------------------|----------------------------------------------------------|
| `gloomy.SetFlags(log.Ldate)`         | [ERROR]: 2018/11/11 Error running Foobar: message          |
| `gloomy.SetFlags(log.Ltime)`         | [ERROR]: 09:42:45 Error running Foobar: message            |
| `gloomy.SetFlags(log.Lmicroseconds)` | [ERROR]: 09:42:50.776015 Error running Foobar: message     |
| `gloomy.SetFlags(log.Llongfile)`     | [ERROR]: /src/main.go:31: Error running Foobar: message    |
| `gloomy.SetFlags(log.Lshortfile)`    | [ERROR]: main.go:31: Error running Foobar: message         |
| `gloomy.SetFlags(log.LUTC)`          | [ERROR]: Error running Foobar: message                     |
| `gloomy.SetFlags(log.LstdFlags)`     | [ERROR]: 2018/11/11 09:43:12 Error running Foobar: message |

```go
func main() {
    lf, err := os.OpenFile(logPath, …, 0660)
    defer gloomy.Init("foo", *verbose, true, lf).Close()
    gloomy.SetFlags(log.LstdFlags)
}
```
