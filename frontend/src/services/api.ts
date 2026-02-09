import type { FeaturedTrim, Vehicle } from '../types';

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Merkezi API Servisi
 * 
 * Tüm ağ istekleri (fetch) burada toplanır.
 * Bileşenler (Components) bu fonksiyonları çağırır.
 * Böylece URL değişirse sadece burayı güncellemek yeterlidir.
 */
export const api = {
    /**
     * Belirli bir marka için araç listesini (Modeller ve Nesiller) getirir.
     * @param brandName Marka adı (örn: "Audi")
     */
    getVehicles: async (brandName: string): Promise<Vehicle[]> => {
        try {
            const query = encodeURIComponent(brandName || 'Audi');
            const response = await fetch(`${API_BASE_URL}/api/vehicles?brand=${query}`);

            if (!response.ok) {
                throw new Error(`API Hatası: ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            console.error('getVehicles hatası:', error);
            throw error;
        }
    },

    /**
     * Ana sayfa için öne çıkan popüler modelleri getirir.
     */
    getFeaturedTrims: async (): Promise<FeaturedTrim[]> => {
        try {
            const response = await fetch(`${API_BASE_URL}/api/featured`);

            if (!response.ok) {
                throw new Error('Featured models yüklenemedi');
            }

            return await response.json();
        } catch (error) {
            console.error('getFeaturedTrims hatası:', error);
            throw error; // Hatayı çağıran yere fırlat ki orada yönetilebilsin (loading/error state)
        }
    },

    /**
     * Tüm markaları getirir.
     */
    getBrands: async (): Promise<any[]> => {
        const response = await fetch(`${API_BASE_URL}/api/brands`);
        if (!response.ok) throw new Error('Failed to load brands');
        return await response.json();
    },

    /**
     * Markaya ait modelleri getirir.
     */
    getModels: async (brandId: number): Promise<any> => {
        const response = await fetch(`${API_BASE_URL}/api/brands/${brandId}/models`);
        if (!response.ok) throw new Error('Failed to load models');
        return await response.json();
    },

    /**
     * Modele ait nesilleri getirir.
     */
    getGenerations: async (modelId: number): Promise<any> => {
        const response = await fetch(`${API_BASE_URL}/api/models/${modelId}/generations`);
        if (!response.ok) throw new Error('Failed to load generations');
        return await response.json();
    },

    /**
     * Nesile ait trimleri (paketleri) getirir.
     */
    getTrims: async (generationId: number): Promise<any> => {
        const response = await fetch(`${API_BASE_URL}/api/generations/${generationId}/trims`);
        if (!response.ok) throw new Error('Failed to load trims');
        return await response.json();
    },

    /**
     * Chatbot ile mesajlaşır.
     */
    sendChatMessage: async (message: string): Promise<{ response: string; error?: string }> => {
        try {
            const response = await fetch(`${API_BASE_URL}/api/chat`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ message }),
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.error || `Chat API Hatası: ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            console.error('sendChatMessage hatası:', error);
            throw error;
        }
    }
};
