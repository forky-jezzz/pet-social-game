import { _decorator, Component, EditBox, Label } from 'cc';
import { visitManager } from '../managers/VisitManager';

const { ccclass, property } = _decorator;

@ccclass('FriendsPanel')
export class FriendsPanel extends Component {
    @property(Label)
    listLabel: Label | null = null;

    @property(EditBox)
    friendIdBox: EditBox | null = null;

    onEnable() {
        void this.refresh();
    }

    async refresh() {
        try {
            const list = await visitManager.listFriends();
            const text = list
                .map((f) => `${f.userId} ${f.nickname} [${f.status}]${f.online ? ' online' : ''}`)
                .join('\n');
            if (this.listLabel) this.listLabel.string = text || '暂无好友';
        } catch (e) {
            console.error(e);
        }
    }

    async onClickRequest() {
        const id = Number(this.friendIdBox?.string || 0);
        if (!id) return;
        await visitManager.requestFriend(id);
        await this.refresh();
    }

    async onClickAccept() {
        const id = Number(this.friendIdBox?.string || 0);
        if (!id) return;
        await visitManager.acceptFriend(id);
        await this.refresh();
    }
}
