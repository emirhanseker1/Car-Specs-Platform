// Configuration for specific generation image styling
// Keys should be lowercase generation codes or "brand-code" for specific overrides.
export const GENERATION_STYLES: Record<string, string> = {
    // BMW
    'f40': 'scale-[1.5] group-hover:scale-[1.58]',
    'f20': 'scale-[1.5] group-hover:scale-[1.6]',
    'g20 lci': 'scale-[1.28] group-hover:scale-[1.3]',
    'g20': 'scale-[1.2] translate-x-[5px] translate-y-[5px] group-hover:scale-[1.25]',
    'f30 lci': 'scale-[1.6] translate-x-[0px] translate-y-[-5px] group-hover:scale-[1.68]',
    'f30': 'scale-[1.85] translate-x-[0px] translate-y-[70px] group-hover:scale-[1.92]',
    'e90 lci': 'scale-[1.95] translate-x-[30px] translate-y-[-32px] group-hover:scale-[2.02]',
    'e90': 'scale-[1.05] translate-y-[10px] group-hover:scale-[1.1]',
    'g60': 'scale-[2.8] translate-y-[-80px] group-hover:scale-[2.88]',
    'g30 lci': 'scale-[2.1] translate-y-[-50px] translate-x-[-30px] group-hover:scale-[2.2]',
    'g30': 'scale-[1.85] translate-y-[-20px] translate-x-[20px] group-hover:scale-[1.92]',
    'f10 lci': 'scale-[1] group-hover:scale-[1.04]',
    'f10': 'scale-[1.15]  group-hover:scale-[1.22]',
    'e60': 'scale-[1.15] translate-y-[5px] translate-x-[0px] group-hover:scale-[1.21]',

    // Audi A3
    '8y': 'scale-[1.2] group-hover:scale-[1.22]',
    '8v1': 'scale-[1.6] group-hover:scale-[1.68]',
    '8v2': 'scale-[1.4] translate-y-[10px] group-hover:scale-[1.48]',
    '8v': 'scale-[1.3] group-hover:scale-[1.4]',
    '8p': 'scale-[1] group-hover:scale-[1.08]',

    // Audi A4 (Specific)
    'audi-b7': 'scale-[1.4] group-hover:scale-[1.45]', // Placeholder style

    // Volkswagen Passat (Specific)
    'volkswagen-b8.5': 'scale-[2] translate-x-[0px] translate-y-[10px] group-hover:scale-[2.1]',
    'volkswagen-b8': 'scale-[1.85] translate-x-[0px] translate-y-[-30px] group-hover:scale-[1.92]',
    'volkswagen-b7': 'scale-[2.3] translate-x-[50px] translate-y-[-30px] group-hover:scale-[2.4]',
    'volkswagen-b6': 'scale-[1.2] translate-x-[0px] translate-y-[10px] group-hover:scale-[1.28]',
    'volkswagen-b5.5': 'scale-[1.6] translate-x-[0px] translate-y-[-30px] group-hover:scale-[1.67]',

    //VW Golf
    'mk8': 'scale-[1.7] translate-x-[0px] translate-y-[-10px] group-hover:scale-[1.78]',
    'mk7.5': 'scale-[1.65] translate-x-[0px] translate-y-[-10px] group-hover:scale-[1.72]',
    'mk7': 'scale-[1.8] translate-x-[10px] translate-y-[20px] group-hover:scale-[1.88]',
    'mk6': 'scale-[2] translate-x-[-60px] translate-y-[-60px] group-hover:scale-[2.1]',
    'mk5': 'scale-[1.8] translate-x-[0px] translate-y-[-20px] group-hover:scale-[1.67]',
};

export function getGenerationStyle(genCode: string, brandName?: string) {
    const code = genCode.toLowerCase();
    const brand = brandName?.toLowerCase();

    // Try brand-specific key first
    if (brand && GENERATION_STYLES[`${brand}-${code}`]) {
        return GENERATION_STYLES[`${brand}-${code}`];
    }

    // Fallback to code only
    if (GENERATION_STYLES[code]) {
        return GENERATION_STYLES[code];
    }

    return 'group-hover:scale-105';
}

// Force re-scan
