package zapLogger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	logger         = new(Logger)
	opt            = new(Option)
	ws             zapcore.WriteSyncer
	debugConsoleWS = zapcore.Lock(os.Stdout)
)

func GetLogger() *zap.SugaredLogger {
	return logger.log.Sugar()
}

func Init(opts ...ModOptions) (err error) {
	if logger.inited {
		return nil
	}
	for _, item := range opts {
		item(opt)
	}
	opt.fixup()
	logger.Opt = opt
	if opt.Development {
		logger.zapConfig = zap.NewDevelopmentConfig()
		logger.zapConfig.EncoderConfig.EncodeTime = timeEncoder
		logger.zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger.zapConfig.EncoderConfig.EncodeCaller = callEncoder
	} else {
		logger.zapConfig = zap.NewProductionConfig()
		logger.zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		logger.zapConfig.EncoderConfig.EncodeTime = timeEncoder
		logger.zapConfig.EncoderConfig.EncodeCaller = callEncoder
	}
	logger.zapConfig.DisableStacktrace = true
	logger.zapConfig.Level.SetLevel(opt.Level)
	if err = logger.init(); err != nil {
		return
	}
	logger.inited = true
	fmt.Println("[NewLogger] success")
	return nil
}

type Logger struct {
	Opt       *Option
	inited    bool
	log       *zap.Logger
	zapConfig zap.Config
}

func (l *Logger) init() error {
	l.setSyncers()
	var err error
	l.log, err = l.zapConfig.Build(l.cors())
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) setSyncers() {
	filePath := filepath.Join(opt.LogDir, fmt.Sprintf("%s.log", opt.FileName))
	f := func(filePath string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    l.Opt.MaxSize,
			MaxAge:     l.Opt.MaxAge,
			MaxBackups: l.Opt.MaxBackups,
			LocalTime:  true,
			Compress:   false,
		})
	}
	ws = f(filePath)
}

func (l *Logger) cors() zap.Option {

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})
	var cores []zapcore.Core
	if l.Opt.Development {
		encoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(encoder, debugConsoleWS, errPriority),
			zapcore.NewCore(encoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(encoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(encoder, debugConsoleWS, debugPriority),
		}...)
	} else {
		encoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(encoder, ws, errPriority),
			zapcore.NewCore(encoder, ws, warnPriority),
			zapcore.NewCore(encoder, ws, infoPriority),
			zapcore.NewCore(encoder, ws, debugPriority),
		}...)
	}
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func callEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if indexSrc := strings.Index(caller.File, "src/"); indexSrc != -1 {
		if len(caller.File[indexSrc:]) > 5 {
			if index := strings.IndexByte(caller.File[indexSrc+4:], '/'); index != -1 {
				buf := GetBuff()
				buf.AppendString(caller.File[indexSrc+4+index:])
				buf.AppendByte(':')
				buf.AppendInt(int64(caller.Line))
				enc.AppendString(buf.String())
				buf.Free()
				return
			}
		}
	}
	enc.AppendString(caller.TrimmedPath())
}
