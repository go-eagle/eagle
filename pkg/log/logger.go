package log

import (
	"io"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/1024casts/snake/pkg/util"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A global variable so that logger functions can be directly accessed
var logger *zap.SugaredLogger

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

//Logger is our contract for the logger
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
}

// InitLogger 初始化logger
func NewLogger() *zap.SugaredLogger {
	encoder := getJSONEncoder()

	// 注意：如果多个文件，最后一个会是全的，前两个可能会丢日志
	infoFilename := viper.GetString("logger.logger_file")
	infoWrite := getLogWriterWithTime(infoFilename)
	warnFilename := viper.GetString("logger.logger_warn_file")
	warnWrite := getLogWriterWithTime(warnFilename)
	errorFilename := viper.GetString("logger.logger_error_file")
	errorWrite := getLogWriterWithTime(errorFilename)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWrite), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWrite), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWrite), errorLevel),
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("ip", util.GetLocalIP()), zap.String("app", viper.GetString("name")))
	// 构造日志
	logger = zap.New(core, caller, development, filed).Sugar()

	return logger
}

func getJSONEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "timestamp",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		NameKey:       "app",
		CallerKey:     "file",
		StacktraceKey: "stacktrace",
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 按时间(小时)进行切割
func getLogWriterWithTime(filename string) io.Writer {
	logFullPath := filename
	hook, err := rotatelogs.New(
		logFullPath+".%Y%m%d%H",                                                // 时间格式使用shell的date时间格式
		rotatelogs.WithLinkName(logFullPath),                                   // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(viper.GetUint("logger.log_backup_count")), // 文件最大保存份数
		rotatelogs.WithRotationTime(time.Hour),                                 // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}
	return hook
}

// Debug logger
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info logger
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn logger
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error logger
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal logger
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Debugf logger
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args)
}

// Infof logger
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args)
}

// Warnf logger
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Errorf logger
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatalf logger
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Panicf logger
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}
