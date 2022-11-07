package internal

import (
	"fmt"
	"log"
)

// ILogger is standard logger interface
type ILogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type stdLogger struct{}

func NewStdLogger() ILogger {
	log.SetPrefix(fmt.Sprintf("[%s]:   ", ToolName))
	log.SetFlags(log.Lmsgprefix)
	return &stdLogger{}
}

func (*stdLogger) Fatal(v ...interface{})                 { log.Fatal(v...) }
func (*stdLogger) Fatalf(format string, v ...interface{}) { log.Fatalf(format, v...) }
func (*stdLogger) Print(v ...interface{})                 { log.Print(v...) }
func (*stdLogger) Println(v ...interface{})               { log.Println(v...) }
func (*stdLogger) Printf(format string, v ...interface{}) { log.Printf(format, v...) }
