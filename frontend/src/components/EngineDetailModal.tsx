import React, { useState } from 'react';
import { X, Lightbulb, AlertTriangle, Wrench, Settings, TrendingUp, Star } from 'lucide-react';
import type { Engine, EngineProblem, EngineSolution, TechnologyFeature } from '../types/engine';
import { parseProblems, parseSolutions, parseTechnologyFeatures, getSeverityBadgeColor } from '../types/engine';

interface EngineDetailModalProps {
    engine: Engine | null;
    onClose: () => void;
}

export default function EngineDetailModal({ engine, onClose }: EngineDetailModalProps) {
    if (!engine) return null;

    const [activeTab, setActiveTab] = useState<'info' | 'problems' | 'solutions'>('info');

    const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    const problems: EngineProblem[] = parseProblems(engine);
    const solutions: EngineSolution[] = parseSolutions(engine);
    const techFeatures: TechnologyFeature[] = parseTechnologyFeatures(engine);

    // Determine fuel type color
    const getFuelTypeColor = (fuelType?: string) => {
        if (!fuelType) return 'from-gray-500 to-gray-600';
        const fuel = fuelType.toLowerCase();
        if (fuel.includes('diesel')) return 'from-amber-500 to-orange-600';
        if (fuel.includes('petrol') || fuel.includes('gasoline')) return 'from-blue-500 to-indigo-600';
        if (fuel.includes('hybrid')) return 'from-green-500 to-emerald-600';
        if (fuel.includes('electric')) return 'from-purple-500 to-violet-600';
        return 'from-gray-500 to-gray-600';
    };

    const productionStatus = engine.productionEndYear
        ? `${engine.productionStartYear || '?'} - ${engine.productionEndYear}`
        : `${engine.productionStartYear || '?'} - Devam Ediyor`;

    const isInProduction = !engine.productionEndYear;

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
                    <div className={`absolute inset-0 bg-gradient-to-br ${getFuelTypeColor(engine.fuelType)} opacity-20`} />
                    <div className="absolute inset-0 bg-gradient-to-b from-transparent via-[#0a0a0a]/50 to-[#0a0a0a]" />

                    {/* Background Big Text */}
                    <div className="absolute -right-10 -top-20 text-[200px] font-black text-white/5 select-none pointer-events-none tracking-tighter leading-none">
                        {engine.code}
                    </div>

                    <div className="relative px-8 pt-20 pb-12 sm:px-12">
                        <div className="flex flex-col md:flex-row items-end gap-6 mb-8">
                            <div>
                                <div className="flex items-center gap-3 mb-4">
                                    {engine.manufacturer && (
                                        <span className="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border border-white/10 bg-white/5 text-white">
                                            {engine.manufacturer}
                                        </span>
                                    )}
                                    {engine.fuelType && (
                                        <span className={`px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border bg-gradient-to-r ${getFuelTypeColor(engine.fuelType)} text-white`}>
                                            {engine.fuelType}
                                        </span>
                                    )}
                                    {engine.aspiration && (
                                        <span className="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border border-white/10 bg-white/5 text-gray-300">
                                            {engine.aspiration}
                                        </span>
                                    )}
                                </div>
                                <h2 className="text-6xl sm:text-7xl font-black text-white tracking-tighter mb-2 drop-shadow-2xl">
                                    {engine.code}
                                </h2>
                                <p className="text-xl text-gray-400 font-light max-w-2xl leading-relaxed">
                                    {engine.name || 'Motor Bilgileri'}
                                </p>
                            </div>
                        </div>

                        {/* Interactive Stats Bar */}
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            {engine.displacementCC && (
                                <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                    <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Motor Hacmi</div>
                                    <div className="text-2xl font-black text-white">{engine.displacementCC} <span className="text-sm font-medium text-gray-500">cc</span></div>
                                </div>
                            )}
                            {engine.cylinders && (
                                <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                    <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Silindir</div>
                                    <div className="text-2xl font-black text-white">{engine.cylinders} <span className="text-sm font-medium text-gray-500">{engine.cylinderLayout || ''}</span></div>
                                </div>
                            )}
                            <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Üretim</div>
                                <div className={`flex items-center gap-2 font-bold text-sm ${isInProduction ? 'text-emerald-400' : 'text-gray-400'}`}>
                                    {productionStatus}
                                </div>
                            </div>
                            {engine.reliabilityRating && (
                                <div className="bg-white/5 backdrop-blur-md rounded-2xl p-4 border border-white/5">
                                    <div className="text-[10px] text-gray-500 uppercase font-bold mb-1">Güvenilirlik</div>
                                    <div className="flex items-center gap-1">
                                        {[...Array(10)].map((_, i) => (
                                            <Star
                                                key={i}
                                                className={`w-3 h-3 ${i < (engine.reliabilityRating || 0) ? 'text-yellow-500 fill-yellow-500' : 'text-gray-600'}`}
                                            />
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                {/* Tabs */}
                <div className="shrink-0 px-8 sm:px-12 border-b border-white/5">
                    <div className="flex gap-1">
                        <button
                            onClick={() => setActiveTab('info')}
                            className={`px-6 py-3 font-bold text-sm transition-all relative ${activeTab === 'info'
                                    ? 'text-white'
                                    : 'text-gray-500 hover:text-gray-300'
                                }`}
                        >
                            Bilgiler
                            {activeTab === 'info' && (
                                <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-white rounded-t-full" />
                            )}
                        </button>
                        <button
                            onClick={() => setActiveTab('problems')}
                            className={`px-6 py-3 font-bold text-sm transition-all relative ${activeTab === 'problems'
                                    ? 'text-white'
                                    : 'text-gray-500 hover:text-gray-300'
                                }`}
                        >
                            Sorunlar {problems.length > 0 && `(${problems.length})`}
                            {activeTab === 'problems' && (
                                <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-white rounded-t-full" />
                            )}
                        </button>
                        <button
                            onClick={() => setActiveTab('solutions')}
                            className={`px-6 py-3 font-bold text-sm transition-all relative ${activeTab === 'solutions'
                                    ? 'text-white'
                                    : 'text-gray-500 hover:text-gray-300'
                                }`}
                        >
                            Çözümler {solutions.length > 0 && `(${solutions.length})`}
                            {activeTab === 'solutions' && (
                                <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-white rounded-t-full" />
                            )}
                        </button>
                    </div>
                </div>

                {/* Scrollable Content */}
                <div className="overflow-y-auto px-8 pb-12 sm:px-12 flex-1 no-scrollbar">

                    {/* Info Tab */}
                    {activeTab === 'info' && (
                        <div className="pt-8 space-y-8">
                            {/* Description */}
                            {engine.description && (
                                <div className="text-lg text-gray-300 leading-relaxed font-light">
                                    <p>{engine.description}</p>
                                </div>
                            )}

                            {/* Technical Specs */}
                            <div className="grid md:grid-cols-2 gap-6">
                                <div className="relative overflow-hidden rounded-3xl bg-white/5 border border-white/10 p-6">
                                    <div className="absolute -right-10 -bottom-10 opacity-10">
                                        <Settings className="w-32 h-32 text-gray-500" />
                                    </div>
                                    <h3 className="text-xl font-bold text-white mb-4 flex items-center gap-3 relative z-10">
                                        <div className="bg-white/10 p-2 rounded-lg">
                                            <Settings className="w-5 h-5 text-gray-300" />
                                        </div>
                                        Teknik Özellikler
                                    </h3>
                                    <div className="space-y-3 relative z-10">
                                        {engine.engineType && (
                                            <div className="flex justify-between items-center text-sm">
                                                <span className="text-gray-500">Motor Tipi:</span>
                                                <span className="text-white font-bold">{engine.engineType}</span>
                                            </div>
                                        )}
                                        {engine.valvesPerCylinder && (
                                            <div className="flex justify-between items-center text-sm">
                                                <span className="text-gray-500">Supap/Silindir:</span>
                                                <span className="text-white font-bold">{engine.valvesPerCylinder}</span>
                                            </div>
                                        )}
                                    </div>
                                </div>

                                {/* Technology Features */}
                                {techFeatures.length > 0 && (
                                    <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-blue-500/10 to-indigo-500/10 border border-blue-500/20 p-6">
                                        <h3 className="text-xl font-bold text-blue-300 mb-4 flex items-center gap-3">
                                            <div className="bg-blue-500/20 p-2 rounded-lg">
                                                <TrendingUp className="w-5 h-5 text-blue-400" />
                                            </div>
                                            Teknolojiler
                                        </h3>
                                        <div className="space-y-2">
                                            {techFeatures.map((feature, i) => (
                                                <div key={i} className="flex gap-3 text-sm text-blue-200/70 bg-blue-500/5 p-3 rounded-xl border border-blue-500/10">
                                                    <div className="w-1.5 h-1.5 rounded-full bg-blue-500 mt-2 shrink-0" />
                                                    <div>
                                                        <div className="font-bold text-blue-300">{feature.name}</div>
                                                        {feature.description && (
                                                            <div className="text-xs text-blue-200/50 mt-1">{feature.description}</div>
                                                        )}
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>

                            {/* Maintenance Notes */}
                            {engine.maintenanceNotes && (
                                <div className="relative overflow-hidden rounded-3xl bg-emerald-950/10 border border-emerald-500/20 p-6">
                                    <h3 className="text-xl font-bold text-emerald-300 mb-4 flex items-center gap-3">
                                        <div className="bg-emerald-500/20 p-2 rounded-lg">
                                            <Lightbulb className="w-5 h-5 text-emerald-400" />
                                        </div>
                                        Bakım Notları
                                    </h3>
                                    <p className="text-emerald-100/80 leading-relaxed">
                                        {engine.maintenanceNotes}
                                    </p>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Problems Tab */}
                    {activeTab === 'problems' && (
                        <div className="pt-8 space-y-4">
                            {problems.length > 0 ? (
                                problems.map((problem, i) => (
                                    <div key={i} className="relative overflow-hidden rounded-3xl bg-rose-950/10 border border-rose-500/20 p-6">
                                        <div className="flex items-start gap-4">
                                            <div className="p-2 rounded-lg bg-rose-500/10 text-rose-500 shrink-0">
                                                <AlertTriangle className="w-5 h-5" />
                                            </div>
                                            <div className="flex-1">
                                                <div className="flex items-center gap-3 mb-2">
                                                    <h4 className="text-lg font-bold text-rose-100">{problem.title}</h4>
                                                    <span className={`px-2 py-1 rounded-full text-[10px] font-black uppercase border ${getSeverityBadgeColor(problem.severity)}`}>
                                                        {problem.severity === 'high' ? 'Yüksek' : problem.severity === 'medium' ? 'Orta' : 'Düşük'}
                                                    </span>
                                                </div>
                                                <p className="text-rose-200/70 leading-relaxed">{problem.description}</p>
                                            </div>
                                        </div>
                                    </div>
                                ))
                            ) : (
                                <div className="text-center py-12 text-gray-500">
                                    <AlertTriangle className="w-16 h-16 mx-auto mb-4 opacity-20" />
                                    <p>Henüz bilinen sorun kaydedilmemiş</p>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Solutions Tab */}
                    {activeTab === 'solutions' && (
                        <div className="pt-8 space-y-4">
                            {solutions.length > 0 ? (
                                solutions.map((solution, i) => (
                                    <div key={i} className="relative overflow-hidden rounded-3xl bg-emerald-950/10 border border-emerald-500/20 p-6">
                                        <div className="flex items-start gap-4">
                                            <div className="p-2 rounded-lg bg-emerald-500/10 text-emerald-500 shrink-0">
                                                <Wrench className="w-5 h-5" />
                                            </div>
                                            <div className="flex-1">
                                                <h4 className="text-lg font-bold text-emerald-100 mb-2">{solution.problemTitle}</h4>
                                                <p className="text-emerald-200/70 leading-relaxed mb-3">{solution.solution}</p>
                                                {solution.estimatedCost && (
                                                    <div className="inline-block px-3 py-1 rounded-full bg-emerald-500/10 text-emerald-400 text-sm font-bold border border-emerald-500/20">
                                                        Tahmini Maliyet: {solution.estimatedCost}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                ))
                            ) : (
                                <div className="text-center py-12 text-gray-500">
                                    <Wrench className="w-16 h-16 mx-auto mb-4 opacity-20" />
                                    <p>Henüz çözüm önerisi kaydedilmemiş</p>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
