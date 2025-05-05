package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

type tag struct {
	Name  string
	Color string
}

var (
	info      = tag{"INFO", "\033[34m"}
	warn      = tag{"WARN", "\033[33m"}
	err       = tag{"ERROR", "\033[31m"}
	debug     = tag{"DEBUG", "\033[90m"}
	updater   = tag{"UPDATER", "\033[35m"}
	installer = tag{"INSTALLER", "\033[36m"}
	deployer  = tag{"DEPLOYER", "\033[95m"}
	system    = tag{"SYSTEM", "\033[32m"}
)

const reset = "\033[0m"

func log(t tag, message string, args ...interface{}) {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		file = "???"
	}
	shortFile := filepath.Base(file)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(message, args...)

	fmt.Printf("[%s] %s[%s]%s [%s] %s\n", timestamp, t.Color, t.Name, reset, shortFile, formatted)
}

func Info(message string, args ...interface{})  { log(info, message, args...) }
func Warn(message string, args ...interface{})  { log(warn, message, args...) }
func Debug(message string, args ...interface{}) { log(debug, message, args...) }
func Error(message string, args ...interface{}) error {
	log(err, message, args...)
	return fmt.Errorf(message, args...)
}

func Updater(message string, args ...interface{})   { log(updater, message, args...) }
func Installer(message string, args ...interface{}) { log(installer, message, args...) }
func Deployer(message string, args ...interface{})  { log(deployer, message, args...) }
func System(message string, args ...interface{})    { log(system, message, args...) }
