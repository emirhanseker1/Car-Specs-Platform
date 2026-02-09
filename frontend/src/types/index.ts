export interface Brand {
    id: number;
    name: string;
    logo_url?: string;
}

export interface Model {
    id: number;
    name: string;
    brand_id?: number;
    brand?: Brand;
}

export interface Vehicle {
    id: number;
    brand: string;
    model: string;
    generation: string;
    image_url: string;
    engine_options?: string[];
    model_id?: number;
    // Helper property often used in UI logic
    generationCount?: number;
}

export interface FeaturedTrim {
    id: number;
    name: string;
    year: number;
    image_url?: string;
    fuel_type?: string;
    transmission_type?: string;
    power_hp?: number;
    model?: {
        id: number;
        name: string;
        brand?: {
            id: number;
            name: string;
        };
    };
}

export interface Spec {
    category: string;
    name: string;
    value: string;
}



export interface Trim {
    id: number;
    name: string;
    image_url?: string;
    // Common specs used directly in Quick Stats
    acceleration_0_100?: number;
    power_hp?: number;
    torque_nm?: number;
    curb_weight_kg?: number;
    transmission_type?: string;
    // Allow dynamic access for other specs
    [key: string]: any;
}

export interface DetailResponse {
    vehicle: Vehicle;
    trims: Trim[];
}
