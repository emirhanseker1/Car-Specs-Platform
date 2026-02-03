import { useState } from 'react';
import { ArrowLeft, Info, Zap, Settings, TrendingUp, CheckCircle2, XCircle, CarFront } from 'lucide-react';
import { Link } from 'react-router-dom';

interface TransmissionSection {
    id: string;
    title: string;
    icon: any;
    description: string;
    mechanism: string;
    pros: string[];
    cons: string[];
    examples: string[];
    image: string;
    color: string;
    videoUrl?: string;
}

const TRANSMISSION_TYPES: TransmissionSection[] = [
    {
        id: 'manual',
        title: 'Manuel Şanzıman',
        icon: Settings,
        description: 'Sürücünün kavrama pedalı ve vites kolu ile doğrudan kontrol ettiği geleneksel şanzıman sistemi.',
        mechanism: 'Manuel şanzıman, sürücünün kavrama pedalına basarak motor ile şanzıman arasındaki bağlantıyı kestiği ve vites kolunu hareket ettirerek farklı dişli oranlarını seçtiği bir sistemdir. Senkronizörler, viteslerin yumuşak bir şekilde geçişini sağlar.',
        pros: [
            'Tam sürücü kontrolü ve bağlantı hissi',
            'Dayanıklılık ve uzun ömür',
            'Düşük bakım ve onarım maliyeti',
            'Daha hafif yapı',
            'Yakıt ekonomisi (doğru kullanımda)'
        ],
        cons: [
            'Öğrenme eğrisi gerektirir',
            'Yoğun trafikte yorucu olabilir',
            'Vites değişimleri daha yavaş',
            'Yanlış kullanımda aşınma riski'
        ],
        examples: [
            'Honda Civic Type R (6 İleri)',
            'Mazda MX-5 Miata',
            'Ford Focus ST',
            'Porsche 911 GT3 (7 İleri)',
            'Toyota GR86'
        ],
        image: '/images/guides/manual_transmission_1769182954404.png',
        color: 'from-slate-700 to-slate-900',
        videoUrl: 'https://www.youtube.com/embed/JtUX0YLD_48'
    },
    {
        id: 'automatic',
        title: 'Otomatik Şanzıman',
        icon: Zap,
        description: 'Hidrolik tork konvertörü ve planetary gear setleri kullanarak otomatik vites değişimi sağlayan sistem.',
        mechanism: 'Tork konvertörü, motor ile şanzıman arasında hidrolik bir bağlantı oluşturur. İçindeki özel sıvı, motorun gücünü tekerleklere aktarır. Planetary gear setleri ve hidrolik valf gövdesi, farklı vites oranlarını otomatik olarak seçer.',
        pros: [
            'Maksimum konfor ve kullanım kolaylığı',
            'Yumuşak ve kesintisiz vites geçişleri',
            'Trafikte rahat kullanım',
            'Geniş model yelpazesinde bulunabilirlik',
            'Modern versiyonlarda yüksek verimlilik'
        ],
        cons: [
            'Manuel kadar doğrudan kontrol hissi yok',
            'Daha ağır ve karmaşık yapı',
            'Bakım maliyeti daha yüksek',
            'Eski modellerde yakıt tüketimi fazla olabilir'
        ],
        examples: [
            'BMW ZF 8HP (8 İleri)',
            'Mercedes-Benz 9G-Tronic',
            'Toyota Aisin 8-Speed',
            'Lexus 10-Speed Automatic',
            'Ford 10R80 (10 İleri)'
        ],
        image: '/images/guides/automatic_transmission_1769182969695.png',
        color: 'from-indigo-700 to-indigo-900',
        videoUrl: 'https://www.youtube.com/embed/LdtXy9By3po'
    },
    {
        id: 'dct',
        title: 'Çift Kavramalı (DCT)',
        icon: TrendingUp,
        description: 'İki ayrı kavrama ve vites seti kullanarak yıldırım hızında vites değişimi sunan modern sistem.',
        mechanism: 'DCT, iki ayrı manuel şanzımanın iç içe geçmiş halidir. Bir kavrama tek vitesleri (1,3,5,7), diğeri çift vitesleri (2,4,6) kontrol eder. Bir sonraki vites önceden hazırlandığı için geçiş milisaniyeler sürer. Mekatronik ünite tüm işlemleri elektronik olarak yönetir.',
        pros: [
            'Çok hızlı vites değişimleri (0.2 saniyeden az)',
            'Mükemmel yakıt ekonomisi',
            'Sportif sürüş performansı',
            'Manuel moda geçiş imkanı',
            'Güç kaybı minimum'
        ],
        cons: [
            'Düşük hızlarda ve trafikte sarsıntılı olabilir',
            'Isınma problemleri (özellikle kuru kavramalı)',
            'Bakım ve onarım maliyeti yüksek',
            'Mekatronik arızaları pahalı',
            'Öğrenme eğrisi gerektirir'
        ],
        examples: [
            'Volkswagen DSG (6/7 İleri)',
            'Porsche PDK (7 İleri)',
            'Hyundai/Kia DCT',
            'Renault EDC',
            'Ford PowerShift (6 İleri)'
        ],
        image: '/images/guides/dct_transmission_1769182986472.png',
        color: 'from-blue-800 to-cyan-900',
        videoUrl: 'https://www.youtube.com/embed/0y8s8sL70pQ'
    },
    {
        id: 'cvt',
        title: 'CVT (Değişken)',
        icon: Info,
        description: 'Kasnak ve kayış sistemi ile sonsuz vites oranı sunan, yakıt ekonomisine odaklanan teknoloji.',
        mechanism: 'CVT, sabit dişliler yerine değişken çaplı iki kasnak ve bunları birbirine bağlayan çelik kayış kullanır. Kasnakların çapı hidrolik olarak değiştirildiğinde, vites oranı sürekli ve kademesiz olarak ayarlanır. Bu sayede motor her zaman en verimli devirde çalışır.',
        pros: [
            'Sonsuz vites oranı - en iyi yakıt ekonomisi',
            'Yumuşak ve kesintisiz ivmelenme',
            'Basit mekanik yapı (az parça)',
            'Sessiz çalışma',
            'Şehir içi kullanımda ideal'
        ],
        cons: [
            '"Lastik bant etkisi" - motor sesi sabit kalır',
            'Sportif sürüş hissi zayıf',
            'Yüksek tork kapasitesi sınırlı',
            'Kayış ömrü sınırlı olabilir',
            'Bazı sürücüler alışamayabilir'
        ],
        examples: [
            'Toyota Corolla Hybrid (e-CVT)',
            'Nissan X-Trail / Qashqai',
            'Subaru WRX (Lineartronic)',
            'Honda Civic CVT',
            'Mitsubishi Outlander'
        ],
        image: '/images/guides/cvt_transmission_1769183002788.png',
        color: 'from-teal-700 to-teal-900',
        videoUrl: 'https://www.youtube.com/embed/bz6LBCj6W-c'
    }
];

