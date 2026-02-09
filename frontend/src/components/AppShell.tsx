import { useLocation, useOutlet, Link } from 'react-router-dom';
import { AnimatePresence, motion } from 'framer-motion';
import TopNav from './TopNav';
import BackgroundSlider from './BackgroundSlider';

export default function AppShell() {
  const location = useLocation();
  const currentOutlet = useOutlet();

  // Pages that share the full-screen immersive background
  const transparentLayoutRoutes = ['/', '/search', '/guides', '/guides/transmission', '/guides/engine', '/compare', '/about', '/guides/transmission/dsg'];
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
        className={`app-header ${isTransparentLayout ? 'absolute' : 'fixed'} top-0 left-0 right-0 z-50 transition-all duration-300 ${isTransparentLayout
          ? 'bg-transparent py-6'
          : 'bg-slate-950/80 backdrop-blur-md py-4 border-b border-white/5 shadow-lg shadow-black/5'
          }`}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <TopNav />
        </div>
      </header>

      {/* Page Content Transition */}
      {location.pathname === '/' && (
        /* DSG Rehberi Floating Button - Moved to Shell to avoid animation shift */
        <Link
          to="/guides/transmission/dsg"
          className="fixed left-6 top-1/3 z-40 group"
        >
          <div className="flex items-center gap-3 pl-2 pr-4 py-2 bg-black/20 backdrop-blur-md rounded-full border border-white/10 hover:bg-white/5 transition-all duration-300 group-hover:scale-105">
            <div className="w-8 h-8 bg-white/10 rounded-full flex items-center justify-center group-hover:bg-white/20 transition-colors">
              <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div className="hidden md:block">
              <span className="font-medium text-white text-sm">DSG Rehberi</span>
            </div>
            <svg className="w-4 h-4 text-white/50 group-hover:text-white group-hover:translate-x-0.5 transition-all hidden md:block" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </Link>
      )}

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
              {currentOutlet}
            </main>
          ) : (
            // Standard Page Content - With Top Padding for Header
            <main className="pt-28 pb-10 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto">
              <div className="rounded-3xl bg-white/95 backdrop-blur border border-white/10 shadow-2xl shadow-black/20 overflow-hidden">
                <div className="p-6 sm:p-8 lg:p-10">
                  {currentOutlet}
                </div>
              </div>
            </main>
          )}
        </motion.div>
      </AnimatePresence>
    </div>
  );
}
