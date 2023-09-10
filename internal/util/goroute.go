package util

import (
	"runtime/debug"

	"github.com/usiot/gbserver/internal/logger"
)

func Recover(f func()) {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrOp().
				Any("op=PANIC||err=", r).
				Bytes("||stack=", RmBSpace(debug.Stack())).
				Done()
		}
	}()

	f()
}

func Go(f func()) { go Recover(f) }
