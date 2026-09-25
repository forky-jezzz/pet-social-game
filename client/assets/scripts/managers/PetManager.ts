import { http } from '../net/HttpClient';
import { PetInfo } from '../net/types';
import { authManager } from './AuthManager';

export class PetManager {
    async create(name: string) {
        const pet = await http.post<PetInfo>('/api/v1/pet', { name });
        authManager.pet = pet;
        return pet;
    }

    async feed(foodCode = 'food_kibble') {
        const pet = await http.post<PetInfo>('/api/v1/pet/feed', { foodCode });
        authManager.pet = pet;
        return pet;
    }

    async outfit(hat?: string, clothes?: string) {
        const body: Record<string, string> = {};
        if (hat) body.hat = hat;
        if (clothes) body.clothes = clothes;
        const pet = await http.post<PetInfo>('/api/v1/pet/outfit', body);
        authManager.pet = pet;
        return pet;
    }

    parseOutfit(pet: PetInfo | null) {
        if (!pet?.outfitJson) return { hat: 'hat_basic', clothes: 'clothes_basic' };
        try {
            return JSON.parse(pet.outfitJson) as { hat: string; clothes: string };
        } catch {
            return { hat: 'hat_basic', clothes: 'clothes_basic' };
        }
    }
}

export const petManager = new PetManager();
