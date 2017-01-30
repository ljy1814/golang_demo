package mylog

import (
	"bytes"
	"demo/display"
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

var count = 0

func TestConsoleLogger(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()

	logger := &ConsoleLogger{}
	logger.SetDefaultInvokingNumber()
	expectedInvokingNumber := Position(1)
	currentInvokingNumber := logger.GetInvokingNumber()
	if currentInvokingNumber != expectedInvokingNumber {
		t.Errorf("The current invoking number %d should be %d\n", currentInvokingNumber, expectedInvokingNumber)
	}
	testLogger(t, logger)
}

func TestLogManager(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			t.Errorf("Fatal Error: %s\n", err)
		}
	}()
	loggers := &LogManager{loggers: []Logger{&ConsoleLogger{position: 2}}}
	//	NewLogger(loggers)
	for _, logger := range loggers.loggers {
		testLogger(t, logger)
	}
}

func testLogger(t *testing.T, logger Logger) {
	var format string
	var content string
	var logContent string

	format = ""
	logContent = "<Error>"
	content = logger.Error(logContent)
	checkContent(t, getErrorLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = " haha Errorf"
	content = logger.Errorf(format, logContent)
	checkContent(t, getErrorLogTag(), content, format, logContent)

	format = ""
	logContent = "<Errorln>"
	content = logger.Errorf(format, logContent)
	checkContent(t, getErrorLogTag(), content, format, logContent)

	format = ""
	logContent = "<Info>"
	content = logger.Info(logContent)
	checkContent(t, getInfoLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Infof"
	content = logger.Infof(format, logContent)
	checkContent(t, getInfoLogTag(), content, format, logContent)

	format = ""
	logContent = "<Infoln>"
	content = logger.Infof(format, logContent)
	checkContent(t, getInfoLogTag(), content, format, logContent)

	format = ""
	logContent = "<Panic>"
	content = logger.Panic(logContent)
	checkContent(t, getPanicLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Panicf"
	content = logger.Panicf(format, logContent)
	checkContent(t, getPanicLogTag(), content, format, logContent)

	format = ""
	logContent = "<Panicln>"
	content = logger.Panicf(format, logContent)
	checkContent(t, getPanicLogTag(), content, format, logContent)

	format = ""
	logContent = "<Fatal>"
	content = logger.Fatal(logContent)
	checkContent(t, getFatalLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Fatalf"
	content = logger.Fatalf(format, logContent)
	checkContent(t, getFatalLogTag(), content, format, logContent)

	format = ""
	logContent = "<Fatalln>"
	content = logger.Fatalf(format, logContent)
	checkContent(t, getFatalLogTag(), content, format, logContent)

	format = ""
	logContent = "<Warn>"
	content = logger.Warn(logContent)
	checkContent(t, getWarnLogTag(), content, format, logContent)

	format = "<%s>"
	logContent = "Warnf"
	content = logger.Warnf(format, logContent)
	checkContent(t, getWarnLogTag(), content, format, logContent)

	format = ""
	logContent = "<Warnln>"
	content = logger.Warnf(format, logContent)
	checkContent(t, getWarnLogTag(), content, format, logContent)
}

func checkContent(t *testing.T, logTag LogTag, content string, format string, logContents ...interface{}) {
	var prefixBuffer bytes.Buffer
	prefixBuffer.WriteString(logTag.Prefix())
	prefixBuffer.WriteString(" demo/mylog.testLogger : (logger_test.go:")
	prefix := prefixBuffer.String()
	//	fmt.Printf("prefix : %s\n", prefix)
	var suffixBuffer bytes.Buffer
	suffixBuffer.WriteString(") - ")
	if len(format) == 0 {
		suffixBuffer.WriteString(fmt.Sprint(logContents...))
	} else {
		suffixBuffer.WriteString(fmt.Sprintf(format, logContents...))
	}
	suffix := suffixBuffer.String()
	//	fmt.Printf("suffix : %s\n", suffix)
	//	fmt.Printf("content : %s\n", content)
	if !strings.HasPrefix(content, prefix) {
		t.Errorf("The content %q should has prefix %q!", content, prefix)
	}
	if !strings.HasSuffix(content, suffix) {
		t.Errorf("The content %q should has suffix %q!", content, suffix)
	}
}

func TestFileName(t *testing.T) {
	fmt.Println(time.Now())
	fmt.Println(time.Now().Year())
	fmt.Println(time.Now().Month())
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	display.Display("Time", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0])
}

func TestFileError(t *testing.T) {
	var fl FileLogger
	fl.Error("xxx")
	fl.Error("sssssssssssssss")
	fl.Errorf("%s %d", "error format", 7)
}
func TestLogType(t *testing.T) {
	content := "------------------------"
	writeContent("error", &content)
}
