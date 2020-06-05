package log

import (
   "log"
   "os"
   "github.com/comail/colog"
)

func LogSetUp() {
  colog.SetDefaultLevel(colog.LInfo)
  colog.SetMinLevel(colog.LDebug)
  colog.SetFormatter(&colog.StdFormatter{
    Colors: true,
    Flag:   log.Ldate | log.Ltime,
  })
  colog.Register()
}

func DebugLog(v ...interface{}) {
  log.Print("debug: ", v)
}

func InfoLog(v ...interface{}) {
  log.Print(v)
}

func ErrorLog(v ...interface{}) {
  log.Printf("error: %#v", v)
}

func FatalLog(v ...interface{}) {
  log.Printf("alert: %#v", v)
  os.Exit(1)
}