export interface UserInfo {
    id: number;
    openId: string;
    platform: string;
    nickname: string;
    avatar: string;
}

export interface PetInfo {
    id: number;
    userId: number;
    name: string;
    level: number;
    hunger: number;
    mood: number;
    outfitJson: string;
}

export interface LoginResult {
    token: string;
    user: UserInfo;
    pet?: PetInfo | null;
}

export interface FriendView {
    userId: number;
    nickname: string;
    avatar: string;
    status: string;
    online: boolean;
}

export interface ChatMsg {
    id?: number;
    channel: string;
    fromId: number;
    fromName: string;
    toId?: number;
    content: string;
    ts: number;
}

export interface VisitHome {
    host: UserInfo;
    pet?: PetInfo | null;
    online: boolean;
}
