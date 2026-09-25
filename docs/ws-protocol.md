# WebSocket 协议

## 连接

```
ws://127.0.0.1:8080/ws?token=<JWT>
```

鉴权失败返回 HTTP 401，不升级。

## 客户端 → 服务端

### 心跳

```json
{ "type": "ping" }
```

### 发消息

```json
{ "type": "chat.send", "channel": "world", "content": "你好" }
```

```json
{ "type": "chat.send", "channel": "friend", "toId": 2, "content": "嗨" }
```

- `content` 最长 200 字  
- 私聊需双方已是好友  
- 发送前走内容安全；不通过不落库  

## 服务端 → 客户端

### 消息

```json
{
  "type": "chat.msg",
  "id": 1,
  "channel": "world",
  "fromId": 1,
  "fromName": "Alice",
  "toId": 0,
  "content": "你好",
  "ts": 1710000000000
}
```

### 拒绝

```json
{ "type": "chat.rejected", "reason": "content_risk" }
```

Mock 拦截词：`违禁`、`赌博`、`色情`。

### 心跳应答 / 错误

```json
{ "type": "pong" }
```

```json
{ "type": "error", "msg": "..." }
```

## 运维说明

- 服务端每 25s 发 WebSocket Ping；客户端建议每 20s 发业务 `ping`  
- 多实例依赖 Redis PubSub：`pubsub:chat:world`、`pubsub:chat:dm`
