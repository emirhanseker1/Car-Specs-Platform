import { FC, SVGProps } from 'react';

interface CarIconProps extends SVGProps<SVGSVGElement> {
    bodyStyle: string;
}

const BaseIcon: FC<SVGProps<SVGSVGElement>> = (props) => (
    <svg viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg" {...props}>
        {props.children}
    </svg>
);

// Paths (Simplified silhouettes)
const PATHS: Record<string, string> = {
    // Sedan (Generic Car)
    'sedan': 'M18.92 6.01C18.72 5.42 18.16 5 17.5 5h-11c-.66 0-1.21.42-1.42 1.01L3 12v8c0 .55.45 1 1 1h1c.55 0 1-.45 1-1v-1h12v1c0 .55.45 1 1 1h1c.55 0 1-.45 1-1v-8l-2.08-5.99zM6.5 16c-.83 0-1.5-.67-1.5-1.5S5.67 13 6.5 13s1.5.67 1.5 1.5S7.33 16 6.5 16zm11 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zM5 11l1.5-4.5h11L19 11H5z',

    // SUV (Taller, boxy)
    'suv': 'M20 8h-3V4c0-.55-.45-1-1-1H8c-.55 0-1 .45-1 1v4H4c-1.1 0-2 .9-2 2v10h2c0 1.66 1.34 3 3 3s3-1.34 3-3h6c0 1.66 1.34 3 3 3s3-1.34 3-3h2v-8c0-1.1-.9-2-2-2zm-3-2v4H9V6h8zm-9.5 13c-.83 0-1.5-.67-1.5-1.5S6.67 16 7.5 16s1.5.67 1.5 1.5S8.33 19 7.5 19zm9 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z',

    // Pickup (Truck bed)
    'pickup': 'M18 7c-.55 0-1 .45-1 1v2H7V8c0-.55-.45-1-1-1H2v9h2c0 1.66 1.34 3 3 3s3-1.34 3-3h6c0 1.66 1.34 3 3 3s3-1.34 3-3h2v-9h-6zM5.5 16c-.83 0-1.5-.67-1.5-1.5S4.67 13 5.5 13s1.5.67 1.5 1.5S6.33 16 5.5 16zm13 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z',

    // Hatchback (Short rear - simplified)
    'hatchback': 'M19 8h-3V5H8v3H4v11h2c0 1.66 1.34 3 3 3s3-1.34 3-3h6c0 1.66 1.34 3 3 3s3-1.34 3-3h2v-9c0-1.1-.9-2-2-2zM7.5 17c-.83 0-1.5-.67-1.5-1.5S6.67 14 7.5 14s1.5.67 1.5 1.5S8.33 17 7.5 17zm11 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zM17 10h-6V7h6v3z',

    // Convertibte / Coupe (Low profile)
    'coupe': 'M19 9h-3V6c0-.55-.45-1-1-1H9c-.55 0-1 .45-1 1v3H4c-1.1 0-2 .9-2 2v8h2c0 1.66 1.34 3 3 3s3-1.34 3-3h6c0 1.66 1.34 3 3 3s3-1.34 3-3h2v-9c0-1.1-.9-2-2-2zm-3-2v3H8V7h8zM7.5 17c-.83 0-1.5-.67-1.5-1.5S6.67 14 7.5 14s1.5.67 1.5 1.5S8.33 17 7.5 17zm11 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z',
};

// Aliases
const MAP: Record<string, string> = {
    'sedan': 'sedan',
    'saloon': 'sedan',
    'suv': 'suv',
    'crossover': 'suv',
    'pickup': 'pickup',
    'single cab': 'pickup',
    'crew cab': 'pickup',
    'station wagon': 'hatchback', // Use hatchback/wagon approximate
    'wagon': 'hatchback',
    'hatchback': 'hatchback',
    'coupe': 'coupe',
    'convertible': 'coupe',
    'cabrio': 'coupe',
    'roadster': 'coupe',
};

export const CarBodyIcon: FC<CarIconProps> = ({ bodyStyle, ...props }) => {
    const key = (bodyStyle || '').toLowerCase().trim();
    let type = 'sedan'; // Default

    // Find matching type
    for (const [k, v] of Object.entries(MAP)) {
        if (key.includes(k)) {
            type = v;
            break;
        }
    }

    const path = PATHS[type] || PATHS['sedan'];

    return (
        <BaseIcon {...props}>
            <path d={path} />
        </BaseIcon>
    );
};
