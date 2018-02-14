package logtool

import (
	"EventFlow/common/config"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type logdata struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Step      string    `json:"step"`
	Mode      string    `json:"mode"`
	ErrorLine string    `json:"errorline,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type sysLogger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	fatal *log.Logger
}

const (
	flagLog   = log.Ldate | log.Ltime
	flagerror = log.Ldate | log.Ltime
	tagDebug  = "debug"
	tagInfo   = "info"
	tagWarn   = "warn"
	tagError  = "error"
	tagFatal  = "fatal"
)

var logFilePath string
var logFilePathWithoutExt string
var currentFileIO *os.File
var fileMutex *sync.Mutex
var logChannel chan logdata
var logger sysLogger

const (
	logChannelBufferSize          = 50
	checkFileSizeIntervalInMinute = 5
)

func init() {

	currentePath, _ := os.Getwd()
	logFolder := fmt.Sprintf("%s/%s", currentePath, config.Config.Log.Path)
	logFilePathWithoutExt = fmt.Sprintf("%s/logfile", logFolder)
	logFilePath = fmt.Sprintf("%s.log", logFilePathWithoutExt)

	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		err := os.Mkdir(logFolder, os.ModePerm)

		if err != nil {
			log.Printf("[tool][log] create log directory failed: %v", err)
			return
		}
	}

	logger = sysLogger{
		debug: log.New(nil, "", flagLog),
		info:  log.New(nil, "", flagLog),
		warn:  log.New(nil, "", flagLog),
		err:   log.New(nil, "", flagLog),
		fatal: log.New(nil, "", flagLog)}

	fileMutex = &sync.Mutex{}
	logChannel = make(chan logdata, logChannelBufferSize)

	checkLogFileSize()
	createLogger()
	createCheckLogFileSizeRoutine()
	createLogWriterRoutine()
}

//Debug write debug log
func Debug(step, mode, message string) {
	writeLogData(tagDebug, step, mode, message)
}

//Info write info log
func Info(step, mode, message string) {
	writeLogData(tagInfo, step, mode, message)
}

//Warn write warning log
func Warn(step, mode, message string) {
	writeLogData(tagWarn, step, mode, message)
}

//Error write error log
func Error(step, mode, message string) {
	writeLogData(tagError, step, mode, message)
}

//Fatal write fatal log
func Fatal(step, mode, message string) {
	writeLogData(tagFatal, step, mode, message)
}

func createLogger() {

	logIO, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		log.Printf("[tool][log] open/create log file failed: %v", err)
		return
	}

	currentFileIO = logIO
	logWriter := []io.Writer{logIO}
	errWriter := []io.Writer{os.Stderr, logIO}

	if config.Config.Log.Verbose {
		logWriter = append(logWriter, os.Stdout)
	}

	logger.debug.SetOutput(io.MultiWriter(logWriter...))
	logger.info.SetOutput(io.MultiWriter(logWriter...))
	logger.warn.SetOutput(io.MultiWriter(logWriter...))
	logger.err.SetOutput(io.MultiWriter(errWriter...))
	logger.fatal.SetOutput(io.MultiWriter(errWriter...))
}

func createCheckLogFileSizeRoutine() {
	ticker := time.NewTicker(time.Minute * checkFileSizeIntervalInMinute)

	go func() {
		for {
			<-ticker.C
			checkLogFileSize()
		}
	}()
}

func checkLogFileSize() {

	fileInfo, err := os.Stat(logFilePath)

	if os.IsNotExist(err) {
		log.Print("[tool][log] log file not exist")
		return

	} else if err != nil {
		log.Printf("[tool][log] check log file state failed: %v", err)
		return
	}

	if fileInfo.Size() >= int64(config.Config.Log.MaxSizeKB*1024) {

		currentTime := time.Now()
		newLogFileName := fmt.Sprintf("%s_%s.log", logFilePathWithoutExt, currentTime.Format("20060102150405"))

		fileMutex.Lock()
		defer fileMutex.Unlock()

		//close current log file
		if currentFileIO != nil {
			if err := currentFileIO.Close(); err != nil {
				log.Printf("[tool][log] close current log file failed: %v", err)
				return
			}
		}

		//rename current log file
		err := os.Rename(logFilePath, newLogFileName)

		if err != nil {
			log.Printf("[tool][log] rename current log file failed: %v", err)
		} else {
			//create new log file
			createLogger()
		}
	}
}

func createLogWriterRoutine() {

	go func() {
		for {
			logData := <-logChannel
			fileMutex.Lock()

			var log *log.Logger

			switch logData.Level {
			case tagDebug:
				log = logger.debug
			case tagInfo:
				log = logger.info
			case tagWarn:
				log = logger.warn
			case tagError:
				log = logger.err
			case tagFatal:
				log = logger.fatal
			}

			var logMsg string

			if config.Config.Log.IncludeCaller && (logData.Level == tagError || logData.Level == tagFatal) {
				logMsg = fmt.Sprintf("[%s] [%s][%s] [%s] %s\r\n", logData.Level, logData.Step, logData.Mode, logData.ErrorLine, logData.Message)
			} else {
				logMsg = fmt.Sprintf("[%s] [%s][%s] %s\r\n", logData.Level, logData.Step, logData.Mode, logData.Message)
			}

			if err := log.Output(0, logMsg); err != nil {
				log.Printf("[tool][log] write log failed: %v", err)
			}

			fileMutex.Unlock()
		}
	}()
}

func writeLogData(level, step, mode, message string) {

	logLevel := config.Config.Log.Level

	switch logLevel {
	case tagDebug:
	case tagInfo:
		if level == tagDebug {
			return
		}
	case tagWarn:
		if level == tagDebug || level == tagInfo {
			return
		}
	case tagError:
		if level != tagError && level != tagFatal {
			return
		}
	case tagFatal:
		if level != tagFatal {
			return
		}
	}

	logData := logdata{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Step:      step,
		Mode:      mode}

	if config.Config.Log.IncludeCaller && (level == tagError || level == tagFatal) {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			shortFile := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					shortFile = file[i+1:]
					break
				}
			}
			logData.ErrorLine = fmt.Sprintf("%s:%d", shortFile, line)

		} else {
			logData.ErrorLine = "???:0"
		}
	}

	logChannel <- logData
}
