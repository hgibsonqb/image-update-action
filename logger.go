package main

import (
	"fmt"
	"github.com/go-logr/logr"
	"os"
	"text/tabwriter"
)

type logSink struct {
	name      string
	keyValues map[string]interface{}
	writer    *tabwriter.Writer
}

var _ logr.LogSink = &logSink{}

func (_ *logSink) Init(info logr.RuntimeInfo) {}

func (_ logSink) Enabled(level int) bool {
	return true
}

func (l logSink) Info(level int, msg string, kvs ...interface{}) {
	fmt.Fprintf(l.writer, "%s\t%s\t", l.name, msg)
	for k, v := range l.keyValues {
		fmt.Fprintf(l.writer, "%s: %+v  ", k, v)
	}
	for i := 0; i < len(kvs); i += 2 {
		fmt.Fprintf(l.writer, "%s: %+v  ", kvs[i], kvs[i+1])
	}
	fmt.Fprintf(l.writer, "\n")
	l.writer.Flush()
}

func (l logSink) Error(err error, msg string, kvs ...interface{}) {
	kvs = append(kvs, "error", err)
	l.Info(0, msg, kvs...)
}

func (l logSink) WithName(name string) logr.LogSink {
	return &logSink{
		name:      l.name + "." + name,
		keyValues: l.keyValues,
	}
}

func (l logSink) WithValues(kvs ...interface{}) logr.LogSink {
	newMap := make(map[string]interface{}, len(l.keyValues)+len(kvs)/2)
	for k, v := range l.keyValues {
		newMap[k] = v
	}
	for i := 0; i < len(kvs); i += 2 {
		newMap[kvs[i].(string)] = kvs[i+1]
	}
	return &logSink{
		name:      l.name,
		keyValues: newMap,
	}
}

func NewLogger() logr.Logger {
	sink := &logSink{
		writer: tabwriter.NewWriter(os.Stdout, 40, 8, 2, '\t', 0),
	}
	return logr.New(sink)
}
