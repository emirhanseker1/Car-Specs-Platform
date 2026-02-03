import { Outlet, useLocation } from 'react-router-dom';
import { AnimatePresence, motion } from 'framer-motion';
import TopNav from './TopNav';
import BackgroundSlider from './BackgroundSlider';

export default function AppShell() {
  const location = useLocation();

  // Pages that share the full-screen immersive background
  const transparentLayoutRoutes = ['/', '/search', '/guides', '/guides/transmission', '/guides/engine', '/compare', '/about'];
  const isTransparentLayout = transparentLayoutRoutes.includes(location.pathname) ||
    location.pathname.startsWith('/brand/') ||
    location.pathname.startsWith('/models/') ||
    location.pathname.startsWith('/generations/') ||
    location.pathname.startsWith('/trims/');


  return (
    <div className={`min-h-screen overflow-hidden selection:bg-primary/30 ${!isTransparentLayout ? 'bg-gradient-to-b from-slate-950 via-slate-950 to-background' : ''}`}>

      {/* Shared Background Slider for specific routes */}
      {isTransparentLayout && <BackgroundSlider />}

      {/* Persistent Header */}
      <header
        className={`${isTransparentLayout ? 'absolute' : 'fixed'} top-0 left-0 right-0 z-50 transition-all duration-300 ${isTransparentLayout
          ? 'bg-transparent py-6'
          : 'bg-slate-950/80 backdrop-blur-md py-4 border-b border-white/5 shadow-lg shadow-black/5'
          }`}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <TopNav />
        </div>
      </header>

      {/* Page Content Transition */}
      <AnimatePresence mode="wait">
        <motion.div
          key={location.pathname}
          initial={{ opacity: 0, y: 20, scale: 0.98 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          exit={{ opacity: 0, y: -20, scale: 0.98 }}
          transition={{ duration: 0.4, ease: [0.22, 1, 0.36, 1] }} // Custom ease for premium feel
          className="min-h-screen w-full relative z-10"
        >
          {isTransparentLayout ? (
            // Full Width Content (Home, Search, Guides)
            <main className="w-full">
              <Outlet />
            </main>
          ) : (
            // Standard Page Content - With Top Padding for Header
            <main className="pt-28 pb-10 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto">
              <div className="rounded-3xl bg-white/95 backdrop-blur border border-white/10 shadow-2xl shadow-black/20 overflow-hidden">
                <div className="p-6 sm:p-8 lg:p-10">
                  <Outlet />
                </div>
              </div>
            </main>
          )}
        </motion.div>
      </AnimatePresence>
    </div>
  );
}
