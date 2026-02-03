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
        title: 'Motor Hacmi',
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
        title: 'Emme Tipi',
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
                            className="inline-flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors mb-8 bg-white/5 hover:bg-white/10 px-4 py-2 rounded-full border border-white/5 w-fit"
                        >
                            <ArrowLeft className="w-4 h-4" />
                            Rehberlere Dön
                        </Link>

                        <div className="max-w-3xl">
                            <h1 className="text-4xl sm:text-6xl font-black text-white mb-6 drop-shadow-2xl">
                                Motor Terimleri Sözlüğü
                            </h1>
                            <p className="text-xl text-slate-300 leading-relaxed font-light">
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
                                className="group cursor-pointer bg-[#1e293b]/60 backdrop-blur-md rounded-3xl p-6 border border-white/10 hover:shadow-2xl hover:border-primary/50 hover:bg-[#1e293b]/80 transition-all duration-300 flex flex-col items-center text-center relative overflow-hidden"
                            >
                                <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>

                                {/* Icon */}
                                <div className="relative z-10 w-24 h-24 mx-auto mb-6 rounded-full bg-white/5 border border-white/10 flex items-center justify-center overflow-hidden group-hover:scale-110 transition-transform duration-300 shadow-lg">
                                    <img
                                        src={term.icon}
                                        alt={term.title}
                                        className="w-full h-full object-cover"
                                    />
                                </div>

                                {/* Title */}
                                <h3 className="relative z-10 text-xl font-bold text-white mb-3 group-hover:text-primary transition-colors">
                                    {term.title}
                                </h3>

                                {/* Short Description */}
                                <p className="relative z-10 text-sm text-slate-400 leading-relaxed mb-6 flex-grow">
                                    {term.shortDesc}
                                </p>

                                {/* Read More */}
                                <div className="relative z-10 flex items-center justify-center text-sm font-bold text-primary uppercase tracking-wider mt-auto">
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
                            <div className="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-4">
                                <div className="bg-[#1e293b] rounded-3xl max-w-3xl w-full max-h-[90vh] overflow-y-auto shadow-2xl border border-white/10 relative">
                                    {/* Close Button Absolute */}
                                    <button
                                        onClick={() => setSelectedTerm(null)}
                                        className="absolute top-4 right-4 z-20 bg-black/40 hover:bg-black/60 text-white p-2 rounded-full transition-colors backdrop-blur-sm border border-white/10"
                                    >
                                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                        </svg>
                                    </button>

                                    {/* Header */}
                                    <div className="bg-gradient-to-b from-primary/20 to-[#1e293b] p-10 text-white flex flex-col items-center text-center border-b border-white/5">
                                        <div className="w-28 h-28 bg-white/10 rounded-full flex items-center justify-center overflow-hidden mb-6 backdrop-blur-md border border-white/10 shadow-2xl">
                                            <img
                                                src={term.icon}
                                                alt={term.title}
                                                className="w-full h-full object-cover"
                                            />
                                        </div>
                                        <div className="text-sm font-bold opacity-70 uppercase tracking-widest mb-2 text-primary">Teknik Terim</div>
                                        <h3 className="text-3xl sm:text-4xl font-black mb-4">{term.title}</h3>
                                        <p className="text-slate-300 text-lg max-w-lg">{term.shortDesc}</p>
                                    </div>

                                    {/* Content */}
                                    <div className="p-8 sm:p-10 space-y-8">
                                        {/* Full Description */}
                                        <div className="bg-black/20 rounded-2xl p-6 border border-white/5">
                                            <h4 className="text-lg font-bold text-white mb-3">Detaylı Açıklama</h4>
                                            <p className="text-slate-300 leading-relaxed text-lg">
                                                {term.fullDesc}
                                            </p>
                                        </div>

                                        {/* Examples */}
                                        {term.examples && term.examples.length > 0 && (
                                            <div>
                                                <h4 className="text-lg font-bold text-white mb-4 flex items-center gap-2">
                                                    Araba Dünyasından Örnekler
                                                </h4>
                                                <div className="grid gap-3">
                                                    {term.examples.map((example, idx) => (
                                                        <div
                                                            key={idx}
                                                            className="bg-white/5 rounded-xl px-5 py-4 border border-white/5 text-slate-200 font-medium flex items-center gap-3"
                                                        >
                                                            <span className="w-2 h-2 rounded-full bg-primary"></span>
                                                            {example}
                                                        </div>
                                                    ))}
                                                </div>
                                            </div>
                                        )}

                                        {/* Close Button Bottom */}
                                        <div className="flex justify-center pt-4">
                                            <button
                                                onClick={() => setSelectedTerm(null)}
                                                className="px-10 py-3 bg-white text-slate-900 hover:bg-slate-200 font-bold rounded-xl transition-colors shadow-lg"
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
        </div>
    );
}
