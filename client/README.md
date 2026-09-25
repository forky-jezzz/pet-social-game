# Pet Social Game Client (Cocos Creator 3.8+)

## Open in Creator

1. Install [Cocos Creator 3.8+](https://www.cocos.com/creator-download)
2. Dashboard → Open → select this `client/` folder (or create a blank 3.8 project and copy `assets/scripts/` in)
3. Creator will generate `.meta` files automatically

## Scene binding (Main.scene)

Create a simple UI and attach scripts:

| Node | Component | Bindings |
|------|-----------|----------|
| Canvas | `GameBootstrap` | — |
| LoginPanel | `LoginPanel` | statusLabel, homeRoot → HomePanel node |
| HomePanel | `HomePanel` | petNameLabel, statsLabel, outfitLabel; Button click → onClickFeed / Treat / Outfit* |
| ChatPanel | `ChatPanel` | logLabel, inputBox, peerIdBox |
| FriendsPanel | `FriendsPanel` | listLabel, friendIdBox |
| VisitPanel | `VisitPanel` | hostIdBox, infoLabel |

Buttons call the corresponding `onClick*` methods via Click Events.

## Config

Edit `assets/scripts/config/GameConfig.ts`:

- `apiBase` / `wsBase` → your server
- `platform`: `auto` | `alipay` | `wechat` | `dev`

## Dev login

Without Alipay/WeChat SDK, set `platform: 'dev'`. Server must have `auth.allow_dev_login: true`.

## Spine

Reserve a child node under the pet sprite named `SpineRoot` for future Spine skins; MVP uses Label stats only.
