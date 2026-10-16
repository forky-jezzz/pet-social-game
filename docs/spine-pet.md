# 宠物 Spine 动作与表情规范

适用：Cocos Creator 3.8+ `sp.Skeleton`（Spine 3.8 / 4.x 按导出版本选扩展）。

## 骨骼建议

| 骨骼 | 用途 |
|------|------|
| `root` | 整体位移 |
| `body` | 呼吸起伏 |
| `head` | 点头、表情跟随 |
| `eye_L` / `eye_R` | 眨眼（缩放 Y 或换附件） |
| `tail` | 摇尾巴 |
| `hat_slot` / `clothes_slot` | 换装附件 |

## 动画列表（一期）

| 动画名 | 循环 | 时长建议 | 说明 |
|--------|------|----------|------|
| `idle_breath` | 是 | 2.0～2.5s | 身体轻微上下 + 胸口起伏，默认常驻 |
| `blink` | 否 | 0.15～0.25s | 双眼合开；由逻辑随机触发 |
| `tail_wag` | 否或短循环 | 0.8～1.2s | 尾巴左右摆；idle 中间歇播放 |
| `happy` | 否 | 0.8s | 喂食/抚摸后；耳朵或身体轻跳 |
| `eat` | 否 | 1.2s | 喂养时低头啃食 |
| `hungry` | 可循环 | 1.5s | hunger 低：耷拉 / 慢眨眼 |
| `sleepy` | 可循环 | 2.0s | mood 低或夜间：半眼 + 慢呼吸 |
| `petted` | 否 | 0.6s | 被串门抚摸：满足眯眼 |

### Track 分配

| Track | 内容 |
|-------|------|
| 0 | 主状态：`idle_breath` / `hungry` / `sleepy` / `eat` |
| 1 | 叠加：`blink`、`tail_wag`（勿打断 track0 呼吸） |
| 2 | 一次性表情：`happy` / `petted` |

## 客户端触发规则

```
启动 → play idle_breath (track0, loop)
每 2.5～5s 随机 → blink (track1)
每 4～8s 随机 → tail_wag (track1)（blink 中则跳过）
喂养成功 → eat (track0) → 结束回 idle 或 happy
抚摸/零食 → petted 或 happy (track2)
hunger < 30 → 主状态切 hungry
mood < 30 → 主状态切 sleepy（优先于 hungry 时取更糟状态）
```

## 换装

- Skin 名与服务端 `outfit_json` 对应，例如：
  - hat: `hat_basic` / `hat_star`
  - clothes: `clothes_basic` / `clothes_cape`
- 或单一 Default skin，运行时 `setAttachment('hat', 'hat_star')`

## 导出检查清单

- [ ] 造型偏动森式圆润；纹理双线性过滤，避免粗糙大像素感
- [ ] atlas 预乘与引擎 `premultipliedAlpha` 一致
- [ ] 动画名与上表完全一致（供 `PetActor` 调用）
- [ ] 文件放入 `client/assets/spine/pet/`（需美术产出后提交）
- [ ] 禁止使用《动物森友会》等商业游戏提取素材
