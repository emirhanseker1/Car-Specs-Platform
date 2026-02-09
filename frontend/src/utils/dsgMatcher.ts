import type { DSGVariant } from '../data/dsgVariants';
import type { Trim } from '../types';
import { DSG_VARIANTS as VARIANTS_DATA } from '../data/dsgVariants';

export const detectTransmissionVariant = (modelNameInput: string, trim: Trim, brandNameInput?: string): DSGVariant | null => {
    // 1. Normalize Inputs
    const brandName = (brandNameInput || '').toLowerCase();
    const modelName = (modelNameInput || '').toLowerCase();
    // const generation = (vehicle.generation || '').toLowerCase(); // Future use
    const trimName = (trim.name || '').toLowerCase();
    const transmissionType = (trim.transmission_type || '').toLowerCase();

    // 0. VAG Group Check
    // DSG/S-Tronic is specific to VAG.
    // Allow if brand is empty (fallback) or explicitly one of the VAG brands.
    const vagBrands = ['volkswagen', 'vw', 'audi', 'seat', 'skoda', 'porsche', 'cupra'];
    if (brandName && !vagBrands.some(b => brandName.includes(b))) {
        return null;
    }

    // Parse Specs
    const torque = trim.torque_nm || 0;
    // const hp = trim.power_hp || 0; // Future use
    const year = trim.year || new Date().getFullYear();
    const gears = parseInt(transmissionType.match(/\d/)?.[0] || '0') || 7; // Default to 7 if unknown but automatic

    // 0. Explicit Match (if data already has the code)
    const explicitMatch = VARIANTS_DATA.find(v => transmissionType.includes(v.code.toLowerCase()));
    if (explicitMatch) return explicitMatch;

    // Skip Manuals
    if (transmissionType.includes('manual') || transmissionType.includes('düz')) return null;

    // --- Heuristic Rules ---

    // Rule 1: Torque & Dry Clutch (DQ200)
    // Applies to: 1.0 TSI, 1.2 TSI, 1.4 TSI, 1.5 TSI, 1.6 TDI
    // Low torque (< 250 Nm), Transverse engine
    if (torque > 0 && torque <= 250) {
        // Exclude Hybrids (DQ400e)
        if (trimName.includes('e-tron') || trimName.includes('gte') || trimName.includes('hybrid')) {
            return VARIANTS_DATA.find(v => v.code === 'DQ400e') || null;
        }
        // Exclude older 6-speed automatics (non-DSG) if needed, but modern small engines are mostly DQ200
        return VARIANTS_DATA.find(v => v.code === 'DQ200') || null;
    }

    // Rule 2: High Torque Transverse (DQ250 vs DQ381 vs DQ500)
    // Applies to: 2.0 TDI, 2.0 TSI, S3, Golf R, Leon Cupra
    if (torque > 250) {
        const isTransverse =
            modelName.includes('a3') || modelName.includes('q2') || modelName.includes('q3') ||
            modelName.includes('golf') || modelName.includes('passat') ||
            modelName.includes('leon') || modelName.includes('octavia') ||
            modelName.includes('superb') || modelName.includes('kodiaq') ||
            modelName.includes('tiguan') || modelName.includes('t-roc');

        if (isTransverse) {
            // Rule 2a: Very High Torque / Commercial (DQ500)
            // Tiguan BiTDI, RS3, TTRS, Transporter
            if (torque >= 500 || modelName.includes('rs3') || modelName.includes('tt rs') || trimName.includes('bitdi') || modelName.includes('transporter')) {
                return VARIANTS_DATA.find(v => v.code === 'DQ500') || null;
            }

            // Rule 2b: Newer 7-Speed Wet (DQ381) -> Approx 2017+
            if (year >= 2017) {
                if (gears === 7) return VARIANTS_DATA.find(v => v.code === 'DQ381') || null;
            }

            // Rule 2c: Older High Torque (DQ250)
            // If year < 2017 and torque > 250, it's almost certainly DQ250 (6-speed)
            // even if data doesn't explicitly say "6 gears" (defaulted to 7).
            // DQ500 is already handled above for very high torque/specific models.
            if (year < 2017) {
                return VARIANTS_DATA.find(v => v.code === 'DQ250') || null;
            }

            // Fallback: If newer but not matched above (e.g. 2018 model with unknown gears)
            // Default to DQ381 for modern high-torque transverse
            return VARIANTS_DATA.find(v => v.code === 'DQ381') || null;
        }
    }

    // Rule 3: Longitudinal Engines (Audi A4, A5, A6, Q5, Q7...)
    // These use DL501 (Older) or DL382 (Newer) or ZF8HP (high torque/conventional)
    const isLongitudinal = !modelName.includes('a1') && !modelName.includes('a3') && !modelName.includes('q2') && !modelName.includes('q3') && !modelName.includes('tt');

    if (isLongitudinal) {
        // Rule 3a: ZF8HP (8-Speed Tiptronic)
        // Used in A6/A7/A8/Q7/Q8 (3.0 TDI/TFSI and up) and RS4/RS5
        // Usually 8 gears
        if (gears === 8 || transmissionType.includes('tiptronic')) {
            return VARIANTS_DATA.find(v => v.code === 'ZF8HP') || null;
        }

        // Rule 3b: S-Tronic 7-Speed Longitudinal
        // DL382 (Newer, ~2015/2016+ B9 chassis) vs DL501 (Older, B8 chassis)
        if (gears === 7) {
            // A4 B9 (2016+), A5 F5 (2017+), Q5 FY (2017+) -> DL382
            // A4 B8, A5 8T, Q5 8R -> DL501
            // Simple year cutoff
            if (year >= 2016) {
                return VARIANTS_DATA.find(v => v.code === 'DL382') || null;
            } else {
                return VARIANTS_DATA.find(v => v.code === 'DL501') || null;
            }
        }
    }

    return null;
};
