// orderbook/log_service.go
package orderbook

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var logService *LogService

type LogService struct {
	mu     sync.Mutex
	logger *log.Logger
}

func InitializeLogService(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	logService = &LogService{
		logger: log.New(file, "ORDERBOOK ", log.LstdFlags|log.Lshortfile),
	}
	return nil
}

func Log(message string) {
	if logService != nil {
		logService.mu.Lock()
		defer logService.mu.Unlock()
		logService.logger.Println(message)
	}
}
