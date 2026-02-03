import { useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

interface Trim {
    id: number;
    name: string;
    year: number;
    generation: string;
    fuel_type?: string;
    displacement_cc?: number;
    power_hp?: number;
    image_url?: string;
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

export default function EngineSelect() {
    const navigate = useNavigate();
    const { brand, modelId } = useParams();
    const [trims, setTrims] = useState<Trim[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!modelId) return;

        setLoading(true);
        fetch(`${API_BASE_URL}/api/models/${modelId}/trims`)
            .then((res) => {
                if (!res.ok) throw new Error('Failed to fetch trims');
                return res.json();
            })
            .then((data) => {
                setTrims(data);
                setLoading(false);
            })
            .catch((err) => {
                console.error(err);
                setError('Failed to load vehicle data.');
                setLoading(false);
            });
    }, [modelId]);

    const engineGroups = useMemo(() => {
        const groups = new Map<string, { name: string; count: number; minYear: number; maxYear: number }>();

        for (const trim of trims) {
            let engineLabel = 'Other';

            // 1. Try to extract common patterns from name (e.g., "1.0 TFSI", "1.6 TDI", "35 TFSI")
            // Regex for displacement (1.0, 2.0) followed by type usually
            const nameMatch = trim.name.match(/(\d\.\d\s+(?:TFSI|TDI|FSI|TSI|MPI|e-tron|g-tron)|(?:30|35|40|45|50|55)\s+(?:TFSI|TDI|TFSIe)|RS3|S3)/i);

            if (nameMatch) {
                engineLabel = nameMatch[0].toUpperCase();
            } else if (trim.displacement_cc && trim.fuel_type) {
                // Fallback to cc + fuel
                const vol = (trim.displacement_cc / 1000).toFixed(1);
                engineLabel = `${vol}L ${trim.fuel_type}`;
            } else {
                engineLabel = trim.name.split(' ').slice(1, 3).join(' '); // Rough fallback
            }

            const existing = groups.get(engineLabel);
            if (existing) {
                groups.set(engineLabel, {
                    name: engineLabel,
                    count: existing.count + 1,
                    minYear: Math.min(existing.minYear, trim.year),
                    maxYear: Math.max(existing.maxYear, trim.year),
                });
            } else {
                groups.set(engineLabel, {
                    name: engineLabel,
                    count: 1,
                    minYear: trim.year,
                    maxYear: trim.year,
                });
            }
        }

        return Array.from(groups.values()).sort((a, b) => a.name.localeCompare(b.name, undefined, { numeric: true }));
    }, [trims]);

    const goBack = () => {
        navigate('/');
    };

    if (loading) return <div className="p-8 text-center text-gray-500">Loading engines...</div>;
    if (error) return <div className="p-8 text-center text-red-500">{error}</div>;

    return (
        <div className="min-h-screen">
            <div className="max-w-7xl mx-auto px-4 py-8 pt-32">
                <div className="flex items-center gap-4 mb-8">
                    <button
                        onClick={goBack}
                        className="inline-flex items-center justify-center w-10 h-10 rounded-full bg-white/10 text-white hover:bg-white/20 transition-all border border-white/5"
                    >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                        </svg>
                    </button>
                    <div>
                        <h1 className="text-3xl font-bold text-white capitalize">
                            {brand} Motor Seçimi
                        </h1>
                        <p className="text-gray-400 mt-1">İncelemek istediğiniz motor tipini seçin</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {engineGroups.map((group) => (
                        <Link
                            key={group.name}
                            to={`/selection/brand/${brand}/model/${modelId}/engine/${encodeURIComponent(group.name)}`}
                            className="block group"
                        >
                            <div className="h-full bg-black/40 backdrop-blur-md border border-white/10 rounded-xl p-6 shadow-lg hover:shadow-xl hover:border-primary/50 transition-all">
                                <div className="flex justify-between items-start mb-4">
                                    <h3 className="text-xl font-bold text-white group-hover:text-primary transition-colors">
                                        {group.name}
                                    </h3>
                                    <span className="inline-flex items-center justify-center rounded-full bg-primary/20 border border-primary/30 px-3 py-1 text-xs font-semibold text-primary-light">
                                        {group.count} varyasyon
                                    </span>
                                </div>

                                <div className="space-y-2">
                                    <div className="flex items-center text-sm text-gray-400">
                                        <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                        </svg>
                                        <span>Yıllar: {group.minYear === group.maxYear ? group.minYear : `${group.minYear} - ${group.maxYear}`}</span>
                                    </div>
                                </div>

                                <div className="mt-6 flex items-center text-sm font-medium text-primary group-hover:translate-x-1 transition-transform">
                                    Jenerasyon Seç
                                    <svg className="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                    </svg>
                                </div>
                            </div>
                        </Link>
                    ))}
                </div>
            </div>
        </div>
    );
}
