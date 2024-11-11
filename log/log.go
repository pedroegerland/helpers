package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func StartZapCtx() error {
	return startZapCtxWithLevel(zapcore.DebugLevel)
}

func startZapCtxWithLevel(logLevel zapcore.Level) error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level.SetLevel(logLevel)

	logTemp, err := config.Build()
	if err != nil {
		return err
	}

	_ = zap.ReplaceGlobals(logTemp)
	return nil
}

func ProcessAndReturn(input interface{}) error {

	switch convertedInput := input.(type) {
	// TODO add more Types
	case string:
		zap.L().Info(convertedInput)
	case error:
		zap.L().Error(convertedInput.Error())
		return convertedInput
	default:
		return nil
	}

	return nil
}
