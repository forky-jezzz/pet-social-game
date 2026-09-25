import { _decorator, Component, EditBox, Label } from 'cc';
import { visitManager } from '../managers/VisitManager';

const { ccclass, property } = _decorator;

@ccclass('VisitPanel')
export class VisitPanel extends Component {
    @property(EditBox)
    hostIdBox: EditBox | null = null;

    @property(Label)
    infoLabel: Label | null = null;

    private hostId = 0;

    async onClickVisit() {
        const id = Number(this.hostIdBox?.string || 0);
        if (!id) return;
        try {
            const home = await visitManager.visit(id);
            this.hostId = id;
            const pet = home.pet;
            const petText = pet
                ? `${pet.name} Lv.${pet.level} 饱食${pet.hunger} 心情${pet.mood}`
                : '对方还没有宠物';
            if (this.infoLabel) {
                this.infoLabel.string = `${home.host.nickname} 的家\n${petText}\n${home.online ? '在线' : '离线'}`;
            }
        } catch (e) {
            if (this.infoLabel) this.infoLabel.string = (e as Error).message;
        }
    }

    async onClickPet() {
        if (!this.hostId) return;
        const pet = await visitManager.interact(this.hostId, 'pet');
        if (this.infoLabel) this.infoLabel.string = `互动成功 心情=${pet.mood}`;
    }

    async onClickGift() {
        if (!this.hostId) return;
        const pet = await visitManager.interact(this.hostId, 'gift');
        if (this.infoLabel) this.infoLabel.string = `送礼成功 饱食=${pet.hunger} 心情=${pet.mood}`;
    }
}
