import type { Engine } from '../types/engine';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export const engineService = {
    // Get all engines with optional filtering
    async getAll(filters?: {
        search?: string;
        fuelType?: string;
        manufacturer?: string;
    }): Promise<Engine[]> {
        const params = new URLSearchParams();
        if (filters?.search) params.append('search', filters.search);
        if (filters?.fuelType) params.append('fuelType', filters.fuelType);
        if (filters?.manufacturer) params.append('manufacturer', filters.manufacturer);

        const url = `${API_BASE_URL}/engines${params.toString() ? `?${params.toString()}` : ''}`;
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error('Failed to fetch engines');
        }

        return response.json();
    },

    // Get engine by ID
    async getById(id: number): Promise<Engine> {
        const response = await fetch(`${API_BASE_URL}/engines/${id}`);

        if (!response.ok) {
            throw new Error('Failed to fetch engine');
        }

        return response.json();
    },

    // Get engine by code
    async getByCode(code: string): Promise<Engine> {
        const response = await fetch(`${API_BASE_URL}/engines/code/${code}`);

        if (!response.ok) {
            throw new Error('Failed to fetch engine');
        }

        return response.json();
    },

    // Create new engine
    async create(engine: Partial<Engine>): Promise<Engine> {
        const response = await fetch(`${API_BASE_URL}/engines`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(engine),
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Failed to create engine');
        }

        return response.json();
    },

    // Update engine
    async update(id: number, engine: Partial<Engine>): Promise<Engine> {
        const response = await fetch(`${API_BASE_URL}/engines/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(engine),
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Failed to update engine');
        }

        return response.json();
    },

    // Delete engine
    async delete(id: number): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/engines/${id}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            throw new Error('Failed to delete engine');
        }
    },
};
