package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ALog struct {
	writerLv    int
	level       int
	logCnt      int
	logFileName string
	fileSize    int64
	fileStream  *os.File
	info        *log.Logger
	ch          chan func()
}

const (
	InfoLevel = iota
)

func NewLogger(logFileName string, level int, chanSize int, tCnt int, flag int) *ALog {
	logger := new(ALog)
	logger.level = level
	logger.writerLv = InfoLevel
	logger.logCnt = 20
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("OpenFile fail")
		log.Fatal(err)
	}
	logger.fileStream = logFile

	logger.logFileName = logFileName
	logger.fileSize = 100000000
	logger.info = log.New(logFile, "[INFO]", flag)
	logger.ch = make(chan func(), chanSize)
	for i := 0; i < tCnt; i++ {
		go logger.PrintLog()
	}
	return logger
}

func (aLog *ALog) ResetOutput() {
	aLog.fileStream.Close()
	logFile, err := os.OpenFile(aLog.logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("OpenFile fail")
		log.Fatal(err)
	}
	aLog.fileStream = logFile

	aLog.info.SetOutput(logFile)

}

func (aLog *ALog) RollFile() {
	var preFileName string
	for i := aLog.logCnt; i >= 1; i-- {
		j := i - 1
		curFileName := fmt.Sprintf("%s_%d.log", aLog.logFileName, i)
		if j == 0 {
			preFileName = aLog.logFileName
		} else {
			preFileName = fmt.Sprintf("%s_%d.log", aLog.logFileName, j)
		}

		_, err := os.Stat(curFileName)
		if err == nil {
			os.Remove(curFileName)
		}

		_, err = os.Stat(preFileName)
		if err == nil {
			os.Rename(preFileName, curFileName)
		}
	}
}

func (aLog *ALog) PrintLog() {
	for function := range aLog.ch {
		fi, err := os.Stat(aLog.logFileName)
		if err == nil {
			if fi.Size() > aLog.fileSize {
				aLog.RollFile()
				aLog.ResetOutput()
			}
		}
		function()
	}
}

func (aLog *ALog) Info(format string, v ...interface{}) {

	timeFormatted := time.Now().Local().Format(time.UnixDate)

	if InfoLevel > aLog.level {
		return
	}
	aLog.ch <- func() {
		aLog.info.Output(2, fmt.Sprintf(fmt.Sprintf("%s\n", timeFormatted), v...))
		aLog.info.Output(2, fmt.Sprintf(format, v...))

	}
}
