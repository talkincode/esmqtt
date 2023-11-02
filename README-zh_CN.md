# Esmqtt

Esmqtt 是一个 可以通过 MQTT 协议订阅接收消息并转发到 Elasticsearch 的工具。

## 功能特性

- [x] 支持 MQTT 协议 V3, V4
- [x] 将 MQTT 消息转发到 Elasticsearch
- [x] 支持消息规则路由， 自定义消息主题对应指定的 Elasticsearch 索引
- [ ] 支持消息持久化

## 快速开始

### 直接安装 esmqtt 服务

```bash     

go install github.com/talkincode/esmqtt

esmqtt -install

```   

在 linux 系统中使用 systemd 管理服务

`systemctl <start | stop | restart> esmqtt`

### docker 部署 esmqtt 服务

[docker-compose.yml](./docker-compose.yml)

```yaml

```bash

docker-compose up -d

```


### 配置

- yaml

>  /etc/esmqtt.yml 文件

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

- 环境变量

> .env 文件

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

### 消息规则

>  /var/esmqtt/rules.json 文件


```json
[
  {
    "topic": "testnode/elastic/message/create",
    "index": "testnode_message",
    "spliter": "day"
  }
  
]
```


- 程序将会在启动时加载该文件，如果文件不存在，程序将会自动订阅主题 `elastic/message/create`。索引取决于消息中指定的 `index` 字段。
- 如果消息中和规则中都没有指定 `index` 字段值，消息将被忽略。 
- 在使用规则时， 如果消息中没有指定 `index` 字段值，将会使用规则中的 `index` 字段值。 `spliter` 的值为 "year | month | day | hour "(或者为空), 根据日期规则生成索引名，
比如 `testnode_message_2021-01-01`， `testnode_message_2021-01`， `testnode_message_2021`。


### 消息模型

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

- `data.id` 指定消息转发到的文档ID，如果消息中没有指定，将会使用内部生成的UUID。
- `data.index` 指定消息转发到的索引，如果消息中没有指定，将会使用规则中的 `index` 字段值， 如果规则也未指定，消息将被忽略。
- `data.payload` 为自定义对象，如果为空，消息将会被忽略。
- `data.timestamp` 为消息时间戳，如果为空，将会使用当前时间戳。
- `data.vector` 为消息向量，如果为空，将会使用空数组(该字段针对GPT模型设计，维度为1536，特定场景使用，可为空)。