package core

import "sync"

var logChan = make(chan string, 100)
var mu sync.Mutex

func AddLog(msg string) {
	mu.Lock()
	defer mu.Unlock()
	select {
	case logChan <- msg:
	default:
	}
}

func GetLogChan() chan string {
	return logChan
}