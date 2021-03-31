package logger

import (
	"io"
	"log"
)

var (
	I *log.Logger
	W *log.Logger
	E *log.Logger
	D *log.Logger
)


func Init(wInfo,wWarn,wErr,wDbg io.Writer) {
	I = log.New(wInfo,"[I]\t",log.Ldate | log.Ltime)
	W = log.New(wInfo,"[W]\t",log.Ldate | log.Ltime)
	E = log.New(wInfo,"[E]\t",log.Ldate | log.Ltime)
	D = log.New(wInfo,"[D]\t",log.Ldate | log.Ltime | log.Lshortfile)
}