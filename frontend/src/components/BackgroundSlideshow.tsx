import { useState, useEffect } from 'react';

interface BackgroundSlideshowProps {
    images: string[];
    interval?: number;
    overlayOpacity?: number;
}

export default function BackgroundSlideshow({
    images,
    interval = 5000,
    overlayOpacity = 0.8
}: BackgroundSlideshowProps) {
    const [currentIndex, setCurrentIndex] = useState(0);

    useEffect(() => {
        if (images.length <= 1) return;

        const timer = setInterval(() => {
            setCurrentIndex((prev) => (prev + 1) % images.length);
        }, interval);

        return () => clearInterval(timer);
    }, [images.length, interval]);

    return (
        <div className="fixed inset-0 z-0 overflow-hidden pointer-events-none">
            {images.map((image, index) => (
                <div
                    key={image}
                    className={`absolute inset-0 bg-cover bg-center transition-opacity duration-1000 ease-in-out`}
                    style={{
                        backgroundImage: `url(${image})`,
                        opacity: index === currentIndex ? 1 : 0,
                        zIndex: index === currentIndex ? 1 : 0
                    }}
                />
            ))}
            <div
                className="absolute inset-0 bg-black backdrop-blur-[2px]"
                style={{ opacity: overlayOpacity, zIndex: 10 }}
            />
        </div>
    );
}
