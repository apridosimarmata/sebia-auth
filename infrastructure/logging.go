package infrastructure

import "log"

var (
	logger = log.Default()
)

func Log(any string) {
	logger.Println(any)
}
