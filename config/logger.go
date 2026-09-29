package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

func NewLogger() *Logger {
	_ = os.MkdirAll("logs", 0755)
	f, err := os.OpenFile(filepath.Join("logs", "app.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &Logger{logger: log.Default()}
	}
	return &Logger{file: f, logger: log.New(f, "", 0)}
}
func (l *Logger) Request(entry map[string]any) {
	b, _ := json.Marshal(entry)
	l.logger.Println(string(b))
	log.Println(string(b))
}
func (l *Logger) Close() {
	if l.file != nil {
		_ = l.file.Close()
	}
}
