export interface DSGVariant {
    code: string;
    name: string;
    gears: number | string;
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
    productionYears: string;
    status: string;
    history: string;
}

export const DSG_VARIANTS: DSGVariant[] = [
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
        borderColor: 'border-red-500/50',
        productionYears: '2008 - Günümüz',
        status: 'Hala Üretimde',
        history: 'Golf Mk6 ve Polo döneminde yaygınlaştı. 250 Nm tork limiti olan küçük hacimli motorlarda (1.0 TSI, 1.5 TSI, 1.6 TDI) kullanılmaya devam ediyor. "Kuru" kavrama olduğu için ısınma ve mekatronik arızalarıyla nam salmıştır ancak revizyonlarla (Gen 2, Gen 3) iyileştirildi.'
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
        borderColor: 'border-emerald-500/50',
        productionYears: '2003 - ~2017/2018',
        status: 'Kısmen Devam Ediyor / Yerini DQ381 Aldı',
        history: 'İlk olarak Golf Mk4 R32 ve Audi TT 3.2 ile piyasaya sürüldü. Dünyanın ilk seri üretim çift kavramalı şanzımanıdır. 400 Nm tork dayanımı ile efsaneleşmiş, en dayanıklı DSG\'lerden biridir. Artık ana modellerde kullanılmıyor ancak bazı spesifik pazarlarda devam ediyor.'
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
        borderColor: 'border-blue-500/50',
        productionYears: '2017 - Günümüz',
        status: 'Hala Üretimde',
        history: 'Golf 7.5 ve Arteon ile tanıtıldı. DQ250\'nin yerini aldı. CO2 emisyonlarını düşürmek için sürtünmesi azaltılmış, yağ pompası ve yatakları geliştirilmiş versiyondur. Şu anki Golf 8, Leon, Octavia vb. modellerin 2.0 TDI veya yüksek güçlü benzinli versiyonlarında standarttır.'
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
        borderColor: 'border-amber-500/50',
        productionYears: '2009/2010 - Günümüz',
        status: 'Hala Üretimde',
        history: 'VW Transporter T5 GP ve Tiguan ile başladı. Audi RS3, TT RS ve ticari araçlar gibi yüksek tork (600 Nm+) gerektiren araçlarda kullanılıyor.'
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
        borderColor: 'border-purple-500/50',
        productionYears: '2008 - ~2015/2016',
        status: 'Üretimden Kalktı / Yerini DL382 Aldı',
        history: 'Audi Q5 ve A4 B8 kasalar ile başladı. Yüksek tork (550-600 Nm) dayanımı vardı ancak mekatronik kart arızaları sık görüldüğü için daha verimli olan DL382 ile değiştirildi. Çok yüksek torklu (RS) modellerde ise yerini ZF8HP\'ye bıraktı.'
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
        borderColor: 'border-cyan-500/50',
        productionYears: '2014/2015 - Günümüz',
        status: 'Hala Üretimde',
        history: 'Audi A6 Ultra ve A4 B9 ile başladı. Audi\'nin güncel "Ultra" teknolojili, önden çekiş veya Quattro modellerinde (A4, A5, A6 2.0 motorlar) kullanılan standart şanzımandır.'
    },
    {
        code: 'ZF8HP',
        name: '8-İleri Tork Konvertörlü (Tiptronic)',
        gears: 8,
        clutchType: 'torque_converter',
        clutchTypeLabel: 'Tork Konvertör',
        maxTorque: 1000,
        description: 'Audi buna da Tiptronic der ama aslında efsanevi ZF şanzımanıdır. DSG değildir, tam otomatiktir. Dünyanın en iyi otomatik şanzımanı kabul edilir.',
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
        borderColor: 'border-yellow-500/50',
        productionYears: '2008 - Günümüz',
        status: 'Hala Üretimde (3. ve 4. Jenerasyon)',
        history: 'İlk olarak BMW 7 Serisi F01\'de kullanıldı. Otomotiv dünyasının en başarılı şanzımanı kabul edilir. VAG grubunda Audi RS6, RS7, Q7, Q8, VW Amarok ve VW Touareg gibi yüksek torklu araçlarda kullanılır. BMW\'nin neredeyse tüm gamında bulunur.'
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
        borderColor: 'border-teal-500/50',
        productionYears: '2014 - Günümüz',
        status: 'Hala Üretimde',
        history: 'Golf GTE ve Audi A3 e-tron ile başladı. Elektrik motorunun şanzıman içine entegre edildiği (P2 hibrit mimarisi) modeldir. Passat GTE, Superb iV gibi PHEV araçlarda kullanılır.'
    },
    {
        code: 'TIPTRONIC',
        name: 'TORK KONVERTÖRLÜ TAM OTOMATİK',
        gears: '6 / 8',
        clutchType: 'torque_converter',
        clutchTypeLabel: 'Tork Konvertör',
        maxTorque: 1000,
        description: 'VAG grubunun tork konvertörlü şanzıman ailesidir. Motor ile şanzıman arasında mekanik kavrama yerine hidrolik güç aktarımı kullanılır. Bu sayede sarsıntısız geçişler ve yüksek tork dayanımı sağlar.',
        chronicProblems: [
            'Solenoid Valf Gövdesi (Beyin): Yüksek kilometrede tıkanıklık sonucu vuruntulu geçişler.',
            'Tork Konvertörü: "Lock-up" (Kilitleme) balatasının aşınması sonucu devir dalgalanması.'
        ],
        maintenanceTips: [
            '"Ömürlük yağ" efsanesine inanmayın; her 60.000 - 80.000 km\'de bir şanzıman yağı ve filtresi değişmelidir.',
            'Park ederken şanzımana yük binmemesi için önce N, sonra El Freni, en son P konumuna alınmalıdır.'
        ],
        clutchInterval: '300.000 km+ (Revizyon Öncesi)',
        smartTip: 'Sıkışık trafikte ısınma sorunu yaşatmaz. Tork konvertörü sayesinde yokuş kalkışlarında geri kaydırma yapmaz ve \'kavrama bitmesi\' gibi riskler, kuru kavramalı şanzımanlara göre çok daha düşüktür.',
        color: 'from-slate-600 to-gray-600',
        borderColor: 'border-slate-500/50',
        productionYears: '2003 - Günümüz',
        status: 'Hala Kullanımda (ZF / Aisin)',
        history: 'VAG grubunun tam otomatik şanzımanlara verdiği isimdir. Üç ana varyasyonu vardır:\n1. Aisin 6 İleri (09G): 2003\'ten günümüze eski Polo/Golf ve bazı pazarlarda kullanılır.\n2. Aisin 8 İleri (AQ300/450): 2017\'den günümüze Golf 8, Tiguan vb. araçların (özellikle ABD pazarı) bazı motorlarında DSG yerine tercih edilir.\n3. ZF 8 İleri (8HP): Audi RS, Touareg ve Amarok gibi yüksek torklu araçlarda kullanılır.'
    }
];

