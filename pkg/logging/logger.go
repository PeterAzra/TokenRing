package logging

import (
	"log"

	"github.com/go-errors/errors"
)

var LoggingLevel = 0

func Information(message string, v ...any) {
	if LoggingLevel == 0 {
		log.Printf(message, v...)
	}
}

func Warning(message string, v ...any) {
	if LoggingLevel < 2 {
		log.Printf(message, v...)
	}

}

func Error(err error, message string, v ...any) {
	if LoggingLevel < 3 {
		log.Printf(message, v...)
		var errorCheck *errors.Error
		if errors.As(err, &errorCheck) {
			log.Println(err.(*errors.Error).ErrorStack())
		} else {
			log.Println(err)
		}
	}
}
