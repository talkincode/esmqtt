# Esmqtt

Esmqtt is a tool that can subscribe to receive messages through the MQTT protocol and forward them to Elasticsearch.

[中文](./README-zh_CN.md)

## Features

- [x] Support for MQTT Protocol V3, V4
- [x] Forward MQTT messages to Elasticsearch.
- [x] Support message rule routing, customize message subject to correspond to the specified Elasticsearch indexes.
- [ ] Support message persistence

## Quick start

### Install the esmqtt service

```bash     

go install github.com/talkincode/esmqtt

esmqtt -install

```   

Use systemd to manage services on a Linux system

`systemctl <start | stop | restart> esmqtt`

### docker Deploy the ESMQTT service

[docker-compose.yml](./docker-compose.yml)

```yaml

```bash

docker-compose up -d

```


### Configuration

- yaml

>  /etc/esmqtt.yml 

```yaml
appid: esmqtt
location: Asia/Shanghai
workdir: /var/esmqtt
debug: true
logger:
  mode: development
  console_enable: true
  loki_enable: false
  file_enable: true
  filename: /var/esmqtt/esmqtt.log
  queue_size: 4096
  loki_api: http://127.0.0.1:3100
  loki_user: esmqtt
  loki_pwd: esmqtt
  loki_job: esmqtt
  metrics_storage: /var/esmqtt/data/metrics
  metrics_history: 168
mqtt:
  server: ""
  username: ""
  password: ""
  debug: false
elastic:
  server: http://127.0.0.1:9200
  api_key: ""
  username: elastic
  password: elastic
  debug: false
  
```

- ENVIRONMENT VARIABLE

> .env 

```bash

ESMQTT_SYSTEM_WORKER_DIR=/tmp/esmqtt
ESMQTT_SYSTEM_DEBUG=true

ESMQTT_MQTT_SERVER=tcp://127.0.0.1:1883
ESMQTT_MQTT_USERNAME=esmqtt
ESMQTT_MQTT_PASSWORD=
ESMQTT_MQTT_DEBUG=true

ESMQTT_ELASTIC_SERVER=https://localhost:9200
ESMQTT_ELASTIC_APIKEY=
ESMQTT_ELASTIC_USERNAME=elastic
ESMQTT_ELASTIC_PASSWORD=elastic
ESMQTT_ELASTIC_DEBUG=true

ESMQTT_LOGGER_JOB=esmqtt
ESMQTT_LOGGER_SERVER=
ESMQTT_LOGGER_USERNAME=esmqtt
ESMQTT_LOGGER_PASSWORD=
ESMQTT_LOGGER_MODE=development
ESMQTT_LOGGER_LOKI_ENABLE=false
ESMQTT_LOGGER_FILE_ENABLE=true

```

### Message rules

>  /var/esmqtt/rules.json 


```json
[
  {
    "topic": "testnode/elastic/message/create",
    "index": "testnode_message",
    "spliter": "day"
  }
  
]
```


- The program will load the file at startup, and if it does not exist, the program will automatically subscribe to the topic `elastic/message/create`. The index depends on the `index` field specified in the message.
- If neither the `index` field value is specified in the message nor in the rule, the message will be ignored. 
- When using a rule, if no `index` field value is specified in the message, the `index` field value in the rule will be used. The value of `spliter` is "year | month | day | hour" (or null), and the index name is generated according to the date rule.
For example `testnode_message_2021-01-01`, `testnode_message_2021-01`, `testnode_message_2021`.


### Message model

```json
{
  "data": {
     "id": "1312313142",
     "index": "testindex",
     "payload": {
       "name": "test", "value": 100
     },
     "vector": [],
     "timestamp": 1698827090749
  }
}
```

- `data.id` specifies the document ID that the message is forwarded to, if not specified in the message, the internally generated UUID will be used.
- `data.index` specifies the index to which the message will be forwarded, if not specified in the message, the value of the `index` field in the rule will be used, if not specified in the rule, the message will be ignored.
- `data.payload` is a custom object, if it is empty, the message will be ignored.
- `data.timestamp` is the message timestamp, if it is empty, the current timestamp will be used.
- `data.vector` is the message vector, if it is empty, an empty array will be used (this field is designed for GPT model, with dimension 1536, for specific scenarios, it can be empty).

