package logging

import (
	"fmt"
	"time"

	"novels-spider/conf"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", conf.RuntimePath, conf.Log["path"])
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s-%s.%s",
		conf.Log["name"],
		time.Now().Format(conf.Log["time_format"]),
		conf.Log["ext"],
	)
}
