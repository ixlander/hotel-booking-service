package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Init(isDev bool) {
	var lg *zap.Logger
	var err error

	if isDev {
		lg, err = zap.NewDevelopment()
	} else {
		lg, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	Log = lg.Sugar()
}
