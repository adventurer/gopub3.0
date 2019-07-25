package mlog

import (
	"log"
	"os"
)

var Mlog *log.Logger

func init() {
	Mlog = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
}
