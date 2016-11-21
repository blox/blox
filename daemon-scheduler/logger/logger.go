package logger

import (
	"os"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	defaultLogFile    = "/var/output/logs/daemon.log"
	logFileEnvVarName = "DS_LOG_FILE"

	defaultLogLevel    = "info"
	logLevelEnvVarName = "DS_LOG_LEVEL"
)

// InitLogger initializes and configures the logger
func InitLogger() error {
	logger, err := log.LoggerFromConfigAsString(loggerConfig())
	if err != nil {
		return errors.Wrap(err, "Could not load logger config")
	}

	err = log.ReplaceLogger(logger)
	if err != nil {
		return errors.Wrap(err, "Could not replace logger")
	}

	return nil
}

func loggerConfig() string {
	return `
	<seelog type="asyncloop" minlevel="` + logLevel() + `">
		<outputs formatid="main">
			<console/>
		    	<rollingfile filename="` + logFile() + `" type="date"
			     datepattern="2006-01-02-15" archivetype="none" maxrolls="72" />
			-->
	    </outputs>
	    <formats>
	        <format id="main" format="%UTCDate(2006-01-02T15:04:05Z07:00) [%LEVEL] %Msg%n" />
	    </formats>
	</seelog>
	`
}

func logFile() string {
	logFile := os.Getenv(logFileEnvVarName)
	if logFile == "" {
		logFile = defaultLogFile
	}
	return logFile
}

func logLevel() string {
	levels := map[string]string{
		"debug": "debug",
		"info":  "info",
		"warn":  "warn",
		"error": "error",
		"crit":  "critical",
		"none":  "off",
	}
	level, ok := levels[os.Getenv(logLevelEnvVarName)]
	if ok {
		return level
	}
	return defaultLogLevel
}
