package conf

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	CIPHER_TEXT_PREFIX = "@ciphered@"
)

// 全局MySQL 客户端实例
var (
	config *Config
	db     *sql.DB
)

// 要想获取配置, 单独提供函数
// 全局Config对象获取函数
// func C() *Config {
// 	return config
// }

// Config 应用配置
type Config struct {
	App   *app   `toml:"app"`
	Log   *log   `toml:"log"`
	MySQL *mysql `toml:"mysql"`
}

func NewConfig() *Config {
	return &Config{
		App:   newDefaultAPP(),
		Log:   newDefaultLog(),
		MySQL: newDefaultMySQL(),
	}
}

type app struct {
	Name       string `toml:"name" env:"APP_NAME"`
	EncryptKey string `toml:"encrypt_key" env:"APP_ENCRYPT_KEY"`
	HTTP       *http  `toml:"http"`
	GRPC       *grpc  `toml:"grpc"`
}

func newDefaultAPP() *app {
	return &app{
		Name:       "cmdb",
		EncryptKey: "defualt app encrypt key",
		HTTP:       newDefaultHttp(),
		GRPC:       newDefaultGrpc(),
	}
}

type http struct {
	Host string `toml:"host" env:"HTTP_HOST"`
	Port string `toml:"port" env:"HTTP_PORT"`
}

func (h *http) GetAddr() string {
	return fmt.Sprintf("%s:%s", h.Host, h.Port)
}

func newDefaultHttp() *http {
	return &http{
		Host: "127.0.0.1",
		Port: "8850",
	}
}

type grpc struct {
	Host string `toml:"host" env:"GRPC_HOST"`
	Port string `toml:"port" env:"GRPC_PORT"`
	EnableSSL bool   `toml:"enable_ssl" env:"GRPC_ENABLE_SSL"`
	CertFile  string `toml:"cert_file" env:"GRPC_CERT_FILE"`
	KeyFile   string `toml:"key_file" env:"GRPC_KEY_FILE"`
}

func (g *grpc) GetAddr() string {
	return fmt.Sprintf("%s:%s", g.Host, g.Port)
}

func newDefaultGrpc() *grpc {
	return &grpc{
		Host: "127.0.0.1",
		Port: "8850",
	}
}

type mysql struct {
	Host        string `toml:"host" env:"MYSQL_HOST"`
	Port        string `toml:"port" env:"MYSQL_PORT"`
	UserName    string `toml:"username" env:"MYSQL_USERNAME"`
	Password    string `toml:"password" env:"MYSQL_PASSWORD"`
	Database    string `toml:"database" env:"MYSQL_DATABASE"`
	MaxOpenConn int    `toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"`
	MaxIdleConn int    `toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"`
	MaxLifeTime int    `toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"`
	MaxIdleTime int    `toml:"max_idle_time" env:"MYSQL_MAX_IDLE_TIME"`
	lock        sync.Mutex
}

func newDefaultMySQL() *mysql {
	return &mysql{
		Host:     "127.0.0.1",
		Port:     "3306",
		UserName: "root",
		Password: "Murphy",
		Database: "test",
	}
}

func (m *mysql) getDBConnection() (*sql.DB, error) {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s>error, %s", dataSourceName, err.Error())
	}

	if m.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(m.MaxIdleConn))
	}

	if m.MaxLifeTime > 0 {
		db.SetConnMaxLifetime(time.Duration(m.MaxLifeTime))
	}

	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetMaxOpenConns(m.MaxOpenConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dataSourceName, err.Error())
	}

	return db, nil
}

func (m *mysql) GetDB() *sql.DB {
	// 加载全局数据库实例

	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		db, err := m.getDBConnection()
		if err != nil {
			panic(err.Error())
		}
		return db
	}

	return db
}

type log struct {
	Filename   string `toml:"path_dir" env:"LOG_PATH_DIR"`       // 日志文件路径
	MaxSize    int    `toml:"max_size" env:"LOG_MAX_SIZE"`       // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int    `toml:"max_backups" env:"LOG_MAX_BACKUPS"` // 日志文件最多保存多少个备份
	MaxAge     int    `toml:"max_age" env:"LOG_MAX_AGE"`         // 文件最多保存多少天
	Compress   bool   `toml:"compress" env:"LOG_COMPRESS"`       // 是否压缩
	Level      string `toml:"level" env:"LOG_LEVEL"`
}

func newDefaultLog() *log {
	return &log{
		Filename:   "./logs/spikeProxy1.log", // 日志文件路径
		MaxSize:    128,                      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                       // 日志文件最多保存多少个备份
		MaxAge:     7,                        // 文件最多保存多少天
		Compress:   true,                     // 是否压缩
		Level:      "debug",
	}
}

func (l *log) LoadGloabalLogger() {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	hook := lumberjack.Logger{
		Filename:   l.Filename,   // 日志文件路径
		MaxSize:    l.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: l.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     l.MaxAge,     // 文件最多保存多少天
		Compress:   l.Compress,   // 是否压缩
	}

	// 设置日志级别
	atomicLevel := l.setLevel()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	// filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	logger := zap.New(core, caller, development)
	logger.Info("log 初始化成功",
		zap.String("当前日志级别为", l.Level),
	)
	// logger.Info("无法获取网址",
	// 	zap.String("url", "http://www.baidu.com"),
	// 	zap.Int("attempt", 3),
	// 	zap.Duration("backoff", time.Second))
	zap.ReplaceGlobals(logger)
}
