import { useState } from 'react';
import { ArrowLeft, AlertTriangle, Wrench, Lightbulb, ChevronDown, ChevronRight } from 'lucide-react';
import { Link } from 'react-router-dom';

interface DSGVariant {
    code: string;
    name: string;
    gears: number;
    clutchType: string;
    clutchTypeLabel: string;
    maxTorque: number;
    description: string;
    chronicProblems: string[];
    maintenanceTips: string[];
    clutchInterval: string;
    smartTip: string;
    color: string;
    borderColor: string;
}

const DSG_VARIANTS: DSGVariant[] = [
    {
        code: 'DQ200',
        name: '7 İleri Kuru Kavrama S-Tronic/DSG',
        gears: 7,
        clutchType: 'dry',
        clutchTypeLabel: 'Kuru Kavrama',
        maxTorque: 250,
        description: 'VAG grubunun en çok tartışılan, en yaygın ve düşük torklu motorlarda (250 Nm torka kadar) kullandığı şanzımandır. Çift kavramalı, 7 ileri vitesli, "kuru" tip bir şanzımandır.',
        chronicProblems: [
            'Mekatronik Arızası: Şanzımanın beyni ve hidrolik ünitesidir.',
            'Kavrama Titremesi (Silkeleme): Özellikle 1. vitesten 2\'ye geçerken.',
            'Isınma: Yoğun dur-kalk trafikte "şanzıman aşırı ısındı" uyarısı.'
        ],
        maintenanceTips: [
            'Trafikte manuel moda alıp 1. viteste sabit gitmek ömrü uzatır.',
            'Yokuşlarda aracı "Auto Hold" veya frene basarak sabit tutun.'
        ],
        clutchInterval: '60.000 - 120.000 km',
        smartTip: 'Yoğun dur-kalk trafikte şanzımanın ısınmaması için aracı sık sık "N" konumuna alın.',
        color: 'from-red-600 to-rose-600',
        borderColor: 'border-red-500/50'
    },
    {
        code: 'DQ250',
        name: '6 İleri Yağlı Kavrama S-Tronic/DSG',
        gears: 6,
        clutchType: 'wet',
        clutchTypeLabel: 'Yağlı Kavrama',
        maxTorque: 400,
        description: 'Daha güçlü motorlarda (2.0 TDI, 2.0 TFSI vb.) kullanılan, DQ200\'e göre çok daha dayanıklı olan abisidir. 6 ileri vitesli, "yağlı" tip şanzımandır.',
        chronicProblems: [
            'Volant (Flywheel) Sesi: Çift kütleli volant zamanla boşluk yapabilir.',
            'Mekatronik Solenoidleri: Vites geçişlerinde vuruntu yaparsa solenoid kirli olabilir.'
        ],
        maintenanceTips: [
            'Her 60.000 km\'de bir şanzıman yağı ve filtresi mutlaka değişmeli.',
            'Motor ve şanzıman yağı ısınmadan dip gaz yapılmamalı.'
        ],
        clutchInterval: '150.000 km+',
        smartTip: 'Kuru kavramaya göre çok daha dayanıklıdır. 60.000 km\'de bir yağ değişimi kritik öneme sahiptir.',
        color: 'from-emerald-600 to-green-600',
        borderColor: 'border-emerald-500/50'
    },
    {
        code: 'DQ381',
        name: '7 İleri Yağlı Kavrama S-Tronic/DSG',
        gears: 7,
        clutchType: 'wet',
        clutchTypeLabel: 'Yağlı Kavrama',
        maxTorque: 430,
        description: 'DQ250\'nin yerini alan, güncel ve optimize edilmiş versiyondur (2017 sonrası). 7 ileri, yağlı kavramadır. 420-430 Nm torklara dayanabilir.',
        chronicProblems: [
            'Erken dönem üretimlerinde yardımcı hidrolik pompa arızaları.',
            'Yazılımsal kararsızlıklar, genelde güncelleme ile çözülür.'
        ],
        maintenanceTips: [
            '60.000 km yağ değişimi hayati önem taşır.',
            'Akü zayıflarsa şanzıman hataları verebilir.'
        ],
        clutchInterval: '150.000 km+',
        smartTip: 'DQ250\'nin geliştirilmiş versiyonudur. Düzenli yağ değişimi ile çok uzun ömürlüdür.',
        color: 'from-blue-600 to-indigo-600',
        borderColor: 'border-blue-500/50'
    },
    {
        code: 'DQ500',
        name: '7-İleri Yağlı Kavrama - Ağır Hizmet Tipi',
        gears: 7,
        clutchType: 'wet',
        clutchTypeLabel: 'Ağır Hizmet Yağlı',
        maxTorque: 600,
        description: 'VAG grubunun "tank" lakaplı, en sağlam şanzımanıdır. Ticari araçlarda (Transporter) ve RS modellerinde kullanılır. Hem yük taşımaya hem pist performansına dayanır.',
        chronicProblems: [
            'En Sorunsuzu: Grubun en az arıza yapan şanzımanıdır.',
            'DPF Rejenerasyonu: Dizel modellerde rölantide kararsız kalabilir (arıza değildir).',
            'Volant Sesi: Yüksek km ticari araçlarda volant sesi yapabilir.'
        ],
        maintenanceTips: [
            '60.000 km periyodunda yağ değişimi aksatılmamalıdır.',
            'Ağır yük çekiliyorsa yağ değişim aralığı 40.000 km\'ye çekilmelidir.'
        ],
        clutchInterval: '200.000 km+',
        smartTip: 'Grubun en sağlam şanzımanıdır. Yağ değişimi ile 200.000 km üzeri sorunsuz kullanım mümkündür.',
        color: 'from-amber-600 to-orange-600',
        borderColor: 'border-amber-500/50'
    },
    {
        code: 'DL501',
        name: '7-İleri Yağlı Kavrama - Boyuna (0B5)',
        gears: 7,
        clutchType: 'wet',
        clutchTypeLabel: 'Çift Hazneli Yağlı',
        maxTorque: 550,
        description: 'Audi\'nin (A4, A5, Q5) boyuna motorlu araçların şanzımanıdır. İki ayrı yağ haznesine (mekatronik ve dişli grubu) sahiptir.',
        chronicProblems: [
            'Mekatronik Kart Arızası: Isınmadan dolayı sık arızalanır.',
            'Vites Pozisyon Sensörü: Sensör arızası nedeniyle şanzıman vites şaşırabilir.'
        ],
        maintenanceTips: [
            'Çift Yağ Değişimi: Hem hidrolik hem dişli tarafın yağı değişmelidir (Çok Önemli).',
            'Mekatronik tamir takımları ile revizyon daha ekonomiktir.'
        ],
        clutchInterval: '100.000 - 150.000 km',
        smartTip: 'ÇİFT yağ haznesi vardır. Servis sadece dış yağı değiştirirse yetersizdir. İkisini de değiştirtin.',
        color: 'from-purple-600 to-fuchsia-600',
        borderColor: 'border-purple-500/50'
    },
    {
        code: 'DL382',
        name: '7-İleri Yağlı Kavrama - Yeni Boyuna',
        gears: 7,
        clutchType: 'wet',
        clutchTypeLabel: 'Yeni Nesil Yağlı',
        maxTorque: 450,
        description: 'Sorunlu DL501\'in yerini alan modern şanzıman. Quattro Ultra ile uyumlu, sürtünmesi azaltılmış, ekonomi odaklıdır.',
        chronicProblems: [
            'Basınç Akümülatörü: Nadiren basınç kaybı yaşanabilir.',
            'DL501\'e göre çok daha güvenilirdir.'
        ],
        maintenanceTips: [
            'Sık dur-kalk yapılan yerlerde Start-Stop\'u kapatmak yağ pompası ömrünü uzatır.',
            '60.000 km yağ değişimi kuralı geçerlidir.'
        ],
        clutchInterval: '150.000 km+',
        smartTip: 'DL501\'in sorunlarını çözen yeni nesil. Quattro Ultra ile uyumlu ve güvenilirdir.',
        color: 'from-cyan-600 to-sky-600',
        borderColor: 'border-cyan-500/50'
    },
    {
        code: 'ZF8HP',
        name: '8-İleri Tork Konvertörlü (Tiptronic)',
        gears: 8,
        clutchType: 'torque_converter',
        clutchTypeLabel: 'Tork Konvertör',
        maxTorque: 1000,
        description: 'Audi buna da Tiptronic der ama aslında efsanevi ZF şanzımandır. DSG değildir, tam otomatiktir. Dünyanın en iyi otomatik şanzımanı kabul edilir.',
        chronicProblems: [
            'Neredeyse Yok: Çok sağlamdır.',
            'Solenoidler: 250k+ km\'de vuruntu yaparsa solenoid değişimi gerekebilir.'
        ],
        maintenanceTips: [
            '"Ömürlük Yağ" yalanına inanma. 80-100k km\'de yağ ve filtre değişirse 500k km yapar.',
            'ZF Lifeguard Fluid kullanılmalıdır.'
        ],
        clutchInterval: '300.000 km+ (Revizyon)',
        smartTip: '"Ömürlük yağ" iddiasına inanmayın, 80-100k\'da yağ değişimi ile yarım milyon km dayanır.',
        color: 'from-yellow-600 to-amber-600',
        borderColor: 'border-yellow-500/50'
    },
    {
        code: 'DQ400e',
        name: '6-İleri Yağlı Kavrama - Hibrit (PHEV)',
        gears: 6,
        clutchType: 'wet',
        clutchTypeLabel: 'Hibrit Yağlı',
        maxTorque: 400,
        description: 'GTE ve e-tron modelleri için tasarlanmıştır. DSG içine entegre elektrik motoru ve üçüncü bir kavrama bulunur.',
        chronicProblems: [
            'Sızdırmazlık: Elektrik motoru ile şanzıman arasında sızıntı (nadir).',
            'Mekatronik: Karmaşık yapısı nedeniyle arızalar maliyetlidir.'
        ],
        maintenanceTips: [
            'Bakımı sadece yetkili veya hibrit sertifikalı servisler yapmalıdır.',
            'Yüksek voltajlı sistemle iç içedir, dikkat gerektirir.'
        ],
        clutchInterval: '150.000 km+',
        smartTip: 'Hibrit sistemle entegre çalıştığı için kavrama daha az aşınır. Uzman servis şarttır.',
        color: 'from-teal-600 to-emerald-600',
        borderColor: 'border-teal-500/50'
    },
    {
        code: 'TIPTRONIC',
        name: 'Eski Tip Tork Konvertörlü',
        gears: 0,
        clutchType: 'torque_converter',
        clutchTypeLabel: 'Tork Konvertör',
        maxTorque: 0,
        description: 'Eski kasalarda görülen klasik otomatik şanzıman. Mekanik kavrama plakası yoktur. Vites geçişleri pürüzsüzdür.',
        chronicProblems: [
            'Türbin (Konvertör) Arızası: Devir dalgalanması veya titreme.',
            'Valf gövdesi aşınırsa vites geçişlerinde sert vurma.'
        ],
        maintenanceTips: [
            'Her 60.000-80.000 km\'de bir yağ değişimi şanzıman ömrünü ikiye katlar.',
            'Soğuk motorla ani gaz açmaktan kaçının.'
        ],
        clutchInterval: '200.000 km+',
        smartTip: 'Fiziksel kavrama balatası olmadığı için çok uzun ömürlüdür.',
        color: 'from-slate-600 to-gray-600',
        borderColor: 'border-slate-500/50'
    }
];

