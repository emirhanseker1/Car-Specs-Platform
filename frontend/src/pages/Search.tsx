import { useEffect, useMemo, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';

// New API response types matching backend
interface Brand {
    id: number;
    name: string;
    country?: string;
}

interface Model {
    id: number;
    brand_id: number;
    name: string;
    body_style?: string;
    brand?: Brand;
}

interface Trim {
    id: number;
    model_id: number;
    name: string;
    year: number;
    generation?: string;
    fuel_type?: string;
    power_hp?: number;
    torque_nm?: number;
    acceleration_0_100?: number;
    transmission_type?: string;
    drivetrain?: string;
    displacement_cc?: number;
    top_speed_kmh?: number;
    fuel_consumption_combined?: number;
    image_url?: string;
    model?: Model;
}

interface Facets {
    fuel_types?: string[];
    transmissions?: string[];
    min_hp?: number;
    max_hp?: number;
    min_year?: number;
    max_year?: number;
}

interface SearchResponse {
    results: Trim[];
    facets: Facets;
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || 'http://localhost:8080';

function Search() {
    const [searchParams, setSearchParams] = useSearchParams();

    const initialBrand = (searchParams.get('brand') || '').trim();
    const initialModel = (searchParams.get('model') || '').trim();
    const initialFuel = (searchParams.get('fuel_type') || '').trim();
    const initialTransmission = (searchParams.get('transmission') || '').trim();
    const initialYear = (searchParams.get('year') || '').trim();

    const [draftBrand, setDraftBrand] = useState(initialBrand);
    const [draftModel, setDraftModel] = useState(initialModel);
    const [draftFuel, setDraftFuel] = useState(initialFuel);
    const [draftTransmission, setDraftTransmission] = useState(initialTransmission);
    const [draftYear, setDraftYear] = useState(initialYear);

    const [data, setData] = useState<SearchResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const applied = useMemo(() => {
        return {
            brand: initialBrand,
            model: initialModel,
            fuel: initialFuel,
            transmission: initialTransmission,
            year: initialYear,
        };
    }, [initialBrand, initialModel, initialFuel, initialTransmission, initialYear]);

    useEffect(() => {
        const qs = new URLSearchParams();
        if (applied.brand) qs.set('brand', applied.brand);
        if (applied.model) qs.set('model', applied.model);
        if (applied.fuel) qs.set('fuel_type', applied.fuel);
        if (applied.transmission) qs.set('transmission', applied.transmission);
        if (applied.year) qs.set('year', applied.year);

        setLoading(true);
        setError(null);
        fetch(`${API_BASE_URL}/api/search?${qs.toString()}`)
            .then(async (res) => {
                const contentType = res.headers.get('content-type') || '';
                if (!res.ok) {
                    const text = await res.text().catch(() => '');
                    throw new Error(text || `Request failed (${res.status})`);
                }
                if (!contentType.includes('application/json')) {
                    const text = await res.text();
                    throw new Error(`Expected JSON, got: ${text.substring(0, 100)}`);
                }
                return res.json();
            })
            .then((json: SearchResponse) => {
                setData(json);
                setLoading(false);
            })
            .catch((err) => {
                setError(err.message);
                setLoading(false);
            });
    }, [applied]);

    const apply = () => {
        const newParams = new URLSearchParams();
        if (draftBrand) newParams.set('brand', draftBrand);
        if (draftModel) newParams.set('model', draftModel);
        if (draftFuel) newParams.set('fuel_type', draftFuel);
        if (draftTransmission) newParams.set('transmission', draftTransmission);
        if (draftYear) newParams.set('year', draftYear);
        setSearchParams(newParams);
    };

    const reset = () => {
        setDraftBrand('');
        setDraftModel('');
        setDraftFuel('');
        setDraftTransmission('');
        setDraftYear('');
        setSearchParams(new URLSearchParams());
    };

    const results = data?.results || [];
    const facets = data?.facets || {};

    return (
        <div className="min-h-screen">
            <div className="max-w-7xl mx-auto px-4 py-8 pt-32">
                {/* Header */}
                <div className="mb-8">
                    <h1 className="text-4xl font-bold text-white mb-2">Vehicle Search</h1>
                    <p className="text-gray-300">Find detailed specifications for your vehicle</p>
                </div>

                {/* Filters */}
                <div className="bg-black/40 backdrop-blur-md border border-white/10 rounded-2xl shadow-lg p-6 mb-8">
                    <h2 className="text-xl font-semibold text-white mb-4">Filters</h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">Brand</label>
                            <input
                                type="text"
                                value={draftBrand}
                                onChange={(e) => setDraftBrand(e.target.value)}
                                placeholder="e.g., BMW"
                                className="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white placeholder:text-gray-500 focus:ring-2 focus:ring-primary/50 focus:border-transparent outline-none transition-all"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">Model</label>
                            <input
                                type="text"
                                value={draftModel}
                                onChange={(e) => setDraftModel(e.target.value)}
                                placeholder="e.g., 3 Series"
                                className="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white placeholder:text-gray-500 focus:ring-2 focus:ring-primary/50 focus:border-transparent outline-none transition-all"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">Fuel Type</label>
                            <select
                                value={draftFuel}
                                onChange={(e) => setDraftFuel(e.target.value)}
                                className="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white focus:ring-2 focus:ring-primary/50 focus:border-transparent outline-none transition-all appearance-none"
                            >
                                <option value="" className="bg-slate-900 text-white">All</option>
                                <option value="Benzin" className="bg-slate-900 text-white">Benzin</option>
                                <option value="Dizel" className="bg-slate-900 text-white">Dizel</option>
                                <option value="Hybrid" className="bg-slate-900 text-white">Hybrid</option>
                                <option value="Elektrik" className="bg-slate-900 text-white">Elektrik</option>
                            </select>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">Transmission</label>
                            <select
                                value={draftTransmission}
                                onChange={(e) => setDraftTransmission(e.target.value)}
                                className="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white focus:ring-2 focus:ring-primary/50 focus:border-transparent outline-none transition-all appearance-none"
                            >
                                <option value="" className="bg-slate-900 text-white">All</option>
                                <option value="Manual" className="bg-slate-900 text-white">Manual</option>
                                <option value="Automatic" className="bg-slate-900 text-white">Automatic</option>
                                <option value="CVT" className="bg-slate-900 text-white">CVT</option>
                            </select>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">Year</label>
                            <input
                                type="number"
                                value={draftYear}
                                onChange={(e) => setDraftYear(e.target.value)}
                                placeholder="2023"
                                className="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-lg text-white placeholder:text-gray-500 focus:ring-2 focus:ring-primary/50 focus:border-transparent outline-none transition-all"
                            />
                        </div>
                    </div>
                    <div className="flex gap-3 mt-6">
                        <button
                            onClick={apply}
                            className="px-6 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg transition-colors font-medium shadow-lg shadow-primary/20"
                        >
                            Apply Filters
                        </button>
                        <button
                            onClick={reset}
                            className="px-6 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg transition-colors font-medium border border-white/5"
                        >
                            Reset
                        </button>
                    </div>
                </div>

                {/* Results */}
                {loading && (
                    <div className="text-center py-12">
                        <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                        <p className="mt-4 text-gray-300">Loading vehicles...</p>
                    </div>
                )}

                {error && (
                    <div className="bg-red-500/10 backdrop-blur border border-red-500/20 rounded-lg p-4 mb-6">
                        <p className="text-red-200">Error: {error}</p>
                    </div>
                )}

                {!loading && !error && (
                    <>
                        <div className="mb-6">
                            <p className="text-gray-300">
                                Found <span className="font-semibold text-white">{results.length}</span> vehicle{results.length !== 1 ? 's' : ''}
                            </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                            {results.map((trim) => (
                                <div key={trim.id} className="bg-black/40 backdrop-blur-md border border-white/10 rounded-xl shadow-lg hover:shadow-xl hover:border-primary/50 transition-all overflow-hidden group">
                                    {/* Image Container with Fixed Aspect Ratio */}
                                    <div className="relative h-48 bg-white/5 overflow-hidden">
                                        {trim.image_url ? (
                                            <img
                                                src={trim.image_url}
                                                alt={`${trim.model?.brand?.name || ''} ${trim.model?.name || ''} ${trim.name}`}
                                                className="w-full h-full object-cover object-center group-hover:scale-110 transition-transform duration-500"
                                                onError={(e) => {
                                                    e.currentTarget.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="400" height="300"%3E%3Crect fill="%231a1a1a" width="400" height="300"/%3E%3Ctext fill="%234b5563" font-family="sans-serif" font-size="18" x="50%25" y="50%25" text-anchor="middle" dominant-baseline="middle"%3ENo Image%3C/text%3E%3C/svg%3E';
                                                }}
                                            />
                                        ) : (
                                            <div className="w-full h-full flex items-center justify-center">
                                                <svg className="w-16 h-16 text-white/20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                                </svg>
                                            </div>
                                        )}
                                    </div>

                                    {/* Card Content */}
                                    <div className="p-6">
                                        <h3 className="text-xl font-bold text-white mb-1 group-hover:text-primary transition-colors">
                                            {trim.model?.brand?.name} {trim.model?.name}
                                        </h3>
                                        <p className="text-gray-400 mb-4">{trim.name}</p>

                                        <div className="space-y-2 text-sm">
                                            <div className="flex justify-between">
                                                <span className="text-gray-400">Year:</span>
                                                <span className="font-medium text-white">{trim.year}</span>
                                            </div>
                                            {trim.power_hp && (
                                                <div className="flex justify-between">
                                                    <span className="text-gray-400">Power:</span>
                                                    <span className="font-medium text-white">{trim.power_hp} HP</span>
                                                </div>
                                            )}
                                            {trim.fuel_type && (
                                                <div className="flex justify-between">
                                                    <span className="text-gray-400">Fuel:</span>
                                                    <span className="font-medium text-white">{trim.fuel_type}</span>
                                                </div>
                                            )}
                                            {trim.transmission_type && (
                                                <div className="flex justify-between">
                                                    <span className="text-gray-400">Transmission:</span>
                                                    <span className="font-medium text-white">{trim.transmission_type}</span>
                                                </div>
                                            )}
                                            {trim.acceleration_0_100 && (
                                                <div className="flex justify-between">
                                                    <span className="text-gray-400">0-100 km/h:</span>
                                                    <span className="font-medium text-white">{trim.acceleration_0_100}s</span>
                                                </div>
                                            )}
                                        </div>

                                        <Link
                                            to={`/vehicle/${trim.id}`}
                                            className="mt-4 block w-full text-center px-4 py-2 bg-white/10 text-white rounded-lg hover:bg-primary hover:text-white transition-all font-medium border border-white/5"
                                        >
                                            View Details
                                        </Link>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {results.length === 0 && (
                            <div className="text-center py-12 bg-white rounded-xl shadow-md">
                                <p className="text-slate-600 text-lg">No vehicles found matching your criteria.</p>
                                <p className="text-slate-500 mt-2">Try adjusting your filters.</p>
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
}

export default Search;
