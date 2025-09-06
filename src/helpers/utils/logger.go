package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ColorizedLogger struct {
	useColor bool
}

func FormatDate(t time.Time) string {
	return t.Format("15:04:05 01/02/2006")
}

var colorCodes = map[string]func(a ...any) string{
	"info":    color.New(color.FgBlue).SprintFunc(),
	"verbose": color.New(color.FgCyan).SprintFunc(),
	"warn":    color.New(color.FgYellow).SprintFunc(),
	"error":   color.New(color.FgRed).SprintFunc(),
	"http":    color.New(color.FgMagenta).SprintFunc(),
	"silly":   color.New(color.FgGreen).SprintFunc(),
}

func (l *ColorizedLogger) log(level, message string) {
	timestamp := FormatDate(time.Now())
	levelUpper := strings.ToUpper(level)

	colorFunc, exists := colorCodes[level]
	if !exists {
		colorFunc = color.New(color.Reset).SprintFunc()
	}

	var logMessage string
	if l.useColor {
		logMessage = fmt.Sprintf("%s - [%s]: %s\n", colorFunc(timestamp), colorFunc(levelUpper), colorFunc(message))
	} else {
		logMessage = fmt.Sprintf("[%s]: [%s] | %s\n", timestamp, levelUpper, message)
	}

	if _, err := os.Stdout.WriteString(logMessage); err != nil {
		log.Printf("Failed To Write Stdout: %v", err)
		return
	}
}

func NewColorizedLogger(useColor bool) *ColorizedLogger {
	return &ColorizedLogger{useColor: useColor}
}

func (l *ColorizedLogger) Info(message string)    { l.log("info", message) }
func (l *ColorizedLogger) Verbose(message string) { l.log("verbose", message) }
func (l *ColorizedLogger) Warn(message string)    { l.log("warn", message) }
func (l *ColorizedLogger) Http(message string)    { l.log("http", message) }
func (l *ColorizedLogger) Silly(message string)   { l.log("silly", message) }
func (l *ColorizedLogger) Error(message string)   { l.log("error", message) }
