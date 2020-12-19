package log

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

// zapLogger logger struct
type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

// newZapLogger new zap logger
func newZapLogger(cfg *Config) (Logger, error) {
	encoder := getJSONEncoder()

	var cores []zapcore.Core
	var options []zap.Option
	// 设置初始化字段
	option := zap.Fields(zap.String("ip", ip.GetLocalIP()), zap.String("app", cfg.Name))
	options = append(options, option)

	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.FatalLevel
	})

	writers := strings.Split(cfg.Writers, ",")
	for _, w := range writers {
		if w == WriterStdOut {
			core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
			cores = append(cores, core)
		}
		if w == WriterFile {
			infoFilename := cfg.LoggerFile
			infoWrite := getLogWriterWithTime(cfg, infoFilename)
			warnFilename := cfg.LoggerWarnFile
			warnWrite := getLogWriterWithTime(cfg, warnFilename)
			errorFilename := cfg.LoggerErrorFile
			errorWrite := getLogWriterWithTime(cfg, errorFilename)

			infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl <= zapcore.InfoLevel
			})
			warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				stacktrace := zap.AddStacktrace(zapcore.WarnLevel)
				options = append(options, stacktrace)
				return lvl == zapcore.WarnLevel
			})
			errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				stacktrace := zap.AddStacktrace(zapcore.ErrorLevel)
				options = append(options, stacktrace)
				return lvl >= zapcore.ErrorLevel
			})

			core := zapcore.NewCore(encoder, zapcore.AddSync(infoWrite), infoLevel)
			cores = append(cores, core)
			core = zapcore.NewCore(encoder, zapcore.AddSync(warnWrite), warnLevel)
			cores = append(cores, core)
			core = zapcore.NewCore(encoder, zapcore.AddSync(errorWrite), errorLevel)
			cores = append(cores, core)
		}
		if w != WriterFile && w != WriterStdOut {
			core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
			cores = append(cores, core)
			allWriter := getLogWriterWithTime(cfg, cfg.LoggerFile)
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

	return &zapLogger{sugaredLogger: logger}, nil
}

// getJSONEncoder
func getJSONEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		NameKey:        "app",
		CallerKey:      "file",
		StacktraceKey:  "trace",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriterWithTime 按时间(小时)进行切割
func getLogWriterWithTime(cfg *Config, filename string) io.Writer {
	logFullPath := filename
	rotationPolicy := cfg.LogRollingPolicy
	backupCount := cfg.LogBackupCount
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
	l.sugaredLogger.Debug(args...)
}

// Info logger
func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

// Warn logger
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

// Error logger
func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}
