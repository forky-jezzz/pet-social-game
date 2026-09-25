export interface IPlatform {
    name: 'alipay' | 'wechat' | 'dev';
    login(): Promise<{ code: string; nickname?: string; avatar?: string }>;
    getStorage(key: string): string | null;
    setStorage(key: string, value: string): void;
    toast(msg: string): void;
}

declare const my: any;
declare const wx: any;

export class DevPlatform implements IPlatform {
    name: 'dev' = 'dev';
    async login() {
        let code = this.getStorage('dev_code');
        if (!code) {
            code = `u${Date.now().toString(36)}`;
            this.setStorage('dev_code', code);
        }
        return { code, nickname: `Dev_${code.slice(-4)}` };
    }
    getStorage(key: string) {
        try {
            return localStorage.getItem(key);
        } catch {
            return null;
        }
    }
    setStorage(key: string, value: string) {
        try {
            localStorage.setItem(key, value);
        } catch {
            /* ignore */
        }
    }
    toast(msg: string) {
        console.log('[toast]', msg);
    }
}

export class AlipayPlatform implements IPlatform {
    name: 'alipay' = 'alipay';
    login(): Promise<{ code: string; nickname?: string; avatar?: string }> {
        return new Promise((resolve, reject) => {
            my.getAuthCode({
                scopes: 'auth_user',
                success: (res: { authCode: string }) => resolve({ code: res.authCode }),
                fail: (err: unknown) => reject(err),
            });
        });
    }
    getStorage(key: string) {
        try {
            return my.getStorageSync({ key })?.data ?? null;
        } catch {
            return null;
        }
    }
    setStorage(key: string, value: string) {
        my.setStorageSync({ key, data: value });
    }
    toast(msg: string) {
        my.showToast({ content: msg });
    }
}

export class WeChatPlatform implements IPlatform {
    name: 'wechat' = 'wechat';
    login(): Promise<{ code: string; nickname?: string; avatar?: string }> {
        return new Promise((resolve, reject) => {
            wx.login({
                success: (res: { code: string }) => resolve({ code: res.code }),
                fail: (err: unknown) => reject(err),
            });
        });
    }
    getStorage(key: string) {
        try {
            return wx.getStorageSync(key) || null;
        } catch {
            return null;
        }
    }
    setStorage(key: string, value: string) {
        wx.setStorageSync(key, value);
    }
    toast(msg: string) {
        wx.showToast({ title: msg, icon: 'none' });
    }
}

export function createPlatform(hint: 'auto' | 'alipay' | 'wechat' | 'dev' = 'auto'): IPlatform {
    if (hint === 'dev') return new DevPlatform();
    if (hint === 'alipay') return new AlipayPlatform();
    if (hint === 'wechat') return new WeChatPlatform();
    if (typeof my !== 'undefined' && my.getAuthCode) return new AlipayPlatform();
    if (typeof wx !== 'undefined' && wx.login) return new WeChatPlatform();
    return new DevPlatform();
}
