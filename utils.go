package gormutils

import (
	"github.com/xiaotushaoxia/errx"
)

func failedTo(err error, opt string) error {
	return errx.Wrap(err, "failed to "+opt)
}

func failedTof(err error, format string, a ...any) error {
	return errx.Wrapf(err, "failed to "+format, a...)
}
