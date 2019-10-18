package stat

import (
	"log"
	"time"
)

type logStat struct {
	errCount   int
	warnCount  int
	infoCount  int
	debugCount int
	minTime    time.Time
	maxTime    time.Time
}

func (l *logStat) addLevel(level string) {
	switch level {
	case "error":
		l.errCount++
	case "warning":
		l.warnCount++
	case "info":
		l.infoCount++
	case "debug":
		l.debugCount++
	default:
		log.Printf("Unkown log level: %s", level)
	}
}

func (l *logStat) addTime(t time.Time) {
	if l.minTime.IsZero() && l.maxTime.IsZero() {
		l.minTime = t
		l.maxTime = t
		return
	}
	if l.minTime.After(t) {
		l.minTime = t
	}
	if l.maxTime.Before(t) {
		l.maxTime = t
	}
}
