package logger

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"
)

import (
	perrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger Logger

// WsLogger is logger struct
type WsLogger struct {
	mutex sync.Mutex
	Logger
	dynamicLevel zap.AtomicLevel
	// disable presents the logger state. if disable is true, the logger will write nothing
	// the default value is false
	disable bool
}

// Logger
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
}

func init() {
	// only use in test case, so just load default config
	if logger == nil {
		InitLogger(nil)
	}
}

// InitLog load from config path
func InitLog(logConfFile string) error {
	if logConfFile == "" {
		InitLogger(nil)
		return perrors.New("log configure file name is nil")
	}
	if path.Ext(logConfFile) != ".yml" {
		InitLogger(nil)
		return perrors.New(fmt.Sprintf("log configure file name %s suffix must be .yml", logConfFile))
	}

	_, err := ioutil.ReadFile(logConfFile)
	if err != nil {
		InitLogger(nil)
		return perrors.New(fmt.Sprintf("ioutil.ReadFile file:%s, error:%v", logConfFile, err))
	}

	conf := &zap.Config{}
	// err = yaml.UnmarshalYML(confFileStream, conf)
	if err != nil {
		InitLogger(nil)
		return perrors.New(fmt.Sprintf("[Unmarshal]init logger error: %v", err))
	}

	InitLogger(conf)

	return nil
}

func InitLogger(conf *zap.Config) {
	var zapLoggerConfig zap.Config
	if conf == nil {
		zapLoggerConfig = zap.NewDevelopmentConfig()
		zapLoggerEncoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
	} else {
		zapLoggerConfig = *conf
	}
	zapLogger, _ := zapLoggerConfig.Build(zap.AddCallerSkip(1))
	// logger = zapLogger.Sugar()
	logger = &WsLogger{Logger: zapLogger.Sugar(), dynamicLevel: zapLoggerConfig.Level}
}

func SetLogger(log Logger) {
	logger = log
}

func GetLogger() Logger {
	return logger
}

func SetLoggerLevel(level string) bool {
	if l, ok := logger.(OpsLogger); ok {
		l.SetLoggerLevel(level)
		return true
	}
	return false
}

type OpsLogger interface {
	Logger
	// SetLoggerLevel function as name
	SetLoggerLevel(level string)
}

// SetLoggerLevel ...
func (dpl *WsLogger) SetLoggerLevel(level string) {
	l := new(zapcore.Level)
	l.Set(level)
	dpl.dynamicLevel.SetLevel(*l)
}
