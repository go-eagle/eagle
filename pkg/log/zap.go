package log

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-eagle/eagle/pkg/utils"
)

const (
	// WriterConsole console输出
	WriterConsole = "console"
	// WriterFile 文件输出
	WriterFile = "file"

	logSuffix      = ".log"
	warnLogSuffix  = "_warn.log"
	errorLogSuffix = "_error.log"

	// defaultSkip zapLogger 包装了一层 zap.Logger，默认要跳过一层
	defaultSkip = 1
)

const (
	// RotateTimeDaily 按天切割
	RotateTimeDaily = "daily"
	// RotateTimeHourly 按小时切割
	RotateTimeHourly = "hourly"
)

var (
	hostname string
	logDir   string
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

// Prevent data race from occurring during zap.AddStacktrace
var zapStacktraceMutex sync.Mutex

func getLoggerLevel(cfg *Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Level]
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
func newZapLogger(cfg *Config, opts ...Option) (*zap.Logger, error) {
	for _, opt := range opts {
		opt(cfg)
	}
	return buildLogger(cfg, defaultSkip), nil
}

// newLoggerWithCallerSkip new logger with caller skip
func newLoggerWithCallerSkip(cfg *Config, skip int, opts ...Option) (Logger, error) {
	for _, opt := range opts {
		opt(cfg)
	}
	return &zapLogger{sugarLogger: buildLogger(cfg, defaultSkip+skip).Sugar()}, nil
}

func buildLogger(cfg *Config, skip int) *zap.Logger {
	logDir = cfg.LoggerDir
	if strings.HasSuffix(logDir, "/") {
		logDir = strings.TrimRight(logDir, "/")
	}

	var encoderCfg zapcore.EncoderConfig
	if cfg.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Encoding == WriterConsole {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var cores []zapcore.Core
	var options []zap.Option
	// init option
	hostname, _ = os.Hostname()
	option := zap.Fields(
		zap.String("ip", utils.GetLocalIP()),
		zap.String("app_id", cfg.ServiceName),
		zap.String("instance_id", hostname),
	)
	options = append(options, option)

	writers := strings.Split(cfg.Writers, ",")
	for _, w := range writers {
		switch w {
		case WriterConsole:
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
		case WriterFile:
			// info
			cores = append(cores, getInfoCore(encoder, cfg))

			// warning
			core, option := getWarnCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}

			// error
			core, option = getErrorCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}
		default:
			// console
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
			// file
			cores = append(cores, getAllCore(encoder, cfg))
		}
	}

	combinedCore := zapcore.NewTee(cores...)

	// 开启开发模式，堆栈跟踪
	if !cfg.DisableCaller {
		caller := zap.AddCaller()
		options = append(options, caller)
	}

	// 跳过文件调用层数
	addCallerSkip := zap.AddCallerSkip(skip)
	options = append(options, addCallerSkip)

	// 构造日志
	return zap.New(combinedCore, options...)
}

func getAllCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	allWriter := getLogWriterWithTime(cfg, GetLogFile(cfg.Filename, logSuffix))
	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.FatalLevel
	})

	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(allWriter),
		FlushInterval: cfg.FlushInterval,
	}
	return zapcore.NewCore(encoder, asyncWriter, allLevel)
}

func getInfoCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	infoWriter := getLogWriterWithTime(cfg, GetLogFile(cfg.Filename, logSuffix))
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})

	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(infoWriter),
		FlushInterval: cfg.FlushInterval,
	}
	return zapcore.NewCore(encoder, asyncWriter, infoLevel)
}

func getWarnCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	warnWriter := getLogWriterWithTime(cfg, GetLogFile(cfg.Filename, warnLogSuffix))
	var stacktrace zap.Option
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			zapStacktraceMutex.Lock()
			stacktrace = zap.AddStacktrace(zapcore.WarnLevel)
			zapStacktraceMutex.Unlock()
		}
		return lvl == zapcore.WarnLevel
	})

	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(warnWriter),
		FlushInterval: cfg.FlushInterval,
	}
	return zapcore.NewCore(encoder, asyncWriter, warnLevel), stacktrace
}

func getErrorCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	errorFilename := GetLogFile(cfg.Filename, errorLogSuffix)
	errorWriter := getLogWriterWithTime(cfg, errorFilename)
	var stacktrace zap.Option
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			zapStacktraceMutex.Lock()
			stacktrace = zap.AddStacktrace(zapcore.ErrorLevel)
			zapStacktraceMutex.Unlock()
		}
		return lvl >= zapcore.ErrorLevel
	})

	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(errorWriter),
		FlushInterval: cfg.FlushInterval,
	}
	return zapcore.NewCore(encoder, asyncWriter, errorLevel), stacktrace
}

// getLogWriterWithTime 按时间(小时)进行切割
func getLogWriterWithTime(cfg *Config, filename string) io.Writer {
	logFullPath := filename
	rotationPolicy := cfg.LogRollingPolicy
	backupCount := cfg.LogBackupCount
	// 默认
	var (
		rotateDuration time.Duration
		// 时间格式使用shell的date时间格式
		timeFormat string
	)
	if rotationPolicy == RotateTimeHourly {
		rotateDuration = time.Hour
		timeFormat = ".%Y%m%d%H"
	} else if rotationPolicy == RotateTimeDaily {
		rotateDuration = time.Hour * 24
		timeFormat = ".%Y%m%d"
	}
	hook, err := rotatelogs.New(
		logFullPath+timeFormat,
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
