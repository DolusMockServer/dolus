package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	writers []io.Writer
	logFile string
}

var Log *Logger 

type LoggerWebSocketClient struct {
	logger *Logger
	conn   *websocket.Conn
}

func (l *Logger) RegisterWebSocketClient(conn *websocket.Conn, lines int) {
	client := &LoggerWebSocketClient{conn: conn, logger: l}
	if lines > 0 {
		l.sendLogs(client, int(lines))
	}
	l.Debugf("New client with address %s connected.", conn.RemoteAddr().String())
	if l.GetLevel() >= logrus.DebugLevel {	// give time for logrus to write to stream before adding the client as a writer
		time.Sleep(time.Millisecond * 250)
	}
	l.addWriter(client)

}

func (l *Logger) addWriter(w io.Writer) {
	l.writers = append(l.writers, w)
	l.SetOutput(io.MultiWriter(l.writers...))
}
func (l *Logger) removeWriter(w io.Writer) {
	i := 0
	for ; i < len(l.writers); i++ {
		if l.writers[i] == w {
			break
		}
	}
	if i >= len(l.writers) {
		return
	}
	l.writers = append(l.writers[:i], l.writers[i+1:]...)
	l.SetOutput(io.MultiWriter(l.writers...))

}

func (c *LoggerWebSocketClient) Write(p []byte) (n int, err error) {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(p)); err != nil {
		c.conn.Close()
		c.logger.removeWriter(c)
		c.logger.Debugf("Client with address %s was removed from writers as connection was closed.", c.conn.RemoteAddr().String())
		return 0, fmt.Errorf("error sending message to client %s: %s", c.conn.RemoteAddr().String(), err.Error())
	}
	return len(p), nil

}

func (l *Logger) GetLogStream(lines int) (string, error) {
	data, err := readLastNLines(l.logFile, lines)
	if err != nil {
		return "", fmt.Errorf("failed to read last %d lines of log file: %s", lines, err.Error())
	}
	return data, nil

}
func (l *Logger) sendLogs(client *LoggerWebSocketClient, lines int) {
	logs, err := l.GetLogStream(lines)
	if err != nil {
		client.Write([]byte(err.Error()))
	}

	client.Write([]byte(logs))

}

type textFormatter struct{}

func (f *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func (l *Logger) SetTextFormatter() {
	l.Logger.SetFormatter(&textFormatter{})
}

func (l *Logger) SetDefaultFormatter() {
	l.Logger.SetFormatter(utcFormatter{&logrus.TextFormatter{
		DisableTimestamp: false, // Disable the timestamp
		DisableQuote:     true,
	}})
}

func NewLogger(logFile string) *Logger {

	logger := &Logger{
		Logger:  logrus.New(),
		logFile: logFile,
	}

	// Set lock so that when errors occur when can reset logger Output writer
	logger.SetNoLock()

	file, err := os.OpenFile(logger.logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err == nil {
		logger.addWriter(io.MultiWriter(os.Stdout, file))
	} else {
		logger.Info("Failed to log to file, using default stderr")
		// need to set a variable to tell registered channels can't fetch early logs
	}

	return logger

}

type utcFormatter struct {
	logrus.Formatter
}

func (u utcFormatter) Format(le *logrus.Entry) ([]byte, error) {
	le.Time = le.Time.UTC()
	return u.Formatter.Format(le)
}

func readLastNLines(filename string, n int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Read lines into a circular buffer to keep the last N lines.
	circularBuffer := make([]string, n)
	currentIndex := 0
	lineCount := 0
	for scanner.Scan() {
		circularBuffer[currentIndex] = scanner.Text()
		currentIndex = (currentIndex + 1) % n
		lineCount++
	}

	// Read the circular buffer in the correct order.
	for i := 0; i < min(n, lineCount); i++ {
		lines = append(lines, circularBuffer[i])
		// currentIndex = (currentIndex + 1) % n
	}

	return strings.Join(lines, "\n"), nil
}
