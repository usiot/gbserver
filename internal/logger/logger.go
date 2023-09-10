package logger

import (
	"io"
	"os"

	"github.com/inysc/hog"
	"github.com/usiot/gbserver/internal/config"
)

var (
	w  io.Writer
	lg = hog.DefaultLogger
)

func Init(cfg *config.Log) {
	w = io.MultiWriter(os.Stdout, &hog.LoggerFile{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  false,
		Compress:   false,
	})

	lg = hog.New(uint8(cfg.Level), w)
	lg.AddSkip(2)
}

func Trace(format string, args ...any) { lg.Trace().Msgf(format, args...) }
func Debug(format string, args ...any) { lg.Debug().Msgf(format, args...) }
func Info(format string, args ...any)  { lg.Info().Msgf(format, args...) }
func Warn(format string, args ...any)  { lg.Warn().Msgf(format, args...) }
func Error(format string, args ...any) { lg.Error().Msgf(format, args...) }
func Fatal(format string, args ...any) { lg.Fatal().Msgf(format, args...) }

func DbgOp() hog.Event { return lg.Debug() }
func InOp() hog.Event  { return lg.Info() }
func ErrOp() hog.Event { return lg.Error() }
func Op() hog.Event    { return lg.Op() }

func Writer() io.Writer { return w }
