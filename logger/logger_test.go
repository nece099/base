package logger

import (
	"errors"
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	log := LogInit()

	fmt.Printf("%+v", log)

	log.Error(errors.New("test"))
}
