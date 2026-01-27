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

interface GolfEngineLink {
    label: string;
    vehicleId: number;
}

type GenerationCard =
    | { kind: 'vehicle'; vehicle: VehicleListItem }
    | {
          kind: 'golf_group';
          mk: string;
          image_url: string;
          generation_meta?: VehicleListItem['generation_meta'];
          engine_links: GolfEngineLink[];
          defaultVehicleId: number;
      };

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function ModelGenerations() {
    const navigate = useNavigate();
    const { brandName, modelName } = useParams();

    const decodedBrandName = brandName ? decodeURIComponent(brandName) : '';
    const decodedModelName = modelName ? decodeURIComponent(modelName) : '';

    const [vehicles, setVehicles] = useState<VehicleListItem[]>([]);
    const [loading, setLoading] = useState(true);

    const goBack = () => {
        if (window.history.length > 1) {
            navigate(-1);
            return;
        }
        navigate('/', { replace: true });
    };

    const inferGolfMk = (v: VehicleListItem) => {
        const isGolf = decodedBrandName.trim().toLowerCase() === 'volkswagen' && decodedModelName.trim().toLowerCase() === 'golf';
        if (!isGolf) return '';
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

    const generations = useMemo<GenerationCard[]>(() => {
        const filtered = vehicles.filter((v) => (v.model || '').trim().toLowerCase() === decodedModelName.toLowerCase());
        const sorted = filtered.slice().sort((a, b) => {
            const ay = (a.generation_meta?.start_year ?? parseInt(a.generation, 10)) || 0;
            const by = (b.generation_meta?.start_year ?? parseInt(b.generation, 10)) || 0;
            return by - ay;
        });

        const isGolf = decodedBrandName.trim().toLowerCase() === 'volkswagen' && decodedModelName.trim().toLowerCase() === 'golf';
        if (!isGolf) return sorted.map((v) => ({ kind: 'vehicle', vehicle: v }));

        const mkOrder: Record<string, number> = { Mk8: 3, 'Mk7.5': 2, Mk7: 1 };

        const grouped = new Map<
            string,
            {
                mk: string;
                defaultVehicleId: number;
                image_url: string;
                generation_meta?: VehicleListItem['generation_meta'];
                optionToVehicle: Map<string, number>;
            }
        >();

        const pickMin = (a?: number, b?: number) => {
            const c = [a, b].filter((x): x is number => typeof x === 'number' && Number.isFinite(x) && x > 0);
            if (!c.length) return undefined;
            return Math.min(...c);
        };

        const pickMax = (a?: number, b?: number) => {
            const c = [a, b].filter((x): x is number => typeof x === 'number' && Number.isFinite(x) && x > 0);
            if (!c.length) return undefined;
            return Math.max(...c);
        };

        for (const v of sorted) {
            const mk = inferGolfMk(v);
            if (!mk) continue;
            const current = grouped.get(mk);

            const options = Array.isArray(v.engine_options) ? v.engine_options : [];

            if (!current) {
                const optionToVehicle = new Map<string, number>();
                for (const opt of options) {
                    const label = (opt || '').trim();
                    if (!label) continue;
                    if (!optionToVehicle.has(label)) optionToVehicle.set(label, v.id);
                }
                grouped.set(mk, {
                    mk,
                    defaultVehicleId: v.id,
                    image_url: v.image_url,
                    generation_meta: v.generation_meta,
                    optionToVehicle,
                });
                continue;
            }

            // Merge option -> representative vehicle id mapping
            for (const opt of options) {
                const label = (opt || '').trim();
                if (!label) continue;
                if (!current.optionToVehicle.has(label)) current.optionToVehicle.set(label, v.id);
            }

            // Prefer having an image
            if (!current.image_url && v.image_url) current.image_url = v.image_url;

            current.generation_meta = {
                ...(current.generation_meta || {}),
                start_year: pickMin(current.generation_meta?.start_year, v.generation_meta?.start_year),
                end_year: pickMax(current.generation_meta?.end_year, v.generation_meta?.end_year),
                is_facelift: (current.generation_meta?.is_facelift || false) || (v.generation_meta?.is_facelift || false),
            };

            grouped.set(mk, current);
        }

        const out: GenerationCard[] = Array.from(grouped.values())
            .sort((a, b) => (mkOrder[b.mk] || 0) - (mkOrder[a.mk] || 0))
            .map((g) => {
                const engine_links: GolfEngineLink[] = Array.from(g.optionToVehicle.entries())
                    .sort((a, b) => a[0].localeCompare(b[0]))
                    .map(([label, vehicleId]) => ({ label, vehicleId }));

                return {
                    kind: 'golf_group',
                    mk: g.mk,
                    image_url: g.image_url,
                    generation_meta: g.generation_meta,
                    engine_links,
                    defaultVehicleId: g.defaultVehicleId,
                };
            });

        return out;
    }, [decodedBrandName, decodedModelName, vehicles]);

    const formatGeneration = (mk: string, meta?: VehicleListItem['generation_meta'], generationFallback?: string) => {
        const start = meta?.start_year;
        const end = meta?.end_year;
        if (mk && start && end && start !== end) return `${mk} (${start}–${end})`;
        if (mk && start && end && start === end) return `${mk} (${start})`;
        if (mk && start && !end) return `${mk} (${start}–)`;
        if (mk && !start && end) return `${mk} (–${end})`;
        if (mk) return mk;
        if (start && end && start !== end) return `${start}–${end}`;
        if (start && end && start === end) return `${start}`;
        const g = (generationFallback || '').trim();
        if (!g) return '—';
        if (g.toLowerCase().includes('http://') || g.toLowerCase().includes('https://')) return '—';
        return g;
    };

    if (loading) {
        return <div className="rounded-2xl bg-white border border-border p-8 text-center text-text-muted shadow-sm">Loading generations...</div>;
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
                        ← {decodedBrandName}
                    </button>

                    <button
                        type="button"
                        onClick={() => navigate('/', { replace: true })}
                        className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1.5 text-xs font-semibold text-text-main hover:bg-background hover:border-primary/30 hover:text-primary transition-colors"
                    >
                        Ana Sayfa
                    </button>
                </div>

                <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-text-main">
                    {decodedBrandName} {decodedModelName}
                </h1>
                <p className="text-sm text-text-muted">Nesil / yıl seç.</p>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                {generations.map((card) => {
                    if (card.kind === 'vehicle') {
                        const v = card.vehicle;
                        const mk = inferGolfMk(v);
                        return (
                            <Link
                                key={v.id}
                                to={`/vehicle/${v.id}/powertrain`}
                                className="rounded-2xl bg-white border border-border shadow-sm hover:shadow-card transition-shadow overflow-hidden group"
                            >
                                <div className="relative h-44 bg-background flex items-center justify-center p-4">
                                    {v.image_url ? (
                                        <img
                                            src={v.image_url}
                                            alt={`${v.brand} ${v.model}`}
                                            className="w-full h-full object-contain group-hover:scale-105 transition-transform duration-300"
                                        />
                                    ) : (
                                        <div className="text-text-muted text-sm">No Image</div>
                                    )}
                                </div>
                                <div className="p-5 border-t border-border/50 space-y-2">
                                    <div className="flex items-center justify-between gap-2">
                                        <div className="text-base font-bold text-text-main">
                                            {formatGeneration(mk, v.generation_meta, v.generation)}
                                        </div>
                                        {v.generation_meta?.is_facelift && mk !== 'Mk7.5' ? (
                                            <span className="inline-flex items-center rounded-full bg-background px-3 py-1 text-xs font-semibold text-text-muted">
                                                Facelift
                                            </span>
                                        ) : null}
                                    </div>
                                    <div className="space-y-2">
                                        <div className="text-xs font-semibold text-text-muted uppercase tracking-wider">Motor seçenekleri</div>
                                        {Array.isArray(v.engine_options) && v.engine_options.length > 0 ? (
                                            <div className="flex flex-wrap gap-2">
                                                {v.engine_options.slice(0, 10).map((opt) => (
                                                    <span
                                                        key={opt}
                                                        className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1 text-xs font-semibold text-text-main"
                                                    >
                                                        {opt}
                                                    </span>
                                                ))}
                                                {v.engine_options.length > 10 ? (
                                                    <span className="inline-flex items-center rounded-full bg-background border border-border px-3 py-1 text-xs font-semibold text-text-muted">
                                                        +{v.engine_options.length - 10}
                                                    </span>
                                                ) : null}
                                            </div>
                                        ) : (
                                            <div className="text-sm text-text-muted">—</div>
                                        )}
                                    </div>
                                </div>
                            </Link>
                        );
                    }

                    const mk = card.mk;
                    return (
                        <div
                            key={mk}
                            className="rounded-2xl bg-white border border-border shadow-sm hover:shadow-card transition-shadow overflow-hidden group"
                        >
                            <div className="relative h-44 bg-background flex items-center justify-center p-4">
                                {card.image_url ? (
                                    <img
                                        src={card.image_url}
                                        alt={`${decodedBrandName} ${decodedModelName} ${mk}`}
                                        className="w-full h-full object-contain group-hover:scale-105 transition-transform duration-300"
                                    />
                                ) : (
                                    <div className="text-text-muted text-sm">No Image</div>
                                )}
                            </div>
                            <div className="p-5 border-t border-border/50 space-y-3">
                                <div className="flex items-center justify-between gap-2">
                                    <div className="text-base font-bold text-text-main">
                                        {formatGeneration(mk, card.generation_meta, '')}
                                    </div>
                                    <Link
                                        to={`/brand/${encodeURIComponent(decodedBrandName)}/model/${encodeURIComponent(decodedModelName)}/mk/${encodeURIComponent(mk)}`}
                                        className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1 text-xs font-semibold text-text-main hover:bg-background hover:border-primary/30 hover:text-primary transition-colors"
                                    >
                                        Detaylar
                                    </Link>
                                </div>

                                <div className="space-y-2">
                                    <div className="text-xs font-semibold text-text-muted uppercase tracking-wider">Motor seçenekleri</div>
                                    {card.engine_links.length > 0 ? (
                                        <div className="flex flex-wrap gap-2">
                                            {card.engine_links.slice(0, 12).map((e) => (
                                                <Link
                                                    key={`${mk}-${e.label}-${e.vehicleId}`}
                                                    to={`/vehicle/${e.vehicleId}/powertrain`}
                                                    className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1 text-xs font-semibold text-text-main hover:bg-background hover:border-primary/30 hover:text-primary transition-colors"
                                                >
                                                    {e.label}
                                                </Link>
                                            ))}
                                            {card.engine_links.length > 12 ? (
                                                <span className="inline-flex items-center rounded-full bg-background border border-border px-3 py-1 text-xs font-semibold text-text-muted">
                                                    +{card.engine_links.length - 12}
                                                </span>
                                            ) : null}
                                        </div>
                                    ) : (
                                        <div className="text-sm text-text-muted">—</div>
                                    )}
                                </div>
                            </div>
                        </div>
                    );
                })}
            </div>

            {generations.length === 0 ? (
                <div className="text-center py-12 bg-white border border-border rounded-2xl text-text-muted shadow-sm">
                    Bu model için nesil/yıl kaydı bulunamadı.
                </div>
            ) : null}
        </div>
    );
}

export default ModelGenerations;
