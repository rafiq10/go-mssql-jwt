package mylog

import (
	"fmt"
	"log"
	"os"
)

type MyLog struct {
	infoLog *log.Logger
	warnLog *log.Logger
	errLog  *log.Logger
}

func (l *MyLog) Init() {
	infoLogFile, err := os.Create("info.log")
	if err != nil {
		fmt.Errorf("Error while creating info.log; myLog.init()%v=", err)
	}
	warnLogFile, err := os.Create("warn.log")
	if err != nil {
		fmt.Errorf("Error while creating warn.log; myLog.init()%v=", err)
	}
	errLogFile, err := os.Create("error.log")
	if err != nil {
		fmt.Errorf("Error while creating error.log; myLog.init()%v=", err)
	}
	l.infoLog.SetOutput(infoLogFile)
	l.warnLog.SetOutput(warnLogFile)
	l.errLog.SetOutput(errLogFile)
	l.infoLog = log.New(infoLogFile, "INFO: ", log.LstdFlags)
	l.warnLog = log.New(warnLogFile, "WARN: ", log.LstdFlags)
	l.errLog = log.New(errLogFile, "ERR: ", log.LstdFlags)
}

func (l *MyLog) Info(input ...interface{}) {
	l.infoLog.Println(input...)
	log.Printf("Info: %s", input)
}
func (l *MyLog) Warn(input ...interface{}) {
	l.warnLog.Println(input...)
	log.Printf("Warning: %s", input)
}
func (l *MyLog) Err(input ...interface{}) {
	l.errLog.Println(input...)
	log.Printf("Error: %s", input)
}
