import { _decorator, Component } from 'cc';
import { GameConfig } from '../config/GameConfig';
import { authManager } from '../managers/AuthManager';
import { chatManager } from '../managers/ChatManager';

const { ccclass } = _decorator;

/**
 * Attach to a persistent root node (e.g. Canvas).
 * Bootstraps managers; panels handle UI interactions.
 */
@ccclass('GameBootstrap')
export class GameBootstrap extends Component {
    start() {
        console.log('[PetSocial] api=', GameConfig.apiBase, 'platform=', authManager.platform.name);
        chatManager.start();
        // If token cached, refresh profile silently.
        if (authManager.platform.getStorage('pet_token')) {
            authManager.refreshMe().catch(() => {
                /* need re-login */
            });
        }
    }
}
