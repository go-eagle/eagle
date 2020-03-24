package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/1024casts/snake/util"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() *zap.Logger {
	encoder := getJsonEncoder()

	// 注意：如果多个文件，最后一个会是全的，前两个可能会丢日志
	infoFilename := viper.GetString("log.logger_file")
	infoWrite := getLogWriterWithTime(infoFilename)
	warnFilename := viper.GetString("log.logger_warn_file")
	warnWrite := getLogWriterWithTime(warnFilename)
	errorFilename := viper.GetString("log.logger_error_file")
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
	logger = zap.New(core, caller, development, filed)

	return logger
}

func getJsonEncoder() zapcore.Encoder {
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
		logFullPath+".%Y%m%d%H",                                             // 时间格式使用shell的date时间格式
		rotatelogs.WithLinkName(logFullPath),                                // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(viper.GetUint("log.log_backup_count")), // 文件最大保存份数
		rotatelogs.WithRotationTime(time.Hour),                              // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Debug(msg string, args ...zap.Field) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...zap.Field) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	logger.Error(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	logger.Fatal(msg, args...)
}

func Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.Info(message)
}
