import { _decorator, Component, Label } from 'cc';
import { authManager } from '../managers/AuthManager';
import { petManager } from '../managers/PetManager';
import { PetActor } from '../pet/PetActor';

const { ccclass, property } = _decorator;

@ccclass('HomePanel')
export class HomePanel extends Component {
    @property(Label)
    petNameLabel: Label | null = null;

    @property(Label)
    statsLabel: Label | null = null;

    @property(Label)
    outfitLabel: Label | null = null;

    @property(PetActor)
    petActor: PetActor | null = null;

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
        this.petActor?.applyStats(pet.hunger, pet.mood);
        this.petActor?.applyOutfit(outfit.hat, outfit.clothes);
    }

    async onClickFeed() {
        try {
            this.petActor?.playEat();
            await petManager.feed('food_kibble');
            this.refresh();
        } catch (e) {
            console.error(e);
        }
    }

    async onClickTreat() {
        try {
            this.petActor?.playEat(() => this.petActor?.playHappy());
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
            this.petActor?.playHappy();
        } catch (e) {
            console.error(e);
        }
    }

    async onClickOutfitCape() {
        try {
            await petManager.outfit(undefined, 'clothes_cape');
            this.refresh();
            this.petActor?.playHappy();
        } catch (e) {
            console.error(e);
        }
    }

    private setStats(msg: string) {
        if (this.statsLabel) this.statsLabel.string = msg;
    }
}
