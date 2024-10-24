package log

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type options struct {
	logLevel   zerolog.Level
	outputPath string
}

type Option func(opts *options)

func NewLogger(opts ...Option) (zerolog.Logger, error) {
	dopts := &options{
		logLevel:   zerolog.InfoLevel,
		outputPath: "",
	}
	for i := range opts {
		opts[i](dopts)
	}
	return newLogger(dopts)
}

func newLogger(opts *options) (zerolog.Logger, error) {
	level := opts.logLevel

	if opts.outputPath == "" {
		return zerolog.New(os.Stdout).Level(level), nil
	}

	outputPath := fmt.Sprintf("%s/outputs.log", opts.outputPath)
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o755)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := zerolog.New(file).Level(level)
	return logger, nil
}

func WithLevel(level zerolog.Level) Option {
	return func(opts *options) {
		opts.logLevel = level
	}
}

func WithOuput(output string) Option {
	return func(opts *options) {
		opts.outputPath = output
	}
}
