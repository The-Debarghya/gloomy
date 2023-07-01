package gologger

import (
	"io"
	"log/syslog"
)

func setup(src string) (io.Writer, io.Writer, io.Writer, error) {
	const facility = syslog.LOG_USER
	il, err := syslog.New(facility|syslog.LOG_NOTICE, src)
	if err != nil {
		return nil, nil, nil, err
	}
	wl, err := syslog.New(facility|syslog.LOG_WARNING, src)
	if err != nil {
		return nil, nil, nil, err
	}
	el, err := syslog.New(facility|syslog.LOG_ERR, src)
	if err != nil {
		return nil, nil, nil, err
	}
	return il, wl, el, nil
}
