import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { ArrowRight, Package } from 'lucide-react';

interface FeaturedTrim {
    id: number;
    name: string;
    year: number;
    image_url?: string;
    fuel_type?: string;
    transmission_type?: string;
    power_hp?: number;
    model?: {
        id: number;
        name: string;
        brand?: {
            id: number;
            name: string;
        };
    };
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function Home() {
    const [featuredTrims, setFeaturedTrims] = useState<FeaturedTrim[]>([]);
    const [featuredLoading, setFeaturedLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Fetch featured trims
    useEffect(() => {
        setFeaturedLoading(true);
        fetch(`${API_BASE_URL}/api/featured`)
            .then(async (res) => {
                if (!res.ok) throw new Error('Failed to fetch featured');
                return res.json();
            })
            .then((data: FeaturedTrim[]) => {
                setFeaturedTrims(data || []);
                setFeaturedLoading(false);
            })
            .catch((err) => {
                console.error('Failed to fetch featured:', err);
                setFeaturedTrims([]);
                setFeaturedLoading(false);
                setError('Popüler modeller yüklenemedi');
            });
    }, []);

    const retryAllFetches = () => {
        setError(null);
        window.location.reload();
    };

    return (
        <div className="relative min-h-screen text-white">

            {/* Content Overlay */}

            {/* Content Overlay */}
            <div className="relative z-10 flex flex-col min-h-screen">

                {/* Hero Section - Full Screen Centered */}
                <div className="flex-grow flex items-center justify-center min-h-screen px-4 py-20">
                    <div className="max-w-4xl w-full text-center space-y-8">
                        <h1 className="text-6xl md:text-8xl font-black tracking-tighter drop-shadow-2xl">
                            ARAÇLARIN
                            <span className="block text-primary">TÜM DETAYLARI</span>
                        </h1>
                        <p className="text-xl md:text-2xl text-gray-200 font-light max-w-2xl mx-auto drop-shadow-lg">
                            Türkiye'nin en kapsamlı araç teknik bilgi platformu ile ilgilendiğiniz aracın tüm özelliklerini keşfedin.
                        </p>

                        <div className="flex flex-wrap justify-center gap-4 pt-8">
                            {['Şanzıman', 'Motor', 'Performans', 'Boyutlar'].map((tag) => (
                                <span key={tag} className="px-6 py-2.5 rounded-full bg-white/5 backdrop-blur-sm border border-white/10 text-sm font-semibold hover:bg-white/10 transition-all cursor-default will-change-transform">
                                    {tag}
                                </span>
                            ))}
                        </div>
                    </div>
                </div>

                {/* DSG Rehberi Floating Button - Left Side */}
                <Link
                    to="/guides/transmission/dsg"
                    className="fixed left-6 top-1/3 z-40 group"
                >
                    <div className="flex items-center gap-3 pl-2 pr-4 py-2 bg-black/20 backdrop-blur-md rounded-full border border-white/10 hover:bg-white/5 transition-all duration-300 group-hover:scale-105">
                        <div className="w-8 h-8 bg-white/10 rounded-full flex items-center justify-center group-hover:bg-white/20 transition-colors">
                            <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            </svg>
                        </div>
                        <div className="hidden md:block">
                            <span className="font-medium text-white text-sm">DSG Rehberi</span>
                        </div>
                        <svg className="w-4 h-4 text-white/50 group-hover:text-white group-hover:translate-x-0.5 transition-all hidden md:block" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    </div>
                </Link>

                {/* Popular Models Section - At the bottom */}
                <div className="w-full bg-gradient-to-t from-black via-black/80 to-transparent pt-32 pb-12 px-6 sm:px-8 lg:px-12">
                    {/* Error Section */}
                    {error && (
                        <div className="mb-8 rounded-2xl bg-red-500/10 backdrop-blur border border-red-500/20 p-5">
                            <div className="flex flex-wrap items-center justify-between gap-3">
                                <div className="text-red-200 text-sm">
                                    <span className="font-bold block mb-1">Hata Oluştu</span>
                                    {error}
                                </div>
                                <button
                                    onClick={retryAllFetches}
                                    className="px-4 py-2 bg-red-500/20 hover:bg-red-500/30 rounded-full text-xs font-bold text-red-200 transition-colors"
                                >
                                    Tekrar Dene
                                </button>
                            </div>
                        </div>
                    )}

                    <div className="flex items-center justify-between mb-8">
                        <h2 className="text-3xl font-bold tracking-tight">Popüler Modeller</h2>
                        <Link
                            to="/search"
                            className="text-sm font-bold text-primary hover:text-white transition-colors flex items-center gap-2"
                        >
                            Tümünü Gör <ArrowRight className="w-4 h-4" />
                        </Link>
                    </div>

                    {featuredLoading ? (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                            {[1, 2, 3, 4].map(i => (
                                <div key={i} className="aspect-video rounded-3xl bg-white/5 animate-pulse" />
                            ))}
                        </div>
                    ) : featuredTrims.length === 0 ? (
                        <div className="text-center py-12 bg-white/5 rounded-3xl backdrop-blur-sm border border-white/10">
                            <p className="text-gray-400">Henüz öne çıkan model bulunmuyor.</p>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                            {featuredTrims.map((trim) => (
                                <Link
                                    key={trim.id}
                                    to={`/vehicle/${trim.id}`}
                                    className="group relative overflow-hidden rounded-3xl bg-white/5 border border-white/10 hover:border-primary/50 hover:bg-white/10 transition-all duration-300"
                                >
                                    {/* Image */}
                                    <div className="aspect-video overflow-hidden bg-white/5">
                                        {trim.image_url ? (
                                            <img
                                                src={trim.image_url}
                                                alt={`${trim.model?.brand?.name} ${trim.name}`}
                                                className="w-full h-full object-contain p-2 group-hover:scale-105 transition-transform duration-700"
                                                onError={(e) => {
                                                    (e.target as HTMLImageElement).style.display = 'none';
                                                }}
                                            />
                                        ) : (
                                            <div className="w-full h-full flex items-center justify-center bg-white/5">
                                                <Package className="w-12 h-12 text-white/20" />
                                            </div>
                                        )}
                                        {/* Overlay Gradient on Image */}
                                        <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-80 group-hover:opacity-100 transition-opacity" />
                                    </div>

                                    {/* Content */}
                                    <div className="absolute bottom-0 left-0 right-0 p-5 translate-y-2 group-hover:translate-y-0 transition-transform duration-300">
                                        <div className="flex items-center justify-between mb-1">
                                            <span className="text-xs font-bold px-2 py-1 rounded bg-primary/20 text-primary uppercase tracking-wider">
                                                {trim.model?.brand?.name}
                                            </span>
                                            <span className="text-xs font-medium text-gray-400">{trim.year}</span>
                                        </div>
                                        <h3 className="text-lg font-bold text-white leading-tight mb-2 truncate">
                                            {trim.model?.name} {trim.name}
                                        </h3>
                                        <div className="flex items-center gap-3 text-xs text-gray-300 opacity-0 group-hover:opacity-100 transition-opacity duration-300 delay-75">
                                            <span className="font-semibold">{trim.power_hp} HP</span>
                                            {trim.fuel_type && (
                                                <>
                                                    <span className="w-1 h-1 rounded-full bg-gray-500" />
                                                    <span className="font-medium text-gray-200">{trim.fuel_type}</span>
                                                </>
                                            )}
                                        </div>
                                    </div>
                                </Link>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

export default Home;
