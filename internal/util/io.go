package util

import (
	"io"
	"log"
)

func CloseWithLog(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("warn: close failed: %v", err)
	}
}
