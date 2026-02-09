import { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { AnimatePresence } from 'framer-motion';
import {
    ArrowLeft,
    Timer,
    Zap,
    Gauge,
    Weight,
    Download,
    Plus,
    ChevronDown,
    ChevronUp,
    CarFront,
    Info,
    Settings,
} from 'lucide-react';
import { api } from '../services/api';
import type { DetailResponse, Spec, Trim, Vehicle } from '../types';
import { DSG_VARIANTS } from '../data/dsgVariants';
import DSGVariantModal from '../components/DSGVariantModal';
import { detectTransmissionVariant } from '../utils/dsgMatcher';

// Helper to determine progress bar percentage for specs
const getSpecPercentage = (name: string, valueStr: string): number => {
    const val = parseFloat(valueStr.replace(/[^0-9.]/g, ''));
    if (isNaN(val)) return 0;

    const n = name.toLowerCase();
    // Approximate max values for visualization
    if (n.includes('power') || n.includes('hp') || n.includes('ps')) return Math.min((val / 800) * 100, 100);
    if (n.includes('torque') || n.includes('nm')) return Math.min((val / 1000) * 100, 100);
    if (n.includes('speed') || n.includes('km/h')) return Math.min((val / 350) * 100, 100); // 350 kmh max
    if (n.includes('acceleration') || n.includes('0-100')) {
        // Inverse: lower is better. 2s is 100%, 15s is 0%
        return Math.max(0, Math.min(((15 - val) / 13) * 100, 100));
    }
    if (n.includes('consumption') || n.includes('l/100')) {
        // Inverse: lower is better. 3L is 100%, 20L is 0%
        return Math.max(0, Math.min(((20 - val) / 17) * 100, 100));
    }
    if (n.includes('displacement') || n.includes('cc')) return Math.min((val / 6000) * 100, 100);

    return 0;
};

function VehicleDetails() {
    const navigate = useNavigate();
    const { brandName, modelName, generationCode, trimIndex } = useParams();
    const [data, setData] = useState<DetailResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedTrimId, setSelectedTrimId] = useState<number | null>(null);
    const [openCategories, setOpenCategories] = useState<Record<string, boolean>>({});
    const [expandedVariantCode, setExpandedVariantCode] = useState<string | null>(null);

    const goBack = () => navigate(-1);

    useEffect(() => {
        setError(null);
        setLoading(true);

        const fetchData = async () => {
            try {
                if (brandName && modelName && generationCode && trimIndex) {
                    // 1. Resolve Brand
                    const brands = await api.getBrands();
                    const brand = brands.find((b: any) => b.name.toLowerCase() === brandName.toLowerCase());
                    if (!brand) throw new Error(`Brand '${brandName}' not found`);

                    // 2. Resolve Model
                    const modelsData = await api.getModels(brand.id);
                    const models = modelsData?.value || [];
                    const model = models.find((m: any) => m.name.toLowerCase() === modelName.toLowerCase());
                    if (!model) throw new Error(`Model '${modelName}' not found`);

                    // 3. Resolve Generation
                    const gensData = await api.getGenerations(model.id);
                    const generations = !gensData ? [] : (Array.isArray(gensData) ? gensData : (gensData.value || []));
                    const generation = generations.find((g: any) => g.code.toLowerCase() === generationCode.toLowerCase());
                    if (!generation) throw new Error(`Generation '${generationCode}' not found`);

                    // 4. Fetch Trims
                    const trimsData = await api.getTrims(generation.id);
                    const trims = !trimsData ? [] : (Array.isArray(trimsData) ? trimsData : (trimsData.value || []));

                    if (trims.length === 0) throw new Error('No trims found for this generation');

                    // 5. Select Trim by Index
                    const index = parseInt(trimIndex, 10) - 1; // 1-based to 0-based
                    if (isNaN(index) || index < 0 || index >= trims.length) {
                        throw new Error(`Invalid trim index: ${trimIndex}`);
                    }
                    const selectedTrim = trims[index];

                    // Construct Data Object
                    const vehicleData: Vehicle = {
                        id: model.id,
                        brand: brand.name,
                        model: model.name,
                        generation: generation.code,
                        image_url: generation.image_url || selectedTrim.image_url
                    };

                    setData({
                        vehicle: vehicleData,
                        trims: trims
                    });
                    setSelectedTrimId(selectedTrim.id);
                }
            } catch (err) {
                console.error('Failed to fetch details:', err);
                const msg = err instanceof Error ? err.message : 'Failed to load vehicle details.';
                setError(msg);
            } finally {
                setLoading(false);
            }
        };
        fetchData();

    }, [brandName, modelName, generationCode, trimIndex]);

    const toggleCategory = (category: string) => {
        setOpenCategories(prev => ({ ...prev, [category]: !prev[category] }));
    };

    const addToCompare = (trim: Trim, specsToStore: Spec[]) => {
        if (!data) return;
        const current = JSON.parse(localStorage.getItem('compareList') || '[]');
        if (current.length >= 2) {
            alert('You can only compare 2 vehicles. Remove one to add another.');
            return;
        }
        if (!current.find((t: any) => t.id === trim.id)) {
            const title = `${data.vehicle.brand} ${data.vehicle.model} • ${trim.name}`;
            const itemToStore = {
                id: trim.id,
                title,
                vehicleId: data.vehicle.id,
                brand: data.vehicle.brand,
                modelName: data.vehicle.model,
                trimName: trim.name,
                specs: specsToStore,
            };
            localStorage.setItem('compareList', JSON.stringify([...current, itemToStore]));
            alert('Added to comparison!');
        } else {
            alert('Already in comparison list.');
        }
    };

    const trims = data?.trims || [];

    const inferCategoryFromName = (name: string) => {
        const n = (name || '').trim().toLowerCase();
        if (!n) return 'General';
        const has = (...needles: string[]) => needles.some((x) => n.includes(x));
        if (has('engine', 'motor', 'cylinder', 'silindir', 'displacement', 'capacity', 'cc', 'tsi', 'tdi', 'power', 'hp', 'kW', 'torque', 'nm')) return 'Engine & Performance';
        if (has('transmission', 'gearbox', 'vites', 'dsg', 'clutch')) return 'Transmission';
        if (has('0-100', 'acceleration', 'speed', 'performance')) return 'Engine & Performance';
        if (has('consumption', 'fuel', 'l/100', 'emission', 'co2')) return 'Efficiency';
        if (has('length', 'width', 'height', 'wheelbase', 'weight', 'mass', 'luggage', 'tank')) return 'Dimensions & Capacities';
        if (has('body', 'hatchback', 'sedan', 'doors', 'seats')) return 'Body';
        if (has('suspen', 'drive', 'brake', 'tire', 'wheel')) return 'Chassis';
        return 'General';
    };

    const currentTrim = useMemo(() => {
        if (!trims.length) return null;
        return trims.find((t) => t.id === selectedTrimId) || trims[0];
    }, [selectedTrimId, trims]);

    const normalizedSpecs = useMemo<Spec[]>(() => {
        if (!currentTrim) return [];
        let list: any[] = [];
        const knownKeys = ['id', 'name', 'year', 'generation', 'model_id', 'generation_id', 'created_at', 'updated_at', 'image_url', 'msrp_price', 'currency', 'market', 'is_facelift'];

        if (currentTrim && typeof currentTrim === 'object') {
            Object.entries(currentTrim).forEach(([k, v]) => {
                if (knownKeys.includes(k) || typeof v === 'object') return;
                if (v === null || v === 0 || v === '') return;

                const name = k.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
                const valStr = String(v);
                const category = inferCategoryFromName(name);

                let suffix = '';
                if (k.includes('_mm')) suffix = ' mm';
                else if (k.includes('_kg')) suffix = ' kg';
                else if (k.includes('_l')) suffix = ' L';
                else if (k.includes('_hp')) suffix = ' HP';
                else if (k.includes('_kw')) suffix = ' kW';
                else if (k.includes('_nm')) suffix = ' Nm';
                else if (k.includes('_cc')) suffix = ' cc';

                const displayName = name.replace(/( Mm| Kg| L| Hp| Kw| Nm| Cc| Kmh| 0 100| Comb| City| Hwy)$/i, '');

                list.push({
                    category,
                    name: displayName,
                    value: valStr + suffix
                });
            });
        }
        return list;
    }, [currentTrim]);

    const groupedSpecs = useMemo(() => {
        const grouped = normalizedSpecs.reduce((acc, spec) => {
            const cat = spec.category || 'General';
            if (!acc[cat]) acc[cat] = [];
            acc[cat].push(spec);
            return acc;
        }, {} as Record<string, Spec[]>);

        if (Object.keys(grouped).length > 0 && Object.keys(openCategories).length === 0) {
            setOpenCategories({ [Object.keys(grouped)[0]]: true });
        }

        return grouped;
    }, [normalizedSpecs]);

    const heroImageUrl = useMemo(() => {
        const raw = (data?.vehicle?.image_url) || (currentTrim as any)?.image_url;
        if (!raw) return '';
        if (raw.includes('/pictures/thumbs/')) return raw.replace(/\/pictures\/thumbs\/\d+px\//, '/pictures/');
        return raw;
    }, [data?.vehicle, currentTrim]);

    const quickStats = useMemo(() => {
        if (!currentTrim) return null;
        const t = currentTrim as any;
        return {
            accel: t.acceleration_0_100 ? `${t.acceleration_0_100}s` : '-',
            power: t.power_hp ? `${t.power_hp} HP` : '-',
            torque: t.torque_nm ? `${t.torque_nm} Nm` : '-',
            weight: t.curb_weight_kg ? `${t.curb_weight_kg} kg` : '-',
        };
    }, [currentTrim]);

    useEffect(() => {
        if (expandedVariantCode) {
            document.body.classList.add('modal-open');
        } else {
            document.body.classList.remove('modal-open');
        }
        return () => document.body.classList.remove('modal-open');
    }, [expandedVariantCode]);

    if (loading) return (
        <div className="min-h-screen bg-[#0f172a] flex items-center justify-center">
            <div className="text-blue-400 animate-pulse">Initializing Dashboard...</div>
        </div>
    );

    if (error) return (
        <div className="min-h-screen bg-[#0f172a] flex items-center justify-center p-4">
            <div className="bg-black/60 backdrop-blur-xl p-8 rounded-2xl border border-red-500/20 text-center max-w-md">
                <div className="text-red-400 font-bold mb-2">System Error</div>
                <div className="text-slate-400 text-sm mb-6">{error}</div>
                <button onClick={goBack} className="px-6 py-2 bg-slate-700 text-white rounded-lg hover:bg-slate-600 transition-colors">
                    Return
                </button>
            </div>
        </div>
    );

    if (!data) return null;

    return (
        <div className="min-h-screen text-slate-200 font-sans selection:bg-blue-500/30 pt-24">
            <div className="border-b border-white/10 bg-transparent mb-6">
                <div className="max-w-7xl mx-auto px-4 h-12 flex items-center gap-4">
                    <button onClick={goBack} className="p-2 hover:bg-white/10 rounded-full transition-colors text-slate-400 hover:text-white">
                        <ArrowLeft size={20} />
                    </button>
                    <div className="text-sm font-medium text-slate-400">
                        Vehicle Database <span className="mx-2">/</span> {data.vehicle.brand} <span className="mx-2">/</span> {data.vehicle.model}
                    </div>
                </div>
            </div>

            <main className="max-w-7xl mx-auto px-4 py-8 space-y-8">
                <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
                    {/* LEFT COLUMN: VISUALS */}
                    <div className="lg:col-span-6 xl:col-span-7 flex flex-col gap-4">
                        <div className="relative aspect-[16/10] bg-black/20 backdrop-blur-md rounded-2xl overflow-hidden border border-white/10 shadow-2xl group">
                            <div className="absolute inset-0 z-10 opacity-20 pointer-events-none"
                                style={{ backgroundImage: 'linear-gradient(rgba(59, 130, 246, 0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(59, 130, 246, 0.1) 1px, transparent 1px)', backgroundSize: '40px 40px' }}>
                            </div>

                            {heroImageUrl ? (
                                <img
                                    src={heroImageUrl}
                                    alt={data.vehicle.model}
                                    className="w-full h-full object-contain p-8 group-hover:scale-105 transition-transform duration-700 ease-out z-0 relative"
                                />
                            ) : (
                                <div className="w-full h-full flex items-center justify-center text-slate-600">
                                    <CarFront size={64} strokeWidth={1} />
                                </div>
                            )}

                            <div className="absolute top-4 left-4 z-20 flex flex-col gap-2">
                                <span className="bg-blue-500/10 backdrop-blur-md border border-blue-500/20 text-blue-400 text-xs font-bold px-3 py-1 rounded-full">
                                    {data.vehicle.generation || 'GEN'}
                                </span>
                            </div>
                        </div>

                        <div className="grid grid-cols-4 gap-3">
                            {[heroImageUrl, heroImageUrl, heroImageUrl].map((img, i) => (
                                <div key={i} className={`aspect-video rounded-lg overflow-hidden border ${i === 0 ? 'border-primary' : 'border-white/10 opacity-50 hover:opacity-100'} cursor-pointer transition-all bg-black/20 backdrop-blur-md`}>
                                    {img && <img src={img} className="w-full h-full object-cover" />}
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* RIGHT COLUMN: DATA */}
                    <div className="lg:col-span-6 xl:col-span-5 flex flex-col gap-6">
                        <div>
                            <h1 className="text-4xl font-bold text-white tracking-tight mb-2">
                                {data.vehicle.brand} {data.vehicle.model}
                            </h1>
                            <div className="flex items-center gap-3 text-sm text-slate-400">
                                <span className="px-2 py-0.5 bg-slate-800 rounded border border-slate-700">Chassis: {data.vehicle.generation || 'N/A'}</span>
                                <span className="w-1 h-1 bg-slate-600 rounded-full"></span>
                                <span>Platform: MQB Evo (Ref)</span>
                            </div>
                        </div>

                        <div className="grid grid-cols-4 gap-3">
                            {[
                                { icon: Timer, label: '0-100', value: quickStats?.accel, color: 'text-blue-400' },
                                { icon: Zap, label: 'Power', value: quickStats?.power, color: 'text-orange-400' },
                                { icon: Gauge, label: 'Torque', value: quickStats?.torque, color: 'text-emerald-400' },
                                { icon: Weight, label: 'Weight', value: quickStats?.weight, color: 'text-purple-400' },
                            ].map((stat, i) => (
                                <div key={i} className="bg-black/40 backdrop-blur-md rounded-xl p-3 border border-white/5 flex flex-col items-center justify-center text-center gap-2 group hover:bg-black/60 transition-all duration-300">
                                    <stat.icon size={20} className={stat.color} />
                                    <div>
                                        <div className="text-[10px] text-slate-500 uppercase font-semibold">{stat.label}</div>
                                        <div className="text-sm font-bold text-slate-200">{stat.value}</div>
                                    </div>
                                </div>
                            ))}
                        </div>

                        <div className="space-y-3">
                            <div className="flex items-center justify-between text-xs text-slate-400 uppercase font-semibold tracking-wider">
                                Transmission
                            </div>

                            {(() => {
                                const matchingVariant = currentTrim ? detectTransmissionVariant(data.vehicle.model, currentTrim, data.vehicle.brand) : null;

                                if (matchingVariant) {
                                    return (
                                        <button
                                            onClick={() => setExpandedVariantCode(matchingVariant.code)}
                                            className="w-full relative overflow-hidden group bg-gradient-to-br from-slate-900 to-black border border-white/10 rounded-xl hover:border-primary/50 transition-all duration-300 text-left"
                                        >
                                            <div className={`absolute inset-0 bg-gradient-to-br ${matchingVariant.color} opacity-10 group-hover:opacity-20 transition-opacity`} />

                                            <div className="relative p-4 flex items-center justify-between">
                                                <div className="flex items-center gap-4">
                                                    <div className={`w-10 h-10 rounded-lg flex items-center justify-center bg-gradient-to-br ${matchingVariant.color} shadow-lg ring-1 ring-white/10`}>
                                                        <Settings className="w-5 h-5 text-white" />
                                                    </div>
                                                    <div>
                                                        <div className="text-white font-bold text-lg leading-none mb-1">{matchingVariant.code}</div>
                                                        <div className="text-slate-400 text-xs font-medium">{matchingVariant.name}</div>
                                                    </div>
                                                </div>

                                                <div className="flex items-center gap-2 text-xs font-bold text-primary bg-primary/10 px-3 py-1.5 rounded-full border border-primary/20 group-hover:bg-primary/20 transition-colors">
                                                    <span>INCELE</span>
                                                    <Settings size={14} className="group-hover:rotate-90 transition-transform duration-500" />
                                                </div>
                                            </div>
                                        </button>
                                    );
                                }

                                const txValue = (currentTrim as any)?.transmission_type || 'Automatic';
                                return (
                                    <div className="bg-black/40 backdrop-blur-md border border-white/10 rounded-xl p-4 flex items-center gap-4">
                                        <div className="w-10 h-10 rounded-lg bg-slate-800 flex items-center justify-center text-slate-400">
                                            <Settings size={20} />
                                        </div>
                                        <div>
                                            <div className="text-slate-200 font-bold">{txValue}</div>
                                            <div className="text-slate-500 text-xs">Transmission</div>
                                        </div>
                                    </div>
                                );
                            })()}
                        </div>

                        <div className="flex gap-3 pt-2">
                            <button
                                onClick={() => currentTrim && addToCompare(currentTrim, normalizedSpecs)}
                                className="flex-1 bg-primary/80 hover:bg-primary backdrop-blur-md border border-primary/50 text-white font-bold h-12 rounded-xl flex items-center justify-center gap-2 transition-all shadow-lg shadow-orange-500/10 active:scale-95"
                            >
                                <Plus size={20} />
                                Compare Vehicle
                            </button>
                            <button className="flex-1 bg-white/5 hover:bg-white/10 backdrop-blur-md border border-white/10 hover:border-white/20 text-slate-300 font-semibold h-12 rounded-xl flex items-center justify-center gap-2 transition-all">
                                <Download size={20} />
                                Tech Sheet
                            </button>
                        </div>
                    </div>
                </div>

                <div className="space-y-6">
                    <h2 className="text-2xl font-bold text-white flex items-center gap-3">
                        <Info size={24} className="text-blue-400" />
                        Technical Specifications
                    </h2>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        {Object.entries(groupedSpecs).map(([category, specs]) => {
                            const isOpen = openCategories[category];
                            return (
                                <div key={category} className="bg-black/40 backdrop-blur-md border border-white/10 rounded-2xl overflow-hidden h-fit transform transition-all hover:border-white/20">
                                    <button
                                        onClick={() => toggleCategory(category)}
                                        className="w-full flex items-center justify-between p-5 bg-transparent hover:bg-white/5 transition-colors"
                                    >
                                        <span className="font-bold text-lg text-slate-200">{category}</span>
                                        {isOpen ? <ChevronUp size={20} className="text-slate-400" /> : <ChevronDown size={20} className="text-slate-500" />}
                                    </button>

                                    {isOpen && (
                                        <div className="p-5 space-y-4">
                                            {specs.map((spec, i) => {
                                                const percentage = getSpecPercentage(spec.name, spec.value);
                                                return (
                                                    <div key={i} className="group">
                                                        <div className="flex justify-between items-baseline mb-1">
                                                            <span className="text-sm text-slate-400">{spec.name}</span>
                                                            <span className="text-sm font-medium text-slate-200">{spec.value}</span>
                                                        </div>
                                                        {percentage > 0 && (
                                                            <div className="h-1.5 w-full bg-white/10 rounded-full overflow-hidden">
                                                                <div
                                                                    className="h-full rounded-full transition-all duration-1000 ease-out bg-gradient-to-r from-primary to-orange-400"
                                                                    style={{ width: `${percentage}%` }}
                                                                ></div>
                                                            </div>
                                                        )}
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </div>
            </main>

            <AnimatePresence>
                {expandedVariantCode && (
                    <DSGVariantModal
                        key="dsg-modal"
                        variant={DSG_VARIANTS.find(v => v.code === expandedVariantCode)!}
                        onClose={() => setExpandedVariantCode(null)}
                    />
                )}
            </AnimatePresence>
        </div>
    );
}

export default VehicleDetails;
