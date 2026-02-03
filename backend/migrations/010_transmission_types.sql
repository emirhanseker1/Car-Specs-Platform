-- Transmission Types Table for detailed transmission information
CREATE TABLE IF NOT EXISTS transmission_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT UNIQUE NOT NULL,  -- DQ200, DQ250, DQ381, TIPTRONIC
    name TEXT NOT NULL,         -- Full display name
    type TEXT NOT NULL,         -- DSG, Automatic, etc.
    gears INTEGER,              -- Number of gears
    clutch_type TEXT,           -- dry, wet, torque_converter
    max_torque_nm INTEGER,      -- Maximum torque capacity
    description TEXT,           -- Technical description
    chronic_problems TEXT,      -- JSON array of common problems
    maintenance_tips TEXT,      -- JSON array of maintenance tips
    clutch_interval_km TEXT,    -- e.g. "60000-120000"
    smart_tip TEXT,             -- Quick user-facing tip
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert transmission type data
INSERT OR REPLACE INTO transmission_types (code, name, type, gears, clutch_type, max_torque_nm, description, chronic_problems, maintenance_tips, clutch_interval_km, smart_tip) VALUES
(
    'DQ200',
    '7-İleri Kuru Kavrama S-Tronic/DSG',
    'DSG',
    7,
    'dry',
    250,
    'VAG grubunun en çok tartışılan, en yaygın ve düşük torklu motorlarda kullandığı şanzımandır. Çift kavramalı, 7 ileri vitesli, "kuru" tip bir şanzımandır. Kavrama plakaları manuel vitesli araçlardaki gibi hava ile soğutulur. Hafiftir ve yakıt tüketimine katkısı pozitiftir.',
    '["Mekatronik Arızası: Şanzımanın beyni ve hidrolik ünitesidir. Yüksek basınç tüpü gevşeyebilir veya elektronik kart arızalanabilir.", "Kavrama Titremesi (Silkeleme): Özellikle 1. vitesten 2''ye geçerken araçta titreme hissedilmesi, kavramanın bittiğine veya ısındığına işarettir.", "Isınma: Yoğun dur-kalk trafikte yeterince soğuyamazsa ''şanzıman aşırı ısındı'' uyarısı verip kendini korumaya alabilir."]',
    '["Sıkışık trafikte manuel moda alıp 1. viteste sabit gitmek ömrü uzatır.", "Yokuşlarda aracı gaza basarak değil, Auto Hold veya frene tam basarak sabit tutun.", "Gaza çok az basıp aracı askıda tutmak kavramayı ''zımparalar''."]',
    '60000-120000',
    'Bu araç kuru kavrama şanzımana sahiptir. Yoğun dur-kalk trafikte şanzımanın ısınmaması için aracı sık sık ''N'' konumuna almanız veya manuel modda kullanmanız tavsiye edilir.'
),
(
    'DQ250',
    '6-İleri Yağlı Kavrama S-Tronic/DSG',
    'DSG',
    6,
    'wet',
    400,
    'Daha güçlü motorlarda (2.0 TDI, 2.0 TFSI vb.) kullanılan, DQ200''e göre çok daha dayanıklı olan şanzımandır. 6 ileri vitesli, yağlı (wet clutch) tip şanzımandır. Kavrama plakaları özel bir şanzıman yağının içinde döner. Bu yağ hem soğutmayı sağlar hem de sürtünmeyi optimize eder.',
    '["Volant (Flywheel) Sesi: Çift kütleli volant zamanla boşluk yapabilir, rölantide ''takır takır'' metal sesi duyulabilir.", "Mekatronik Solenoidleri: Vites geçişlerinde vuruntu (kütleme) yaparsa solenoid valflerin kirlenmiş veya bozulmuş olma ihtimali yüksektir."]',
    '["Her 60.000 km''de bir şanzıman yağı ve filtresi mutlaka orijinaliyle değişmelidir. Yağ kirlenirse mekatroniği bozar.", "Motor ve şanzıman yağı ısınmadan (ilk 10-15 dk) Launch Control veya dip gaz yapılmamalıdır."]',
    '150000+',
    'Bu araç yağlı kavrama DSG şanzımana sahiptir. Kuru kavramaya göre çok daha dayanıklıdır. 60.000 km''de bir yağ değişimi şanzıman ömrü için kritik öneme sahiptir.'
),
(
    'DQ381',
    '7-İleri Yağlı Kavrama S-Tronic/DSG',
    'DSG',
    7,
    'wet',
    430,
    'DQ250''nin yerini alan, 2017 sonrası ve güncel kasalarda kullanılan optimize edilmiş versiyonudur. 7 ileri vitesli, yağlı kavramadır. DQ250''ye göre daha düşük sürtünmeli yağ kullanır, daha yeni bir yağ pompası sistemine sahiptir ve emisyon odaklı geliştirilmiştir.',
    '["Erken dönem üretimlerinde yardımcı hidrolik pompa arızaları rapor edilmiştir.", "Yazılımsal kararsızlıklar (vites seçiminde gecikme) görülebilir, genelde güncelleme ile çözülür."]',
    '["60.000 km veya 120.000 km aralığında yağ değişimi hayati önem taşır.", "Start-Stop sistemiyle entegre çalıştığı için akü voltajına hassastır. Akü zayıflarsa şanzıman hataları verebilir."]',
    '150000+',
    'Bu araç güncel nesil yağlı kavrama DSG şanzımana sahiptir. DQ250''nin geliştirilmiş versiyonudur. Düzenli yağ değişimi ile çok uzun ömürlüdür.'
),
(
    'TIPTRONIC',
    'Tork Konvertörlü Tam Otomatik',
    'Automatic',
    NULL,
    'torque_converter',
    NULL,
    'Eski kasalarda (8L ve erken 8P) veya çok yüksek torklu lüks Audi modellerinde görülür. Mekanik bir kavrama plakası yoktur. Motor gücünü şanzımana sıvı (yağ) basıncıyla (tork konvertörü) iletir. Vites geçişleri DSG kadar hızlı değildir ama çok daha pürüzsüzdür.',
    '["Türbin (Konvertör) Arızası: Sabit hızda giderken devir saatinde dalgalanma veya araçta titreme yapabilir.", "Vuruntu: Valf gövdesi içindeki kanallar aşınırsa vites geçişlerinde sert vurma yapabilir."]',
    '["''Ömürlük yağ'' efsanesine inanmayın. Her 60.000-80.000 km''de bir yağ değişimi şanzımanın ömrünü ikiye katlar.", "Soğuk motorla ani gaz açmaktan kaçının."]',
    '200000+',
    'Bu araç geleneksel tork konvertörlü otomatik şanzımana sahiptir. Fiziksel kavrama balatası olmadığı için çok uzun ömürlüdür. Düzenli yağ değişimi önemlidir.'
);
