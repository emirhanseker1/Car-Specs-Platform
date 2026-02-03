import { useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

interface Vehicle {
    id: number;
    brand: string;
    model: string;
    generation: string;
    image_url: string;
    engine_options?: string[];
    model_id?: number;
}

interface ModelCard {
    brand: string;
    model: string;
    model_id?: number;
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
        const brandQuery = decodedBrandName || 'Audi';
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
                model_id: v.model_id,
                image_url: undefined as string | undefined,
                generationCount: 0,
                engine_options: [] as string[],
            };

            if (!current.model_id && v.model_id) current.model_id = v.model_id;

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

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                    <p className="mt-4 text-gray-300">Modeller yükleniyor...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen">
            <div className="max-w-7xl mx-auto px-4 py-8 pt-32">
                {/* Header */}
                <div className="mb-12">
                    <button
                        type="button"
                        onClick={goBack}
                        className="mb-4 inline-flex items-center gap-2 px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg transition-all backdrop-blur-sm border border-white/10"
                    >
                        ← Geri
                    </button>
                    <h1 className="text-5xl font-bold text-white mb-4">
                        {decodedBrandName || 'Modeller'}
                    </h1>
                    <p className="text-gray-300 text-lg">
                        Model seçerek nesilleri görüntüleyin
                    </p>
                </div>

                {models.length === 0 ? (
                    <div className="text-center py-12 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl">
                        <p className="text-gray-300 text-lg">
                            {decodedBrandName || 'Bu marka'} için model bulunamadı.
                        </p>
                    </div>
                ) : (
                    <div className="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 max-w-6xl mx-auto">
                        {models.map((m) => (
                            <Link
                                key={m.model}
                                to={`/brand/${brandName}/${m.model}`}
                                className="group bg-black/40 backdrop-blur-md border border-white/10 rounded-xl shadow-lg hover:shadow-2xl hover:border-primary/50 transition-all overflow-hidden"
                            >
                                {/* Model Image */}
                                <div className="relative h-40 bg-white/5 flex items-center justify-center p-6">
                                    {m.image_url ? (
                                        <img
                                            src={m.image_url}
                                            alt={m.model}
                                            className="w-full h-full object-contain group-hover:scale-110 transition-transform duration-300"
                                        />
                                    ) : m.model === 'A3' ? (
                                        <img
                                            src="/images/audi-a3.png"
                                            alt="Audi A3"
                                            className="w-full h-full object-contain scale-[1.6] group-hover:scale-[1.7] transition-transform duration-300"
                                        />
                                    ) : (
                                        <div className="text-center">
                                            <div className="text-4xl font-bold text-white/30 group-hover:text-primary/50 transition-colors">
                                                {m.model.charAt(0)}
                                            </div>
                                        </div>
                                    )}
                                </div>

                                {/* Model Name & Info */}
                                <div className="p-5 border-t border-white/10">
                                    <h3 className="text-2xl font-bold text-white text-center group-hover:text-primary transition-colors mb-2">
                                        {m.model}
                                    </h3>
                                    <div className="text-sm text-gray-400 text-center">
                                        {m.generationCount} nesil
                                    </div>
                                </div>

                                {/* Hover Effect */}
                                <div className="absolute inset-0 bg-gradient-to-t from-primary/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
                            </Link>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}

export default ModelList;
