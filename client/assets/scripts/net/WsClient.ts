import { GameConfig } from '../config/GameConfig';
import { http } from './HttpClient';

type Handler = (msg: Record<string, unknown>) => void;

export class WsClient {
    private ws: WebSocket | null = null;
    private handlers = new Map<string, Set<Handler>>();
    private reconnectMs = 1500;
    private closedByUser = false;
    private heartbeatTimer: ReturnType<typeof setInterval> | null = null;

    connect() {
        this.closedByUser = false;
        const token = http.getToken();
        if (!token) {
            throw new Error('login first');
        }
        if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) {
            return;
        }
        const url = `${GameConfig.wsBase}/ws?token=${encodeURIComponent(token)}`;
        this.ws = new WebSocket(url);
        this.ws.onopen = () => {
            this.reconnectMs = 1500;
            this.startHeartbeat();
        };
        this.ws.onmessage = (ev) => {
            try {
                const msg = JSON.parse(String(ev.data)) as Record<string, unknown>;
                const type = String(msg.type || '');
                this.emit(type, msg);
                this.emit('*', msg);
            } catch {
                /* ignore */
            }
        };
        this.ws.onclose = () => {
            this.stopHeartbeat();
            this.ws = null;
            if (!this.closedByUser) {
                setTimeout(() => this.connect(), this.reconnectMs);
                this.reconnectMs = Math.min(this.reconnectMs * 1.5, 15000);
            }
        };
        this.ws.onerror = () => {
            /* onclose will retry */
        };
    }

    close() {
        this.closedByUser = true;
        this.stopHeartbeat();
        this.ws?.close();
        this.ws = null;
    }

    on(type: string, handler: Handler) {
        if (!this.handlers.has(type)) this.handlers.set(type, new Set());
        this.handlers.get(type)!.add(handler);
        return () => this.handlers.get(type)?.delete(handler);
    }

    send(payload: Record<string, unknown>) {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            throw new Error('ws not connected');
        }
        this.ws.send(JSON.stringify(payload));
    }

    private emit(type: string, msg: Record<string, unknown>) {
        this.handlers.get(type)?.forEach((h) => h(msg));
    }

    private startHeartbeat() {
        this.stopHeartbeat();
        this.heartbeatTimer = setInterval(() => {
            try {
                this.send({ type: 'ping' });
            } catch {
                /* ignore */
            }
        }, 20000);
    }

    private stopHeartbeat() {
        if (this.heartbeatTimer) {
            clearInterval(this.heartbeatTimer);
            this.heartbeatTimer = null;
        }
    }
}

export const wsClient = new WsClient();
