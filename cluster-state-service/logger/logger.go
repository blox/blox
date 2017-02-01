// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package logger

import (
	"os"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	defaultLogFile    = "/var/output/logs/css.log"
	logFileEnvVarName = "CSS_LOG_FILE"

	defaultLogLevel    = "info"
	logLevelEnvVarName = "CSS_LOG_LEVEL"
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
	<!-- TODO: only errors go into the error.log -->
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
