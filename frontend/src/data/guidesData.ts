export interface GuideSection {
    title: string;
    content: string;
    imagePlaceholder?: string; // Hint for what image should go here
}

export interface GuideArticle {
    id: string;
    title: string;
    summary: string;
    sections: GuideSection[];
}

export const GUIDES: GuideArticle[] = [
    {
        id: 'transmission',
        title: 'Şanzıman Dünyası: Otomatik, DCT, CVT',
        summary: 'Tork konvertörlü tam otomatikler, çift kavramalı (DCT) sistemler ve sürekli değişken (CVT) şanzımanlar arasındaki farklar neler?',
        sections: [
            {
                title: 'Tork Konvertörlü (Tam Otomatik)',
                content: 'Geleneksel otomatik şanzıman olarak da bilinir. Motor ile tekerlekler arasında mekanik bir bağ yerine hidrolik bir "sıvı kavraması" (tork konvertörü) kullanır. Bu sayede vites geçişleri son derece yumuşaktır ve dur-kalk trafikte konforludur. ZF 8HP ve Aisin şanzımanlar bu teknolojinin zirvesidir.',
                imagePlaceholder: 'torque_converter_diagram'
            },
            {
                title: 'Çift Kavramalı (DCT / DSG / EDC)',
                content: 'İki adet manuel şanzımanın iç içe geçmiş hali gibidir. Biri tek (1,3,5), diğeri çift (2,4,6) vitesleri kontrol eder. Vites geçişleri milisaniyeler sürer, performans ve yakıt ekonomisi sunar. Ancak dur-kalk trafikte ısınabilir ve "kararsız" davranabilir. VW DSG ve Renault EDC en bilinenleridir.',
                imagePlaceholder: 'dual_clutch_cutaway'
            },
            {
                title: 'CVT (Sürekli Değişken)',
                content: 'Dişli yerine kasnak ve çelik kayış kullanır. Sonsuz sayıda vites oranına sahiptir. Motoru sürekli en verimli devirde tutar, bu yüzden yakıt ekonomisinde liderdir. Ancak "lastik bant etkisi" denilen, hızlanırken motor sesinin sabit kalması hissi bazı sürücüleri rahatsız edebilir. Toyota ve Nissan sıkça kullanır.',
                imagePlaceholder: 'cvt_pulley_system'
            }
        ]
    },
    {
        id: 'engine',
        title: 'Motor Terimleri Sözlüğü',
        summary: 'Beygir gücü, tork, sıkıştırma oranı... Teknik veriler ne anlama geliyor ve aracınızı nasıl etkiliyor?',
        sections: [
            {
                title: 'Beygir Gücü (HP) vs Tork (Nm)',
                content: 'Basitçe: Beygir gücü (HP) otomobilin ne kadar "hızlı gidebileceğini" (son hız), Tork (Nm) ise ne kadar "çabuk hızlanabileceğini" ve yük altında çekiş gücünü belirler. Dizel motorlar yüksek tork üretirken, Benzinli motorlar genellikle daha yüksek devir çevirip yüksek beygir üretir.',
                imagePlaceholder: 'hp_vs_torque_graph'
            },
            {
                title: 'Turbo Besleme vs Atmosferik',
                content: 'Turbo, egzoz gazını kullanarak motora daha fazla hava basar ve küçük hacimden büyük güç (ve tork) elde edilmesini sağlar. Atmosferik motorlar ise daha doğrusal güç üretir, gaz tepkileri daha keskindir ancak genellikle daha fazla yakıt tüketir ve düşük devirde cansızdır.',
                imagePlaceholder: 'turbocharger_diagram'
            },
            {
                title: 'Silindir Hacmi ve Konfigürasyon',
                content: 'Inline-4 (I4) en yaygın tiptir. V6 ve V8 motorlar daha pürüzsüz ve güçlüdür. Boxer motorlar (Subaru, Porsche) ise ağırlık merkezini aşağı çeker. Hacim arttıkça genellikle tüketim artar ama motor ömrü ve pürüzsüzlük de artar.',
                imagePlaceholder: 'engine_configurations'
            }
        ]
    },
    {
        id: 'chronicles',
        title: 'Kronikler ve Çözümleri',
        summary: 'Her güzelin bir kusuru vardır. Popüler araçların sık karşılaşılan "kronik" sorunları ve bunların nedenleri.',
        sections: [
            {
                title: 'Nedir Bu "Kronik"?',
                content: 'Bir araç modelinin üretim hatası, parça kalitesi veya tasarım tercihi nedeniyle belirli bir kilometrede neredeyse tüm kullanıcılarında çıkardığı ortak arızadır. Örneğin: E90 kasa BMW 3 serisinin kapı kollarının erimesi veya Golf 7 DSG mekatronik kart arızası.',
                imagePlaceholder: 'mechanic_checking_car'
            },
            {
                title: 'Platformumuzda Ne Yapabilirsiniz?',
                content: 'Yeni "Sorun Bildir" özelliğimizle, yaşadığınız arızayı veritabanımıza ekleyebilirsiniz. Eğer bir çözüm bulduysanız (örn: "X parçasını yan sanayi ile değiştirmeyin"), bunu da paylaşarak topluluğa katkıda bulunabilirsiniz.',
            }
        ]
    },
    {
        id: 'platforms',
        title: 'Platform Kardeşliği',
        summary: 'Aynı altyapıyı kullanan araçlar ve parça uyumluluğu hakkında bilmeniz gerekenler.',
        sections: [
            {
                title: 'Modüler Platformlar (MQB, EMP2, CLAR)',
                content: 'Otomotiv devleri maliyeti düşürmek için devasa "lego" setleri kullanır. VW Golf, Audi A3, Seat Leon ve Skoda Octavia aslında "MQB" platformunda aynı şasidir. Motor, şanzıman ve elektroniklerin %70"i ortaktır.',
                imagePlaceholder: 'car_chassis_platform'
            },
            {
                title: 'Parça Uyumluluğu',
                content: 'Audi logosu basılı bir su pompası 5000 TL iken, aynı üreticinin (örn: INA veya Bosch) kutusundaki aynı parça 1500 TL olabilir. Platform kardeşliğini bilmek, bakım maliyetlerinizi %60"a kadar düşürebilir.',
                imagePlaceholder: 'spare_parts_comparison'
            }
        ]
    }
];