export default function TransmissionGuide() {
    const [expandedSection, setExpandedSection] = useState<string | null>(null);

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

            <div className="relative z-10">
                                    
                {/* Hero Section */}
                <div className="relative border-b border-white/10 bg-black/20 backdrop-blur-sm">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 pt-32">
                        <Link
                            to="/guides"
                            className="inline-flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors mb-8 bg-white/5 hover:bg-white/10 px-4 py-2 rounded-full w-fit border border-white/5"
                        >
                            <ArrowLeft className="w-4 h-4" />
                            Rehberlere Dön
                        </Link>

                        <div className="max-w-3xl">
                            <h1 className="text-lg sm:text-6xl font-black text-white mb-6 drop-shadow-2xl">
                                Şanzıman Dünyası
                            </h1>
                            <p className="text-lg text-slate-300 leading-relaxed mb-8 font-light">
                                Otomobilinizin kalbinden tekerleklerine gücü ileten kritik sistem: Şanzıman.
                                Manuel'den CVT'ye, her sistemin nasıl çalıştığını, avantajlarını ve dezavantajlarını keşfedin.
                            </p>
                            <div className="flex flex-wrap gap-3 text-sm">
                                <div className="flex items-center gap-2 bg-[#1e293b]/80 backdrop-blur px-4 py-2 rounded-xl border border-white/10 text-white">
                                    <Settings className="w-4 h-4 text-primary" />
                                    <span className="font-medium">4 Ana Tip</span>
                                </div>
                                <div className="flex items-center gap-2 bg-[#1e293b]/80 backdrop-blur px-4 py-2 rounded-xl border border-white/10 text-white">
                                    <Zap className="w-4 h-4 text-primary" />
                                    <span className="font-medium">Detaylı Mekanizma</span>
                                </div>
                                <div className="flex items-center gap-2 bg-[#1e293b]/80 backdrop-blur px-4 py-2 rounded-xl border border-white/10 text-white">
                                    <TrendingUp className="w-4 h-4 text-primary" />
                                    <span className="font-medium">Örnek Modeller</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Main Content */}
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                    {/* Introduction */}
                    <div className="bg-[#1e293b]/60 backdrop-blur-md rounded-2xl p-5 border border-white/10 shadow-lg mb-12">
                        <h2 className="text-lg font-bold text-white mb-4">Şanzıman Nedir?</h2>
                        <div className="prose prose-invert max-w-none text-slate-300">
                            <p className="leading-relaxed mb-4">
                                Şanzıman (transmisyon), motorun ürettiği gücü ve torku tekerleklere ileten, aynı zamanda
                                farklı hız ve yük koşullarına göre optimize eden mekanik bir sistemdir. Motorlar genellikle
                                belirli bir devir aralığında en verimli çalışır; şanzıman bu devir aralığını koruyarak
                                aracın farklı hızlarda hareket etmesini sağlar.
                            </p>
                            <p className="leading-relaxed">
                                Farklı vites oranları sayesinde, düşük hızlarda yüksek tork (çekiş gücü) ve yüksek hızlarda
                                düşük motor devri (yakıt ekonomisi) elde edilir. Modern otomobillerde manuel, otomatik,
                                çift kavramalı (DCT) ve sürekli değişken (CVT) olmak üzere dört ana şanzıman tipi bulunur.
                            </p>
                        </div>
                    </div>



                    {/* Transmission Types Grid */}
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
                        {TRANSMISSION_TYPES.map((transmission, index) => {
                            const Icon = transmission.icon;

                            return (
                                <div
                                    key={transmission.id}
                                    onClick={() => setExpandedSection(transmission.id)}
                                    className={`group cursor-pointer rounded-2xl bg-gradient-to-br ${transmission.color} p-5 text-white shadow-xl border border-white/10 hover:border-white/30 hover:scale-105 transition-all duration-300 relative overflow-hidden`}
                                >
                                    <div className="absolute inset-0 bg-black/20 group-hover:bg-transparent transition-colors"></div>

                                    <div className="relative z-10">
                                    
                                        <div className="w-14 h-14 bg-white/10 rounded-2xl flex items-center justify-center backdrop-blur-sm mb-4 border border-white/10">
                                            <Icon className="w-7 h-7" />
                                        </div>

                                        <div className="text-xs font-bold uppercase tracking-wider opacity-70 mb-2">Tip {index + 1}</div>
                                        <h3 className="text-lg font-bold mb-3 leading-tight">{transmission.title}</h3>
                                        <p className="text-white/80 text-sm leading-relaxed mb-4">{transmission.description}</p>

                                        {/* DSG Guide Button for DCT */}
                                        {transmission.id === 'dct' && (
                                            <Link
                                                to="/guides/transmission/dsg"
                                                onClick={(e) => e.stopPropagation()}
                                                className="flex items-center gap-2 bg-white/10 hover:bg-white/20 text-white text-xs font-bold px-4 py-2 rounded-lg border border-white/10 mb-4 transition-all w-full justify-center group/btn"
                                            >
                                                <Settings className="w-3.5 h-3.5 group-hover/btn:rotate-90 transition-transform" />
                                                VAG DSG Rehberi
                                            </Link>
                                        )}

                                        <div className="flex items-center text-sm font-semibold mt-auto">
                                            Detayları Gör
                                            <svg className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                            </svg>
                                        </div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>

                    {/* Detail Modal */}
                    {expandedSection && (() => {
                        const transmission = TRANSMISSION_TYPES.find(t => t.id === expandedSection);
                        if (!transmission) return null;
                        const Icon = transmission.icon;

                        return (
                            <div className="fixed inset-0 z-[9999] bg-[#0f172a] overflow-y-auto animate-in fade-in zoom-in-95 duration-300">
                                {/* Sticky Header */}
                                <div className="sticky top-0 z-[100] bg-[#0f172a]/90 backdrop-blur-xl border-b border-white/5 shadow-2xl">
                                    <div className="max-w-[1920px] mx-auto px-5 py-3 flex items-center justify-between">
                                        <div className="flex items-center gap-4">
                                            <div className={`w-12 h-10 rounded-lg flex items-center justify-center bg-gradient-to-br ${transmission.color} shadow-lg ring-1 ring-white/10`}>
                                                <Icon className="w-4 h-4 text-white" />
                                            </div>
                                            <div>
                                                <div className="flex items-center gap-2">
                                                    <h3 className="text-lg font-bold text-white leading-tight">{transmission.title}</h3>
                                                    <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400 bg-slate-800/50 px-2 py-0.5 rounded-md border border-slate-700">Teknik İnceleme</span>
                                                </div>
                                                <div className="text-xs font-medium text-slate-400 hidden sm:block">Detaylı mekanizma ve özellikler</div>
                                            </div>
                                        </div>
                                        <button
                                            onClick={() => setExpandedSection(null)}
                                            className="group flex items-center gap-2 bg-white/5 hover:bg-white/10 px-5 py-2.5 rounded-lg border border-white/10 transition-all hover:border-white/20 active:scale-95"
                                        >
                                            <span className="text-sm font-bold text-slate-300 group-hover:text-white">KAPAT</span>
                                            <div className="bg-white/10 rounded-full p-0.5 group-hover:bg-white/20 transition-colors">
                                                <XCircle className="w-4 h-4 text-slate-300 group-hover:text-white" />
                                            </div>
                                        </button>
                                    </div>
                                </div>

                                {/* Dashboard Content */}
                                <div className="max-w-[1920px] mx-auto px-5 py-8 md:py-12">
                                    <div className="grid grid-cols-1 xl:grid-cols-12 gap-5 md:gap-10 items-start">
                                        
                                        {/* LEFT COLUMN (7/12) */}
                                        <div className="xl:col-span-9 space-y-8">
                                            {/* Hero Description */}
                                            <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-slate-900 via-slate-900 to-[#0B1120] border border-white/5 p-5 md:p-5 shadow-2xl">
                                                <div className={`absolute top-0 right-0 w-[600px] h-[600px] bg-gradient-to-br ${transmission.color} opacity-[0.1] blur-[100px] rounded-full pointer-events-none -translate-y-1/2 translate-x-1/4`} />
                                                <h2 className="relative text-lg md:text-lg font-black text-white mb-6 leading-tight tracking-tight">
                                                    {transmission.title} <span className="text-slate-500">Nedir?</span>
                                                </h2>
                                                <p className="relative text-base md:text-lg text-slate-300 leading-relaxed font-light max-w-4xl">
                                                    {transmission.description}
                                                </p>
                                            </div>

                                            {/* Mechanism Card */}
                                            <div className="bg-blue-900/5 border border-blue-500/10 rounded-2xl p-5 md:p-10 relative overflow-hidden group">
                                                <div className="absolute top-0 right-0 w-96 h-96 bg-blue-500/10 blur-[120px] rounded-full pointer-events-none group-hover:bg-blue-500/15 transition-all duration-700" />
                                                <h4 className="relative text-lg font-bold text-blue-100 mb-6 flex items-center gap-3">
                                                    <div className="bg-blue-500/20 p-2.5 rounded-lg ring-1 ring-blue-500/30">
                                                        <Settings className="w-4 h-4 text-blue-400" />
                                                    </div>
                                                    Çalışma Mekanizması
                                                </h4>
                                                <p className="relative text-blue-100/80 leading-relaxed text-base md:text-lg font-light">
                                                    {transmission.mechanism}
                                                </p>
                                            </div>

                                            {/* Visuals Grid */}
                                            <div className="space-y-6">
                                                {/* Video */}
                                                {transmission.videoUrl && (
                                                    <div className="space-y-4">
                                                        <h4 className="text-base font-bold text-slate-200 flex items-center gap-3 px-2">
                                                            <Info className="w-4 h-4 text-amber-400" />
                                                            Nasıl Çalışır?
                                                        </h4>
                                                        <div className="rounded-2xl overflow-hidden border border-white/10 bg-black shadow-2xl ring-1 ring-white/5 group">
                                                            <div className="relative aspect-video">
                                                                <iframe
                                                                    className="absolute top-0 left-0 w-full h-full opacity-90 group-hover:opacity-100 transition-opacity"
                                                                    src={transmission.videoUrl}
                                                                    title="Video"
                                                                    frameBorder="0"
                                                                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                                                                    allowFullScreen
                                                                ></iframe>
                                                            </div>
                                                        </div>
                                                    </div>
                                                )}

                                                {/* Image */}
                                                
                                                </div>
                                            </div>
                                        </div>

                                        {/* RIGHT COLUMN (5/12) */}
                                        <div className="xl:col-span-3 space-y-8 sticky top-32">
                                            
                                            {/* Pros & Cons */}
                                            <div className="grid gap-5">
                                                {/* Advantages */}
                                                <div className="bg-emerald-950/20 border border-emerald-500/10 rounded-2xl p-5">
                                                    <h4 className="text-lg font-bold text-emerald-400 mb-6 flex items-center gap-3">
                                                        <div className="bg-emerald-500/10 p-2 rounded-lg ring-1 ring-emerald-500/20">
                                                           <CheckCircle2 className="w-4 h-4" />
                                                        </div>
                                                        Avantajlar
                                                    </h4>
                                                    <div className="grid gap-3">
                                                        {transmission.pros.map((pro, idx) => (
                                                            <div key={idx} className="bg-emerald-500/5 border border-emerald-500/10 rounded-lg p-4 flex items-start gap-4 hover:bg-emerald-500/10 transition-colors">
                                                                <CheckCircle2 className="w-4 h-4 text-emerald-500 flex-shrink-0 mt-0.5" />
                                                                <span className="text-emerald-100/90 font-medium">{pro}</span>
                                                            </div>
                                                        ))}
                                                    </div>
                                                </div>

                                                {/* Disadvantages */}
                                                <div className="bg-rose-950/10 border border-rose-500/10 rounded-2xl p-5">
                                                    <h4 className="text-lg font-bold text-rose-400 mb-6 flex items-center gap-3">
                                                        <div className="bg-rose-500/10 p-2 rounded-lg ring-1 ring-rose-500/20">
                                                           <XCircle className="w-4 h-4" />
                                                        </div>
                                                        Dezavantajlar
                                                    </h4>
                                                    <div className="grid gap-3">
                                                        {transmission.cons.map((con, idx) => (
                                                            <div key={idx} className="bg-rose-500/5 border border-rose-500/10 rounded-lg p-4 flex items-start gap-4 hover:bg-rose-500/10 transition-colors">
                                                                <XCircle className="w-4 h-4 text-rose-500 flex-shrink-0 mt-0.5" />
                                                                <span className="text-rose-100/90 font-medium">{con}</span>
                                                            </div>
                                                        ))}
                                        </div>
                                    </div>

                                            {/* Examples */}
                                            <div className="bg-slate-900/50 border border-white/5 rounded-2xl p-5">
                                                <h4 className="text-lg font-bold text-indigo-300 mb-6 flex items-center gap-3">
                                                    <div className="bg-indigo-500/10 p-2 rounded-lg ring-1 ring-indigo-500/20">
                                                        <CarFront className="w-4 h-4 text-indigo-400" />
                                                    </div>
                                                    Popüler Örnekler
                                                </h4>
                                                <div className="flex flex-wrap gap-3">
                                                    {transmission.examples.map((example, idx) => (
                                                        <div
                                                            key={idx}
                                                            className="bg-slate-800/50 hover:bg-indigo-600/20 border border-white/5 hover:border-indigo-500/30 rounded-full px-5 py-3 text-sm text-slate-300 hover:text-white font-medium transition-all cursor-default"
                                                        >
                                                            {example}
                                                        </div>
                                                    ))}
                                                </div>
                                            </div>

                                        </div>

                                    </div>
                                </div>
                            </div>
                        );
                    })()}

{/* Comparison Section */}
                    <div className="mt-12 bg-[#1e293b]/60 backdrop-blur-md rounded-2xl p-5 border border-white/10 shadow-lg">
                        <h2 className="text-lg font-bold text-white mb-6">Karşılaştırma ve Seçim Rehberi</h2>

                        <div className="rounded-2xl overflow-hidden border border-white/10 mb-8 bg-black/40">
                            <img
                                src="/images/guides/transmission_comparison_1769183038027.png"
                                alt="Şanzıman Karşılaştırması"
                                className="w-full h-auto opacity-90"
                            />
                        </div>

                        <div className="space-y-6">
                            <div>
                                <h3 className="text-lg font-semibold text-white mb-4">Hangi Kullanım İçin Hangi Şanzıman?</h3>
                                <div className="grid md:grid-cols-2 gap-4">
                                    <div className="bg-blue-500/10 rounded-xl p-4 border border-blue-500/20">
                                        <div className="font-bold text-blue-400 mb-2">🏙️ Şehir İçi Kullanım</div>
                                        <div className="text-sm text-blue-200">Otomatik veya CVT - Konfor ve yakıt ekonomisi</div>
                                    </div>
                                    <div className="bg-purple-500/10 rounded-xl p-4 border border-purple-500/20">
                                        <div className="font-bold text-purple-400 mb-2">🏁 Sportif Sürüş</div>
                                        <div className="text-sm text-purple-200">DCT veya Manuel - Hız ve kontrol</div>
                                    </div>
                                    <div className="bg-green-500/10 rounded-xl p-4 border border-green-500/20">
                                        <div className="font-bold text-green-400 mb-2">🛣️ Uzun Yol</div>
                                        <div className="text-sm text-green-200">Otomatik (8+ ileri) - Konfor ve verimlilik</div>
                                    </div>
                                    <div className="bg-orange-500/10 rounded-xl p-4 border border-orange-500/20">
                                        <div className="font-bold text-orange-400 mb-2">💰 Düşük Maliyet</div>
                                        <div className="text-sm text-orange-200">Manuel - Bakım ve yakıt tasarrufu</div>
                                    </div>
                                </div>
                            </div>

                            <div className="bg-white/5 rounded-xl p-5 border border-white/10">
                                <h4 className="font-bold text-white mb-3">💡 Uzman Tavsiyesi</h4>
                                <p className="text-sm text-slate-400 leading-relaxed">
                                    Şanzıman seçimi tamamen kullanım amacınıza bağlıdır. Yoğun şehir trafiğinde her gün
                                    kullanacaksanız otomatik veya CVT konforlu olacaktır. Sportif sürüş ve performans
                                    arıyorsanız DCT veya manuel tercih edilebilir. Modern otomatik şanzımanlar (8-10 ileri)
                                    artık yakıt ekonomisinde de manuel kadar verimli olabilmektedir.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
