import { Target, Info, Database, Layers, CheckCircle2, ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function About() {
  return (
    <div className="relative min-h-screen font-sans text-slate-200">
      {/* Background Image & Overlay */}
      <div className="fixed inset-0 z-0">
        <div
          className="absolute inset-0 bg-cover bg-center"
          style={{ backgroundImage: 'url(/hero-2.jpg)' }}
        />
        <div className="absolute inset-0 bg-slate-900/90 backdrop-blur-sm" />
      </div>

      <main className="relative z-10 max-w-7xl mx-auto px-4 py-8 pt-32 space-y-16">

        {/* Hero Section */}
        <div className="text-center max-w-3xl mx-auto space-y-6">
          <div className="inline-flex items-center justify-center p-2 bg-primary/10 rounded-full mb-4 ring-1 ring-primary/30">
            <Info className="w-5 h-5 text-primary mr-2" />
            <span className="text-primary text-sm font-bold uppercase tracking-wide">Hakkımızda</span>
          </div>
          <h1 className="text-5xl md:text-6xl font-black text-white tracking-tight drop-shadow-2xl">
            Türkiye'nin En Kapsamlı <br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-orange-400">
              Araç Teknik Bilgi Platformu
            </span>
          </h1>
          <p className="text-xl text-slate-300 leading-relaxed font-light">
            Araç satın alma süreçlerindeki bilgi kirliliğini ortadan kaldırıyoruz. Doğru, detaylı ve güvenilir teknik verilere tek bir noktadan ulaşın.
          </p>
        </div>

        {/* Cards Section */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Mission Card */}
          <div className="group relative bg-[#1e293b]/60 backdrop-blur-md border border-white/10 rounded-3xl p-8 hover:bg-[#1e293b]/80 transition-all duration-300 hover:shadow-2xl hover:shadow-primary/10 hover:border-primary/30 overflow-hidden">
            <div className="absolute top-0 right-0 p-8 opacity-10 group-hover:opacity-20 transition-opacity">
              <Target size={120} />
            </div>
            <div className="relative z-10">
              <div className="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center mb-6 text-blue-400 grouping-hover:scale-110 transition-transform">
                <Target size={28} strokeWidth={2.5} />
              </div>
              <h2 className="text-2xl font-bold text-white mb-4">Misyon & Amaç</h2>
              <p className="text-slate-300 leading-relaxed">
                İlan sitelerindeki teknik bilgiler genellikle eksik veya tutarsızdır. Amacımız, kullanıcıların ilgilendikleri marka ve modeli seçerek, fabrika verilerine dayalı en doğru bilgilere en hızlı şekilde ulaşmalarını sağlamaktır.
              </p>
              <ul className="mt-6 space-y-3">
                <li className="flex items-center gap-3 text-slate-400">
                  <CheckCircle2 size={18} className="text-blue-500" />
                  <span>Doğrulanmış fabrika verileri</span>
                </li>
                <li className="flex items-center gap-3 text-slate-400">
                  <CheckCircle2 size={18} className="text-blue-500" />
                  <span>Kolay karşılaştırma imkanı</span>
                </li>
                <li className="flex items-center gap-3 text-slate-400">
                  <CheckCircle2 size={18} className="text-blue-500" />
                  <span>Kullanıcı dostu arayüz</span>
                </li>
              </ul>
            </div>
          </div>

          {/* Scope Card */}
          <div className="group relative bg-[#1e293b]/60 backdrop-blur-md border border-white/10 rounded-3xl p-8 hover:bg-[#1e293b]/80 transition-all duration-300 hover:shadow-2xl hover:shadow-purple/10 hover:border-purple-500/30 overflow-hidden">
            <div className="absolute top-0 right-0 p-8 opacity-10 group-hover:opacity-20 transition-opacity">
              <Database size={120} />
            </div>
            <div className="relative z-10">
              <div className="w-12 h-12 bg-purple-500/20 rounded-xl flex items-center justify-center mb-6 text-purple-400 group-hover:scale-110 transition-transform">
                <Layers size={28} strokeWidth={2.5} />
              </div>
              <h2 className="text-2xl font-bold text-white mb-4">Veri Kapsamı</h2>
              <p className="text-slate-300 leading-relaxed">
                Standart özelliklerin ötesine geçiyoruz. Diğer platformlarda bulamayacağınız derin teknik detayları sunuyoruz.
              </p>
              <div className="grid grid-cols-2 gap-4 mt-6">
                {[
                  'Şanzıman Tipi ve Kodu',
                  'Tork Dayanım Değerleri',
                  'Motor Kodları & Aileleri',
                  'Kronik Problemler',
                  'Platform Ortaklıkları',
                  'Parça Uyumlulukları'
                ].map((item, i) => (
                  <div key={i} className="flex items-center gap-2 text-sm text-slate-400 bg-white/5 p-2 rounded-lg border border-white/5">
                    <div className="w-1.5 h-1.5 rounded-full bg-purple-500" />
                    {item}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Call to Action */}
        <div className="relative rounded-3xl overflow-hidden text-center bg-gradient-to-r from-primary/20 via-[#1e293b] to-[#1e293b] border border-white/10 p-12">
          <div className="relative z-10 flex flex-col items-center gap-6">
            <h2 className="text-3xl font-bold text-white">Araçları Keşfetmeye Başlayın</h2>
            <p className="text-slate-300 max-w-xl">
              binlerce araçlık veritabanımızda aradığınız o modeli bulun, teknik detaylarına hakim olun.
            </p>
            <Link
              to="/search"
              className="inline-flex items-center gap-2 bg-primary hover:bg-primary-hover text-white px-8 py-4 rounded-xl font-bold text-lg transition-all shadow-lg shadow-primary/25 hover:scale-105"
            >
              Detaylı Arama Yap <ArrowRight size={20} />
            </Link>
          </div>
        </div>

      </main>
    </div>
  );
}
