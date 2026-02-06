package core

import "sync"

var (
	logChan = make(chan string, 100)
	mu      sync.Mutex
)

// AddLog 向通道发送日志，前端会实时收到
func AddLog(msg string) {
	mu.Lock()
	defer mu.Unlock()
	select {
	case logChan <- msg:
	default:
		// 频道满时丢弃老日志
	}
}

func GetLogChan() chan string {
	return logChan
}