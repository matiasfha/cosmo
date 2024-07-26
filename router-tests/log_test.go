package integration_test

import (
    "bytes"
    "os"
    "testing"
	"strings"
    "encoding/json"

    "github.com/stretchr/testify/require"
    "go.uber.org/zap/zapcore"

    "github.com/wundergraph/cosmo/router/pkg/logging"
)

// logBuffer will capture the log output
var logBuffer bytes.Buffer

func TestLog(t *testing.T) {
    // Redirect standard output to capture log output
    var buf bytes.Buffer
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    // Adjust the logging.New function call as per its actual signature
    logger, err := logging.New(&logging.Config{
        PrettyLogging: false,
        Debug:         false,
        LogLevel:      "info",
        LogFile:       "test.log",
    })
    require.NoError(t, err)

    // Log something
    logger.Info("Test log message")

    w.Close()
    os.Stdout = old
    buf.ReadFrom(r)
    output := buf.String()

    // Verify the log output
    require.Contains(t, output, "Test log message")
}

type LogEntry struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func getLogs() ([]LogEntry, error) {
	var logs []LogEntry
	lines := strings.Split(logBuffer.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}
	return logs, nil
}

func assertLogContains(t *testing.T, level zapcore.Level, msg string) {
	logs, err := getLogs()
	require.NoError(t, err)
	for _, log := range logs {
		if log.Level == level.String() && strings.Contains(log.Msg, msg) {
			return
		}
	}
	t.Errorf("Expected log with level %s and message containing %q", level.String(), msg)
}

