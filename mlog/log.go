package mlog

import (
	"log"
	"os"
	"time"

	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"
)

var Mlog *log.Logger
var CustomLogger context.Handler

func init() {
	os.Mkdir("logs", 0766)
	Mlog = log.New(os.Stdout, "", log.Llongfile|log.LstdFlags)

	CustomLogger = logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query:   true,
		Columns: true,
		// Columns: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})
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
