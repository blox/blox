package logger

import (
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const logConfig = `
<!-- TODO: only errors go into the error.log -->
<seelog type="asyncloop" minlevel="debug">
	<outputs formatid="main">
		<console/>
		<!-- TODO: Enable this when we have container builds. Until then, console 
		logger should do the trick. 
	    	<rollingfile filename="/var/output/logs/handler.log" type="date"
		     datepattern="2006-01-02-15" archivetype="none" maxrolls="72" />
		-->
    </outputs>
    <formats>
        <format id="main" format="%UTCDate(2006-01-02T15:04:05Z07:00) [%LEVEL] %Msg%n" />
    </formats>
</seelog>
`

func InitLogger() error {
	logger, err := log.LoggerFromConfigAsString(logConfig)
	if err != nil {
		return errors.Wrap(err, "Could not load logger config")
	}

	err = log.ReplaceLogger(logger)
	if err != nil {
		return errors.Wrap(err, "Could not replace logger")
	}

	return nil
}
