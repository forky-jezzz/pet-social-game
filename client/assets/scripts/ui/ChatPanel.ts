import { _decorator, Component, EditBox, Label, Node } from 'cc';
import { chatManager } from '../managers/ChatManager';
import { ChatMsg } from '../net/types';

const { ccclass, property } = _decorator;

@ccclass('ChatPanel')
export class ChatPanel extends Component {
    @property(Label)
    logLabel: Label | null = null;

    @property(EditBox)
    inputBox: EditBox | null = null;

    @property(EditBox)
    peerIdBox: EditBox | null = null;

    private lines: string[] = [];
    private unsub: (() => void) | null = null;
    private mode: 'world' | 'friend' = 'world';

    onEnable() {
        chatManager.start();
        this.unsub = chatManager.onMessage((m) => this.append(m));
        void this.load();
    }

    onDisable() {
        this.unsub?.();
        this.unsub = null;
    }

    async load() {
        const peerId = this.mode === 'friend' ? Number(this.peerIdBox?.string || 0) : 0;
        try {
            const list = await chatManager.loadHistory(this.mode, peerId);
            this.lines = list.map((m) => this.format(m));
            this.render();
        } catch (e) {
            console.error(e);
        }
    }

    onClickWorld() {
        this.mode = 'world';
        void this.load();
    }

    onClickFriend() {
        this.mode = 'friend';
        void this.load();
    }

    onClickSend() {
        const content = this.inputBox?.string?.trim() || '';
        if (!content) return;
        if (this.mode === 'world') {
            chatManager.sendWorld(content);
        } else {
            const peerId = Number(this.peerIdBox?.string || 0);
            if (!peerId) return;
            chatManager.sendFriend(peerId, content);
        }
        if (this.inputBox) this.inputBox.string = '';
    }

    private append(m: ChatMsg) {
        if (this.mode === 'world' && m.channel !== 'world') return;
        if (this.mode === 'friend' && m.channel !== 'friend') return;
        this.lines.push(this.format(m));
        if (this.lines.length > 80) this.lines.shift();
        this.render();
    }

    private format(m: ChatMsg) {
        return `[${m.fromName}] ${m.content}`;
    }

    private render() {
        if (this.logLabel) this.logLabel.string = this.lines.join('\n');
    }
}
