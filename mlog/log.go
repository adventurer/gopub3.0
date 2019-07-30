package mlog

import (
	"log"
	"os"
	"time"
)

var Mlog *log.Logger

func init() {
	os.Mkdir("logs", 0766)
	Mlog = log.New(os.Stdout, "", log.Llongfile|log.LstdFlags)
}

func Flog(fileName string, title string, content string) {
	now := time.Now().Format("2006-01-02")
	os.Mkdir("logs/"+fileName, 0766)
	logFile, err := os.OpenFile("logs/"+fileName+"/"+now+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

	debugLog := log.New(logFile, title, log.LstdFlags)
	debugLog.Println(content)
}
