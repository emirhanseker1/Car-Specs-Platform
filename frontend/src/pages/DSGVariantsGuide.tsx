import { useState, useEffect } from 'react';
import { ArrowLeft, AlertTriangle, ChevronRight, TrendingUp } from 'lucide-react';
import { Link } from 'react-router-dom';
import { DSG_VARIANTS } from '../data/dsgVariants';
import DSGVariantModal from '../components/DSGVariantModal';
import { AnimatePresence } from 'framer-motion';

export default function DSGVariantsGuide() {
    const [expandedCode, setExpandedCode] = useState<string | null>(null);

    const handleCardClick = (code: string) => {
        setExpandedCode(expandedCode === code ? null : code);
    };

    // Toggle body class for modal
    useEffect(() => {
        if (expandedCode) {
            document.body.classList.add('modal-open');
        } else {
            document.body.classList.remove('modal-open');
        }
        return () => document.body.classList.remove('modal-open');
    }, [expandedCode]);

    const expandedVariant = DSG_VARIANTS.find(v => v.code === expandedCode) || null;

    return (
        <div className="min-h-screen font-sans text-gray-100">
            {/* Header */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 pt-32">
                <div className="mb-12">
                    <Link
                        to="/guides/transmission"
                        className="inline-flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors mb-6 bg-black/40 hover:bg-black/60 px-4 py-2 rounded-full w-fit border border-white/10 backdrop-blur-md"
                    >
                        <ArrowLeft className="w-4 h-4" />
                        Şanzıman Rehberine Dön
                    </Link>

                    <h1 className="text-3xl sm:text-5xl font-black text-white mb-4 drop-shadow-2xl">
                        VAG DSG/S-Tronic Detayları
                    </h1>
                    <p className="text-lg text-gray-300 max-w-2xl leading-relaxed font-light">
                        DQ200, DQ250, DQ381 ve Tiptronic - Kronik sorunlar, bakım önerileri ve teknik analizler.
                    </p>
                </div>

                {/* Warning Box - Premium Glass */}
                <div className="bg-amber-900/20 border border-amber-500/20 rounded-2xl p-6 mb-10 flex items-start gap-4 backdrop-blur-md shadow-lg">
                    <div className="p-2.5 rounded-xl bg-amber-500/10 text-amber-500 shrink-0">
                        <AlertTriangle className="w-6 h-6" />
                    </div>
                    <div>
                        <h4 className="text-amber-500 font-bold mb-1">Önemli Yasal Uyarı</h4>
                        <p className="text-amber-200/80 text-sm leading-relaxed">
                            Bu sayfada yer alan bilgiler genel teknik verilere ve servis tecrübelerine dayanmaktadır.
                            Aracınızın spesifik durumu, üretim yılı ve yazılım versiyonuna göre değişiklik gösterebilir.
                            Ciddi mekanik sorunlarda mutlaka yetkili veya uzman servise başvurun.
                        </p>
                    </div>
                </div>

                {/* DSG Variants - Grid Layout */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-16">
                    {DSG_VARIANTS.map((variant) => {
                        const isProduction = variant.status.includes('Hala Üretimde') || variant.status.includes('Hala Kullanımda');
                        const isPartial = variant.status.includes('Kısmen');

                        return (
                            <button
                                key={variant.code}
                                onClick={() => handleCardClick(variant.code)}
                                className={`
                                    group relative flex flex-col h-full text-left transition-all duration-500 rounded-[2rem] overflow-hidden min-h-[220px]
                                    ${expandedCode === variant.code ? 'scale-[1.02] shadow-[0_0_50px_rgba(0,0,0,0.5)] z-10' : 'hover:-translate-y-2 hover:shadow-2xl'}
                                `}
                            >
                                {/* Background & Borders */}
                                <div className={`absolute inset-0 bg-gray-900/80 backdrop-blur-xl border border-white/10 transition-colors duration-500 ${expandedCode === variant.code ? 'bg-gray-900/95 border-white/20' : 'group-hover:bg-gray-900/90 group-hover:border-white/20'}`} />

                                {/* Dynamic Gradient Glow */}
                                <div className={`absolute -inset-0.5 bg-gradient-to-br ${variant.color} opacity-0 group-hover:opacity-20 transition-opacity duration-700 blur-2xl`} />

                                {/* Status Indicator Badge */}
                                <div className="absolute top-6 right-6 z-20">
                                    <div className={`
                                        flex items-center gap-2 px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest backdrop-blur-md border shadow-lg
                                        ${isProduction
                                            ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20 shadow-emerald-900/20'
                                            : isPartial
                                                ? 'bg-amber-500/10 text-amber-400 border-amber-500/20 shadow-amber-900/20'
                                                : 'bg-rose-500/10 text-rose-400 border-rose-500/20 shadow-rose-900/20'}
                                    `}>
                                        <div className={`w-1.5 h-1.5 rounded-full ${isProduction ? 'bg-emerald-400 animate-pulse' : isPartial ? 'bg-amber-400' : 'bg-rose-400'}`} />
                                        {isProduction ? 'Aktif' : isPartial ? 'Kısmen' : 'Eski'}
                                    </div>
                                </div>

                                <div className="relative z-10 p-8 flex flex-col h-full">
                                    {/* Header Info */}
                                    <div className="mb-6">
                                        <div className="flex items-center gap-3 mb-3">
                                            <div className={`p-2 rounded-lg bg-white/5 border border-white/10 ${expandedCode === variant.code ? 'text-white' : 'text-gray-400 group-hover:text-white'} transition-colors`}>
                                                <TrendingUp className="w-5 h-5" />
                                            </div>
                                            <span className="text-xs font-bold text-gray-500 uppercase tracking-widest">DSG Serisi</span>
                                        </div>
                                        <h3 className="text-4xl font-black text-white tracking-tighter mb-2 group-hover:scale-105 origin-left transition-transform duration-300">
                                            {variant.code}
                                        </h3>
                                        <div className="h-1 w-12 rounded-full bg-gradient-to-r from-white/20 to-transparent group-hover:w-24 group-hover:from-primary group-hover:to-primary/20 transition-all duration-500" />
                                    </div>

                                    <div className="mt-auto pt-4 border-t border-white/5 group-hover:border-white/10 transition-colors">
                                        <span className="text-xs font-medium text-gray-400 group-hover:text-white transition-colors flex items-center justify-between">
                                            Detayları İncele
                                            <div className="bg-white/10 p-2 rounded-full group-hover:bg-white/20 transition-all group-hover:translate-x-1">
                                                <ChevronRight className="w-4 h-4 text-white" />
                                            </div>
                                        </span>
                                    </div>
                                </div>
                            </button>
                        )
                    })}
                </div>

                {/* Detail Modal */}
                <AnimatePresence>
                    {expandedVariant && (
                        <DSGVariantModal
                            key="dsg-guide-modal"
                            variant={expandedVariant}
                            onClose={() => setExpandedCode(null)}
                        />
                    )}
                </AnimatePresence>

                {/* Back Button */}
                <div className="mt-12 text-center">
                    <Link
                        to="/guides/transmission"
                        className="inline-flex items-center gap-3 px-8 py-4 bg-black/40 hover:bg-black/60 text-white font-bold rounded-2xl transition-all border border-white/10 hover:border-white/20 hover:scale-105 active:scale-95 shadow-lg backdrop-blur-md"
                    >
                        <ArrowLeft className="w-5 h-5" />
                        Şanzıman Rehberine Dön
                    </Link>
                </div>
            </div>
        </div>
    );
}
