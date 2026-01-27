import { useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

interface VehicleListItem {
    id: number;
    brand: string;
    model: string;
    generation: string;
    image_url: string;
    generation_meta?: {
        start_year?: number;
        end_year?: number;
        is_facelift?: boolean;
        market?: string;
    };
    engine_options?: string[];
}

interface EngineLink {
    label: string;
    vehicleId: number;
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function GolfMkGroup() {
    const navigate = useNavigate();
    const { brandName, modelName, mk } = useParams();

    const decodedBrandName = brandName ? decodeURIComponent(brandName) : '';
    const decodedModelName = modelName ? decodeURIComponent(modelName) : '';
    const decodedMk = mk ? decodeURIComponent(mk) : '';

    const [vehicles, setVehicles] = useState<VehicleListItem[]>([]);
    const [loading, setLoading] = useState(true);

    const goBack = () => {
        if (window.history.length > 1) {
            navigate(-1);
            return;
        }
        navigate('/', { replace: true });
    };

    const normalizeMk = (raw: string) => {
        const s = raw.trim().toLowerCase();
        if (s === 'mk7' || s === '7') return 'Mk7';
        if (s === 'mk8' || s === '8') return 'Mk8';
        if (s === 'mk7.5' || s === 'mk7,5' || s === '7.5' || s === '7,5') return 'Mk7.5';
        if (s.includes('7.5') || s.includes('7,5')) return 'Mk7.5';
        if (s.includes('mk7')) return 'Mk7';
        if (s.includes('mk8')) return 'Mk8';
        return raw.trim();
    };

    const inferGolfMk = (v: VehicleListItem) => {
        const meta = v.generation_meta;
        const start = meta?.start_year;
        const end = meta?.end_year;
        const g = (v.generation || '').toLowerCase();

        const hasMk8 = /\bmk\s*8\b/.test(g);
        const hasMk7 = /\bmk\s*7\b/.test(g);

        if (meta?.is_facelift || g.includes('facelift') || g.includes('7.5') || g.includes('7,5')) return 'Mk7.5';
        if ((start && start >= 2020) || (end && end >= 2020) || hasMk8) return 'Mk8';
        if ((start && start >= 2012) || (end && end >= 2012) || hasMk7) return 'Mk7';
        return '';
    };

    useEffect(() => {
        setLoading(true);
        fetch(`${API_BASE_URL}/api/vehicles?brand=${encodeURIComponent(decodedBrandName)}`)
            .then((res) => res.json())
            .then((data) => {
                setVehicles(Array.isArray(data) ? data : []);
                setLoading(false);
            })
            .catch((err) => {
                console.error('Failed to fetch vehicles:', err);
                setVehicles([]);
                setLoading(false);
            });
    }, [decodedBrandName]);

    const group = useMemo(() => {
        const isGolf = decodedBrandName.trim().toLowerCase() === 'volkswagen' && decodedModelName.trim().toLowerCase() === 'golf';
        const targetMk = normalizeMk(decodedMk);
        if (!isGolf) return null;

        const filtered = vehicles.filter((v) => (v.model || '').trim().toLowerCase() === decodedModelName.toLowerCase());
        const inGroup = filtered.filter((v) => inferGolfMk(v) === targetMk);

        const engineMap = new Map<string, number>();
        for (const v of inGroup) {
            const options = Array.isArray(v.engine_options) ? v.engine_options : [];
            for (const opt of options) {
                const label = (opt || '').trim();
                if (!label) continue;
                if (!engineMap.has(label)) engineMap.set(label, v.id);
            }
        }

        const engineLinks: EngineLink[] = Array.from(engineMap.entries())
            .sort((a, b) => a[0].localeCompare(b[0]))
            .map(([label, vehicleId]) => ({ label, vehicleId }));

        const image_url = inGroup.find((v) => v.image_url)?.image_url || '';
        const years = inGroup
            .map((v) => ({ start: v.generation_meta?.start_year, end: v.generation_meta?.end_year }))
            .filter((x) => typeof x.start === 'number' || typeof x.end === 'number');

        const startCandidates = years.map((x) => x.start).filter((x): x is number => typeof x === 'number' && Number.isFinite(x) && x > 0);
        const endCandidates = years.map((x) => x.end).filter((x): x is number => typeof x === 'number' && Number.isFinite(x) && x > 0);

        const start_year = startCandidates.length ? Math.min(...startCandidates) : undefined;
        const end_year = endCandidates.length ? Math.max(...endCandidates) : undefined;

        return {
            mk: targetMk,
            image_url,
            start_year,
            end_year,
            engineLinks,
            vehicleCount: inGroup.length,
        };
    }, [decodedBrandName, decodedMk, decodedModelName, vehicles]);

    const title = useMemo(() => {
        if (!group) return `${decodedBrandName} ${decodedModelName}`.trim();
        if (group.start_year && group.end_year && group.start_year !== group.end_year) return `${decodedBrandName} ${decodedModelName} • ${group.mk} (${group.start_year}–${group.end_year})`;
        if (group.start_year && group.end_year && group.start_year === group.end_year) return `${decodedBrandName} ${decodedModelName} • ${group.mk} (${group.start_year})`;
        if (group.start_year && !group.end_year) return `${decodedBrandName} ${decodedModelName} • ${group.mk} (${group.start_year}–)`;
        if (!group.start_year && group.end_year) return `${decodedBrandName} ${decodedModelName} • ${group.mk} (–${group.end_year})`;
        return `${decodedBrandName} ${decodedModelName} • ${group.mk}`;
    }, [decodedBrandName, decodedModelName, group]);

    if (loading) {
        return <div className="rounded-2xl bg-white border border-border p-8 text-center text-text-muted shadow-sm">Loading...</div>;
    }

    if (!group) {
        return (
            <div className="rounded-2xl bg-white border border-border p-8 text-center text-text-muted shadow-sm">
                Bu sayfa yalnızca Volkswagen Golf Mk grupları için kullanılabilir.
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="space-y-2">
                <div className="flex flex-wrap items-center gap-3">
                    <button
                        type="button"
                        onClick={goBack}
                        className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1.5 text-xs font-semibold text-text-main hover:bg-background hover:border-primary/30 hover:text-primary transition-colors"
                    >
                        ← Geri
                    </button>
                    <button
                        type="button"
                        onClick={() => navigate('/', { replace: true })}
                        className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1.5 text-xs font-semibold text-text-main hover:bg-background hover:border-primary/30 hover:text-primary transition-colors"
                    >
                        Ana Sayfa
                    </button>
                </div>

                <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-text-main">{title}</h1>
                <p className="text-sm text-text-muted">Bu Mk grubundaki tüm motor seçenekleri.</p>
            </div>

            <div className="rounded-2xl bg-white border border-border shadow-sm overflow-hidden">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-0">
                    <div className="bg-background p-6 flex items-center justify-center">
                        {group.image_url ? (
                            <img src={group.image_url} alt={title} className="w-full h-64 object-contain" />
                        ) : (
                            <div className="text-text-muted text-sm">No Image</div>
                        )}
                    </div>
                    <div className="p-6 space-y-4">
                        <div className="text-xs text-text-muted">{group.vehicleCount} kayıt</div>

                        <div className="space-y-2">
                            <div className="text-xs font-semibold text-text-muted uppercase tracking-wider">Motor seçenekleri</div>
                            {group.engineLinks.length > 0 ? (
                                <div className="flex flex-wrap gap-2">
                                    {group.engineLinks.map((e) => (
                                        <Link
                                            key={`${group.mk}-${e.label}-${e.vehicleId}`}
                                            to={`/vehicle/${e.vehicleId}/powertrain`}
                                            className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1 text-xs font-semibold text-text-main hover:bg-background hover:border-primary/30 hover:text-primary transition-colors"
                                        >
                                            {e.label}
                                        </Link>
                                    ))}
                                </div>
                            ) : (
                                <div className="text-sm text-text-muted">—</div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default GolfMkGroup;
