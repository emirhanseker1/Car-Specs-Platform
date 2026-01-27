import { Search, Bell, Settings, User } from 'lucide-react';

export default function Header() {
    return (
        <header className="h-20 bg-white border-b border-border flex items-center justify-between px-8 sticky top-0 z-10">
            <div className="flex-1 max-w-xl">
                <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted w-5 h-5" />
                    <input
                        type="text"
                        placeholder="Search or type"
                        className="w-full pl-10 pr-4 py-2.5 bg-background rounded-lg text-sm text-text-main focus:outline-none focus:ring-2 focus:ring-primary/20 transition-shadow"
                    />
                </div>
            </div>

            <div className="flex items-center gap-6">
                <div className="flex items-center gap-4 text-text-muted">
                    <button className="p-2 hover:bg-background rounded-full transition-colors relative">
                        <Bell className="w-5 h-5" />
                        <span className="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full border-2 border-white"></span>
                    </button>
                    <button className="p-2 hover:bg-background rounded-full transition-colors">
                        <Settings className="w-5 h-5" />
                    </button>
                </div>

                <div className="flex items-center gap-3 pl-6 border-l border-border">
                    <div className="w-10 h-10 bg-gray-200 rounded-full overflow-hidden">
                        {/* Placeholder for user avatar */}
                        <div className="w-full h-full flex items-center justify-center bg-primary/10 text-primary">
                            <User className="w-6 h-6" />
                        </div>
                    </div>
                    <div className="hidden md:block">
                        <p className="text-sm font-medium text-text-main">Admin</p>
                        <p className="text-xs text-text-muted">Admin</p>
                    </div>
                </div>
            </div>
        </header>
    );
}
