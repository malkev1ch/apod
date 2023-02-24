package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var defLogger = New(&Config{
	LogLevel: "info",
	DevMode:  false,
	name:     "default",
})

type Config struct {
	LogLevel string
	DevMode  bool
	name     string
}

func NewConfig(logLevel string, devMode bool) *Config {
	return &Config{LogLevel: logLevel, DevMode: devMode}
}

type Logger interface {
	Sync() error
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(template string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(template string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(template string, keysAndValues ...interface{})
	Err(msg string, err error)
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(template string, keysAndValues ...interface{})
	Printf(template string, args ...interface{})
}

// Application logger
type appLogger struct {
	level       string
	devMode     bool
	name        string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

func New(cfg *Config) Logger {
	l := appLogger{
		level:   cfg.LogLevel,
		devMode: cfg.DevMode,
		name:    cfg.name,
	}
	l.InitLogger()

	return &l
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *appLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func (l *appLogger) InitLogger() {
	logLevel := l.getLoggerLevel()

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderConsole zapcore.EncoderConfig
	if l.devMode {
		encoderConsole = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConsole = zap.NewProductionEncoderConfig()
	}

	encoderConsole.TimeKey = "time"
	encoderConsole.LevelKey = "level"
	encoderConsole.CallerKey = zapcore.OmitKey
	encoderConsole.MessageKey = "message"
	encoderConsole.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConsole.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConsole.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConsole.EncodeName = zapcore.FullNameEncoder
	encoderConsole.EncodeDuration = zapcore.StringDurationEncoder

	consoleEncoder := zapcore.NewJSONEncoder(encoderConsole)

	core := zapcore.NewCore(consoleEncoder, logWriter, zap.NewAtomicLevelAt(logLevel))

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.logger = logger
	l.sugarLogger = logger.Sugar()

	l.logger = l.logger.Named(l.name)
	l.sugarLogger = l.sugarLogger.Named(l.name)
}

func (l *appLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *appLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *appLogger) Debugw(template string, keysAndValues ...interface{}) {
	l.sugarLogger.Errorw(template, keysAndValues...)
}

func (l *appLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *appLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *appLogger) Infow(template string, keysAndValues ...interface{}) {
	l.sugarLogger.Infow(template, keysAndValues...)
}

func (l *appLogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *appLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *appLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

func (l *appLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *appLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *appLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *appLogger) Errorw(template string, keysAndValues ...interface{}) {
	l.sugarLogger.Errorw(template, keysAndValues...)
}

func (l *appLogger) Err(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

func (l *appLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *appLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *appLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *appLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *appLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *appLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

func (l *appLogger) Fatalw(template string, keysAndValues ...interface{}) {
	l.sugarLogger.Fatalw(template, keysAndValues...)
}

func (l *appLogger) Sync() error {
	go l.logger.Sync() // nolint: errcheck
	return l.sugarLogger.Sync()
}
