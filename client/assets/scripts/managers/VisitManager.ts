import { http } from '../net/HttpClient';
import { FriendView, PetInfo, VisitHome } from '../net/types';

export class VisitManager {
    listFriends() {
        return http.get<FriendView[]>('/api/v1/friends');
    }

    requestFriend(friendId: number) {
        return http.post<{ ok: boolean }>('/api/v1/friends/request', { friendId });
    }

    acceptFriend(friendId: number) {
        return http.post<{ ok: boolean }>('/api/v1/friends/accept', { friendId });
    }

    visit(userId: number) {
        return http.get<VisitHome>(`/api/v1/visit/${userId}`);
    }

    interact(userId: number, action: 'pet' | 'gift') {
        return http.post<PetInfo>(`/api/v1/visit/${userId}/interact`, { action });
    }
}

export const visitManager = new VisitManager();
