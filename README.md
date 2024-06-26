# RMAnnounce

自动拉取 RoboMaster 官网的最新公告并通过飞书机器人发送。

## Build

```shell
docker build .
```

## Configuration

### Example

etc/config.yaml

```yaml
# 飞书机器人 Webhook
webhooks:
  - https://open.feishu.cn/open-apis/bot/v2/hook/...
# RM 公告最后一条 ID
lastId: 1708
# 监控更新的页面
monitored_pages:
  - id: 1653 # 页面 ID
    hash: # 页面正文的散列值 第一次运行时请留空
```

## Run

### Docker Run

```shell
docker run -d -v /path/to/config.yaml:/app/etc/config.yaml registry.cn-guangzhou.aliyuncs.com/scutrobot/rm-announce:latest
```

### Docker Compose

```yaml
version: '3'
services:
  rm-announce:
    image: registry.cn-guangzhou.aliyuncs.com/scutrobot/rm-announce:latest
    container_name: rm-announce
    volumes:
      - /path/to/config.yaml:/app/etc/config.yaml
    restart: always
```
