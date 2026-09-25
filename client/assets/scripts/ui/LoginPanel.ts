import { _decorator, Component, Label, Node } from 'cc';
import { authManager } from '../managers/AuthManager';
import { petManager } from '../managers/PetManager';

const { ccclass, property } = _decorator;

@ccclass('LoginPanel')
export class LoginPanel extends Component {
    @property(Label)
    statusLabel: Label | null = null;

    @property(Node)
    homeRoot: Node | null = null;

    async onClickLogin() {
        this.setStatus('登录中...');
        try {
            const data = await authManager.login();
            this.setStatus(`欢迎 ${data.user.nickname}`);
            if (!data.pet) {
                await petManager.create(`${data.user.nickname}的宠物`);
            } else {
                await authManager.refreshMe();
            }
            if (this.homeRoot) this.homeRoot.active = true;
            this.node.active = false;
        } catch (e) {
            this.setStatus(`登录失败: ${(e as Error).message}`);
        }
    }

    private setStatus(msg: string) {
        if (this.statusLabel) this.statusLabel.string = msg;
        else console.log(msg);
    }
}
