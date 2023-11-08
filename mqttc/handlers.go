package mqttc

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/talkincode/esmqtt/app"
	"github.com/talkincode/esmqtt/common"
	"github.com/talkincode/esmqtt/common/zaplog/log"
	"github.com/talkincode/esmqtt/models"
)

func onElasticMessage(c mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	payloadStr := string(payload)
	if !json.Valid([]byte(payloadStr)) {
		//此处可以记录控制日志，或者根据控制命令实现一些功能，比如重启，升级等，目前只是简单的打印
		fmt.Println("The payload is not valid JSON, message skipped.")
		return
	}
	emsg, err := DecodeMessage[models.ElasticMessage](payload)
	if err != nil {
		log.Errorf("onElasticMessage decode error: %s", err)
		return
	}

	var newmsg models.ElasticMessage
	newmsg = emsg.Data

	if newmsg.Payload == nil {
		return
	}

	rule := app.GetTopicRule(msg.Topic())
	if rule != nil && newmsg.Index == "" && rule.Index != "" {
		newmsg.Index = rule.Index
		switch rule.Spliter {
		case "year":
			newmsg.Index = fmt.Sprintf("%s-%s", newmsg.Index, time.Now().UTC().Format("2006"))
		case "month":
			newmsg.Index = fmt.Sprintf("%s-%s", newmsg.Index, time.Now().UTC().Format("2006_01"))
		case "day":
			newmsg.Index = fmt.Sprintf("%s-%s", newmsg.Index, time.Now().UTC().Format("2006_01_02"))
		case "hour":
			newmsg.Index = fmt.Sprintf("%s-%s", newmsg.Index, time.Now().UTC().Format("2006_01_02_15"))
		}
	}

	if newmsg.Index == "" {
		log.Errorf("onElasticMessage index is empty")
		return
	}

	if newmsg.ID == "" {
		newmsg.ID = common.Md5UUID()
	}

	if newmsg.Timestamp == 0 {
		newmsg.Timestamp = time.Now().UTC().UnixMilli()
	}

	app.MsgQueue().PushBack(newmsg)
	msg.Ack()
}
