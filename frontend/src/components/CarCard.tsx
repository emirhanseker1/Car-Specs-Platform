import { Fuel, Gauge, Users, Settings2, Heart } from 'lucide-react';

interface CarProps {
    name: string;
    type: string;
    image: string;
    price: number;
    rating: number;
    reviewCount: number;
    specs: {
        transmission: string;
        range: string;
        class: string;
        seats: number;
    };
}

export default function CarCard({ car }: { car: CarProps }) {
    return (
        <div className="bg-white rounded-2xl p-6 shadow-card hover:shadow-lg transition-shadow border border-transparent hover:border-primary/10 group">
            <div className="flex justify-between items-start mb-4">
                <div>
                    <h3 className="text-lg font-bold text-text-main">{car.name}</h3>
                    <p className="text-sm text-text-muted">{car.type}</p>
                </div>
                <button className="text-text-muted hover:text-red-500 transition-colors">
                    <Heart className="w-5 h-5" />
                </button>
            </div>

            <div className="relative h-40 mb-6 flex items-center justify-center">
                {/* Gradient Background Effect */}
                <div className="absolute inset-x-8 inset-y-4 bg-gradient-to-t from-primary/20 to-transparent rounded-full blur-xl opacity-0 group-hover:opacity-100 transition-opacity" />
                <img
                    src={car.image}
                    alt={car.name}
                    className="w-full h-full object-contain relative z-10"
                />
            </div>

            <div className="grid grid-cols-2 gap-4 mb-6">
                <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Settings2 className="w-4 h-4 text-primary" />
                    <span>{car.specs.transmission}</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Gauge className="w-4 h-4 text-primary" />
                    <span>{car.specs.range}</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Fuel className="w-4 h-4 text-primary" />
                    <span>{car.specs.class}</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-text-muted">
                    <Users className="w-4 h-4 text-primary" />
                    <span>{car.specs.seats} Seat</span>
                </div>
            </div>

            <div className="flex items-center justify-between pt-4 border-t border-border/50">
                <div>
                    <span className="text-xl font-bold text-text-main">${car.price.toFixed(2)}/</span>
                    <span className="text-sm text-text-muted">Day</span>
                </div>
                <button className="bg-primary hover:bg-primary-hover text-white px-6 py-2.5 rounded-lg font-medium transition-colors shadow-lg shadow-primary/25">
                    Rent Now
                </button>
            </div>
        </div>
    );
}
