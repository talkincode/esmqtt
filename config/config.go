package config

import (
	"fmt"
	"os"
	"path"

	"github.com/talkincode/esmqtt/common"
	"github.com/talkincode/esmqtt/common/envutils"
	"gopkg.in/yaml.v3"
)

type MqttConfig struct {
	Server   string `yaml:"server" json:"server"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Debug    bool   `yaml:"debug" json:"debug"`
}

type ElasticConfig struct {
	Server   string `yaml:"server" json:"server"`
	ApiKey   string `yaml:"api_key" json:"api_key"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Debug    bool   `yaml:"debug" json:"debug"`
}

type LogConfig struct {
	Mode           string `yaml:"mode"`
	ConsoleEnable  bool   `yaml:"console_enable"`
	LokiEnable     bool   `yaml:"loki_enable"`
	FileEnable     bool   `yaml:"file_enable"`
	Filename       string `yaml:"filename"`
	QueueSize      int    `yaml:"queue_size"`
	LokiApi        string `yaml:"loki_api"`
	LokiUser       string `yaml:"loki_user"`
	LokiPwd        string `yaml:"loki_pwd"`
	LokiJob        string `yaml:"loki_job"`
	MetricsStorage string `yaml:"metrics_storage"`
	MetricsHistory int    `yaml:"metrics_history"`
}

type AppConfig struct {
	Appid    string `yaml:"appid"`
	Location string `yaml:"location"`
	Workdir  string `yaml:"workdir"`
	Debug    bool   `yaml:"debug"`

	Logger  LogConfig     `yaml:"logger" json:"logger"`
	Mqtt    MqttConfig    `yaml:"mqtt" json:"mqtt"`
	Elastic ElasticConfig `yaml:"elastic" json:"elastic"`
}

func (c *AppConfig) GetLogDir() string {
	return path.Join(c.Workdir, "logs")
}

func (c *AppConfig) GetDataDir() string {
	return path.Join(c.Workdir, "data")
}

func (c *AppConfig) GetPrivateDir() string {
	return path.Join(c.Workdir, "private")
}

func (c *AppConfig) InitDirs() {
	err := os.MkdirAll(path.Join(c.Workdir, "logs"), 0755)
	err = os.MkdirAll(path.Join(c.Workdir, "private"), 0755)
	err = os.MkdirAll(path.Join(c.Workdir, "data"), 0755)
	if err != nil {
		fmt.Println(err)
	}
}

var DefaultAppConfig = &AppConfig{
	Appid:    "esmqtt",
	Location: "Asia/Shanghai",
	Workdir:  "/var/esmqtt",
	Debug:    true,
	Elastic: ElasticConfig{
		Server:   "http://127.0.0.1:9200",
		ApiKey:   "",
		Username: "elastic",
		Password: "elastic",
		Debug:    false,
	},
	Mqtt: MqttConfig{
		Server:   "",
		Username: "",
		Password: "",
		Debug:    false,
	},
	Logger: LogConfig{
		Mode:           "development",
		ConsoleEnable:  true,
		LokiEnable:     false,
		FileEnable:     true,
		Filename:       "/var/esmqtt/esmqtt.log",
		QueueSize:      4096,
		LokiApi:        "http://127.0.0.1:3100",
		LokiUser:       "esmqtt",
		LokiPwd:        "esmqtt",
		LokiJob:        "esmqtt",
		MetricsStorage: "/var/esmqtt/data/metrics",
		MetricsHistory: 24 * 7,
	},
}

func LoadConfig(cfile string) *AppConfig {
	// 开发环境首先查找当前目录是否存在自定义配置文件
	if cfile == "" {
		cfile = "esmqtt.yml"
	}
	if !common.FileExists(cfile) {
		cfile = "/etc/esmqtt.yml"
	}
	cfg := new(AppConfig)
	if common.FileExists(cfile) {
		data := common.Must2(os.ReadFile(cfile))
		common.Must(yaml.Unmarshal(data.([]byte), cfg))
	} else {
		cfg = DefaultAppConfig
	}

	envutils.SetEnvValue("ESMQTT_SYSTEM_WORKER_DIR", &cfg.Workdir)
	envutils.SetEnvBoolValue("ESMQTT_SYSTEM_DEBUG", &cfg.Debug)

	envutils.SetEnvValue("ESMQTT_ELASTIC_SERVER", &cfg.Elastic.Server)
	envutils.SetEnvValue("ESMQTT_ELASTIC_APIKEY", &cfg.Elastic.ApiKey)
	envutils.SetEnvValue("ESMQTT_ELASTIC_USERNAME", &cfg.Elastic.Username)
	envutils.SetEnvValue("ESMQTT_ELASTIC_PASSWORD", &cfg.Elastic.Password)
	envutils.SetEnvBoolValue("ESMQTT_ELASTIC_DEBUG", &cfg.Elastic.Debug)

	envutils.SetEnvValue("ESMQTT_MQTT_SERVER", &cfg.Mqtt.Server)
	envutils.SetEnvValue("ESMQTT_MQTT_USERNAME", &cfg.Mqtt.Username)
	envutils.SetEnvValue("ESMQTT_MQTT_PASSWORD", &cfg.Mqtt.Password)
	envutils.SetEnvBoolValue("ESMQTT_MQTT_DEBUG", &cfg.Mqtt.Debug)

	envutils.SetEnvValue("ESMQTT_LOGGER_JOB", &cfg.Logger.LokiJob)
	envutils.SetEnvValue("ESMQTT_LOGGER_SERVER", &cfg.Logger.LokiApi)
	envutils.SetEnvValue("ESMQTT_LOGGER_USERNAME", &cfg.Logger.LokiUser)
	envutils.SetEnvValue("ESMQTT_LOGGER_PASSWORD", &cfg.Logger.LokiPwd)
	envutils.SetEnvValue("ESMQTT_LOGGER_MODE", &cfg.Logger.Mode)
	envutils.SetEnvBoolValue("ESMQTT_LOGGER_LOKI_ENABLE", &cfg.Logger.LokiEnable)
	envutils.SetEnvBoolValue("ESMQTT_LOGGER_FILE_ENABLE", &cfg.Logger.FileEnable)

	cfg.InitDirs()

	return cfg
}
