import { useState } from 'react';
import { ArrowLeft, Info, Zap, Settings, TrendingUp } from 'lucide-react';
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
        title: 'Otomatik Şanzıman (Tork Konvertörlü)',
        icon: Zap,
        description: 'Hidrolik tork konvertörü ve planetary gear setleri kullanarak otomatik vites değişimi sağlayan sistem.',
        mechanism: 'Tork konvertörü, motor ile şanzıman arasında hidrolik bir bağlantı oluşturur. İçindeki özel sıvı, motorun gücünü tekerleklere aktarır. Planetary gear setleri ve hidrolik valf gövdesi, farklı vites oranlarını otomatik olarak seçer.',
        pros: [
            'Maksimum konfor ve kullanım kolaylığı',
            'Yumuşak ve kesintisiz vites geçişleri',
            'Trafikte rahat kullanım',
            'Geniş model yelpazesinde bulunabilirlik',
            'Modern versiyonlarda yüksek verimlilik (8-10 ileri)'
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
        title: 'Çift Kavramalı Şanzıman (DCT/DSG)',
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
            'Öğrenme eğrisi gerektirir (sürüş tarzı)'
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
        title: 'CVT (Sürekli Değişken Şanzıman)',
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
        <div className="min-h-screen bg-background">
            {/* Hero Section */}
            <div className="relative overflow-hidden bg-gradient-to-br from-primary/10 via-background to-background border-b border-border">
                <div className="absolute inset-0 bg-grid-pattern opacity-5"></div>
                <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
                    <Link
                        to="/guides"
                        className="inline-flex items-center gap-2 text-sm text-text-muted hover:text-primary transition-colors mb-8"
                    >
                        <ArrowLeft className="w-4 h-4" />
                        Rehberlere Dön
                    </Link>

                    <div className="max-w-3xl">
                        <h1 className="text-4xl sm:text-5xl font-bold text-text-main mb-6">
                            Şanzıman Dünyası
                        </h1>
                        <p className="text-lg text-text-muted leading-relaxed mb-8">
                            Otomobilinizin kalbinden tekerleklerine gücü ileten kritik sistem: Şanzıman.
                            Manuel'den CVT'ye, her sistemin nasıl çalıştığını, avantajlarını ve dezavantajlarını keşfedin.
                        </p>
                        <div className="flex flex-wrap gap-4 text-sm">
                            <div className="flex items-center gap-2 bg-white px-4 py-2 rounded-full border border-border">
                                <Settings className="w-4 h-4 text-primary" />
                                <span className="text-text-main font-medium">4 Ana Tip</span>
                            </div>
                            <div className="flex items-center gap-2 bg-white px-4 py-2 rounded-full border border-border">
                                <Zap className="w-4 h-4 text-primary" />
                                <span className="text-text-main font-medium">Detaylı Mekanizma</span>
                            </div>
                            <div className="flex items-center gap-2 bg-white px-4 py-2 rounded-full border border-border">
                                <TrendingUp className="w-4 h-4 text-primary" />
                                <span className="text-text-main font-medium">Örnek Modeller</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Main Content */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                {/* Introduction */}
                <div className="bg-white rounded-3xl p-8 border border-border shadow-sm mb-12">
                    <h2 className="text-2xl font-bold text-text-main mb-4">Şanzıman Nedir?</h2>
                    <div className="prose prose-slate max-w-none">
                        <p className="text-text-muted leading-relaxed mb-4">
                            Şanzıman (transmisyon), motorun ürettiği gücü ve torku tekerleklere ileten, aynı zamanda
                            farklı hız ve yük koşullarına göre optimize eden mekanik bir sistemdir. Motorlar genellikle
                            belirli bir devir aralığında en verimli çalışır; şanzıman bu devir aralığını koruyarak
                            aracın farklı hızlarda hareket etmesini sağlar.
                        </p>
                        <p className="text-text-muted leading-relaxed">
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
                                className={`group cursor-pointer rounded-3xl bg-gradient-to-br ${transmission.color} p-6 text-white hover:shadow-2xl hover:scale-105 transition-all duration-300 relative overflow-hidden`}
                            >
                                <div className="absolute inset-0 bg-white/5 opacity-0 group-hover:opacity-100 transition-opacity"></div>

                                <div className="relative z-10">
                                    <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center backdrop-blur-sm mb-4">
                                        <Icon className="w-7 h-7" />
                                    </div>

                                    <div className="text-xs font-medium opacity-90 mb-2">Tip {index + 1}</div>
                                    <h3 className="text-xl font-bold mb-3 leading-tight">{transmission.title}</h3>
                                    <p className="text-white/90 text-sm leading-relaxed mb-4">{transmission.description}</p>

                                    <div className="flex items-center text-sm font-semibold">
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
                        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
                            <div className="bg-white rounded-3xl max-w-5xl w-full max-h-[90vh] overflow-y-auto shadow-2xl">
                                {/* Header */}
                                <div className={`bg-gradient-to-r ${transmission.color} p-8 text-white sticky top-0 z-10`}>
                                    <div className="flex items-start justify-between">
                                        <div className="flex items-start gap-4">
                                            <div className="w-12 h-12 bg-white/20 rounded-2xl flex items-center justify-center backdrop-blur-sm">
                                                <Icon className="w-6 h-6" />
                                            </div>
                                            <div>
                                                <div className="text-sm font-medium opacity-90 mb-1">Detaylı İnceleme</div>
                                                <h3 className="text-2xl font-bold mb-2">{transmission.title}</h3>
                                                <p className="text-white/90 text-sm max-w-2xl">{transmission.description}</p>
                                            </div>
                                        </div>
                                        <button
                                            onClick={() => setExpandedSection(null)}
                                            className="text-white/80 hover:text-white transition-colors p-2 hover:bg-white/10 rounded-xl"
                                        >
                                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>

                                {/* Content */}
                                <div className="p-8 space-y-8">
                                    {/* Mechanism */}
                                    <div>
                                        <h4 className="text-lg font-semibold text-text-main mb-3 flex items-center gap-2">
                                            <Settings className="w-5 h-5 text-primary" />
                                            Çalışma Mekanizması
                                        </h4>
                                        <p className="text-text-muted leading-relaxed">
                                            {transmission.mechanism}
                                        </p>
                                    </div>

                                    {/* Pros & Cons */}
                                    <div className="grid md:grid-cols-2 gap-6">
                                        {/* Pros */}
                                        <div className="bg-green-50 rounded-2xl p-6 border border-green-200">
                                            <h4 className="text-lg font-semibold text-green-900 mb-4 flex items-center gap-2">
                                                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                                                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                                                </svg>
                                                Avantajlar
                                            </h4>
                                            <ul className="space-y-2">
                                                {transmission.pros.map((pro, idx) => (
                                                    <li key={idx} className="text-sm text-green-800 flex items-start gap-2">
                                                        <span className="text-green-600 mt-0.5">✓</span>
                                                        <span>{pro}</span>
                                                    </li>
                                                ))}
                                            </ul>
                                        </div>

                                        {/* Cons */}
                                        <div className="bg-red-50 rounded-2xl p-6 border border-red-200">
                                            <h4 className="text-lg font-semibold text-red-900 mb-4 flex items-center gap-2">
                                                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                                                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                                                </svg>
                                                Dezavantajlar
                                            </h4>
                                            <ul className="space-y-2">
                                                {transmission.cons.map((con, idx) => (
                                                    <li key={idx} className="text-sm text-red-800 flex items-start gap-2">
                                                        <span className="text-red-600 mt-0.5">✗</span>
                                                        <span>{con}</span>
                                                    </li>
                                                ))}
                                            </ul>
                                        </div>
                                    </div>

                                    {/* Video */}
                                    {transmission.videoUrl && (
                                        <div>
                                            <h4 className="text-lg font-semibold text-text-main mb-3 flex items-center gap-2">
                                                <svg className="w-5 h-5 text-primary" fill="currentColor" viewBox="0 0 20 20">
                                                    <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" />
                                                </svg>
                                                Açıklama Videosu
                                            </h4>
                                            <div className="rounded-2xl overflow-hidden border border-border bg-slate-50">
                                                <div className="relative" style={{ paddingBottom: '56.25%' }}>
                                                    <iframe
                                                        className="absolute top-0 left-0 w-full h-full"
                                                        src={transmission.videoUrl}
                                                        title={`${transmission.title} - Açıklama Videosu`}
                                                        frameBorder="0"
                                                        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                                                        allowFullScreen
                                                    ></iframe>
                                                </div>
                                            </div>
                                            <p className="text-xs text-text-muted mt-2">
                                                📺 Anlatan Adamlar kanalı tarafından hazırlanmıştır
                                            </p>
                                        </div>
                                    )}

                                    {/* Image */}
                                    <div>
                                        <h4 className="text-lg font-semibold text-text-main mb-3 flex items-center gap-2">
                                            <svg className="w-5 h-5 text-primary" fill="currentColor" viewBox="0 0 20 20">
                                                <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
                                            </svg>
                                            Teknik Görsel
                                        </h4>
                                        <div className="rounded-2xl overflow-hidden border border-border bg-slate-50">
                                            <img
                                                src={transmission.image}
                                                alt={transmission.title}
                                                className="w-full h-auto"
                                            />
                                        </div>
                                    </div>

                                    {/* Examples */}
                                    <div>
                                        <h4 className="text-lg font-semibold text-text-main mb-4 flex items-center gap-2">
                                            <TrendingUp className="w-5 h-5 text-primary" />
                                            Popüler Örnekler
                                        </h4>
                                        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-3">
                                            {transmission.examples.map((example, idx) => (
                                                <div
                                                    key={idx}
                                                    className="bg-slate-50 rounded-xl px-4 py-3 border border-slate-200 text-sm text-text-main font-medium"
                                                >
                                                    {example}
                                                </div>
                                            ))}
                                        </div>
                                    </div>

                                    {/* Close Button */}
                                    <div className="flex justify-center pt-4">
                                        <button
                                            onClick={() => setExpandedSection(null)}
                                            className="px-6 py-3 bg-slate-100 hover:bg-slate-200 text-text-main font-semibold rounded-xl transition-colors"
                                        >
                                            Kapat
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    );
                })()}

                {/* Comparison Section */}
                <div className="mt-12 bg-white rounded-3xl p-8 border border-border shadow-sm">
                    <h2 className="text-2xl font-bold text-text-main mb-6">Karşılaştırma ve Seçim Rehberi</h2>

                    <div className="rounded-2xl overflow-hidden border border-border mb-6">
                        <img
                            src="/images/guides/transmission_comparison_1769183038027.png"
                            alt="Şanzıman Karşılaştırması"
                            className="w-full h-auto"
                        />
                    </div>

                    <div className="space-y-6">
                        <div>
                            <h3 className="text-lg font-semibold text-text-main mb-3">Hangi Kullanım İçin Hangi Şanzıman?</h3>
                            <div className="grid md:grid-cols-2 gap-4">
                                <div className="bg-blue-50 rounded-xl p-4 border border-blue-200">
                                    <div className="font-semibold text-blue-900 mb-2">🏙️ Şehir İçi Kullanım</div>
                                    <div className="text-sm text-blue-800">Otomatik veya CVT - Konfor ve yakıt ekonomisi</div>
                                </div>
                                <div className="bg-purple-50 rounded-xl p-4 border border-purple-200">
                                    <div className="font-semibold text-purple-900 mb-2">🏁 Sportif Sürüş</div>
                                    <div className="text-sm text-purple-800">DCT veya Manuel - Hız ve kontrol</div>
                                </div>
                                <div className="bg-green-50 rounded-xl p-4 border border-green-200">
                                    <div className="font-semibold text-green-900 mb-2">🛣️ Uzun Yol</div>
                                    <div className="text-sm text-green-800">Otomatik (8+ ileri) - Konfor ve verimlilik</div>
                                </div>
                                <div className="bg-orange-50 rounded-xl p-4 border border-orange-200">
                                    <div className="font-semibold text-orange-900 mb-2">💰 Düşük Maliyet</div>
                                    <div className="text-sm text-orange-800">Manuel - Bakım ve yakıt tasarrufu</div>
                                </div>
                            </div>
                        </div>

                        <div className="bg-slate-50 rounded-xl p-6 border border-slate-200">
                            <h4 className="font-semibold text-text-main mb-3">💡 Uzman Tavsiyesi</h4>
                            <p className="text-sm text-text-muted leading-relaxed">
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
    );
}
