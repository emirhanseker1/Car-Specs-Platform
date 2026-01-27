import { useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

interface Vehicle {
    id: number;
    brand: string;
    model: string;
    generation: string;
    image_url: string;
    engine_options?: string[];
}

interface ModelCard {
    brand: string;
    model: string;
    image_url?: string;
    generationCount: number;
    engine_options: string[];
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function ModelList() {
    const navigate = useNavigate();
    const { brandName } = useParams();
    const [vehicles, setVehicles] = useState<Vehicle[]>([]);
    const [loading, setLoading] = useState(true);

    const decodedBrandName = brandName ? decodeURIComponent(brandName) : '';

    const goBack = () => {
        if (window.history.length > 1) {
            navigate(-1);
            return;
        }
        navigate('/', { replace: true });
    };

    useEffect(() => {
        const brandQuery = decodedBrandName || 'Fiat';
        fetch(`${API_BASE_URL}/api/vehicles?brand=${encodeURIComponent(brandQuery)}`)
            .then((res) => res.json())
            .then((data) => {
                setVehicles(data || []);
                setLoading(false);
            })
            .catch((err) => {
                console.error('Failed to fetch vehicles:', err);
                setLoading(false);
            });
    }, [decodedBrandName]);

    const models = useMemo<ModelCard[]>(() => {
        const byModel = new Map<string, ModelCard>();
        for (const v of vehicles) {
            const model = (v.model || '').trim();
            if (!model) continue;
            const key = model.toLowerCase();
            const current = byModel.get(key) || {
                brand: decodedBrandName,
                model,
                image_url: undefined as string | undefined,
                generationCount: 0,
                engine_options: [] as string[],
            };

            current.generationCount += 1;
            if (!current.image_url && v.image_url) current.image_url = v.image_url;

            if (Array.isArray(v.engine_options)) {
                for (const opt of v.engine_options) {
                    const o = (opt || '').trim();
                    if (!o) continue;
                    if (!current.engine_options.includes(o)) current.engine_options.push(o);
                }
            }

            byModel.set(key, current);
        }

        return Array.from(byModel.values()).sort((a, b) => a.model.localeCompare(b.model));
    }, [decodedBrandName, vehicles]);

    if (loading) return <div className="rounded-2xl bg-white border border-border p-8 text-center text-text-muted shadow-sm">Loading vehicles...</div>;

    return (
        <div className="space-y-6">
            <div className="flex items-start justify-between gap-4">
                <div className="space-y-1">
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
                        {decodedBrandName || 'Fiat'}
                    </h1>
                    <p className="text-sm text-text-muted">Model seç, nesil/yıl seç ve motor/şanzıman donanımına göre incele.</p>
                </div>
                <div className="text-xs text-text-muted pt-1">
                    {models.length} model
                </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                {models.map((m) => (
                    <Link
                        key={m.model}
                        to={`/brand/${encodeURIComponent(decodedBrandName || m.brand)}/model/${encodeURIComponent(m.model)}`}
                        className="rounded-2xl bg-white border border-border shadow-sm hover:shadow-card transition-shadow overflow-hidden group"
                    >
                        <div className="relative h-44 bg-background flex items-center justify-center p-4">
                            {m.image_url ? (
                                <img
                                    src={m.image_url}
                                    alt={m.model}
                                    className="w-full h-full object-contain group-hover:scale-105 transition-transform duration-300"
                                />
                            ) : (
                                <div className="text-text-muted text-sm">No Image</div>
                            )}
                        </div>
                        <div className="p-5 border-t border-border/50">
                            <h3 className="text-base font-bold text-text-main leading-snug">
                                {m.model}
                            </h3>
                            <div className="mt-2 text-sm text-text-muted">
                                {m.engine_options.length > 0
                                    ? `Motor seçenekleri: ${m.engine_options.slice(0, 3).join(' • ')}${m.engine_options.length > 3 ? '…' : ''}`
                                    : 'Motor seçenekleri: —'}
                                <div className="mt-1 text-xs text-text-muted">{m.generationCount} nesil/yıl</div>
                            </div>
                        </div>
                    </Link>
                ))}
            </div>

            {models.length === 0 && (
                <div className="text-center py-12 bg-white border border-border rounded-2xl text-text-muted shadow-sm">
                    No models found for {decodedBrandName || 'this brand'}.
                </div>
            )}
        </div>
    );
}

export default ModelList;
