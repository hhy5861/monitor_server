package utils

import (
	"time"
)

func TimesTamp() string {
	return time.Now().Format("2006-01-02 15:04:05")

}

func TimesUnix() int64 {
	return time.Now().Unix()
}
