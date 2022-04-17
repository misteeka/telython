package log

import (
	logPackage "log"
	"os"
)

var InfoLogger *logPackage.Logger
var WarnLogger *logPackage.Logger
var ErrorLogger *logPackage.Logger

func init() {
	InfoLogger = logPackage.New(os.Stdout, "[INFO] ", logPackage.Lshortfile|logPackage.Ldate)
	WarnLogger = logPackage.New(os.Stdout, "[WARN] ", logPackage.Lshortfile|logPackage.Ldate)
	ErrorLogger = logPackage.New(os.Stdout, "[ERROR] ", logPackage.Lshortfile|logPackage.Ldate)
}
