import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

interface Brand {
    id: number;
    name: string;
    logo?: string;
    created_at?: string;
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function Search() {
    const [brands, setBrands] = useState<Brand[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchBrands = async () => {
            try {
                setLoading(true);
                setError(null);
                const response = await fetch(`${API_BASE_URL}/api/brands`);
                if (!response.ok) throw new Error('Markalar yüklenemedi');
                const data = await response.json();
                setBrands(data || []);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchBrands();
    }, []);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                    <p className="mt-4 text-gray-300">Markalar yükleniyor...</p>
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
                        onClick={() => window.location.reload()}
                        className="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg transition-colors"
                    >
                        Yeniden Dene
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen">
            <div className="max-w-7xl mx-auto px-4 py-8 pt-32">
                {/* Header */}
                <div className="mb-12 text-center">
                    <h1 className="text-5xl font-bold text-white mb-4">
                        Araç Markası Seçin
                    </h1>
                    <p className="text-gray-300 text-lg">
                        Marka seçerek modelleri görüntüleyin
                    </p>
                </div>

                {brands.length === 0 ? (
                    <div className="text-center py-12 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl">
                        <p className="text-gray-300 text-lg">Henüz marka bulunmamaktadır.</p>
                    </div>
                ) : (
                    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 max-w-4xl mx-auto">
                        {brands.map((brand) => (
                            <Link
                                key={brand.id}
                                to={`/brand/${encodeURIComponent(brand.name)}`}
                                className="group bg-black/40 backdrop-blur-md border border-white/10 rounded-xl shadow-lg hover:shadow-2xl hover:border-primary/50 transition-all overflow-hidden"
                            >
                                {/* Brand Logo Container */}
                                <div className="relative h-24 bg-white/5 flex items-center justify-center p-4">
                                    {brand.logo ? (
                                        <img
                                            src={brand.logo}
                                            alt={brand.name}
                                            className="w-full h-full object-contain scale-150 group-hover:scale-[1.65] transition-transform duration-300"
                                        />
                                    ) : brand.name === 'Audi' ? (
                                        <img
                                            src="/images/audi-logo.png"
                                            alt="Audi Logo"
                                            className="w-full h-full object-contain scale-150 group-hover:scale-[1.65] transition-transform duration-300"
                                        />
                                    ) : (
                                        <div className="text-center">
                                            <div className="text-6xl font-bold text-white/20 group-hover:text-primary/40 transition-colors">
                                                {brand.name.charAt(0)}
                                            </div>
                                        </div>
                                    )}
                                </div>

                                {/* Brand Name */}
                                <div className="p-4 border-t border-white/10">
                                    <h3 className="text-xl font-bold text-white text-center group-hover:text-primary transition-colors">
                                        {brand.name}
                                    </h3>
                                </div>

                                {/* Hover Effect */}
                                <div className="absolute inset-0 bg-gradient-to-t from-primary/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
                            </Link>
                        ))}
                    </div>
                )}

                {/* Info Section */}
                <div className="mt-16 text-center">
                    <div className="inline-block bg-black/40 backdrop-blur-md border border-white/10 rounded-xl p-8 max-w-2xl">
                        <h2 className="text-2xl font-bold text-white mb-4">
                            Nasıl Çalışır?
                        </h2>
                        <div className="space-y-3 text-gray-300 text-left">
                            <div className="flex items-start gap-3">
                                <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">1</span>
                                <p><strong className="text-white">Marka Seçin:</strong> İstediğiniz araç markasını seçin</p>
                            </div>
                            <div className="flex items-start gap-3">
                                <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">2</span>
                                <p><strong className="text-white">Model Seçin:</strong> Markanın modellerini görüntüleyin</p>
                            </div>
                            <div className="flex items-start gap-3">
                                <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">3</span>
                                <p><strong className="text-white">Nesil Seçin:</strong> Modelin farklı nesillerini inceleyin</p>
                            </div>
                            <div className="flex items-start gap-3">
                                <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">4</span>
                                <p><strong className="text-white">Motor Seçin:</strong> Nesle ait motor seçeneklerini görün</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Search;
