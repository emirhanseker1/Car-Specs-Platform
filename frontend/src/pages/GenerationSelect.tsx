import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

interface Generation {
    id: number;
    model_id: number;
    code: string;
    name?: string;
    start_year: number;
    end_year?: number;
    is_current: boolean;
    trim_count: number;
    image_url?: string;
}

import { getGenerationStyle } from '../utils/generationStyles';

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

export default function GenerationSelect() {
    const navigate = useNavigate();
    const { modelId: routeModelId, brandName, modelName } = useParams();
    const [modelId, setModelId] = useState<string | null>(routeModelId || null);
    const [generations, setGenerations] = useState<Generation[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Fetch model ID if we have brandName and modelName instead
    useEffect(() => {
        if (routeModelId) {
            setModelId(routeModelId);
            return;
        }

        if (!brandName || !modelName) return;

        const fetchModelId = async () => {
            try {
                // Step 1: Get brand ID from brand name
                const brandsResponse = await fetch(`${API_BASE_URL}/api/brands`);
                if (!brandsResponse.ok) {
                    throw new Error('Markalar yüklenemedi');
                }
                const allBrands = await brandsResponse.json();
                const brand = allBrands.find((b: any) =>
                    b.name.toLowerCase() === brandName.toLowerCase()
                );

                if (!brand || !brand.id) {
                    throw new Error('Marka bulunamadı');
                }

                // Step 2: Fetch models using numeric brand ID
                const modelsResponse = await fetch(`${API_BASE_URL}/api/brands/${brand.id}/models`);

                if (!modelsResponse.ok) {
                    console.error('API error:', modelsResponse.status, modelsResponse.statusText);
                    throw new Error('Modeller yüklenemedi');
                }

                const data = await modelsResponse.json();
                const models = data?.value || [];

                // Find matching model
                const model = models.find((m: any) => m.name === modelName);

                if (model && model.id) {
                    setModelId(model.id.toString());
                } else {
                    setError('Model bulunamadı');
                    setLoading(false);
                }
            } catch (err: any) {
                console.error('Error fetching model:', err);
                setError(err.message);
                setLoading(false);
            }
        };

        fetchModelId();
    }, [routeModelId, brandName, modelName]);

    useEffect(() => {
        if (!modelId) return;

        const fetchGenerations = async () => {
            try {
                setLoading(true);
                setError(null);
                const response = await fetch(`${API_BASE_URL}/api/models/${modelId}/generations`);
                if (!response.ok) throw new Error('Nesiller yüklenemedi');
                const data = await response.json();

                // Handle both wrapped and unwrapped responses
                const generationsList = !data ? [] : (Array.isArray(data) ? data : (data.value || []));
                setGenerations(generationsList);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchGenerations();
    }, [modelId]);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                    <p className="mt-4 text-gray-300">Nesiller yükleniyor...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="bg-red-500/10 backdrop-blur border border-red-500/20 rounded-lg p-6 max-w-md">
                    <p className="text-red-200 mb-4">Hata: {error}</p>
                    <button
                        onClick={() => navigate(-1)}
                        className="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg transition-colors"
                    >
                        Geri Dön
                    </button>
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
                        onClick={() => navigate(-1)}
                        className="mb-4 inline-flex items-center gap-2 px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg transition-all backdrop-blur-sm border border-white/10"
                    >
                        ← Geri
                    </button>
                    <h1 className="text-5xl font-bold text-white mb-4">
                        Nesil Seçin
                    </h1>
                    <p className="text-gray-300 text-lg">
                        Hangi nesli incelemek istersiniz?
                    </p>
                </div>

                {generations.length === 0 ? (
                    <div className="text-center py-12 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl">
                        <p className="text-gray-300 text-lg">Nesil bulunamadı.</p>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {generations.map((gen) => (
                            <Link
                                key={gen.id}
                                to={brandName && modelName
                                    ? `/brand/${brandName}/${modelName}/${gen.code}`
                                    : `/generations/${gen.id}/trims`
                                }
                                className="group bg-black/40 backdrop-blur-md border border-white/10 rounded-xl shadow-lg hover:shadow-2xl hover:border-primary/50 transition-[transform,border,shadow] duration-300 overflow-hidden"
                            >
                                {/* Generation Image */}
                                {gen.image_url && (
                                    <div className="relative h-48 overflow-hidden flex items-center justify-center p-4">
                                        <img
                                            src={gen.image_url}
                                            alt={`${gen.code} ${gen.name || ''}`}
                                            className={`max-w-full max-h-full object-contain object-center transition-transform duration-300 ${getGenerationStyle(gen.code, brandName)
                                                }`}
                                            onError={(e) => {
                                                // Hide image if it fails to load
                                                e.currentTarget.style.display = 'none';
                                            }}
                                        />
                                    </div>
                                )}

                                <div className="p-6">
                                    {/* Gen Code & Badge */}
                                    <div className="flex justify-between items-start mb-4">
                                        <div>
                                            <h3 className="text-2xl font-bold text-white group-hover:text-primary transition-colors">
                                                {gen.code}
                                            </h3>
                                            {gen.name && (
                                                <p className="text-sm text-gray-400 mt-1">
                                                    {gen.name}
                                                </p>
                                            )}
                                        </div>
                                        {gen.is_current && (
                                            <span className="inline-flex items-center justify-center rounded-full bg-green-500/20 border border-green-500/30 px-3 py-1 text-xs font-semibold text-green-400">
                                                Güncel
                                            </span>
                                        )}
                                    </div>

                                    {/* Years */}
                                    <div className="mb-4">
                                        <div className="text-gray-300 text-lg">
                                            {gen.end_year
                                                ? `${gen.start_year} - ${gen.end_year}`
                                                : `${gen.start_year} - Günümüz`}
                                        </div>
                                    </div>

                                    {/* Trim Count */}
                                    {gen.trim_count !== undefined && gen.trim_count > 0 && (
                                        <div className="flex items-center gap-2 mb-4 text-sm text-gray-400">
                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                            </svg>
                                            <span>{gen.trim_count} motor seçeneği</span>
                                        </div>
                                    )}

                                    {/* Action */}
                                    <div className="flex items-center gap-2 text-primary font-medium pt-2 border-t border-white/10">
                                        <span>Motor seçeneklerini gör</span>
                                        <svg className="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                        </svg>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
