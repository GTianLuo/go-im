package log

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex

	clientLogWriter io.Writer
	clientInfoLog   *log.Logger
	clientErrorLog  *log.Logger
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
	Fatal  = errorLog.Fatal
	Fatalf = errorLog.Fatalf

	ClientInfo  func(v ...any)
	ClientInfof func(format string, v ...any)

	ClientError  func(v ...any)
	ClientErrorf func(format string, v ...any)
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func InitClientLogFile() {
	// 打开日志文件
	file, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("打开日志文件失败：", err)
	}

	clientLogWriter = file
	clientInfoLog = log.New(clientLogWriter, "[info]", log.LstdFlags|log.Lshortfile)
	clientErrorLog = log.New(clientLogWriter, "[error]", log.LstdFlags|log.Lshortfile)

	ClientInfo = clientInfoLog.Println
	ClientInfof = clientInfoLog.Printf
	ClientError = clientErrorLog.Println
	ClientErrorf = clientInfoLog.Printf
}

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if level > ErrorLevel {
		errorLog.SetOutput(ioutil.Discard)
	}
	if level > InfoLevel {
		infoLog.SetOutput(ioutil.Discard)
	}
}
