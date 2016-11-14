// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
