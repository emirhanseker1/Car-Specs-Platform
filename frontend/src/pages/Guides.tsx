import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, BookOpen, Wrench, AlertTriangle, Layers } from 'lucide-react';
import { GUIDES } from '../data/guidesData';
import type { GuideArticle } from '../data/guidesData';

// Map IDs to Icons
const ICONS: Record<string, any> = {
  'transmission': BookOpen,
  'engine': Wrench,
  'chronicles': AlertTriangle,
  'platforms': Layers,
};

export default function Guides() {
  const navigate = useNavigate();
  const [selectedGuide, setSelectedGuide] = useState<GuideArticle | null>(null);

  if (selectedGuide) {
    return (
      <div className="relative min-h-screen font-sans text-slate-200">
        {/* Background Image & Overlay (Consistent with app) */}
        <div className="fixed inset-0 z-0">
          <div
            className="absolute inset-0 bg-cover bg-center"
            style={{ backgroundImage: 'url(/hero-2.jpg)' }}
          />
          <div className="absolute inset-0 bg-slate-900/90 backdrop-blur-sm" />
        </div>

        <div className="relative z-10 max-w-7xl mx-auto px-4 py-8 pt-32 space-y-6">
          <button
            onClick={() => setSelectedGuide(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors mb-4 bg-white/5 hover:bg-white/10 px-4 py-2 rounded-full border border-white/5 w-fit"
          >
            <ArrowLeft className="w-4 h-4" />
            Rehberlere Dön
          </button>

          <article className="prose prose-invert max-w-none">
            <h1 className="text-4xl font-bold text-white mb-4">{selectedGuide.title}</h1>
            <p className="text-xl text-slate-300 mb-8 leading-relaxed border-l-4 border-primary pl-6">
              {selectedGuide.summary}
            </p>

            <div className="space-y-12">
              {selectedGuide.sections.map((section, idx) => (
                <div key={idx} className="bg-[#1e293b]/60 backdrop-blur-md rounded-3xl p-8 border border-white/10 shadow-lg">
                  <h2 className="text-2xl font-semibold text-white mb-4 flex items-center gap-3">
                    <span className="flex items-center justify-center w-10 h-10 rounded-full bg-primary/20 text-primary text-base font-bold">
                      {idx + 1}
                    </span>
                    {section.title}
                  </h2>
                  <p className="text-slate-300 leading-relaxed text-lg">
                    {section.content}
                  </p>
                  {section.imagePlaceholder && (
                    <div className="mt-8 h-64 bg-black/40 rounded-2xl flex items-center justify-center border-2 border-dashed border-white/10">
                      <div className="text-center">
                        <div className="text-5xl mb-3 opacity-50">📷</div>
                        <span className="text-sm text-slate-500 uppercase tracking-wider font-semibold">
                          {section.imagePlaceholder.replace(/_/g, ' ')}
                        </span>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </article>
        </div>
      </div>
    );
  }

  return (
    <div className="relative min-h-screen font-sans">
      {/* Background Image & Overlay */}
      <div className="fixed inset-0 z-0">
        <div
          className="absolute inset-0 bg-cover bg-center"
          style={{ backgroundImage: 'url(/hero-2.jpg)' }}
        />
        <div className="absolute inset-0 bg-slate-900/90 backdrop-blur-sm" />
      </div>

      <div className="relative z-10 max-w-7xl mx-auto px-4 py-8 pt-32 space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl sm:text-5xl font-bold tracking-tight text-white drop-shadow-lg">Rehberler</h1>
          <p className="text-slate-300 text-xl max-w-2xl font-light">
            Otomobil dünyasının teknik derinliklerine inin. Şanzıman tipleri, motor teknolojileri ve daha fazlasını uzmanından öğrenin.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 pb-12">
          {GUIDES.map((guide) => {
            const Icon = ICONS[guide.id] || BookOpen;
            const isTransmission = guide.id === 'transmission';
            const isEngine = guide.id === 'engine';

            return (
              <div
                key={guide.id}
                onClick={() => {
                  if (isTransmission) {
                    navigate('/guides/transmission');
                  } else if (isEngine) {
                    navigate('/guides/engine');
                  } else {
                    setSelectedGuide(guide);
                  }
                }}
                className="group cursor-pointer rounded-3xl border border-white/10 bg-[#1e293b]/60 backdrop-blur-md p-8 hover:bg-[#1e293b]/80 hover:shadow-2xl hover:border-primary/50 transition-all duration-500 relative overflow-hidden"
              >
                {/* Background Decor */}
                <div className="absolute -right-10 -top-10 text-white/5 opacity-50 group-hover:scale-110 transition-transform duration-700">
                  <Icon size={200} />
                </div>

                <div className="relative z-10 flex flex-col h-full">
                  <div className="w-16 h-16 rounded-2xl bg-white/10 text-primary flex items-center justify-center mb-6 group-hover:bg-primary group-hover:text-white transition-all duration-300 shadow-lg shadow-black/20">
                    <Icon className="w-8 h-8" />
                  </div>

                  <h3 className="text-2xl font-bold text-white mb-3 group-hover:text-primary transition-colors">
                    {guide.title}
                  </h3>

                  <p className="text-base text-slate-300 leading-relaxed mb-8 flex-grow opacity-90">
                    {guide.summary}
                  </p>

                  <div className="flex items-center text-sm font-bold text-primary tracking-wide uppercase">
                    Okumaya Başla
                    <ArrowLeft className="w-4 h-4 ml-2 rotate-180 group-hover:translate-x-2 transition-transform duration-300" />
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
