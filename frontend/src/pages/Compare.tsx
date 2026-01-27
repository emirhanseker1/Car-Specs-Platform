import { useMemo, useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Trophy, X } from 'lucide-react';

interface Spec {
    category: string;
    name: string;
    value: string;
}

interface CompareItem {
    id: number;
    title?: string;
    specs?: unknown;
}

function normalizeSpecs(raw: unknown): Spec[] {
    if (Array.isArray(raw)) return raw as Spec[];
    if (typeof raw === 'string') {
        const trimmed = raw.trim();
        if (!trimmed) return [];
        try {
            const parsed = JSON.parse(trimmed);
            if (Array.isArray(parsed)) return parsed as Spec[];
            if (Array.isArray((parsed as any)?.specs)) return (parsed as any).specs as Spec[];
        } catch {
            return [];
        }
    }
    return [];
}

function parseValue(valStr: string): number | null {
    if (!valStr) return null;
    const match = valStr.match(/(\d+[.,]?\d*)/);
    if (!match) return null;
    const num = parseFloat(match[0].replace(',', '.'));
    return isNaN(num) ? null : num;
}

type MetricType = 'higher-better' | 'lower-better' | 'neutral';

function getMetricType(name: string): MetricType {
    const n = name.toLowerCase();

    // Lower is better
    if (n.includes('acceleration') || n.includes('0-100') || n.includes('consumption') || n.includes('tüketim') || n.includes('emission') || n.includes('co2')) {
        return 'lower-better';
    }

    // Higher is better
    if (n.includes('power') || n.includes('hp') || n.includes('ps') || n.includes('torque') || n.includes('nm') || n.includes('speed') || n.includes('hız') || n.includes('range') || n.includes('menzil') || n.includes('capacity') || n.includes('bagaj') || n.includes('trunk')) {
        return 'higher-better';
    }

    return 'neutral';
}

