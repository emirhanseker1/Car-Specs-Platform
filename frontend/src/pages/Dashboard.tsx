import { Search } from 'lucide-react';
import HeroSection from '../components/HeroSection';
import CarCard from '../components/CarCard';

const MOCK_CARS = [
    {
        name: "BMW 3 Series Sedan",
        type: "The 3 Sedan",
        image: "https://placehold.co/400x200/png?text=BMW+3+Series", // Placeholder
        price: 87.44,
        rating: 4.9,
        reviewCount: 100,
        specs: { transmission: "Automatic", range: "250km", class: "Premium", seats: 4 }
    },
    {
        name: "Mercedes-Benz C-Class",
        type: "The C Sedan",
        image: "https://placehold.co/400x200/png?text=Mercedes+C-Class",
        price: 60.20,
        rating: 4.9,
        reviewCount: 100,
        specs: { transmission: "Automatic", range: "250km", class: "Premium", seats: 4 }
    },
    {
        name: "Audi Q5 SUV",
        type: "The Q5",
        image: "https://placehold.co/400x200/png?text=Audi+Q5",
        price: 97.82,
        rating: 4.9,
        reviewCount: 100,
        specs: { transmission: "Automatic", range: "250km", class: "Premium", seats: 4 }
    },
    {
        name: "Toyota Camry Sedan",
        type: "The Camry",
        image: "https://placehold.co/400x200/png?text=Toyota+Camry",
        price: 45.50,
        rating: 4.8,
        reviewCount: 85,
        specs: { transmission: "Automatic", range: "Unlimited", class: "Mid-size", seats: 5 }
    }
];

export default function Dashboard() {
    return (
        <div className="space-y-8 pb-8">
            <HeroSection />

            <div className="grid grid-cols-12 gap-8">
                {/* Filters Sidebar */}
                <aside className="col-span-12 lg:col-span-3 space-y-8">
                    <div className="bg-white rounded-2xl p-6 shadow-sm space-y-6">
                        <div className="flex items-center justify-between">
                            <h3 className="font-bold text-text-main">Filter Plans</h3>
                            <button className="text-primary text-sm font-medium hover:underline">Reset</button>
                        </div>

                        <div className="relative">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4" />
                            <input type="text" placeholder="Search" className="w-full pl-10 pr-4 py-2 bg-background rounded-lg text-sm/6 focus:outline-none focus:ring-2 focus:ring-primary/20" />
                        </div>

                        <div className="space-y-4">
                            <h4 className="text-sm font-bold text-text-main uppercase tracking-wider">Price & Budget</h4>
                            <div className="h-32 bg-background rounded-lg flex items-center justify-center text-xs text-text-muted">
                                Price Range Chart Placeholder
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="text-xs text-text-muted block mb-1">Min Price</label>
                                    <div className="px-3 py-2 bg-background rounded-lg text-sm font-medium">$120</div>
                                </div>
                                <div>
                                    <label className="text-xs text-text-muted block mb-1">Max Price</label>
                                    <div className="px-3 py-2 bg-background rounded-lg text-sm font-medium">$420</div>
                                </div>
                            </div>
                        </div>

                        <div className="space-y-4">
                            <h4 className="text-sm font-bold text-text-main uppercase tracking-wider">Brand & Model</h4>
                            <div className="space-y-3">
                                {['BMW', 'Mercedes', 'Audi', 'Toyota', 'Honda', 'Nissan'].map(brand => (
                                    <label key={brand} className="flex items-center gap-3 cursor-pointer group">
                                        <input type="checkbox" className="w-4 h-4 rounded border-gray-300 text-primary focus:ring-primary" />
                                        <span className="text-sm text-text-muted group-hover:text-text-main transition-colors">{brand}</span>
                                        <span className="ml-auto text-xs text-text-muted">(12)</span>
                                    </label>
                                ))}
                            </div>
                        </div>
                    </div>
                </aside>

                {/* Car Grid */}
                <div className="col-span-12 lg:col-span-9">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-3 gap-6">
                        {MOCK_CARS.map((car, i) => (
                            <CarCard key={i} car={car} />
                        ))}
                    </div>
                    <div className="mt-8 flex justify-center">
                        <button className="px-6 py-2 bg-white border border-border rounded-lg text-sm font-medium hover:bg-background transition-colors">
                            Show More Cars
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
