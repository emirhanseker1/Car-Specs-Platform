// Engine types for the engine codes feature

export interface Engine {
    id: number;
    code: string;
    name?: string;
    manufacturer?: string;
    engineType?: string;
    fuelType?: string;
    displacementCC?: number;
    cylinders?: number;
    cylinderLayout?: string;
    valvesPerCylinder?: number;
    aspiration?: string;
    description?: string;
    technologyFeatures?: string; // JSON string
    productionStartYear?: number;
    productionEndYear?: number;
    commonProblems?: string; // JSON string
    solutions?: string; // JSON string
    maintenanceNotes?: string;
    reliabilityRating?: number;
    createdAt: string;
    updatedAt: string;
}

export interface EngineProblem {
    title: string;
    description: string;
    severity: 'low' | 'medium' | 'high';
}

export interface EngineSolution {
    problemTitle: string;
    solution: string;
    estimatedCost?: string;
}

export interface TechnologyFeature {
    name: string;
    description?: string;
}

// Helper functions to parse JSON fields
export function parseProblems(engine: Engine): EngineProblem[] {
    if (!engine.commonProblems) return [];
    try {
        return JSON.parse(engine.commonProblems);
    } catch {
        return [];
    }
}

export function parseSolutions(engine: Engine): EngineSolution[] {
    if (!engine.solutions) return [];
    try {
        return JSON.parse(engine.solutions);
    } catch {
        return [];
    }
}

export function parseTechnologyFeatures(engine: Engine): TechnologyFeature[] {
    if (!engine.technologyFeatures) return [];
    try {
        return JSON.parse(engine.technologyFeatures);
    } catch {
        return [];
    }
}

// Severity color functions
export function getSeverityColor(severity: 'low' | 'medium' | 'high'): string {
    switch (severity) {
        case 'high':
            return 'from-rose-500 to-red-600';
        case 'medium':
            return 'from-amber-500 to-orange-600';
        case 'low':
            return 'from-emerald-500 to-green-600';
        default:
            return 'from-gray-500 to-gray-600';
    }
}

export function getSeverityBadgeColor(severity: 'low' | 'medium' | 'high'): string {
    switch (severity) {
        case 'high':
            return 'bg-rose-500/10 text-rose-400 border-rose-500/20';
        case 'medium':
            return 'bg-amber-500/10 text-amber-400 border-amber-500/20';
        case 'low':
            return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20';
        default:
            return 'bg-gray-500/10 text-gray-400 border-gray-500/20';
    }
}