export default function DSGVariantsGuide() {
    const [expandedCode, setExpandedCode] = useState<string | null>(null);

    const handleCardClick = (code: string) => {
        setExpandedCode(expandedCode === code ? null : code);
    };

    const expandedVariant = DSG_VARIANTS.find(v => v.code === expandedCode);

    return (
        <div className="relative min-h-screen font-sans text-slate-200">
            {/* Background */}
            <div className="fixed inset-0 z-0">
                <div className="absolute inset-0 bg-cover bg-center" style={{ backgroundImage: 'url(/hero-2.jpg)' }} />
                <div className="absolute inset-0 bg-slate-900/90 backdrop-blur-sm" />
            </div>

            <div className="relative z-10">
                {/* Header */}
                <div className="relative border-b border-white/10 bg-black/20 backdrop-blur-sm">
                    <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 pt-28">
                        <Link
                            to="/guides/transmission"
                            className="inline-flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors mb-6 bg-white/5 hover:bg-white/10 px-4 py-2 rounded-full w-fit border border-white/5"
                        >
                            <ArrowLeft className="w-4 h-4" />
                            Şanzıman Rehberine Dön
                        </Link>

                        <h1 className="text-3xl sm:text-4xl font-black text-white mb-3">
                            VAG DSG/S-Tronic Detayları
                        </h1>
                        <p className="text-lg text-slate-300 max-w-2xl">
                            DQ200, DQ250, DQ381 ve Tiptronic - Kronik sorunlar ve bakım önerileri.
                        </p>
                        <div className="absolute top-32 right-8 hidden md:flex items-center gap-4 opacity-70 hover:opacity-100 transition-all">
                            <img src="/images/vag/logo1.png" className="h-6 w-auto object-contain grayscale hover:grayscale-0 transition-all" alt="Seat" />
                            <img src="/images/vag/logo2.png" className="h-6 w-auto object-contain grayscale hover:grayscale-0 transition-all" alt="Audi" />
                            <img src="/images/vag/logo3.png" className="h-6 w-auto object-contain grayscale hover:grayscale-0 transition-all" alt="Skoda" />
                            <img src="/images/vag/logo4.png" className="h-6 w-auto object-contain grayscale hover:grayscale-0 transition-all" alt="VW" />
                        </div>
                    </div>
                </div>

                {/* Content */}
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    {/* Warning Box - Compact */}
                    <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-4 mb-6">
                        <div className="flex items-center gap-3">
                            <AlertTriangle className="w-5 h-5 text-slate-400 flex-shrink-0" />
                            <p className="text-slate-300 text-sm">
                                Bu bilgiler genel öneridir. Her aracın durumu farklıdır. Ciddi sorunlarda yetkili servise başvurun.
                            </p>
                        </div>
                    </div>

                    {/* DSG Variants - Grid Layout */}
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-6">
                        {DSG_VARIANTS.map((variant) => (
                            <button
                                key={variant.code}
                                onClick={() => handleCardClick(variant.code)}
                                className={`
                                    relative p-4 rounded-2xl text-left transition-all duration-300 border group h-full
                                    ${expandedCode === variant.code
                                        ? `bg-slate-800 border-2 ${variant.borderColor.replace('/50', '')} shadow-[0_0_20px_rgba(0,0,0,0.3)] scale-105 z-10`
                                        : 'bg-[#1e293b]/80 border-slate-700 hover:border-slate-500 hover:bg-slate-800/80'}
                                `}
                            >
                                {/* Active Indicator Glow */}
                                {expandedCode === variant.code && (
                                    <div className={`absolute inset-0 rounded-2xl bg-gradient-to-br ${variant.color} opacity-10 blur-xl`} />
                                )}

                                <div className="relative z-10 flex flex-col h-full justify-between min-h-[100px]">
                                    <div className="flex items-center justify-between mb-2">
                                        <h3 className={`text-xl font-black ${expandedCode === variant.code ? 'text-white' : 'text-slate-200'}`}>
                                            {variant.code}
                                        </h3>
                                        {expandedCode === variant.code ? (
                                            <div className="bg-white/10 p-1 rounded-full">
                                                <ChevronDown className="w-4 h-4 text-white" />
                                            </div>
                                        ) : (
                                            <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-slate-400" />
                                        )}
                                    </div>

                                    <div className="mt-auto">
                                        <div className="flex items-center gap-2 mb-2">
                                            <div className={`h-1.5 w-1.5 rounded-full bg-gradient-to-r ${variant.color}`} />
                                            <span className={`text-xs font-bold uppercase tracking-wider ${expandedCode === variant.code ? 'text-white/80' : 'text-slate-500 group-hover:text-slate-400'}`}>
                                                {variant.clutchTypeLabel}
                                            </span>
                                        </div>

                                        <div className={`text-xs font-medium transition-colors ${expandedCode === variant.code ? 'text-blue-300' : 'text-slate-600 group-hover:text-blue-400'}`}>
                                            İncelemek için tıkla
                                        </div>
                                    </div>
                                </div>
                            </button>
                        ))}
                    </div>

                    {/* Expanded Detail Panel */}
                    {expandedVariant && (
                        <div className="relative overflow-hidden rounded-3xl border border-slate-700 shadow-2xl animate-fadeIn mt-6">
                            {/* Dynamic Background Glow */}
                            <div className={`absolute top-0 right-0 w-[500px] h-[500px] bg-gradient-to-br ${expandedVariant.color} opacity-20 blur-[100px] rounded-full pointer-events-none -translate-y-1/2 translate-x-1/2`} />
                            <div className="absolute bottom-0 left-0 w-[300px] h-[300px] bg-blue-500/10 blur-[80px] rounded-full pointer-events-none translate-y-1/2 -translate-x-1/2" />

                            <div className="relative bg-[#0f172a]/90 backdrop-blur-xl p-8">
                                {/* Header Section */}
                                <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-6 mb-10 pb-8 border-b border-slate-700/50">
                                    <div>
                                        <div className="flex items-center gap-3 mb-2">
                                            <h2 className="text-4xl md:text-5xl font-black text-white tracking-tight">
                                                {expandedVariant.code}
                                            </h2>
                                            <span className={`px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider border ${expandedVariant.borderColor} bg-white/5 text-white/80`}>
                                                {expandedVariant.clutchTypeLabel}
                                            </span>
                                        </div>
                                        <p className="text-lg text-slate-400 font-medium">{expandedVariant.name}</p>
                                    </div>

                                    {/* Quick Stats */}
                                    <div className="flex gap-4">
                                        <div className="px-5 py-3 rounded-2xl bg-slate-800/50 border border-slate-700/50 text-center min-w-[100px]">
                                            <div className="text-xs text-slate-500 uppercase font-bold tracking-wider mb-1">Maks Tork</div>
                                            <div className="text-xl font-black text-white">
                                                {expandedVariant.maxTorque > 0 ? (
                                                    <span>{expandedVariant.maxTorque} <span className="text-sm font-normal text-slate-400">Nm</span></span>
                                                ) : 'Yüksek'}
                                            </div>
                                        </div>
                                        <div className="px-5 py-3 rounded-2xl bg-slate-800/50 border border-slate-700/50 text-center min-w-[100px]">
                                            <div className="text-xs text-slate-500 uppercase font-bold tracking-wider mb-1">Vites</div>
                                            <div className="text-xl font-black text-white">
                                                {expandedVariant.gears > 0 ? expandedVariant.gears : 'Auto'}
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                {/* Main Grid */}
                                <div className="grid lg:grid-cols-3 gap-8">
                                    {/* Left Column: Description & Smart Tip */}
                                    <div className="lg:col-span-1 space-y-6">
                                        <div className="group relative p-6 rounded-3xl bg-slate-800/30 border border-slate-700/50 hover:bg-slate-800/50 transition-colors">
                                            <h3 className="text-lg font-bold text-white mb-3 flex items-center gap-2">
                                                <div className="w-1 h-6 bg-blue-500 rounded-full" />
                                                Teknik Bakış
                                            </h3>
                                            <p className="text-slate-300 leading-relaxed text-sm">
                                                {expandedVariant.description}
                                            </p>
                                        </div>

                                        <div className="relative overflow-hidden p-6 rounded-3xl bg-gradient-to-br from-blue-900/20 to-indigo-900/20 border border-blue-500/20 group">
                                            <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                                                <Lightbulb className="w-24 h-24 text-blue-400" />
                                            </div>
                                            <div className="relative z-10">
                                                <h3 className="text-lg font-bold text-blue-300 mb-3 flex items-center gap-2">
                                                    <Lightbulb className="w-5 h-5" />
                                                    Uzman İpucu
                                                </h3>
                                                <p className="text-blue-100/90 text-sm leading-relaxed font-medium">
                                                    {expandedVariant.smartTip}
                                                </p>
                                            </div>
                                        </div>
                                    </div>

                                    {/* Right Column: Problems & Maintenance */}
                                    <div className="lg:col-span-2 grid sm:grid-cols-2 gap-6">
                                        {/* Chronic Problems */}
                                        <div className="space-y-4">
                                            <div className="flex items-center gap-3 mb-2">
                                                <div className="p-2 rounded-xl bg-red-500/10 text-red-400">
                                                    <AlertTriangle className="w-5 h-5" />
                                                </div>
                                                <h3 className="text-lg font-bold text-slate-200">Bilinen Durumlar</h3>
                                            </div>

                                            <div className="space-y-3">
                                                {expandedVariant.chronicProblems.map((problem, idx) => (
                                                    <div key={idx} className="p-4 rounded-2xl bg-[#1e1b2e]/50 border border-red-500/10 hover:border-red-500/30 transition-colors group">
                                                        <div className="flex gap-3">
                                                            <div className="mt-1.5 w-1.5 h-1.5 rounded-full bg-red-500/40 group-hover:bg-red-500 group-hover:shadow-[0_0_8px_rgba(239,68,68,0.6)] transition-all flex-shrink-0" />
                                                            <p className="text-sm text-slate-400 group-hover:text-slate-200 transition-colors cursor-default">
                                                                {problem}
                                                            </p>
                                                        </div>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>

                                        {/* Maintenance Tips */}
                                        <div className="space-y-4">
                                            <div className="flex items-center gap-3 mb-2">
                                                <div className="p-2 rounded-xl bg-green-500/10 text-green-400">
                                                    <Wrench className="w-5 h-5" />
                                                </div>
                                                <h3 className="text-lg font-bold text-slate-200">Bakım & Koruma</h3>
                                            </div>

                                            <div className="space-y-3">
                                                {expandedVariant.maintenanceTips.map((tip, idx) => (
                                                    <div key={idx} className="p-4 rounded-2xl bg-[#1b2e25]/30 border border-green-500/10 hover:border-green-500/30 transition-colors group">
                                                        <div className="flex gap-3">
                                                            <div className="mt-1 flex-shrink-0">
                                                                <div className="w-4 h-4 rounded-full bg-green-500/20 flex items-center justify-center group-hover:bg-green-500/40 transition-colors">
                                                                    <div className="w-1.5 h-1.5 rounded-full bg-green-500" />
                                                                </div>
                                                            </div>
                                                            <p className="text-sm text-slate-400 group-hover:text-slate-200 transition-colors cursor-default">
                                                                {tip}
                                                            </p>
                                                        </div>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                {/* Life Expectancy Bar */}
                                <div className="mt-10 pt-8 border-t border-slate-700/50">
                                    <div className="bg-slate-800/30 rounded-3xl p-6 border border-slate-700/50">
                                        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                                            <div className="flex items-center gap-4">
                                                <div className="w-12 h-12 rounded-2xl bg-blue-500/10 flex items-center justify-center text-blue-400">
                                                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                                    </svg>
                                                </div>
                                                <div>
                                                    <h4 className="font-bold text-white text-lg">Beklenen Ömür</h4>
                                                    <p className="text-slate-400 text-sm">Ortalama kavrama/revizyon süresi</p>
                                                </div>
                                            </div>
                                            <div className="text-right">
                                                <div className="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-300">
                                                    {expandedVariant.clutchInterval}
                                                </div>
                                                <div className="w-full h-1.5 bg-slate-700/50 rounded-full mt-2 overflow-hidden">
                                                    <div className="h-full bg-gradient-to-r from-blue-500 to-cyan-400 w-3/4 rounded-full shadow-[0_0_10px_rgba(59,130,246,0.5)]" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Back Button */}
                    <div className="mt-8 text-center">
                        <Link
                            to="/guides/transmission"
                            className="inline-flex items-center gap-2 px-6 py-3 bg-slate-800 hover:bg-slate-700 text-white font-semibold rounded-xl transition-all border border-slate-700"
                        >
                            <ArrowLeft className="w-4 h-4" />
                            Şanzıman Rehberine Dön
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}
