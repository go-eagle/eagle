package log

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/1024casts/snake/config"
	"github.com/1024casts/snake/pkg/net/ip"
)

const (
	// WriterStdOut 标准输出
	WriterStdOut = "stdout"
	// WriterFile 文件输出
	WriterFile = "file"
)

const (
	// RotateTimeDaily 按天切割
	RotateTimeDaily = "daily"
	// RotateTimeHourly 按小时切割
	RotateTimeHourly = "hourly"
)

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// zapLogger logger struct
type zapLogger struct {
	sugarLogger *zap.SugaredLogger
}

// newZapLogger new zap logger
func newZapLogger(cfg *config.Config) (Logger, error) {
	var encoderCfg zapcore.EncoderConfig
	if cfg.App.Mode == "debug" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder

	if cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var cores []zapcore.Core
	var options []zap.Option
	// 设置初始化字段
	hostname, _ := os.Hostname()
	option := zap.Fields(
		zap.String("ip", ip.GetLocalIP()),
		zap.String("app_id", cfg.Logger.Name),
		zap.String("instance_id", hostname),
	)
	options = append(options, option)

	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.FatalLevel
	})

	writers := strings.Split(cfg.Logger.Writers, ",")
	for _, w := range writers {
		switch w {
		case WriterStdOut:
			core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
			cores = append(cores, core)
		case WriterFile:
			cores = append(cores, getInfoCore(encoder, cfg))

			core, option := getWarnCore(encoder, cfg)
			cores = append(cores, core)
			options = append(options, option)

			core, option = getErrorCore(encoder, cfg)
			cores = append(cores, core)
			options = append(options, option)
		default:
			core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
			cores = append(cores, core)
			allWriter := getLogWriterWithTime(cfg, cfg.Logger.LoggerFile)
			core = zapcore.NewCore(encoder, zapcore.AddSync(allWriter), allLevel)
			cores = append(cores, core)
		}
	}

	combinedCore := zapcore.NewTee(cores...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	options = append(options, caller)
	// 开启文件及行号
	development := zap.Development()
	options = append(options, development)
	// 跳过文件调用层数
	addCallerSkip := zap.AddCallerSkip(2)
	options = append(options, addCallerSkip)

	// 构造日志
	logger := zap.New(combinedCore, options...).Sugar()
	if err := logger.Sync(); err != nil {
		logger.Error(err)
	}

	return &zapLogger{sugarLogger: logger}, nil
}

func getInfoCore(encoder zapcore.Encoder, cfg *config.Config) zapcore.Core {
	infoFilename := cfg.Logger.LoggerFile
	infoWrite := getLogWriterWithTime(cfg, infoFilename)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(infoWrite), infoLevel)
}

func getWarnCore(encoder zapcore.Encoder, cfg *config.Config) (zapcore.Core, zap.Option) {
	warnFilename := cfg.Logger.LoggerWarnFile
	warnWrite := getLogWriterWithTime(cfg, warnFilename)
	var stacktrace zap.Option
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		stacktrace = zap.AddStacktrace(zapcore.WarnLevel)
		return lvl == zapcore.WarnLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(warnWrite), warnLevel), stacktrace
}

func getErrorCore(encoder zapcore.Encoder, cfg *config.Config) (zapcore.Core, zap.Option) {
	errorFilename := cfg.Logger.LoggerErrorFile
	errorWrite := getLogWriterWithTime(cfg, errorFilename)
	var stacktrace zap.Option
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		stacktrace = zap.AddStacktrace(zapcore.ErrorLevel)
		return lvl >= zapcore.ErrorLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(errorWrite), errorLevel), stacktrace
}

// getLogWriterWithTime 按时间(小时)进行切割
func getLogWriterWithTime(cfg *config.Config, filename string) io.Writer {
	logFullPath := filename
	rotationPolicy := cfg.Logger.LogRollingPolicy
	backupCount := cfg.Logger.LogBackupCount
	// 默认
	rotateDuration := time.Hour * 24
	if rotationPolicy == RotateTimeHourly {
		rotateDuration = time.Hour
	}
	hook, err := rotatelogs.New(
		logFullPath+".%Y%m%d%H",                     // 时间格式使用shell的date时间格式
		rotatelogs.WithLinkName(logFullPath),        // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(backupCount),   // 文件最大保存份数
		rotatelogs.WithRotationTime(rotateDuration), // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}
	return hook
}

// Debug logger
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Info logger
func (l *zapLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Warn logger
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// Error logger
func (l *zapLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugarLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugarLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugarLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugarLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugarLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugarLogger.Panicf(format, args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugarLogger.With(f...)
	return &zapLogger{newLogger}
}
