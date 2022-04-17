package log

import (
	logPackage "log"
	"os"
)

var InfoLogger *logPackage.Logger
var WarnLogger *logPackage.Logger
var ErrorLogger *logPackage.Logger

func init() {
	InfoLogger = logPackage.New(os.Stdout, "[INFO] ", logPackage.Lshortfile|logPackage.Ltime|logPackage.Ldate|logPackage.Lmicroseconds)
	WarnLogger = logPackage.New(os.Stdout, "[WARN] ", logPackage.Lshortfile|logPackage.Ltime|logPackage.Ldate|logPackage.Lmicroseconds)
	ErrorLogger = logPackage.New(os.Stdout, "[ERROR] ", logPackage.Lshortfile|logPackage.Ltime|logPackage.Ldate|logPackage.Lmicroseconds)
}
