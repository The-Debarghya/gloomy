package gologger

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoggingBeforeInit(t *testing.T) {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stderr = w

	init_logger()

	info := "info log"
	warning := "warning log"
	errL := "error log"
	fatal := "fatal log"

	Info(info)
	Warning(warning)
	Error(errL)

	defaultLogger.output(sFatal, 0, fatal)

	w.Close()
	os.Stderr = old

	var b bytes.Buffer
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b.Write(scanner.Bytes())
	}

	out := b.String()

	for _, txt := range []string{info, warning, errL, fatal} {
		if !strings.Contains(out, txt) {
			t.Errorf("Log output %q doesn't match text: %q", out, txt)
		}
	}
}

func TestInit(t *testing.T) {
	var buf1 bytes.Buffer
	l1 := Init("test1", false, false, &buf1)
	if !reflect.DeepEqual(l1, defaultLogger) {
		t.Fatal("defaultLogger does not match logger returned by Init")
	}

	// Subsequent runs of Init shouldn't change defaultLogger.
	var buf2 bytes.Buffer
	l2 := Init("test2", false, false, &buf2)
	if !reflect.DeepEqual(l1, defaultLogger) {
		t.Error("defaultLogger shouldn't have changed")
	}

	// Check log output.
	l1.Info("logger #1")
	l2.Info("logger #2")
	defaultLogger.Info("default logger")

	tests := []struct {
		out  string
		want int
	}{
		{buf1.String(), 2},
		{buf2.String(), 1},
	}

	for i, tt := range tests {
		got := len(strings.Split(strings.TrimSpace(tt.out), "\n"))
		if got != tt.want {
			t.Errorf("Logger %d mismatched no. of lines, expected %d, got %d", i+1, tt.want, got)
		}
	}
}