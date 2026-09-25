# HTTP API

Base: `http://127.0.0.1:8080`

统一响应：

```json
{ "code": 0, "msg": "ok", "data": {} }
```

`code != 0` 为业务/HTTP 错误。鉴权接口需头：`Authorization: Bearer <token>`。

## POST `/api/v1/auth/login`

```json
{ "platform": "dev|alipay|wechat", "code": "alice", "nickname": "Alice", "avatar": "" }
```

返回：`{ token, user, pet }`

## GET `/api/v1/me`

返回：`{ user, pet, online }`

## POST `/api/v1/pet`

```json
{ "name": "小橘" }
```

每用户仅可创建一次。

## POST `/api/v1/pet/feed`

```json
{ "foodCode": "food_kibble" }
```

可选：`food_treat`。提高 hunger / mood；双高时升级并重置。

## POST `/api/v1/pet/outfit`

```json
{ "hat": "hat_star", "clothes": "clothes_cape" }
```

字段均可选；种子服装见 migration。

## GET `/api/v1/friends`

返回好友（accepted）+ 待处理申请（pending，对方发起）。

## POST `/api/v1/friends/request`

```json
{ "friendId": 2 }
```

## POST `/api/v1/friends/accept`

```json
{ "friendId": 1 }
```

仅被申请人可接受。

## GET `/api/v1/visit/:userId`

串门查看宿主与宠物，写入 `visit_logs(action=view)`。

## POST `/api/v1/visit/:userId/interact`

```json
{ "action": "pet|gift" }
```

`pet` 加心情；`gift` 加饱食与心情。

## GET `/api/v1/chat/history`

Query：`channel=world|friend`，`peerId`（私聊必填），`limit`（默认 50）。

## GET `/healthz`

探活。
