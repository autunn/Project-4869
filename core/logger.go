package core

var logChan = make(chan string, 100)

func AddLog(msg string) {
	select {
	case logChan <- msg:
	default:
	}
}

func GetLogChan() chan string {
	return logChan
}