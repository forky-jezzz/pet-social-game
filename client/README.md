# Cocos Creator 场景绑定（对齐 Design System）

见上级 [docs/art-direction.md](../docs/art-direction.md)。

## Open in Creator

1. Install Cocos Creator 3.8+
2. Open `client/`（或新建 3.8 工程后拷入 `assets/scripts/`）
3. 按示意搭建三个场景：`MainRoom` / `VisitRoom` / `SocialHub`

## MainRoom 节点

| 节点 | 组件 | 说明 |
|------|------|------|
| RoomRoot | — | 地板/墙/窗光/家具分层 |
| PetHost | `sp.Skeleton` + `PetActor` | 中心宠 |
| HUD_TopLeft | — | 头像、宠名、心情/饥饿胶囊 |
| HUD_TopRight | — | 信件、背包、设置 |
| HUD_SideRight | — | 聊天入口、晒图入口 |
| BottomActions | — | 喂养/换装/串门/拍照 |
| FeedSheet / WardrobeDrawer / VisitSheet | 面板脚本 | 滑动菜单、衣柜抽屉、串门选项 |
| HomePanel | `HomePanel` | 绑定 labels + petActor |

## VisitRoom

- 顶栏「XX 的家」+ 加好友/赞  
- PetHost + PetGuest 双 Spine  
- 打扫 / 帮喂食 / 留言 / 特产交换面板  

## SocialHub

- PhotoWall 双列瀑布流  
- ChatDrawer 半透明上滑抽屉（复用 `ChatManager`）  

## 配置

`assets/scripts/config/GameConfig.ts`：`apiBase` / `wsBase` / `platform`

## Spine

资源目录 `assets/spine/pet/`，动画名见 [docs/spine-pet.md](../docs/spine-pet.md)。
