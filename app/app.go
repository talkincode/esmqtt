package app

import (
	"crypto/tls"
	"net/http"
	"path"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/cenkalti/backoff/v4"
	"github.com/coocood/freecache"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/panjf2000/ants/v2"
	_ "github.com/panjf2000/ants/v2"
	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
	"github.com/talkincode/esmqtt/common"
	"github.com/talkincode/esmqtt/common/cziploc"
	"github.com/talkincode/esmqtt/common/deque"
	"github.com/talkincode/esmqtt/common/zaplog"
	"github.com/talkincode/esmqtt/common/zaplog/log"
	"github.com/talkincode/esmqtt/config"
	"github.com/talkincode/esmqtt/models"
	_ "github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
)

var app *Application

var LogNamespace = zap.String("namespace", "esmqtt")

const (
	MqttQueueName = "esmqtt.message"
)

type Application struct {
	teilock   sync.Mutex
	appConfig *config.AppConfig
	taskPool  *ants.Pool
	cache     *freecache.Cache
	sched     *cron.Cron
	ipFind    *cziploc.IpFetch
	esclient  *elasticsearch.Client
	queue     *deque.Deque[models.ElasticMessage]
}

func GApp() *Application {
	return app
}

func GConfig() *config.AppConfig {
	return app.appConfig
}

func MsgQueue() *deque.Deque[models.ElasticMessage] {
	return app.queue
}

func GCache() *freecache.Cache {
	return app.cache
}

func TaskPool() *ants.Pool {
	return app.taskPool
}

func IpFind() *cziploc.IpFetch {
	return app.ipFind
}

func SubmitTask(task func()) {
	app.SubmitTask(task)
}

func InitGlobalApplication(cfg *config.AppConfig) {
	app = NewApplication(cfg)
	app.Init(cfg)
}

func NewApplication(appConfig *config.AppConfig) *Application {
	return &Application{appConfig: appConfig, teilock: sync.Mutex{}}
}

func (a *Application) Init(cfg *config.AppConfig) {
	loc, err := time.LoadLocation(cfg.Location)
	if err != nil {
		log.Error("timezone config error")
	} else {
		time.Local = loc
	}

	err = a.inttElastic()
	common.Must(err)
	log.Info("init elastic success")

	a.queue = deque.New[models.ElasticMessage](100000)

	a.cache = freecache.NewCache(8 * 1024 * 1024)
	zaplog.InitGlobalLogger(zaplog.LogConfig{
		Mode:          cfg.Logger.Mode,
		ConsoleEnable: true,
		LokiEnable:    cfg.Logger.LokiEnable,
		FileEnable:    cfg.Logger.FileEnable,
		Filename:      cfg.Logger.Filename,
		LokiApi:       cfg.Logger.LokiApi,
		LokiUser:      cfg.Logger.LokiUser,
		LokiPwd:       cfg.Logger.LokiPwd,
		LokiJob:       cfg.Logger.LokiJob,
	})
	a.taskPool, err = ants.NewPool(1024)
	common.Must(err)

	a.ipFind = cziploc.NewIpFetch(path.Join(cfg.Workdir, "data/qqwry.dat"))
	a.initJob()

	go a.startPostTask()
}

func (a *Application) SubmitTask(task func()) {
	log.ErrorIf(a.taskPool.Submit(task), LogNamespace)
}

func (a *Application) inttElastic() (err error) {
	retryBackoff := backoff.NewExponentialBackOff()
	a.esclient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses:     []string{a.appConfig.Elastic.Server},
		APIKey:        a.appConfig.Elastic.ApiKey,
		Username:      a.appConfig.Elastic.Username,
		Password:      a.appConfig.Elastic.Password,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // WARNING: This disables TLS verification - only use in trusted environments
			},
		},
	})
	return err
}

func Release() {
	zaplog.Release()
	app.taskPool.Release()
}
