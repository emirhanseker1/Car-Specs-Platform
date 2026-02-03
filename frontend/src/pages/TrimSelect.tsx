import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';

interface Trim {
    id: number;
    name: string;
    power_hp?: number;
    torque_nm?: number;
    acceleration_0_100?: number;
    fuel_type?: string;
    transmission_type?: string;
    transmission_code?: string;
    drivetrain?: string;
    year: number;
    start_year?: number;
    end_year?: number;
}

interface Generation {
    id: number;
    code: string;
    name?: string;
    start_year: number;
    end_year?: number;
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function TrimSelect() {
    const { generationId, brandName, modelName, generationCode } = useParams<{
        generationId: string;
        brandName: string;
        modelName: string;
        generationCode: string;
    }>();
    const navigate = useNavigate();
    const [trims, setTrims] = useState<Trim[]>([]);
    const [generation, setGeneration] = useState<Generation | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);

                let targetGenId = generationId;

                // Resolve generation ID from pretty URL params if needed
                if (!targetGenId && brandName && modelName && generationCode) {
                    // 1. Fetch Brand
                    const brandsRes = await fetch(`${API_BASE_URL}/api/brands`);
                    if (!brandsRes.ok) throw new Error('Markalar yüklenemedi');
                    const brands = await brandsRes.json();
                    const brand = brands.find((b: any) => b.name.toLowerCase() === brandName.toLowerCase());
                    if (!brand) throw new Error('Marka bulunamadı');

                    // 2. Fetch Model
                    const modelsRes = await fetch(`${API_BASE_URL}/api/brands/${brand.id}/models`);
                    if (!modelsRes.ok) throw new Error('Modeller yüklenemedi');
                    const modelsData = await modelsRes.json();
                    const models = modelsData.value || [];
                    const model = models.find((m: any) => m.name.toLowerCase() === modelName.toLowerCase());
                    if (!model) throw new Error('Model bulunamadı');

                    // 3. Fetch Generations
                    const gensRes = await fetch(`${API_BASE_URL}/api/models/${model.id}/generations`);
                    if (!gensRes.ok) throw new Error('Nesiller yüklenemedi');
                    const gensData = await gensRes.json();
                    const generationsList = Array.isArray(gensData) ? gensData : (gensData.value || []);

                    const targetGen = generationsList.find((g: any) => g.code.toLowerCase() === generationCode.toLowerCase());
                    if (!targetGen) throw new Error('Nesil bulunamadı');

                    targetGenId = targetGen.id.toString();
                }

                if (!targetGenId) return; // Should not happen if routing is correct

                // Fetch generation info
                const genRes = await fetch(`${API_BASE_URL}/api/generations/${targetGenId}`);
                if (!genRes.ok) throw new Error('Nesil bilgisi bulunamadı');
                const genData = await genRes.json();
                setGeneration(genData);

                // Fetch trims for this generation
                const trimRes = await fetch(`${API_BASE_URL}/api/generations/${targetGenId}/trims`);
                if (!trimRes.ok) throw new Error('Motor seçenekleri yüklenemedi');
                const trimData = await trimRes.json();

