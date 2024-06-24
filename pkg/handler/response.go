package handler

import (
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/sirupsen/logrus"
)

func newErrorResponse(message string) error {
	logrus.Error(message)
	return gqlerrors.NewFormattedError(message)
}
