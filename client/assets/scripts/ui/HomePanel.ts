import { _decorator, Component, Label } from 'cc';
import { authManager } from '../managers/AuthManager';
import { petManager } from '../managers/PetManager';

const { ccclass, property } = _decorator;

@ccclass('HomePanel')
export class HomePanel extends Component {
    @property(Label)
    petNameLabel: Label | null = null;

    @property(Label)
    statsLabel: Label | null = null;

    @property(Label)
    outfitLabel: Label | null = null;

    onEnable() {
        this.refresh();
    }

    refresh() {
        const pet = authManager.pet;
        if (!pet) {
            this.setStats('尚未创建宠物');
            return;
        }
        if (this.petNameLabel) this.petNameLabel.string = pet.name;
        this.setStats(`Lv.${pet.level}  饱食 ${pet.hunger}  心情 ${pet.mood}`);
        const outfit = petManager.parseOutfit(pet);
        if (this.outfitLabel) this.outfitLabel.string = `帽:${outfit.hat} 衣:${outfit.clothes}`;
    }

    async onClickFeed() {
        try {
            await petManager.feed('food_kibble');
            this.refresh();
        } catch (e) {
            console.error(e);
        }
    }

    async onClickTreat() {
        try {
            await petManager.feed('food_treat');
            this.refresh();
        } catch (e) {
            console.error(e);
        }
    }

    async onClickOutfitHat() {
        try {
            await petManager.outfit('hat_star');
            this.refresh();
        } catch (e) {
            console.error(e);
        }
    }

    async onClickOutfitCape() {
        try {
            await petManager.outfit(undefined, 'clothes_cape');
            this.refresh();
        } catch (e) {
            console.error(e);
        }
    }

    private setStats(msg: string) {
        if (this.statsLabel) this.statsLabel.string = msg;
    }
}
