package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
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
		log.Info("Index already exists.")
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
			"vector":    message.Vector,
			"payload":   objdata,
			"timestamp": message.Timestamp,
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
	var result = make(map[string][]models.ElasticMessage)
	send := func(data map[string][]models.ElasticMessage) error {
		for index, message := range data {
			err := a.BatchPost(index, message)
			if err != nil {
				log.Errorf("batch post elastic message error: %s", err)
				return err
			} else {
				log.Infof("batch post elastic message success: %s %d", index, len(message))
			}
		}
		return nil
	}

	log.Info("Start post task")
	var lastPostTime = time.Now().Unix()
	for {
		msg := MsgQueue().PopFront()
		if msg == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if _, ok := result[msg.Index]; !ok {
			result[msg.Index] = make([]models.ElasticMessage, 0)
		}
		result[msg.Index] = append(result[msg.Index], *msg)

		if len(result) >= 1000 || time.Now().Unix()-lastPostTime >= 1 {
			lastPostTime = time.Now().Unix()
			_ = a.taskPool.Submit(func() {
				err := send(result)
				if err != nil {
					log.Errorf("send message error: %s", err)
				}
			})
			time.Sleep(1 * time.Second)
		}
	}
}
