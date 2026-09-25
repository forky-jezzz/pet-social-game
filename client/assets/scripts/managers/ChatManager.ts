import { http } from '../net/HttpClient';
import { wsClient } from '../net/WsClient';
import { ChatMsg } from '../net/types';

type ChatListener = (msg: ChatMsg) => void;

export class ChatManager {
    private listeners = new Set<ChatListener>();
    private unsub: (() => void) | null = null;

    start() {
        if (this.unsub) return;
        this.unsub = wsClient.on('chat.msg', (raw) => {
            const msg: ChatMsg = {
                id: Number(raw.id || 0),
                channel: String(raw.channel || ''),
                fromId: Number(raw.fromId || 0),
                fromName: String(raw.fromName || ''),
                toId: raw.toId ? Number(raw.toId) : undefined,
                content: String(raw.content || ''),
                ts: Number(raw.ts || Date.now()),
            };
            this.listeners.forEach((l) => l(msg));
        });
        wsClient.on('chat.rejected', (raw) => {
            console.warn('chat rejected', raw.reason);
        });
    }

    onMessage(listener: ChatListener) {
        this.listeners.add(listener);
        return () => this.listeners.delete(listener);
    }

    async loadHistory(channel: 'world' | 'friend', peerId = 0, limit = 50) {
        const q = new URLSearchParams({ channel, limit: String(limit) });
        if (peerId) q.set('peerId', String(peerId));
        return http.get<ChatMsg[]>(`/api/v1/chat/history?${q.toString()}`);
    }

    sendWorld(content: string) {
        wsClient.send({ type: 'chat.send', channel: 'world', content });
    }

    sendFriend(toId: number, content: string) {
        wsClient.send({ type: 'chat.send', channel: 'friend', toId, content });
    }
}

export const chatManager = new ChatManager();
