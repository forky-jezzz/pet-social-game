import { GameConfig } from '../config/GameConfig';

export interface ApiBody<T = unknown> {
    code: number;
    msg: string;
    data?: T;
}

export class HttpClient {
    private token = '';

    setToken(token: string) {
        this.token = token;
    }

    getToken() {
        return this.token;
    }

    async request<T>(method: string, path: string, body?: unknown): Promise<T> {
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
        };
        if (this.token) {
            headers.Authorization = `Bearer ${this.token}`;
        }
        const res = await fetch(`${GameConfig.apiBase}${path}`, {
            method,
            headers,
            body: body === undefined ? undefined : JSON.stringify(body),
        });
        const json = (await res.json()) as ApiBody<T>;
        if (json.code !== 0) {
            throw new Error(json.msg || `http ${res.status}`);
        }
        return json.data as T;
    }

    get<T>(path: string) {
        return this.request<T>('GET', path);
    }

    post<T>(path: string, body?: unknown) {
        return this.request<T>('POST', path, body ?? {});
    }
}

export const http = new HttpClient();
