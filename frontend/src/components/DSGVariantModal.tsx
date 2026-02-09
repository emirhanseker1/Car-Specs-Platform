import React from 'react';
import { X, Lightbulb, AlertTriangle, Wrench, History } from 'lucide-react';
import type { DSGVariant } from '../data/dsgVariants';

interface DSGVariantModalProps {
    variant: DSGVariant | null;
    onClose: () => void;
}

export default function DSGVariantModal({ variant, onClose }: DSGVariantModalProps) {
    if (!variant) return null;

    const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    const isProduction = variant.status.includes('Hala Üretimde') || variant.status.includes('Hala Kullanımda');
    const isPartial = variant.status.includes('Kısmen');

    return (
        <div className="fixed inset-0 z-[9999] flex items-center justify-center p-0 sm:p-4 text-left font-sans" style={{ margin: 0 }}>
            <div
                className="absolute inset-0 bg-black/90 backdrop-blur-xl animate-in fade-in duration-300"
                onClick={handleBackdropClick}
            />

            <div className="relative w-full h-full sm:h-auto sm:max-h-[95vh] sm:max-w-6xl bg-[#0a0a0a] sm:rounded-[2.5rem] shadow-2xl overflow-hidden animate-in zoom-in-95 duration-300 flex flex-col border border-white/10">

                {/* Close Button */}
                <button
                    onClick={onClose}
                    className="absolute top-6 right-6 z-50 p-2 rounded-full bg-black/50 text-white/70 hover:text-white hover:bg-black/80 backdrop-blur-md transition-all border border-white/10"
                >
                    <X className="w-6 h-6" />
                </button>

                {/* Hero Header */}
                <div className="relative shrink-0 overflow-hidden">
                    <div className={`absolute inset-0 bg-gradient-to-br ${variant.color} opacity-20`} />
                    <div className="absolute inset-0 bg-gradient-to-b from-transparent via-[#0a0a0a]/50 to-[#0a0a0a]" />

                    {/* Background Big Text */}
                    <div className="absolute -right-10 -top-20 text-[200px] font-black text-white/5 select-none pointer-events-none tracking-tighter leading-none">
                        {variant.code}
                    </div>

                    <div className="relative px-8 pt-20 pb-12 sm:px-12">
                        <div className="flex flex-col md:flex-row items-end gap-6 mb-8">
                            <div>
                                <div className="flex items-center gap-3 mb-4">
                                    <span className={`px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border bg-white/5 backdrop-blur-sm ${variant.borderColor} text-white`}>
                                        {variant.clutchTypeLabel}
                                    </span>
                                    {variant.gears > 0 && (
                                        <span className="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border border-white/10 bg-white/5 text-gray-300">
                                            {variant.gears} İleri
                                        </span>
                                    )}
                                </div>
                                <h2 className="text-6xl sm:text-7xl font-black text-white tracking-tighter mb-2 drop-shadow-2xl">
                                    {variant.code}
                                </h2>
                                <p className="text-xl text-gray-400 font-light max-w-2xl leading-relaxed">
                                    {variant.name}
                                </p>
                            </div>
                        </div>

                        {/* Interactive Stats Bar */}
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Tork Dayanımı</div>
                                <div className="text-2xl font-black text-white">{variant.maxTorque} <span className="text-sm font-medium text-gray-500">Nm</span></div>
                            </div>
                            <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Durum</div>
                                <div className={`flex items-center gap-2 font-bold ${isProduction ? 'text-emerald-400' : isPartial ? 'text-amber-400' : 'text-rose-400'}`}>
                                    {isProduction ? 'Aktif Üretim' : isPartial ? 'Kısmen Aktif' : 'Üretim Bitti'}
                                </div>
                            </div>
                            <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Üretim Yılı</div>
                                <div className="text-lg font-bold text-white">{variant.productionYears.split('-')[0].trim()} <span className="text-gray-600 mx-1">→</span> {variant.productionYears.split('-')[1]?.trim() || '...'}</div>
                            </div>
                            <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Kavrama Ömrü</div>
                                <div className="text-lg font-bold text-white">{variant.clutchInterval}</div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Scrollable Content */}
                <div className="overflow-y-auto px-8 pb-12 sm:px-12 flex-1 no-scrollbar">

                    <div className="grid lg:grid-cols-12 gap-10">
                        {/* Main Info Column */}
                        <div className="lg:col-span-7 space-y-8">
                            {/* Description */}
                            <div className="text-lg text-gray-300 leading-relaxed font-light">
                                <p>{variant.description}</p>
                            </div>

                            {/* History Section */}
                            <div className="relative overflow-hidden rounded-3xl bg-white/5 border border-white/10 p-8">
                                <div className={`absolute -right-10 -bottom-10 opacity-10`}>
                                    <History className="w-40 h-40 text-gray-500" />
                                </div>
                                <h3 className="text-xl font-bold text-white mb-4 flex items-center gap-3 relative z-10">
                                    <div className="bg-white/10 p-2 rounded-lg">
                                        <History className="w-5 h-5 text-gray-300" />
                                    </div>
                                    Tarihçe ve Gelişim
                                </h3>
                                <p className="text-gray-400 leading-relaxed relative z-10 font-medium whitespace-pre-line">
                                    {variant.history}
                                </p>
                            </div>

                            {/* Expert Tip */}
                            <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-blue-500/10 to-indigo-500/10 border border-blue-500/20 p-8">
                                <h3 className="text-xl font-bold text-blue-300 mb-4 flex items-center gap-3">
                                    <div className="bg-blue-500/20 p-2 rounded-lg">
                                        <Lightbulb className="w-5 h-5 text-blue-400" />
                                    </div>
                                    Uzman Görüşü
                                </h3>
                                <p className="text-blue-100/80 leading-relaxed text-lg italic">
                                    "{variant.smartTip}"
                                </p>
                            </div>
                        </div>

                        {/* Side Info Column */}
                        <div className="lg:col-span-5 space-y-6">

                            {/* Problems */}
                            <div className="bg-rose-950/10 border border-rose-500/20 rounded-3xl p-6">
                                <div className="flex items-center gap-3 mb-6">
                                    <div className="p-2 rounded-lg bg-rose-500/10 text-rose-500">
                                        <AlertTriangle className="w-5 h-5" />
                                    </div>
                                    <h3 className="text-lg font-bold text-rose-100">Kronik Sorunlar</h3>
                                </div>
                                <div className="space-y-3">
                                    {variant.chronicProblems.map((problem, i) => (
                                        <div key={i} className="flex gap-3 text-sm text-rose-200/70 bg-rose-500/5 p-3 rounded-xl border border-rose-500/10">
                                            <div className="w-1.5 h-1.5 rounded-full bg-rose-500 mt-2 shrink-0" />
                                            {problem}
                                        </div>
                                    ))}
                                </div>
                            </div>

                            {/* Maintenance */}
                            <div className="bg-emerald-950/10 border border-emerald-500/20 rounded-3xl p-6">
                                <div className="flex items-center gap-3 mb-6">
                                    <div className="p-2 rounded-lg bg-emerald-500/10 text-emerald-500">
                                        <Wrench className="w-5 h-5" />
                                    </div>
                                    <h3 className="text-lg font-bold text-emerald-100">Bakım Önerileri</h3>
                                </div>
                                <div className="space-y-3">
                                    {variant.maintenanceTips.map((tip, i) => (
                                        <div key={i} className="flex gap-3 text-sm text-emerald-200/70 bg-emerald-500/5 p-3 rounded-xl border border-emerald-500/10">
                                            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 mt-2 shrink-0" />
                                            {tip}
                                        </div>
                                    ))}
                                </div>
                            </div>

                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
