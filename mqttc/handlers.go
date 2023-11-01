package mqttc

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/talkincode/esmqtt/app"
	"github.com/talkincode/esmqtt/common/zaplog/log"
	"github.com/talkincode/esmqtt/models"
)

func onElasticMessage(c mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	// log.Infof("onTeamsdnsClientMetrics %s %s", msg.Topic(), string(payload))
	emsg, err := DecodeMessage[models.ElasticMessage](payload)
	if err != nil {
		log.Errorf("onElasticMessage decode error: %s", err)
		return
	}
	app.MsgQueue().PushBack(emsg.Data)
	msg.Ack()
}
