/** Global game / API endpoints. Change for prod builds. */
export const GameConfig = {
    /** REST base, no trailing slash */
    apiBase: 'http://127.0.0.1:8080',
    /** WebSocket base, no trailing slash */
    wsBase: 'ws://127.0.0.1:8080',
    /** Platform hint: auto | alipay | wechat | dev */
    platform: 'auto' as 'auto' | 'alipay' | 'wechat' | 'dev',
};
