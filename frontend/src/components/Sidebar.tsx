import {
    LayoutGrid,
    Car,
    CalendarDays,
    Wallet,
    ShoppingBag,
    Wrench,
    MessageSquare,
    Settings,
    HelpCircle,
    LogOut
} from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';

export default function Sidebar() {
    const location = useLocation();

    const menuItems = [
        { name: 'Dashboard', icon: LayoutGrid, path: '/' },
        { name: 'Assets', icon: Car, path: '/assets' },
        { name: 'Booking', icon: CalendarDays, path: '/booking' },
        { name: 'Sell Cars', icon: ShoppingBag, path: '/sell' },
        { name: 'Buy Cars', icon: Wallet, path: '/buy' },
        { name: 'Services', icon: Wrench, path: '/services' },
        { name: 'Calendar', icon: CalendarDays, path: '/calendar' },
        { name: 'Messages', icon: MessageSquare, path: '/messages' },
    ];

    const bottomItems = [
        { name: 'Settings', icon: Settings, path: '/settings' },
        { name: 'Help', icon: HelpCircle, path: '/help' },
    ];

    return (
        <div className="w-64 h-screen bg-white border-r border-border flex flex-col fixed left-0 top-0">
            <div className="p-6 flex items-center gap-3">
                <div className="bg-primary p-2 rounded-lg">
                    <Car className="text-white w-6 h-6" />
                </div>
                <h1 className="text-2xl font-bold text-text-main">Motive</h1>
            </div>

            <nav className="flex-1 px-4 py-4 space-y-1 overflow-y-auto">
                <div className="mb-2 px-4 text-xs font-semibold text-text-muted uppercase tracking-wider">
                    Menu
                </div>
                {menuItems.map((item) => {
                    const isActive = location.pathname === item.path;
                    return (
                        <Link
                            key={item.name}
                            to={item.path}
                            className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${isActive
                                    ? 'bg-background text-text-main font-medium'
                                    : 'text-text-muted hover:bg-background hover:text-text-main'
                                }`}
                        >
                            <item.icon className={`w-5 h-5 ${isActive ? 'text-primary' : ''}`} />
                            <span>{item.name}</span>
                        </Link>
                    );
                })}
            </nav>

            <div className="p-4 border-t border-border mt-auto">
                {bottomItems.map((item) => (
                    <Link
                        key={item.name}
                        to={item.path}
                        className="flex items-center gap-3 px-4 py-3 rounded-lg text-text-muted hover:bg-background hover:text-text-main transition-colors"
                    >
                        <item.icon className="w-5 h-5" />
                        <span>{item.name}</span>
                    </Link>
                ))}
                <button className="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-text-muted hover:bg-red-50 hover:text-red-500 transition-colors mt-2">
                    <LogOut className="w-5 h-5" />
                    <span>Log Out</span>
                </button>
            </div>
        </div>
    );
}
