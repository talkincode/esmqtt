package mqttc

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/talkincode/esmqtt/app"
	"github.com/talkincode/esmqtt/common/zaplog/log"
	"go.uber.org/zap"
)

var client mqtt.Client

const (
	MqttElasticMessageCreate = "elastic/message/create"
)

func Start() error {
	mqtt.ERROR = zap.NewStdLog(zap.L())
	return startMqttDaemon()
}

func Publish(qos int, topic string, payload interface{}) error {
	t := client.Publish(topic, byte(qos), false, payload)
	if t.Wait() && t.Error() != nil {
		log.Errorf("publish %s error: %s\n", topic, t.Error())
		return t.Error()
	}

	if app.GConfig().Debug {
		log.Infof("publish %s %s", topic, payload)
	}
	return nil
}

// 连接时订阅主题
func onConnect(client mqtt.Client) {
	log.Info("mqtt connect success")
	if token := client.Subscribe(MqttElasticMessageCreate, 1, onElasticMessage); token.Wait() && token.Error() != nil {
		log.Errorf("onConnect subscribe MqttTeamsdnsInform error: %s", token.Error())
	}
}

// Start 启动守护进程
func startMqttDaemon() error {
	opts := mqtt.NewClientOptions().
		AddBroker(app.GConfig().Mqtt.Server).
		SetClientID(app.GConfig().Mqtt.Username).
		SetUsername(app.GConfig().Mqtt.Username).
		SetPassword(app.GConfig().Mqtt.Password)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectRetryInterval(10 * time.Second)
	opts.SetCleanSession(true)
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.OnConnect = onConnect
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
