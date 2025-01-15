package database

import (
	coreError "github.com/luancpereira/APICheckout/core/errors"
)

type Utils struct{}

func (Utils) CoreErrorDatabase(err error) *coreError.CoreError {
	return coreError.New("error.database", err.Error())
}
