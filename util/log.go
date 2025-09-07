package util

import (
	"fmt"
	"github.com/fatih/color"
	"log"
)

var (
	fgRed    = color.New(color.FgRed)
	fgYellow = color.New(color.FgYellow)
	fgGreen  = color.New(color.FgGreen)
)

func init() {
	fgRed.EnableColor()
	fgYellow.EnableColor()
	fgGreen.EnableColor()
}

func LogError(format string, v ...any) {
	colorStr := fgRed.Sprintf(format, v...)
	log.Println(colorStr)
}

func LogSuccess(format string, v ...any) {
	colorStr := fgGreen.Sprintf(format, v...)
	log.Println(colorStr)
}
func LogWarn(format string, v ...any) {
	colorStr := fgYellow.Sprintf(format, v...)
	log.Println(colorStr)
}

func Log(format string, v ...any) {
	log.Println(fmt.Sprintf(format, v...))
}

func LogFatal(format string, v ...any) {
	log.Fatalln(fmt.Sprintf(format, v...))
}