                // Handle both wrapped and unwrapped responses
                const trimsList = Array.isArray(trimData) ? trimData : (trimData.value || []);
                setTrims(trimsList);
            } catch (err: any) {
                console.error(err);
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [generationId, brandName, modelName, generationCode]);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                    <p className="mt-4 text-gray-300">Motor seçenekleri yükleniyor...</p>
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

    const yearRange = generation?.end_year
        ? `${generation.start_year} - ${generation.end_year}`
        : `${generation?.start_year} - Günümüz`;

    return (
        <div className="min-h-screen">
            <div className="max-w-7xl mx-auto px-4 py-8 pt-32">
                {/* Header */}
                <div className="mb-8">
                    <button
                        onClick={() => navigate(-1)}
                        className="inline-flex items-center text-gray-400 hover:text-white mb-4 transition-colors"
                    >
                        <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                        </svg>
                        Nesillere Dön
                    </button>
                    <h1 className="text-4xl font-bold text-white mb-2">
                        Nesil {generation?.code || ''} - Motor Seçimi
                    </h1>
                    <p className="text-gray-300">{yearRange}</p>
                    <p className="text-gray-400 mt-2">
                        Motor seçeneğinizi seçin
                    </p>
                </div>

                {trims.length === 0 ? (
                    <div className="text-center py-12 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl">
                        <p className="text-gray-300 text-lg">Bu nesil için motor bilgisi bulunamadı.</p>
                        <button
                            onClick={() => navigate(-1)}
                            className="mt-4 px-6 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg transition-colors font-medium"
                        >
                            Geri Dön
                        </button>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {trims.map((trim, index) => (
                            <Link
                                key={trim.id}
                                to={`/brand/${brandName}/${modelName}/${generationCode}/${index + 1}`}
                                className="bg-black/40 backdrop-blur-md border border-white/10 rounded-xl shadow-lg hover:shadow-xl hover:border-primary/50 transition-all overflow-hidden group p-6"
                            >
                                {/* Trim Header */}
                                <div className="mb-4">
                                    <h3 className="text-2xl font-bold text-white mb-1 group-hover:text-primary transition-colors">
                                        {trim.name.replace(/Tfsi/g, 'TFSI').replace(/Tdi/g, 'TDI')}
                                    </h3>
                                    <p className="text-gray-400 text-sm">
                                        {trim.end_year
                                            ? `${trim.start_year || trim.year} - ${trim.end_year}`
                                            : `${trim.start_year || trim.year} - Devam ediyor`}
                                    </p>
                                </div>

                                {/* Specs Grid */}
                                <div className="space-y-3 mb-4">
                                    {trim.fuel_type && (
                                        <div className="flex justify-between items-center pb-2 border-b border-white/10">
                                            <span className="text-gray-400 text-sm">Yakıt</span>
                                            <span className="text-white font-medium">{trim.fuel_type}</span>
                                        </div>
                                    )}
                                    {trim.power_hp && (
                                        <div className="flex justify-between items-center pb-2 border-b border-white/10">
                                            <span className="text-gray-400 text-sm">Güç</span>
                                            <span className="text-white font-medium">{trim.power_hp} HP</span>
                                        </div>
                                    )}
                                    {trim.torque_nm && (
                                        <div className="flex justify-between items-center pb-2 border-b border-white/10">
                                            <span className="text-gray-400 text-sm">Tork</span>
                                            <span className="text-white font-medium">{trim.torque_nm} Nm</span>
                                        </div>
                                    )}
                                    {trim.acceleration_0_100 && (
                                        <div className="flex justify-between items-center pb-2 border-b border-white/10">
                                            <span className="text-gray-400 text-sm">0-100 km/h</span>
                                            <span className="text-white font-medium">{trim.acceleration_0_100}s</span>
                                        </div>
                                    )}
                                    {trim.transmission_type && (
                                        <div className="flex justify-between items-center pb-2 border-b border-white/10">
                                            <span className="text-gray-400 text-sm">Şanzıman</span>
                                            <span className="text-white font-medium">{trim.transmission_type}</span>
                                        </div>
                                    )}
                                    {trim.transmission_code && (
                                        <div className="flex justify-between items-center pb-2 border-b border-white/10">
                                            <span className="text-gray-400 text-sm">Şanzıman Tipi</span>
                                            <span className="text-white font-medium">{trim.transmission_code}</span>
                                        </div>
                                    )}
                                    {trim.drivetrain && (
                                        <div className="flex justify-between items-center pb-2 border-b border-white/10">
                                            <span className="text-gray-400 text-sm">Çekiş</span>
                                            <span className="text-white font-medium">{trim.drivetrain}</span>
                                        </div>
                                    )}
                                </div>

                                {/* Action */}
                                <div className="flex items-center justify-between pt-4 mt-4 border-t border-white/10">
                                    <span className="text-sm text-primary font-medium">Detayları Gör</span>
                                    <svg className="w-5 h-5 text-primary group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                    </svg>
                                </div>
                            </Link>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}

export default TrimSelect;
