import { useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

interface Trim {
    id: number;
    name: string;
    powertrain_meta?: {
        engine_code?: string;
        fuel_type?: string;
        displacement_cc?: number;
        power_hp?: number;
        torque_nm?: number;
        transmission_type?: string;
        gears?: number;
        drive?: string;
    };
}

interface Vehicle {
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
}

interface DetailResponse {
    vehicle: Vehicle;
    trims: Trim[];
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function PowertrainSelect() {
    const navigate = useNavigate();
    const { id } = useParams();

    const [data, setData] = useState<DetailResponse | null>(null);
    const [loading, setLoading] = useState(true);

    const goBack = () => {
        if (window.history.length > 1) {
            navigate(-1);
            return;
        }
        navigate('/', { replace: true });
    };

    useEffect(() => {
        setLoading(true);
        fetch(`${API_BASE_URL}/api/vehicles/${id}`)
            .then((res) => res.json())
            .then((payload) => {
                const parsed = payload as Partial<DetailResponse>;
                if (!parsed || typeof parsed !== 'object' || !parsed.vehicle || !Array.isArray(parsed.trims)) {
                    setData(null);
                    setLoading(false);
                    return;
                }
                setData({ vehicle: parsed.vehicle as Vehicle, trims: parsed.trims as Trim[] });
                setLoading(false);
            })
            .catch((err) => {
                console.error('Failed to fetch powertrains:', err);
                setData(null);
                setLoading(false);
            });
    }, [id]);

    const formatGeneration = (v: Vehicle) => {
        const isGolf = (v.brand || '').trim().toLowerCase() === 'volkswagen' && (v.model || '').trim().toLowerCase() === 'golf';
        const meta = v.generation_meta;
        const start = meta?.start_year;
        const end = meta?.end_year;
        const g = (v.generation || '').toLowerCase();

        const hasMk8 = /\bmk\s*8\b/.test(g);
        const hasMk7 = /\bmk\s*7\b/.test(g);

        const mk = (() => {
            if (!isGolf) return '';
            if (meta?.is_facelift || g.includes('facelift') || g.includes('7.5') || g.includes('7,5')) return 'Mk7.5';
            if ((start && start >= 2020) || (end && end >= 2020) || hasMk8) return 'Mk8';
            if ((start && start >= 2012) || (end && end >= 2012) || hasMk7) return 'Mk7';
            return '';
        })();

        if (mk && start && end && start !== end) return `${mk} (${start}–${end})`;
        if (mk && start && end && start === end) return `${mk} (${start})`;
        if (mk && start && !end) return `${mk} (${start}–)`;
        if (mk && !start && end) return `${mk} (–${end})`;
        if (mk) return mk;
        if (start && end && start !== end) return `${start}–${end}`;
        if (start && end && start === end) return `${start}`;
        const gg = (v.generation || '').trim();
        if (!gg) return '';
        if (gg.toLowerCase().includes('http://') || gg.toLowerCase().includes('https://')) return '';
        return gg;
    };

    const formatPowertrain = (t: Trim) => {
        const pt = t.powertrain_meta;
        if (pt) {
            const pieces: string[] = [];
            if (pt.displacement_cc) {
                pieces.push(`${(pt.displacement_cc / 1000).toFixed(1)}L`);
            }
            if (pt.fuel_type) pieces.push(pt.fuel_type);
            if (pt.power_hp) pieces.push(`${pt.power_hp} hp`);
            if (pt.transmission_type) pieces.push(pt.transmission_type);
            if (pt.drive) pieces.push(pt.drive);
            const label = pieces.filter(Boolean).join(' • ');
            return label || t.name;
        }
        return t.name;
    };

    const grouped = useMemo(() => {
        const trims = data?.trims || [];
        const map = new Map<string, Trim[]>();
        for (const t of trims) {
            const key = formatPowertrain(t);
            map.set(key, [...(map.get(key) || []), t]);
        }
        return Array.from(map.entries()).sort((a, b) => a[0].localeCompare(b[0]));
    }, [data]);

    if (loading) {
        return <div className="rounded-2xl bg-white border border-border p-8 text-center text-text-muted shadow-sm">Loading powertrains...</div>;
    }

    if (!data) {
        return <div className="rounded-2xl bg-white border border-border p-8 text-center text-text-muted shadow-sm">No data found.</div>;
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
                <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-text-main">
                    {data.vehicle.brand} {data.vehicle.model}
                    {formatGeneration(data.vehicle) ? ` • ${formatGeneration(data.vehicle)}` : ''}
                </h1>
                <p className="text-sm text-text-muted">Motor / şanzıman (donanım) seç.</p>
            </div>

            <div className="rounded-2xl bg-white border border-border shadow-sm overflow-hidden">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-0">
                    <div className="bg-background p-6 flex items-center justify-center">
                        {data.vehicle.image_url ? (
                            <img
                                src={data.vehicle.image_url}
                                alt={data.vehicle.model}
                                className="h-auto w-auto max-w-[560px] max-h-[360px] object-contain"
                                loading="lazy"
                            />
                        ) : (
                            <div className="text-text-muted">No Image Available</div>
                        )}
                    </div>
                    <div className="p-6 sm:p-8">
                        <div className="text-xs font-semibold text-text-muted uppercase tracking-wider">Seçenekler</div>
                        <div className="mt-3 space-y-3">
                            {grouped.map(([label, trims]) => (
                                <div key={label} className="rounded-2xl border border-border bg-background/40 p-4">
                                    <div className="flex items-start justify-between gap-3">
                                        <div className="min-w-0">
                                            <div className="text-sm font-bold text-text-main">{label}</div>
                                            <div className="mt-1 text-xs text-text-muted">{trims.length} donanım</div>
                                        </div>
                                    </div>
                                    <div className="mt-3 flex flex-wrap gap-2">
                                        {trims.slice(0, 8).map((t) => (
                                            <Link
                                                key={t.id}
                                                to={`/vehicle/${data.vehicle.id}?trimId=${encodeURIComponent(String(t.id))}`}
                                                className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1 text-xs font-semibold text-text-main hover:border-primary/30 hover:text-primary transition-colors"
                                            >
                                                {t.name}
                                            </Link>
                                        ))}
                                        {trims.length > 8 ? (
                                            <Link
                                                to={`/vehicle/${data.vehicle.id}`}
                                                className="inline-flex items-center rounded-full bg-white border border-border px-3 py-1 text-xs font-semibold text-text-muted hover:text-text-main transition-colors"
                                            >
                                                Tümünü gör
                                            </Link>
                                        ) : null}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default PowertrainSelect;
