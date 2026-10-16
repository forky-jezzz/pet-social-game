# 支付宝社交宠物养成小游戏 · Pet Social Game

Cocos Creator 3.8+ (TypeScript) 客户端 + Go/Gin 服务端。  
MVP 方案 B：登录、喂养、换装、串门、世界/好友聊天、内容安全。

**最快试玩**：起服务后打开浏览器 → [http://127.0.0.1:8080/play/](http://127.0.0.1:8080/play/)

## 目录

```
pet-social-game/
├── play/            # ★ 浏览器试玩页（一期可玩）
├── client/          # Cocos Creator 脚本工程
├── server/          # Go API + WebSocket
└── docs/            # 架构与协议
```

## 5 分钟试玩

### 1. 启动依赖

```bash
cd server
docker compose up -d
```

等待 MySQL / Redis healthy（migration 会自动导入表与种子道具）。

### 2. 启动 API（同时托管试玩页）

```bash
cd server
export GOPROXY=https://goproxy.cn,direct
go run ./cmd/api -config configs/config.yaml
```

### 3. 打开试玩

浏览器访问：**http://127.0.0.1:8080/play/**

建议流程（双人社交）：

1. 窗口 A：账号码 `alice`，开始试玩 → 记下顶部 **ID**
2. 窗口 B（无痕/另一浏览器）：账号码 `bob` → 记下 **ID**
3. 互相在「好友」填对方 ID → 申请 / 接受
4. 「喂狗粮 / 喂零食 / 换装」体验养成
5. 「世界聊天」互发消息；「串门」抚摸/送礼；「好友私聊」互聊
6. 试内容安全：发送含「违禁」的文本会被拦截

## Cocos 客户端

见 [client/README.md](client/README.md)。正式支付宝/微信小游戏发布用 Creator 打开 `client/`，对接同一套 API。

## 文档

- [MVP 方案 B：范围 / 决策 / 计划变更](docs/mvp-plan.md)
- [架构总览](docs/architecture.md)
- [视觉 Design System（家园/串门/社交）](docs/art-direction.md)
  - [场景1 家园](docs/refs/scene-01-main-room.png)
  - [场景2 串门](docs/refs/scene-02-visiting.png)
  - [场景3 社交](docs/refs/scene-03-social-hub.png)
- [宠物 Spine 动作规范](docs/spine-pet.md)
- [HTTP API](docs/api.md)
- [WebSocket 协议](docs/ws-protocol.md)

## 二期（未做）

晒图广场、投票榜、OSS、Spine 正式换装、商城支付。详见 [docs/mvp-plan.md](docs/mvp-plan.md)。
