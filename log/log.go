package myLog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type InitConfig struct {
	logPath     string
	loglevel    string
	maxSize     int
	maxBackups  int
	maxAge      int
	compress    bool
	serviceName string
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

var Logger *zap.Logger

// "info", 200, 30, 90, false, sys.ServerName
func InitZap(co ...ConfigOption) {
	ic := &InitConfig{
		logPath:     "",
		loglevel:    "info",
		maxSize:     200,
		maxBackups:  30,
		maxAge:      90,
		compress:    false,
		serviceName: "",
	}
	for _, ele := range co {
		ele(ic)
	}
	initialize(ic)
}

func initialize(initConfig *InitConfig) {
	if initConfig.logPath == "" {
		panic("日志路径不能为空")
	}
	if initConfig.serviceName == "" {
		panic("服务名称不能为空")
	}
	var level zapcore.Level
	switch initConfig.loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	lp := fmt.Sprintf("%s/%s/%s.log", initConfig.logPath, initConfig.serviceName, initConfig.loglevel)
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   lp,                    // 日志文件路径，默认 os.TempDir()
		MaxSize:    initConfig.maxSize,    // 每个日志文件保存10M，默认 100M
		MaxBackups: initConfig.maxBackups, // 保留30个备份，默认不限
		MaxAge:     initConfig.maxAge,     // 保留7天，默认不限
		Compress:   initConfig.compress,   // 是否压缩，默认不压缩
	}
	//write := zapcore.AddSync(&hook)
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "Logger",
		CallerKey:      "linenum", //显示源码行号
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     TimeEncoder,                    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器 zapcore.FullCallerEncoder
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		//zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook), zapcore.AddSync(os.Stdout)), // 打印到控制台和文件
		//write,
		level,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段,如：添加一个服务器名称
	filed := zap.Fields(zap.String("service", initConfig.serviceName))
	// 构造日志
	Logger = zap.New(core, caller, development, filed)
	//Logger = zap.New(core, filed)
}

type ConfigOption func(*InitConfig)

// 设置日志目录
func WithLogPath(lp string) ConfigOption {
	return func(config *InitConfig) {
		config.logPath = lp
	}
}

// 设置日志级别
func WithLogLevel(ll string) ConfigOption {
	return func(config *InitConfig) {
		config.loglevel = ll
	}
}

// 设置每个日志保存的最大大小
func WithMaxSize(ms int) ConfigOption {
	return func(config *InitConfig) {
		config.maxSize = ms
	}
}

// 设置保存天数
func WithMaxAge(ma int) ConfigOption {
	return func(config *InitConfig) {
		config.maxAge = ma
	}
}

// 设置是否压缩
func WithCompress(c bool) ConfigOption {
	return func(config *InitConfig) {
		config.compress = c
	}
}

// 设置服务名称
func WithServiceName(sn string) ConfigOption {
	return func(config *InitConfig) {
		config.serviceName = sn
	}
}
