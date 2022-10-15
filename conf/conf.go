package conf

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvPROD  = "prod"
	EnvBETA  = "beta"
	EnvDEV   = "dev"
	EnvLOCAL = "local"
	EnvTEST  = "test"
)

var (
	ENV            string
	WorkDir        = "."
	ConfigFile     = ""
	MysqlConfig    *Mysql
	MqttConfig     *Mqtt
	TdengineConfig *Tdengine
	RedisConfig    *Redis
	RabbitMQConfig *RabbitMQ
	ZapConfig      *Zap
)

type Mysql struct {
	Host         string
	Username     string
	Password     string
	Database     string
	Port         int
	MaxIdleConns int  `mapstructure:"max_idle_conns"`
	MaxOpenConns int  `mapstructure:"max_open_conns"`
	ShowSql      bool `mapstructure:"show_sql"`
}
type Mqtt struct {
	Host     string
	Port     int
	Username string
	Password string
	ClientId string
}
type Tdengine struct {
	Host     string
	Username string
	Password string
	Database string
	Port     int
	Keep     int
	Days     int
}
type Redis struct {
	Addr     string
	Password string
	Db       int
}
type RabbitMQ struct {
	Uri      string
	Exchange string
}

type Conf struct {
	Mqtt      Mqtt
	Mysql     Mysql
	Tdengine  Tdengine
	Redis     Redis
	RabbitMQ  RabbitMQ
	Zap       Zap
	AllowList []string `mapstructure:"allow_list"`
	MaxCpu    int      `mapstructure:"max_cpu"`
}

type Config struct {
	Conf Conf
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名

	MaxAge       int  `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 日志留存时间
	ShowLine     bool `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole bool `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
}

// ZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// TransportLevel 根据字符串转化为 zapcore.Level
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// TODO: we should no longer use init() function after remove all handler's integration tests
// ENV=test is for integration tests only, other ENV should call "InitConf" explicitly
func init() {
	if env := os.Getenv("ENV"); env == EnvTEST {
		InitConf()
	}
}

func InitConf() {
	//go test
	if workDir := os.Getenv("APISIX_API_WORKDIR"); workDir != "" {
		WorkDir = workDir
	}

	setupConfig()
	setupEnv()

}

func setupConfig() {
	// setup config file path
	if ConfigFile == "" {
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(WorkDir + "/conf")
	} else {
		viper.SetConfigFile(ConfigFile)
	}

	// load config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Sprintf("fail to find configuration: %s", ConfigFile))
		} else {
			panic(fmt.Sprintf("fail to read configuration: %s, err: %s", ConfigFile, err.Error()))
		}
	}

	// unmarshal config
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("fail to unmarshal configuration: %s, err: %s", ConfigFile, err.Error()))
	}

	if len(config.Conf.Zap.Level) > 0 {
		initZapConfig(config.Conf.Zap)
	}

	if len(config.Conf.Mysql.Host) > 0 {
		initMysqlConfig(config.Conf.Mysql)
	}

	if len(config.Conf.Tdengine.Host) > 0 {
		initTdengine(config.Conf.Tdengine)

	}
	if len(config.Conf.Mqtt.Host) > 0 {
		initMqtt(config.Conf.Mqtt)
	}
	if len(config.Conf.Redis.Addr) > 0 {
		initRedis(config.Conf.Redis)
	}
	if len(config.Conf.RabbitMQ.Uri) > 0 {
		initRabbitMQ(config.Conf.RabbitMQ)
	}

	// set degree of parallelism
	initParallelism(config.Conf.MaxCpu)

}

func setupEnv() {
	ENV = EnvPROD
	if env := os.Getenv("ENV"); env != "" {
		ENV = env
	}
}

func initZapConfig(conf Zap) {
	if conf.Level != "" {
		ZapConfig = &Zap{
			Level:         conf.Level,
			Prefix:        conf.Prefix,
			Format:        conf.Format,
			Director:      conf.Director,
			EncodeLevel:   conf.EncodeLevel,
			StacktraceKey: conf.StacktraceKey,
			MaxAge:        conf.MaxAge,
			ShowLine:      conf.ShowLine,
			LogInConsole:  conf.LogInConsole,
		}
	}
}

func initMysqlConfig(conf Mysql) {
	MysqlConfig = &Mysql{
		Host:         conf.Host,
		Username:     conf.Username,
		Password:     conf.Password,
		Database:     conf.Database,
		Port:         conf.Port,
		MaxIdleConns: conf.MaxIdleConns,
		MaxOpenConns: conf.MaxOpenConns,
		ShowSql:      conf.ShowSql,
	}
}

// initialize parallelism settings
func initParallelism(choiceCores int) {
	if choiceCores < 1 {
		return
	}
	maxSupportedCores := runtime.NumCPU()

	if choiceCores > maxSupportedCores {
		choiceCores = maxSupportedCores
	}
	runtime.GOMAXPROCS(choiceCores)
}

func initTdengine(conf Tdengine) {
	TdengineConfig = &Tdengine{
		Host:     conf.Host,
		Port:     conf.Port,
		Username: conf.Username,
		Password: conf.Password,
	}
}
func initMqtt(conf Mqtt) {
	MqttConfig = &Mqtt{
		Host:     conf.Host,
		Port:     conf.Port,
		Username: conf.Username,
		Password: conf.Password,
	}
}
func initRedis(conf Redis) {
	RedisConfig = &Redis{
		Addr:     conf.Addr,
		Password: conf.Password,
		Db:       conf.Db,
	}
}
func initRabbitMQ(conf RabbitMQ) {
	RabbitMQConfig = &RabbitMQ{
		Uri: conf.Uri,
	}
}
