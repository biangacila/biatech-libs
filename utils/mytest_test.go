package utils

import (
	"errors"
	"testing"
)

func TestMakeLog(t *testing.T) {
	NewLoggerSlogInfo(errors.New("test me here"))
}
