package log

import (
	"io"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func InitLogger() *zap.Logger {
	encoder := getEncoder()

	infoWrite := getLogWriter("./", "snake.log")
	warnWrite := getLogWriter("./", "snake.log.wf")
	errorWrite := getLogWriter("./", "snake.log.err")

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
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

	zapLog = zap.New(core)

	return zapLog
}

func getEncoder() zapcore.Encoder {
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

func getLogWriter(logPath, filename string) io.Writer {
	logFullPath := path.Join(logPath, filename)
	hook, err := rotatelogs.New(
		logFullPath+".%Y%m%d%H",                   // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(logFullPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(7),           // 文件最大保存份数
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic(err)
	}
	return hook
}
