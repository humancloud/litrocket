package handlelog

import (
	"fmt"
	"io"
	"litrocket/common"
	"os"
	"time"
)

const (
	// Directory To Save Log File.
	logdir = "log"
)

// Log File.
var filename = logdir + "/" + time.Now().Format("2006-01-02")

// Handlelog Append Message to log file, If in "debug" mode print log info.
// Level : INFO,WARNING,FATAL.
// INFO  : Some Success Message.
// WARNING : Some Thing Failed.
// FATAL : Fatal error, Server Send Mail to Manager.
func Handlelog(level, msg string) {
	var (
		f   *os.File
		err error
	)

	// Create "log" Directory If The Directory Is Not Exist.
	if _, err = os.Stat(logdir); err != nil && !os.IsExist(err) {
		if err = os.MkdirAll(logdir, os.ModePerm); err != nil {
			return
		}
	}

	// Open Or Create Log file.
	f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	// Write Log To File.
	content := level + ": " + time.Now().Format("2006-01-02 15:04:05") + " " + msg + "\n"
	if _, err = io.WriteString(f, content); err != nil {
		return
	}

	f.Close()

	// Write to Stdout If AppMode is "debug".
	if common.AppMode == "debug" {
		fmt.Printf("%s", content)
	}

	// Send Mail If FaTal.
	if level == "FATAL" {
		SendMail(content)
	}
}

// todo Send Mail To Manager If Fatal error occurred.
func SendMail(msg string) {
	os.Exit(-1)
}
