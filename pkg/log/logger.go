package log

import (
	"io"
	"os"
	"time"

	"github.com/1024casts/snake/util"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() *zap.Logger {
	encoder := getJsonEncoder()

	// 注意：如果多个文件，最后一个会是全的，前两个可能会丢日志
	filename := "./snake"
	infoWrite := getLogWriterWithTime(filename, "")
	warnWrite := getLogWriterWithTime(filename, ".wf")
	errorWrite := getLogWriterWithTime(filename, ".err")

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
	filed := zap.Fields(zap.String("ip", util.GetCurrentIP()), zap.String("app", "ebao-policy"))
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

func getLogWriterWithTime(filename string, suffix string) io.Writer {
	logFullPath := filename + suffix
	hook, err := rotatelogs.New(
		logFullPath+".%Y%m%d%H",                // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(logFullPath),   // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(7),        // 文件最大保存份数
		rotatelogs.WithRotationTime(time.Hour), // 日志切割时间间隔
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
