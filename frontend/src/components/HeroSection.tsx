import { MapPin, Calendar, Clock } from 'lucide-react';

export default function HeroSection() {
    return (
        <div className="bg-white rounded-2xl p-6 lg:p-8 shadow-sm flex flex-col lg:flex-row gap-8 items-center">
            <div className="flex-1 space-y-6">
                <div className="space-y-4 max-w-lg">
                    <h2 className="text-4xl font-bold text-text-main leading-tight">
                        All In One Car Platform
                    </h2>
                    <p className="text-text-muted text-lg">
                        Renting a car gives your freedom, and we'll help you find the best car for you at a great price.
                    </p>
                </div>

                <div className="relative h-48 w-full max-w-md mt-4">
                    {/* Placeholder for Hero Cars Image */}
                    <div className="absolute inset-0 bg-gradient-to-r from-blue-100 to-blue-50 rounded-xl flex items-center justify-center text-blue-300 font-bold text-2xl">
                        Cars Banner Image
                    </div>
                </div>
            </div>

            <div className="flex-1 w-full bg-background/50 rounded-xl p-6 border border-border/50">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Pick-up */}
                    <div className="md:col-span-2 space-y-2">
                        <label className="text-sm font-semibold text-text-main flex items-center gap-2">
                            <span className="w-2 h-2 rounded-full bg-blue-500"></span>
                            Pick-up Location
                        </label>
                        <div className="relative">
                            <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4" />
                            <input type="text" placeholder="Enter City, Airport, or Address" className="w-full pl-10 pr-4 py-3 bg-white rounded-lg text-sm border-none shadow-sm focus:ring-2 focus:ring-primary/20" />
                        </div>
                    </div>

                    {/* Drop-off */}
                    <div className="md:col-span-2 space-y-2">
                        <label className="text-sm font-semibold text-text-main flex items-center gap-2">
                            <span className="w-2 h-2 rounded-full bg-primary"></span>
                            Drop-off Location
                        </label>
                        <div className="relative">
                            <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4" />
                            <input type="text" placeholder="Enter City, Airport, or Address" className="w-full pl-10 pr-4 py-3 bg-white rounded-lg text-sm border-none shadow-sm focus:ring-2 focus:ring-primary/20" />
                        </div>
                    </div>

                    {/* Date/Time Pick-up */}
                    <div className="space-y-2">
                        <label className="text-sm font-bold text-text-main">Pick-up Date</label>
                        <div className="relative">
                            <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4" />
                            <input type="date" className="w-full pl-10 pr-4 py-3 bg-white rounded-lg text-sm border-none shadow-sm focus:ring-2 focus:ring-primary/20" />
                        </div>
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-bold text-text-main">Pick-up Time</label>
                        <div className="relative">
                            <Clock className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4" />
                            <input type="time" className="w-full pl-10 pr-4 py-3 bg-white rounded-lg text-sm border-none shadow-sm focus:ring-2 focus:ring-primary/20" />
                        </div>
                    </div>
                    {/* Date/Time Drop-off - Simplified for brevity in layout, can duplicate above */}
                </div>
            </div>
        </div>
    );
}
