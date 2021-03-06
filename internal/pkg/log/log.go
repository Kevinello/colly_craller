package log

import (
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly/debug"
	goLogger "github.com/phachon/go-logger"
)

type _Logger struct {
	*goLogger.Logger
	counter int32
	start   time.Time
}

func (l *_Logger) Init() error {
	l.counter = 0
	l.start = time.Now()
	l.Logger = goLogger.NewLogger()
	return nil
}

func (l *_Logger) Event(e *debug.Event) {
	i := atomic.AddInt32(&l.counter, 1)
	l.Debugf("[%06d] %d [%06d - %s] %q (%s)\n", i, e.CollectorID, e.RequestID, e.Type, e.Values, time.Since(l.start))
}

var (
	GLogger    = new(_Logger)
	ModuleName = os.Getenv("MODULE_NAME")
	level      = os.Getenv("LOG_LEVEL")
)

func init() {
	GLogger.Init()

	// 日志等级默认为INFO
	logLevel := goLogger.LOGGER_LEVEL_INFO
	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = goLogger.LOGGER_LEVEL_DEBUG
	case "INFO":
		logLevel = goLogger.LOGGER_LEVEL_INFO
	case "ERROR":
		logLevel = goLogger.LOGGER_LEVEL_ERROR
	case "EMERGENCY":
		logLevel = goLogger.LOGGER_LEVEL_EMERGENCY
	case "ALERT":
		logLevel = goLogger.LOGGER_LEVEL_ALERT
	case "CRITICAL":
		logLevel = goLogger.LOGGER_LEVEL_CRITICAL
	case "NOTICE":
		logLevel = goLogger.LOGGER_LEVEL_NOTICE
	case "WARNING":
		logLevel = goLogger.LOGGER_LEVEL_WARNING
	}

	consoleConfig := &goLogger.ConsoleConfig{
		Format: "%timestamp_format% [%level_string%] [%file%:%function%](line %line%): %body%",
	}
	GLogger.Detach("console")
	GLogger.Attach("console", logLevel, consoleConfig)
}
