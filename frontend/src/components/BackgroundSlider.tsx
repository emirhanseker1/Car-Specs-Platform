import { useEffect, useState } from 'react';

// Hero slider images - New High Quality Local Images
const HERO_IMAGES = [
    '/hero/hero-1.jpg',
    '/hero/hero-2.jpg',
    '/hero/hero-3.jpg',
    '/hero/hero-4.jpg',
    '/hero/hero-5.jpg',
];

export default function BackgroundSlider() {
    const [currentSlide, setCurrentSlide] = useState(0);

    // Hero slider auto-advance
    useEffect(() => {
        const interval = setInterval(() => {
            setCurrentSlide((prev) => (prev + 1) % HERO_IMAGES.length);
        }, 5000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="fixed inset-0 z-0 overflow-hidden pointer-events-none">
            {HERO_IMAGES.map((image, index) => (
                <div
                    key={index}
                    className="absolute inset-0 transition-opacity duration-1000 ease-in-out bg-slate-900"
                    style={{
                        opacity: currentSlide === index ? 1 : 0,
                        backgroundImage: `url("${image}")`,
                        backgroundSize: 'cover',
                        backgroundPosition: 'center',
                    }}
                >
                    {/* Minimal dark overlay to ensure text readability without killing clearity */}
                    {/* Reduced opacity from 0.4 to 0.3 to let more "quality" through */}
                    <div className="absolute inset-0 bg-black/30" />
                </div>
            ))}
        </div>
    );
}
