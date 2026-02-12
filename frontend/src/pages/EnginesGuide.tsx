import { useEffect, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import EngineDetailModal from '../components/EngineDetailModal';
import { engineService } from '../services/engineService';
import type { Engine } from '../types/engine';

export default function EnginesGuide() {
    const [engines, setEngines] = useState<Engine[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedEngine, setSelectedEngine] = useState<Engine | null>(null);

    // Background carousel
    const [currentBgIndex, setCurrentBgIndex] = useState(0);
    const backgroundImages = [
        '/hero/hero-1.jpg',
        '/hero/hero-2.jpg',
        '/hero/hero-3.jpg',
        '/hero/hero-4.jpg',
        '/hero/hero-5.jpg',
    ];

    // Rotate background every 5 seconds
    useEffect(() => {
        const interval = setInterval(() => {
            setCurrentBgIndex((prev) => (prev + 1) % backgroundImages.length);
        }, 5000);
        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        const fetchEngines = async () => {
            try {
                setLoading(true);
                setError(null);
                const data = await engineService.getAll();
                setEngines(data || []);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchEngines();
    }, []);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-900">
                <div className="text-center">
                    <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                    <p className="mt-4 text-gray-300">Motorlar yükleniyor...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-900">
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
        <div className="relative min-h-screen">
            {/* Full Screen Animated Background */}
            <div className="fixed inset-0 z-0">
                {backgroundImages.map((img, idx) => (
                    <motion.div
                        key={img}
                        initial={{ opacity: 0 }}
                        animate={{ opacity: idx === currentBgIndex ? 1 : 0 }}
                        transition={{ duration: 1.5, ease: "easeInOut" }}
                        className="absolute inset-0"
                        style={{
                            backgroundImage: `url(${img})`,
                            backgroundSize: 'cover',
                            backgroundPosition: 'center',
                        }}
                    />
                ))}
                {/* Dark overlay for readability */}
                <div className="absolute inset-0 bg-gradient-to-b from-black/40 via-black/50 to-black/60" />
            </div>

            {/* Content */}
            <div className="relative z-10">
                <div className="max-w-7xl mx-auto px-4 py-8 pt-32">
                    {/* Header */}
                    <div className="mb-12 text-center">
                        <h1 className="text-5xl font-bold text-white mb-4">
                            Motor Kodu Seçin
                        </h1>
                        <p className="text-gray-300 text-lg">
                            Motor seçerek detayları görüntüleyin
                        </p>
                    </div>

                    {engines.length === 0 ? (
                        <div className="text-center py-12 bg-black/40 backdrop-blur-md border border-white/10 rounded-xl">
                            <p className="text-gray-300 text-lg">Henüz motor bulunmamaktadır.</p>
                        </div>
                    ) : (
                        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 max-w-4xl mx-auto">
                            {engines.map((engine) => (
                                <button
                                    key={engine.id}
                                    onClick={() => setSelectedEngine(engine)}
                                    className="group bg-black/40 backdrop-blur-md border border-white/10 rounded-xl shadow-lg hover:shadow-2xl hover:border-primary/50 transition-all overflow-hidden"
                                >
                                    {/* Engine Code Container */}
                                    <div className="relative h-24 bg-white/5 flex items-center justify-center p-4">
                                        <div className="text-center">
                                            {/* Main Code */}
                                            <div className="text-3xl font-black text-white/90 group-hover:text-primary transition-colors tracking-tight">
                                                {engine.code}
                                            </div>
                                            {/* Manufacturer */}
                                            {engine.manufacturer && (
                                                <div className="text-[9px] text-gray-500 font-semibold uppercase mt-0.5 tracking-wider">
                                                    {engine.manufacturer}
                                                </div>
                                            )}
                                        </div>

                                        {/* Status Badge */}
                                        {!engine.productionEndYear && (
                                            <div className="absolute top-2 right-2 w-2 h-2 bg-emerald-400 rounded-full animate-pulse" />
                                        )}
                                    </div>

                                    {/* Engine Name */}
                                    <div className="p-4 border-t border-white/10">
                                        <h3 className="text-sm font-bold text-white text-center group-hover:text-primary transition-colors line-clamp-2">
                                            {engine.name || engine.code}
                                        </h3>
                                        {/* Fuel Type */}
                                        {engine.fuelType && (
                                            <p className="text-xs text-gray-400 text-center mt-1">
                                                {engine.fuelType}
                                            </p>
                                        )}
                                    </div>

                                    {/* Hover Effect */}
                                    <div className="absolute inset-0 bg-gradient-to-t from-primary/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
                                </button>
                            ))}
                        </div>
                    )}

                    {/* Info Section */}
                    <div className="mt-16 text-center">
                        <div className="inline-block bg-black/40 backdrop-blur-md border border-white/10 rounded-xl p-8 max-w-2xl">
                            <h2 className="text-2xl font-bold text-white mb-4">
                                Motor Kodları Rehberi
                            </h2>
                            <div className="space-y-3 text-gray-300 text-left">
                                <div className="flex items-start gap-3">
                                    <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">1</span>
                                    <p><strong className="text-white">Motor Seçin:</strong> İlgilendiğiniz motor kodunu seçin</p>
                                </div>
                                <div className="flex items-start gap-3">
                                    <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">2</span>
                                    <p><strong className="text-white">Detayları İnceleyin:</strong> Motor özelliklerini görüntüleyin</p>
                                </div>
                                <div className="flex items-start gap-3">
                                    <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">3</span>
                                    <p><strong className="text-white">Kronik Sorunlar:</strong> Yaygın arızaları ve çözümleri öğrenin</p>
                                </div>
                                <div className="flex items-start gap-3">
                                    <span className="flex-shrink-0 w-8 h-8 bg-primary/20 text-primary rounded-full flex items-center justify-center font-bold">4</span>
                                    <p><strong className="text-white">Güvenilirlik:</strong> Motor güvenilirlik puanlarını kontrol edin</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Modal */}
            <AnimatePresence>
                {selectedEngine && (
                    <EngineDetailModal
                        engine={selectedEngine}
                        onClose={() => setSelectedEngine(null)}
                    />
                )}
            </AnimatePresence>
        </div>
    );
}
