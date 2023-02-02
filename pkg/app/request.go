package app

import (
	"github.com/beego/beego/v2/core/validation"
	"github.com/go-gin-example/pkg/logging"
)

func MakeErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}
