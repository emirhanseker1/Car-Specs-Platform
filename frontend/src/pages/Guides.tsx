import { useState } from 'react';
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
  const [selectedGuide, setSelectedGuide] = useState<GuideArticle | null>(null);

  if (selectedGuide) {
    return (
      <div className="space-y-6">
        <button
          onClick={() => setSelectedGuide(null)}
          className="flex items-center gap-2 text-sm text-text-muted hover:text-primary transition-colors mb-4"
        >
          <ArrowLeft className="w-4 h-4" />
          Rehberlere Dön
        </button>

        <article className="prose prose-slate max-w-none">
          <h1 className="text-3xl font-bold text-text-main mb-4">{selectedGuide.title}</h1>
          <p className="text-lg text-text-muted mb-8 leading-relaxed border-l-4 border-primary pl-4">
            {selectedGuide.summary}
          </p>

          <div className="space-y-12">
            {selectedGuide.sections.map((section, idx) => (
              <div key={idx} className="bg-white rounded-3xl p-8 border border-border shadow-sm">
                <h2 className="text-2xl font-semibold text-text-main mb-4 flex items-center gap-3">
                  <span className="flex items-center justify-center w-8 h-8 rounded-full bg-primary/10 text-primary text-sm font-bold">
                    {idx + 1}
                  </span>
                  {section.title}
                </h2>
                <p className="text-text-muted leading-relaxed">
                  {section.content}
                </p>
                {section.imagePlaceholder && (
                  <div className="mt-6 h-48 bg-slate-100 rounded-xl flex items-center justify-center border-2 border-dashed border-slate-200">
                    <div className="text-center">
                      <div className="text-4xl mb-2">📷</div>
                      <span className="text-xs text-text-muted uppercase tracking-wider font-semibold">
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
    );
  }

  return (
    <div className="space-y-8 pt-32 px-6 sm:px-8 lg:px-12">
      <div className="space-y-2">
        <h1 className="text-3xl sm:text-4xl font-bold tracking-tight text-white">Rehberler</h1>
        <p className="text-gray-300 text-lg">
          Otomobil dünyasının teknik derinliklerine inin. Şanzıman tipleri, motor teknolojileri ve daha fazlası.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pb-12">
        {GUIDES.map((guide) => {
          const Icon = ICONS[guide.id] || BookOpen;
          const isTransmission = guide.id === 'transmission';
          const isEngine = guide.id === 'engine';

          return (
            <div
              key={guide.id}
              onClick={() => {
                if (isTransmission) {
                  window.location.href = '/guides/transmission';
                } else if (isEngine) {
                  window.location.href = '/guides/engine';
                } else {
                  setSelectedGuide(guide);
                }
              }}
              className="group cursor-pointer rounded-3xl border border-white/10 bg-black/40 backdrop-blur-md p-6 hover:shadow-2xl hover:border-primary/30 transition-all duration-300 relative overflow-hidden"
            >
              <div className="absolute top-0 right-0 p-6 opacity-5 group-hover:opacity-10 transition-opacity transform group-hover:scale-110 duration-500">
                <Icon className="w-32 h-32 text-white" />
              </div>

              <div className="relative z-10 flex flex-col h-full">
                <div className="w-12 h-12 rounded-2xl bg-white/10 text-primary flex items-center justify-center mb-4 group-hover:bg-primary group-hover:text-white transition-colors duration-300">
                  <Icon className="w-6 h-6" />
                </div>

                <h3 className="text-xl font-bold text-white mb-2 group-hover:text-primary transition-colors">
                  {guide.title}
                </h3>

                <p className="text-sm text-gray-400 leading-relaxed mb-4 flex-grow">
                  {guide.summary}
                </p>

                <div className="flex items-center text-sm font-semibold text-primary">
                  Okumaya Başla
                  <ArrowLeft className="w-4 h-4 ml-2 rotate-180 group-hover:translate-x-1 transition-transform" />
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
