package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/talkincode/esmqtt/common"
	"github.com/talkincode/esmqtt/common/zaplog/log"
	"github.com/talkincode/esmqtt/models"
)

func (a *Application) createMapping(indexName string) error {
	// 检查索引是否存在
	res, err := a.esclient.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	// 如果索引已存在，直接返回
	if res.StatusCode == 200 {
		return nil
	}

	// 定义映射
	mapping := `{
		"mappings": {
			"properties": {
				"vector": {
					"type": "dense_vector",
					"dims": 1536
				},
				"payload": {
					"type": "object"
				},
				"timestamp": {
					"type": "date",
					"format": "epoch_millis"
				},
				"@timestamp": {
					"type": "date",
					"format": "strict_date_optional_time||epoch_millis"
				},
				"clientid": {
					"type": "keyword"
				},
				"Node": {
					"type": "long"
				}	
			}
		}
	}`

	// 创建索引并设置映射
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err = req.Do(context.Background(), a.esclient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		} else {
			return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
		}
	}

	return nil
}

func (a *Application) BatchPost(indexName string, messages []models.ElasticMessage) error {
	err := a.createMapping(indexName)
	if err != nil {
		return err
	}

	// Prepare bulk request
	var buf bytes.Buffer
	for _, message := range messages {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    message.ID,
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}
		var objdata map[string]interface{}
		err := json.Unmarshal(message.Payload, &objdata)
		if err != nil {
			return err
		}
		body := map[string]interface{}{
			"vector":     message.Vector,
			"payload":    objdata,
			"timestamp":  message.Timestamp,
			"@timestamp": message.AtTime,
			"clientid":   message.ClientID,
			"Node":       message.Node,
		}
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	}

	// Bulk insert data
	res, err := a.esclient.Bulk(bytes.NewReader(buf.Bytes()), a.esclient.Bulk.WithIndex(indexName))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"]))
	}

	return nil
}

func (a *Application) startPostTask() {

	send := func(data map[string][]models.ElasticMessage) error {
		for index, messages := range data {
			err := a.BatchPost(index, messages)
			if err != nil {
				log.Errorf("batch post elastic messages error: %s", err)
				return err
			} else {
				log.Infof("batch post elastic messages success: %s %d", index, len(messages))
			}
		}
		return nil
	}

	log.Info("Start post task")
	var result = make(map[string][]models.ElasticMessage)
	var lastPostTime = time.Now().Unix()
	var counter int32 = 0
	for {

		if len(result) > 0 && counter >= 2000 || time.Now().Unix()-lastPostTime >= 3 {
			lastPostTime = time.Now().Unix()
			atomic.StoreInt32(&counter, 0)
			sresult := common.DeepCopy(result).(map[string][]models.ElasticMessage)
			result = make(map[string][]models.ElasticMessage)
			_ = a.taskPool.Submit(func() {
				err := send(sresult)
				if err != nil {
					log.Errorf("send message error: %s", err)
				}
			})
		}

		msg := MsgQueue().PopFront()
		if msg != nil {
			atomic.AddInt32(&counter, 1)
			if _, ok := result[msg.Index]; !ok {
				result[msg.Index] = make([]models.ElasticMessage, 0)
			}
			result[msg.Index] = append(result[msg.Index], *msg)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
