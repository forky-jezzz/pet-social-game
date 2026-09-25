# Pet Social Game MVP（方案 B）· 产品与计划文档

> 本文档固化一期目标、范围边界、仓库结构、技术决策与计划变更记录。  
> 实现细节见 [architecture.md](./architecture.md)、[api.md](./api.md)、[ws-protocol.md](./ws-protocol.md)。

---

## 1. 目标范围（一期交付）

| 能力 | 说明 | 状态 |
|------|------|------|
| 支付宝登录 | auth code → session JWT；开发可用 `platform=dev` | 已实现（生产 OAuth 换票为骨架） |
| 宠物创建 / 喂养 / 换装 | 每用户一只宠；道具种子数据；喂养分布式锁 | 已实现 |
| 串门 | 访问好友家园、抚摸/送礼互动、visit_logs | 已实现 |
| 世界聊天 + 好友私聊 | WebSocket + Redis PubSub | 已实现 |
| 文本内容安全 | `alipay.security.risk.content.detect`；开发 mock | Mock 默认；Alipay 骨架 |
| HTTP API + MySQL + Redis | 持久化 / 缓存 / 在线状态 / 锁 / PubSub | 已实现 |
| 浏览器试玩页 | `/play/` 不依赖 Cocos，一期可玩 | 已实现（计划增补） |

### 明确不做（二期）

- 晒图广场
- 投票榜
- OSS 上传
- Spine 正式资源管线
- 商城支付

---

## 2. 技术栈

**客户端**

- 引擎：Cocos Creator 3.8+（TypeScript）
- UI & 动画：Cocos GUI + Spine（二期正式接入；一期占位）
- 网络：HTTP/HTTPS + WebSocket
- 适配：支付宝 | 微信小游戏
- 联调：`play/` 静态页（同域托管）

**服务端**

- Go + Gin + GORM + go-redis + gorilla/websocket
- MySQL（用户/宠物/道具持久化）
- Redis（缓存/在线状态/分布式锁/聊天 PubSub）
- 内容安全：接口抽象 + mock / alipay

---

## 3. 仓库结构

```
/opt/infra/workspace/pet-social-game/
├── README.md
├── docs/
│   ├── mvp-plan.md            # 本文：范围与计划
│   ├── architecture.md        # 架构总览与模块
│   ├── api.md
│   └── ws-protocol.md
├── play/                      # 浏览器试玩（一期可玩）
├── server/                    # Go + Gin
│   ├── cmd/api/main.go
│   ├── configs/config.yaml
│   ├── go.mod
│   ├── internal/
│   │   ├── config/
│   │   ├── middleware/        # JWT
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repo/
│   │   ├── ws/
│   │   ├── security/
│   │   ├── model/
│   │   └── pkg/
│   ├── migrations/
│   └── docker-compose.yml
└── client/                    # Cocos Creator 3.8+ TS 骨架
    ├── assets/scripts/
    │   ├── net/
    │   ├── platform/
    │   ├── managers/
    │   ├── ui/
    │   └── config/
    ├── package.json
    └── README.md
```

远程仓库：https://github.com/forky-jezzz/pet-social-game

---

## 4. 架构总览

详见 [architecture.md](./architecture.md) 完整图与数据流。摘要：

```mermaid
flowchart TB
  subgraph clients [Clients]
    Play[Browser_Play]
    Cocos[Cocos_TS]
    AlipaySDK[Alipay_MiniGame_SDK]
    WxSDK[WeChat_MiniGame_SDK]
  end

  subgraph server [Go_Server]
    Gin[Gin_HTTP]
    WS[WS_Hub]
    Svc[Services]
    Sec[ContentSecurity]
  end

  subgraph data [Data]
    MySQL[(MySQL)]
    Redis[(Redis)]
  end

  Play --> Gin
  Play --> WS
  Cocos --> AlipaySDK
  Cocos --> WxSDK
  Cocos -->|HTTPS_REST| Gin
  Cocos -->|WSS_Chat| WS
  Gin --> Svc
  WS --> Svc
  Svc --> MySQL
  Svc --> Redis
  WS --> Redis
  Svc --> Sec
```

---

## 5. 关键默认决策（重点）

| 决策项 | 选择 |
|--------|------|
| MVP 范围 | 方案 B（养成 + 社交聊天，不含晒图广场） |
| 项目路径 | `/opt/infra/workspace/pet-social-game/` |
| 服务端 | Go 1.22+ / Gin / GORM / go-redis / gorilla/websocket |
| 鉴权 | JWT HS256；`platform=dev` 联调登录 |
| 内容安全 | 默认 `mock`；配置密钥后切 `alipay` |
| 客户端正式端 | Cocos 脚本骨架；场景在 Creator 内绑定 |
| 可玩性 | 增加 `play/` + 服务端静态托管，保证一期可试玩 |
| 模块代理 | `GOPROXY=https://goproxy.cn,direct` |

---

## 6. 计划变更记录（Changelog）

| 日期 | 变更 | 原因 |
|------|------|------|
| 2026-09-25 | 选定 MVP **方案 B**（相对纯养成 A / 全功能 C） | 一期需含世界/好友聊天与内容安全，但排除晒图/投票/OSS |
| 2026-09-25 | 目录名 `pet-social-game`；服务端确认 Go+Gin | 与用户确认 |
| 2026-09-25 | 增补 `play/` 浏览器试玩页，并由 API 托管 `/play/` | 无 Cocos 时也能端到端试玩社交与养成 |
| 2026-09-25 | 主文档拆分为 mvp-plan / architecture / api / ws | 范围与架构、协议分离，便于后续改计划只改本文 |
| 2026-09-25 | 推送 GitHub `forky-jezzz/pet-social-game` | 版本管理与协作 |

后续若调整范围（例如提前做晒图），请在本表追加一行，并同步更新「目标范围」表格状态列。

---

## 7. 数据与协议索引

- 表结构：`server/migrations/001_init.sql`
- Redis：`online:{userId}`、`lock:feed:{userId}`、`pubsub:chat:world` / `pubsub:chat:dm`
- HTTP：[api.md](./api.md)
- WebSocket：[ws-protocol.md](./ws-protocol.md)

---

## 8. 验收清单（一期）

- [x] `docker compose up` 后 API `/healthz` 正常
- [x] 开发登录拿 JWT，创建宠物、喂养、换装
- [x] 双用户加好友、串门互动
- [x] 世界聊天 / 私聊经 WS 互通
- [x] 敏感词 mock 拦截 `chat.rejected`
- [x] 浏览器打开 `/play/` 可完成上述流程
- [x] Cocos 脚本与平台适配骨架齐全（场景绑定需 Creator）
