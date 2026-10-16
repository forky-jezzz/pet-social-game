import { _decorator, Component, sp } from 'cc';

const { ccclass, property } = _decorator;

/** Spine animation names — keep in sync with docs/spine-pet.md */
export const PetAnims = {
    IdleBreath: 'idle_breath',
    Blink: 'blink',
    TailWag: 'tail_wag',
    Happy: 'happy',
    Eat: 'eat',
    Hungry: 'hungry',
    Sleepy: 'sleepy',
    Petted: 'petted',
} as const;

enum MainState {
    Idle = 'idle',
    Hungry = 'hungry',
    Sleepy = 'sleepy',
    Eat = 'eat',
}

/**
 * Drives pet Spine micro-motions: breath, blink, tail wag, and reaction clips.
 * Attach to the node that has sp.Skeleton. If skeleton is missing, no-ops safely.
 */
@ccclass('PetActor')
export class PetActor extends Component {
    @property(sp.Skeleton)
    skeleton: sp.Skeleton | null = null;

    @property
    blinkMinSec = 2.5;

    @property
    blinkMaxSec = 5;

    @property
    wagMinSec = 4;

    @property
    wagMaxSec = 8;

    private mainState: MainState = MainState.Idle;
    private blinkLeft = 3;
    private wagLeft = 6;
    private busyOverlay = false;

    onLoad() {
        if (!this.skeleton) {
            this.skeleton = this.getComponent(sp.Skeleton);
        }
    }

    start() {
        this.playIdle();
        this.scheduleNextBlink();
        this.scheduleNextWag();
    }

    update(dt: number) {
        if (!this.hasAnim()) return;
        if (this.mainState === MainState.Eat) return;

        this.blinkLeft -= dt;
        this.wagLeft -= dt;

        if (this.blinkLeft <= 0) {
            this.playOverlay(PetAnims.Blink, false);
            this.scheduleNextBlink();
        }
        if (this.wagLeft <= 0 && !this.busyOverlay) {
            this.playOverlay(PetAnims.TailWag, false);
            this.scheduleNextWag();
        }
    }

    /** Sync visual mood from server pet stats. */
    applyStats(hunger: number, mood: number) {
        if (this.mainState === MainState.Eat) return;
        if (mood < 30) {
            this.setMain(MainState.Sleepy, PetAnims.Sleepy, true);
        } else if (hunger < 30) {
            this.setMain(MainState.Hungry, PetAnims.Hungry, true);
        } else {
            this.setMain(MainState.Idle, PetAnims.IdleBreath, true);
        }
    }

    playIdle() {
        this.setMain(MainState.Idle, PetAnims.IdleBreath, true);
    }

    playEat(onDone?: () => void) {
        if (!this.hasAnim(PetAnims.Eat)) {
            this.playHappy();
            onDone?.();
            return;
        }
        this.setMain(MainState.Eat, PetAnims.Eat, false);
        this.onceComplete(0, () => {
            this.playHappy();
            onDone?.();
        });
    }

    playHappy() {
        this.playOneShot(PetAnims.Happy, 2);
        this.setMain(MainState.Idle, PetAnims.IdleBreath, true);
    }

    playPetted() {
        this.playOneShot(PetAnims.Petted, 2);
    }

    /** Apply outfit skins/attachments when skins exist. */
    applyOutfit(hat?: string, clothes?: string) {
        const sk = this.skeleton;
        if (!sk) return;
        try {
            if (hat) sk.setSkin(hat);
            // If using combined skins, prefer setSlotsToSetupPose after skin change.
            sk.setSlotsToSetupPose();
            if (clothes && typeof (sk as any).setAttachment === 'function') {
                // Optional slot swap when not using full skins.
                try {
                    (sk as any).setAttachment('clothes', clothes);
                } catch {
                    /* skin-only pipeline */
                }
            }
        } catch (e) {
            console.warn('[PetActor] outfit skin missing', hat, clothes, e);
        }
    }

    private setMain(state: MainState, anim: string, loop: boolean) {
        this.mainState = state;
        if (!this.hasAnim(anim)) return;
        this.skeleton!.setAnimation(0, anim, loop);
    }

    private playOverlay(anim: string, loop: boolean) {
        if (!this.hasAnim(anim)) return;
        this.busyOverlay = true;
        this.skeleton!.setAnimation(1, anim, loop);
        this.onceComplete(1, () => {
            this.busyOverlay = false;
        });
    }

    private playOneShot(anim: string, track: number) {
        if (!this.hasAnim(anim)) return;
        this.skeleton!.setAnimation(track, anim, false);
    }

    private onceComplete(track: number, cb: () => void) {
        const sk = this.skeleton;
        if (!sk) {
            cb();
            return;
        }
        const handler = (entry: sp.spine.TrackEntry | null) => {
            if (!entry || entry.trackIndex !== track) return;
            sk.setCompleteListener(null as any);
            cb();
        };
        sk.setCompleteListener(handler as any);
    }

    private scheduleNextBlink() {
        this.blinkLeft = this.rand(this.blinkMinSec, this.blinkMaxSec);
    }

    private scheduleNextWag() {
        this.wagLeft = this.rand(this.wagMinSec, this.wagMaxSec);
    }

    private rand(min: number, max: number) {
        return min + Math.random() * (max - min);
    }

    private hasAnim(name?: string) {
        const sk = this.skeleton;
        if (!sk || !sk.skeletonData) return false;
        if (!name) return true;
        try {
            return !!sk.findAnimation(name);
        } catch {
            return false;
        }
    }
}
