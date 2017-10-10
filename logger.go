package logger
import
	"log/syslog"

func logger () {
	logger, e := syslog.New(syslog.LOG_MAIL, "SenderPolicy")
	if e != nil {
		logger.Err(e.Error())
	}

}