function Compare() {
    const [items, setItems] = useState<CompareItem[]>([]);

    useEffect(() => {
        const stored = localStorage.getItem('compareList');
        if (stored) {
            setItems(JSON.parse(stored));
        }
    }, []);

    const normalizedItems = useMemo(() => {
        return items
            .map((item) => ({
                id: item.id,
                title: item.title || 'Vehicle',
                specs: normalizeSpecs(item.specs),
            }))
            .filter((x) => Number.isFinite(x.id));
    }, [items]);

    const removeItem = (id: number) => {
        const newItems = items.filter(i => i.id !== id);
        setItems(newItems);
        localStorage.setItem('compareList', JSON.stringify(newItems));
    };

    const clearAll = () => {
        setItems([]);
        localStorage.setItem('compareList', JSON.stringify([]));
    };

    // Find all unique spec names to align rows
    const allSpecNames = Array.from(new Set(
        normalizedItems.flatMap(item => item.specs.map(s => s.name))
    ));

    // Calculate ranges for bars
    const specRanges = useMemo(() => {
        const ranges: Record<string, { min: number, max: number }> = {};
        allSpecNames.forEach(name => {
            const values = normalizedItems
                .map(item => item.specs.find(s => s.name === name)?.value)
                .map(v => v ? parseValue(v) : null)
                .filter((v): v is number => v !== null);

            if (values.length > 0) {
                ranges[name] = {
                    min: Math.min(...values),
                    max: Math.max(...values)
                };
            }
        });
        return ranges;
    }, [allSpecNames, normalizedItems]);

    return (
        <div className="relative min-h-screen text-slate-200 font-sans selection:bg-blue-500/30">
            {/* Background Image & Overlay */}
            <div className="fixed inset-0 z-0">
                <div
                    className="absolute inset-0 bg-cover bg-center"
                    style={{ backgroundImage: 'url(/hero-2.jpg)' }}
                ></div>
                <div className="absolute inset-0 bg-slate-900/90 backdrop-blur-sm"></div>
            </div>

            <main className="relative z-10 max-w-7xl mx-auto px-4 py-8 space-y-8 pt-32">
                {/* Header */}
                <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
                    <div className="space-y-1">
                        <h1 className="text-3xl sm:text-4xl font-bold tracking-tight text-white">Karşılaştır</h1>
                        <p className="text-sm text-slate-400">İki donanımı yan yana kıyasla.</p>
                    </div>
                    {normalizedItems.length > 0 && (
                        <div className="flex items-center gap-2">
                            <button
                                type="button"
                                onClick={clearAll}
                                className="rounded-xl bg-white/5 hover:bg-white/10 border border-white/10 px-4 py-2 text-sm font-semibold text-slate-200 transition-colors"
                            >
                                Temizle
                            </button>
                            <Link
                                to="/"
                                className="rounded-xl bg-primary hover:bg-primary-hover text-white px-4 py-2 text-sm font-bold transition-colors shadow-lg shadow-orange-500/20"
                            >
                                Keşfet
                            </Link>
                        </div>
                    )}
                </div>

                {normalizedItems.length === 0 ? (
                    <div className="rounded-3xl bg-[#1e293b]/60 backdrop-blur-md border border-white/10 p-12 text-center space-y-6 max-w-lg mx-auto mt-20">
                        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-white/5 mb-2">
                            <Trophy className="w-8 h-8 text-slate-500" />
                        </div>
                        <div>
                            <h2 className="text-2xl font-bold text-white mb-2">Karşılaştırma Listesi Boş</h2>
                            <p className="text-slate-400">Bir araç detay sayfasından “Add to Compare” butonunu kullanarak listeye araç ekleyebilirsin.</p>
                        </div>
                        <Link to="/" className="inline-flex items-center justify-center rounded-xl bg-primary hover:bg-primary-hover text-white px-8 py-3 font-bold transition-all shadow-lg shadow-orange-500/20 hover:scale-105">
                            Modellere Göz At
                        </Link>
                    </div>
                ) : (
                    <div className="overflow-x-auto bg-[#1e293b]/60 backdrop-blur-md rounded-3xl border border-white/10 shadow-2xl">
                        <table className="w-full text-left border-collapse table-fixed min-w-[600px]">
                            <thead>
                                <tr className="border-b border-white/10 bg-white/5">
                                    <th className="p-5 font-bold text-slate-400 uppercase text-xs tracking-wider w-48">Özellik</th>
                                    {normalizedItems.map(item => (
                                        <th key={item.id} className="p-5 font-bold text-white relative align-top">
                                            <div className="pr-10">
                                                {item.title}
                                            </div>
                                            <button
                                                type="button"
                                                onClick={() => removeItem(item.id)}
                                                className="absolute top-4 right-4 text-slate-500 hover:text-red-400 hover:bg-red-500/10 p-2 rounded-full transition-colors"
                                                aria-label="Remove"
                                            >
                                                <X className="w-4 h-4" />
                                            </button>
                                        </th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                {allSpecNames.map(specName => {
                                    const range = specRanges[specName];
                                    const type = getMetricType(specName);
                                    let bestVal: number | null = null;
                                    if (range && type !== 'neutral') {
                                        bestVal = type === 'higher-better' ? range.max : range.min;
                                    }

                                    return (
                                        <tr key={specName} className="hover:bg-white/5 transition-colors group">
                                            <td className="p-5 text-sm font-medium text-slate-400 break-words w-48 group-hover:text-slate-300 transition-colors">{specName}</td>
                                            {normalizedItems.map(item => {
                                                const spec = item.specs.find(s => s.name === specName);
                                                const valStr = spec?.value;
                                                const valNum = valStr ? parseValue(valStr) : null;

                                                const isWinner = bestVal !== null && valNum !== null && valNum === bestVal;

                                                let barPercent = 0;
                                                if (range && valNum !== null && range.max > 0) {
                                                    barPercent = (valNum / range.max) * 100;
                                                }

                                                return (
                                                    <td key={item.id} className="p-5 relative">
                                                        <div className={`relative z-10 text-sm font-medium flex items-center gap-2 ${isWinner ? 'text-green-400 font-bold' : 'text-slate-200'}`}>
                                                            {valStr || '-'}
                                                            {isWinner && <Trophy className="w-4 h-4 text-yellow-500 fill-yellow-500" />}
                                                        </div>
                                                        {/* Visual Bar Background */}
                                                        {barPercent > 0 && (
                                                            <div
                                                                className="absolute bottom-0 left-0 h-0.5 bg-blue-500/30 transition-all duration-500"
                                                                style={{ width: `${Math.min(barPercent, 100)}%` }}
                                                            />
                                                        )}
                                                    </td>
                                                );
                                            })}
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                )}
            </main>
        </div>
    );
}

export default Compare;
