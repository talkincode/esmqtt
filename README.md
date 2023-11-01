# Esmqtt

Esmqtt 是一个 可以通过 MQTT 协议订阅接收消息并转发到 Elasticsearch 的工具。


## 功能特性

- [x] 支持 MQTT 协议 V3, V4
- [x] 将 MQTT 消息转发到 Elasticsearch
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