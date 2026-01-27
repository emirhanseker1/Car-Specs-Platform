import { NavLink } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Car, BookOpen, GitCompare, Info, Search } from 'lucide-react';

const navItems = [
  { to: '/', label: 'Keşfet', icon: Car, end: true },
  { to: '/search', label: 'Ara', icon: Search },
  { to: '/compare', label: 'Karşılaştır', icon: GitCompare },
  { to: '/guides', label: 'Rehberler', icon: BookOpen },
  { to: '/about', label: 'Hakkında', icon: Info },
];

export default function TopNav() {
  return (
    <div className="flex items-center justify-between gap-6">
      <div className="flex items-center gap-3">
        <div className="h-10 w-10 rounded-2xl bg-white/10 border border-white/10 flex items-center justify-center">
          <Car className="h-5 w-5 text-white" />
        </div>
        <div>
          <div className="text-white font-bold leading-tight">Car Specs</div>
          <div className="text-xs text-white/60">Türkiye araç teknik bilgi platformu</div>
        </div>
      </div>

      <nav className="hidden md:flex items-center gap-2 rounded-2xl bg-white/10 border border-white/10 p-2">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end as any}
            className="relative px-4 py-2 text-sm font-semibold transition-colors"
          >
            {({ isActive }) => (
              <>
                {isActive && (
                  <motion.span
                    layoutId="nav-pill"
                    className="absolute inset-0 bg-white rounded-xl shadow-sm"
                    transition={{ type: "spring", stiffness: 300, damping: 30 }}
                  />
                )}
                <div className={`relative z-10 flex items-center gap-2 ${isActive ? 'text-slate-900' : 'text-white/80 hover:text-white'}`}>
                  <item.icon className="h-4 w-4" />
                  <span>{item.label}</span>
                </div>
              </>
            )}
          </NavLink>
        ))}
      </nav>

      <div className="md:hidden rounded-2xl bg-white/10 border border-white/10 px-3 py-2 text-xs font-semibold text-white/80">
        Menü
      </div>
    </div>
  );
}
