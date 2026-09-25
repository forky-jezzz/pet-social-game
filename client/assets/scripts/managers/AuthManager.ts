import { createPlatform, IPlatform } from '../platform/IPlatform';
import { GameConfig } from '../config/GameConfig';
import { http } from '../net/HttpClient';
import { LoginResult, PetInfo, UserInfo } from '../net/types';
import { wsClient } from '../net/WsClient';

const TOKEN_KEY = 'pet_token';

export class AuthManager {
    platform: IPlatform;
    user: UserInfo | null = null;
    pet: PetInfo | null = null;

    constructor() {
        this.platform = createPlatform(GameConfig.platform);
        const cached = this.platform.getStorage(TOKEN_KEY);
        if (cached) http.setToken(cached);
    }

    async login() {
        const cred = await this.platform.login();
        const data = await http.post<LoginResult>('/api/v1/auth/login', {
            platform: this.platform.name,
            code: cred.code,
            nickname: cred.nickname || '',
            avatar: cred.avatar || '',
        });
        http.setToken(data.token);
        this.platform.setStorage(TOKEN_KEY, data.token);
        this.user = data.user;
        this.pet = data.pet || null;
        wsClient.connect();
        return data;
    }

    async refreshMe() {
        const data = await http.get<{ user: UserInfo; pet?: PetInfo | null }>('/api/v1/me');
        this.user = data.user;
        this.pet = data.pet || null;
        return data;
    }
}

export const authManager = new AuthManager();
