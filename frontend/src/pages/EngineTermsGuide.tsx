import { useState } from 'react';
import { ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';

interface EngineTerm {
    id: string;
    title: string;
    shortDesc: string;
    fullDesc: string;
    icon: string;
    examples?: string[];
}

const ENGINE_TERMS: EngineTerm[] = [
    {
        id: 'horsepower',
        title: 'Beygir Gücü (HP/PS)',
        shortDesc: 'Motorun maksimum hızda ürettiği güç',
        fullDesc: 'Beygir gücü, motorun belirli bir devirde (genellikle maksimum devirde) ürettiği gücü ifade eder. Aracın ne kadar hızlı gidebileceğini, maksimum hızını ve yüksek devirlerdeki performansını belirler. 1 HP = 745.7 Watt\'tır. Örneğin 200 HP\'lik bir motor, 149 kW güç üretir.',
        icon: '/images/guides/horsepower_icon_1769194949851.png',
        examples: ['Civic Type R: 320 HP', 'Golf GTI: 245 HP', 'Porsche 911: 450+ HP']
    },
    {
        id: 'torque',
        title: 'Tork (Nm)',
        shortDesc: 'Motorun dönme kuvveti ve çekiş gücü',
        fullDesc: 'Tork, motorun ürettiği dönme kuvvetidir ve Newton-metre (Nm) ile ölçülür. Aracın ne kadar hızlı ivmelenebileceğini, yokuş tırmanma kabiliyetini ve yük taşıma kapasitesini belirler. Yüksek tork, özellikle düşük devirlerde güçlü çekiş sağlar. Dizel motorlar genellikle benzinli motorlara göre daha yüksek tork üretir.',
        icon: '/images/guides/torque_icon_1769194961203.png',
        examples: ['Dizel Motor: 400 Nm @ 1750 rpm', 'Benzinli Turbo: 350 Nm @ 2500 rpm', 'Elektrikli: 600+ Nm @ 0 rpm']
    },
    {
        id: 'turbo',
        title: 'Turbo Besleme',
        shortDesc: 'Egzoz gazıyla motora daha fazla hava basma sistemi',
        fullDesc: 'Turboşarj, egzoz gazlarının enerjisini kullanarak bir türbin döndürür ve bu türbin motora daha fazla hava pompalar. Daha fazla hava = daha fazla yakıt = daha fazla güç demektir. Küçük hacimli motorlardan büyük güç elde edilmesini sağlar. Ancak "turbo lag" denilen gecikme yaşanabilir. Modern motorlarda twin-scroll turbo ve değişken geometrili turbo ile bu sorun minimize edilmiştir.',
        icon: '/images/guides/turbo_icon_1769194973896.png',
        examples: ['1.5 TSI: 150 HP (Turbo)', '1.5 NA: 110 HP (Atmosferik)', 'Bi-Turbo V8: 600+ HP']
    },
    {
        id: 'displacement',
        title: 'Motor Hacmi (cc/L)',
        shortDesc: 'Silindirlerin toplam hacmi',
        fullDesc: 'Motor hacmi, tüm silindirlerin toplam hacmini ifade eder ve cc (santimetreküp) veya litre (L) ile ölçülür. 1000 cc = 1.0 L\'dir. Daha büyük hacim genellikle daha fazla güç ve tork anlamına gelir, ancak yakıt tüketimi de artar. Modern turbo motorlar, küçük hacimle yüksek performans sunarak "downsizing" trendini başlatmıştır.',
        icon: '/images/guides/displacement_icon_1769195000365.png',
        examples: ['1.0 TSI: 999 cc', '2.0 TDI: 1968 cc', '5.0 V8: 4999 cc']
    },
    {
        id: 'cylinders',
        title: 'Silindir Sayısı ve Dizilimi',
        shortDesc: 'Motorun silindir konfigürasyonu',
        fullDesc: 'Silindir sayısı ve dizilimi, motorun performansını, titreşimini ve ses karakteristiğini etkiler. Inline-4 (I4) en yaygın olanıdır. V6 ve V8 daha pürüzsüz çalışır. Boxer motorlar (Subaru, Porsche) ağırlık merkezini düşürür. Inline-6 (I6) mükemmel denge sunar. Daha fazla silindir = daha pürüzsüz çalışma, ancak daha karmaşık ve pahalı.',
        icon: '/images/guides/cylinder_icon_1769194987304.png',
        examples: ['I4: Golf, Civic, Corolla', 'V6: Camry, Passat', 'V8: Mustang, M3', 'Boxer: Subaru WRX']
    },
    {
        id: 'compression',
        title: 'Sıkıştırma Oranı',
        shortDesc: 'Pistonun hava-yakıt karışımını sıkıştırma derecesi',
        fullDesc: 'Sıkıştırma oranı, pistonun en alt noktadaki hacmin en üst noktadaki hacme oranıdır. Yüksek sıkıştırma oranı (10:1 ve üzeri) daha verimli yanma ve daha iyi performans sağlar, ancak yüksek oktan yakıt gerektirir. Dizel motorlar çok yüksek sıkıştırma oranına sahiptir (16:1 - 20:1) ve bu nedenle buji olmadan ateşleme yapabilir.',
        icon: '/images/guides/compression_icon_1769195013762.png',
        examples: ['Benzinli NA: 10:1 - 13:1', 'Benzinli Turbo: 9:1 - 10:1', 'Dizel: 16:1 - 20:1']
    },
    {
        id: 'valvetrain',
        title: 'Supap Sistemi (SOHC/DOHC)',
        shortDesc: 'Supapları kontrol eden mekanizma',
        fullDesc: 'SOHC (Single Overhead Camshaft): Tek eksantrik mil, genellikle 2 supap/silindir. Basit ve ekonomik. DOHC (Double Overhead Camshaft): Çift eksantrik mil, genellikle 4 supap/silindir. Daha iyi havalandırma, daha yüksek performans. VVT (Variable Valve Timing) teknolojisi ile supap zamanlaması optimize edilir.',
        icon: '/images/guides/horsepower_icon_1769194949851.png',
        examples: ['SOHC: Eski motorlar', 'DOHC: Modern performans motorları', 'DOHC VVT-i: Toyota']
    },
    {
        id: 'aspiration',
        title: 'Emme Tipi (NA/Turbo/Supercharger)',
        shortDesc: 'Motorun hava alma yöntemi',
        fullDesc: 'Naturally Aspirated (NA): Atmosferik basınçla hava alır, doğrusal güç, keskin gaz tepkisi. Turbocharged: Egzoz gazıyla çalışan türbin, yüksek güç, turbo lag olabilir. Supercharged: Motor tarafından mekanik olarak çalıştırılan kompresör, anlık tepki, sürekli güç tüketimi. Twin-Turbo: İki turbo, daha geniş devir aralığında güç.',
        icon: '/images/guides/turbo_icon_1769194973896.png',
        examples: ['NA: Honda Civic Type R (eski)', 'Turbo: Golf GTI', 'Supercharged: Jaguar V6']
    },
    {
        id: 'fueltype',
        title: 'Yakıt Tipi',
        shortDesc: 'Motorun kullandığı yakıt türü',
        fullDesc: 'Benzin: Yüksek devir, keskin tepki, daha temiz egzoz. Dizel: Yüksek tork, ekonomik, dayanıklı, ancak daha gürültülü. Hybrid: Benzin + elektrik motor, maksimum verimlilik. Plug-in Hybrid (PHEV): Şarj edilebilir batarya, elektrikli menzil. Full Electric (EV): Sadece elektrik, sıfır emisyon, anlık tork.',
        icon: '/images/guides/displacement_icon_1769195000365.png',
        examples: ['Benzin: 95/98 oktan', 'Dizel: Motorin', 'Hybrid: Toyota Prius', 'EV: Tesla Model 3']
    }
];

export default function EngineTermsGuide() {
    const [selectedTerm, setSelectedTerm] = useState<string | null>(null);

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
                            Motor Terimleri Sözlüğü
                        </h1>
                        <p className="text-lg text-text-muted leading-relaxed">
                            Beygir gücü, tork, turbo, silindir hacmi... Teknik terimlerin ne anlama geldiğini
                            ve aracınızın performansını nasıl etkilediğini keşfedin.
                        </p>
                    </div>
                </div>
            </div>

            {/* Main Content */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                {/* Terms Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {ENGINE_TERMS.map((term) => (
                        <div
                            key={term.id}
                            onClick={() => setSelectedTerm(term.id)}
                            className="group cursor-pointer bg-white rounded-3xl p-6 border border-border hover:shadow-card hover:border-primary/30 transition-all duration-300"
                        >
                            {/* Icon */}
                            <div className="w-20 h-20 mx-auto mb-4 rounded-2xl bg-slate-50 flex items-center justify-center overflow-hidden">
                                <img
                                    src={term.icon}
                                    alt={term.title}
                                    className="w-16 h-16 object-contain"
                                />
                            </div>

                            {/* Title */}
                            <h3 className="text-lg font-bold text-text-main mb-2 text-center group-hover:text-primary transition-colors">
                                {term.title}
                            </h3>

                            {/* Short Description */}
                            <p className="text-sm text-text-muted text-center leading-relaxed mb-4">
                                {term.shortDesc}
                            </p>

                            {/* Read More */}
                            <div className="flex items-center justify-center text-sm font-semibold text-primary">
                                Detayları Gör
                                <svg className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                </svg>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Detail Modal */}
                {selectedTerm && (() => {
                    const term = ENGINE_TERMS.find(t => t.id === selectedTerm);
                    if (!term) return null;

                    return (
                        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
                            <div className="bg-white rounded-3xl max-w-3xl w-full max-h-[90vh] overflow-y-auto shadow-2xl">
                                {/* Header */}
                                <div className="bg-gradient-to-r from-slate-700 to-slate-900 p-8 text-white sticky top-0 z-10">
                                    <div className="flex items-start justify-between">
                                        <div className="flex items-start gap-4">
                                            <div className="w-16 h-16 bg-white rounded-2xl flex items-center justify-center p-2">
                                                <img
                                                    src={term.icon}
                                                    alt={term.title}
                                                    className="w-full h-full object-contain"
                                                />
                                            </div>
                                            <div>
                                                <div className="text-sm font-medium opacity-90 mb-1">Motor Terimi</div>
                                                <h3 className="text-2xl font-bold mb-2">{term.title}</h3>
                                                <p className="text-white/90 text-sm">{term.shortDesc}</p>
                                            </div>
                                        </div>
                                        <button
                                            onClick={() => setSelectedTerm(null)}
                                            className="text-white/80 hover:text-white transition-colors p-2 hover:bg-white/10 rounded-xl"
                                        >
                                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>

                                {/* Content */}
                                <div className="p-8 space-y-6">
                                    {/* Full Description */}
                                    <div>
                                        <h4 className="text-lg font-semibold text-text-main mb-3">Detaylı Açıklama</h4>
                                        <p className="text-text-muted leading-relaxed">
                                            {term.fullDesc}
                                        </p>
                                    </div>

                                    {/* Examples */}
                                    {term.examples && term.examples.length > 0 && (
                                        <div>
                                            <h4 className="text-lg font-semibold text-text-main mb-3">Örnekler</h4>
                                            <div className="grid gap-2">
                                                {term.examples.map((example, idx) => (
                                                    <div
                                                        key={idx}
                                                        className="bg-slate-50 rounded-xl px-4 py-3 border border-slate-200 text-sm text-text-main"
                                                    >
                                                        {example}
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}

                                    {/* Close Button */}
                                    <div className="flex justify-center pt-4">
                                        <button
                                            onClick={() => setSelectedTerm(null)}
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
            </div>
        </div>
    );
}
