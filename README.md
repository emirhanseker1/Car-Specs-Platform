# 🚗 Araçların Tüm Detayları - Car Specs Platform

Araçların detaylı teknik özelliklerini görüntülemek için tasarlanmış kapsamlı bir platform. Modern Go backend, React frontend ve yapay zeka destekli chatbot ile güçlendirilmiştir.

![Car Specs Platform](frontend/public/project-screenshot.png)

## ✨ Özellikler

### 🔍 Araç Veritabanı ve Arama
- **Detaylı Araç Veritabanı**: Motor, şanzıman, boyutlar ve daha fazlası dahil olmak üzere kapsamlı teknik özellikler
- **Gelişmiş Arama**: Marka, model, yıl, yakıt tipi ve şanzıman türüne göre filtreleme
- **Karşılaştırma Aracı**: Farklı araç trim'lerinin yan yana karşılaştırılması
- **Nesil Bazlı Organizasyon**: Araçlar nesil kodları ile organize edilmiş (örn: Audi A3 8P, 8V, 8Y)

### 🤖 Yapay Zeka Destekli Chatbot
- **Otomotiv Uzman Asistanı**: Google Gemini 1.5 Flash ile güçlendirilmiş AI chatbot
- **Gerçek Zamanlı Web Araması**: Tavily Search API entegrasyonu ile güncel bilgilere erişim
- **Domain Özel Bilgi**: Otomotiv mühendisliği ve servis teknisyeni uzmanlığı
- **Akıllı Doğruluk Kontrolleri**: Motor tipine özel validasyon (benzinli/dizel ayrımı)
- **Türkçe Destek**: Tam Türkçe dil desteği ile yanıtlar

### 🎨 Modern Kullanıcı Arayüzü
- **Sürükleyici Tasarım**: Glassmorphism efektleri ve yüksek kaliteli görseller ile modern dark mode
- **Etkileşimli Rehberler**: 
  - Şanzıman Rehberi (Otomatik, Manuel, DSG, CVT)
  - Motor Terimleri Sözlüğü
  - DSG Varyantları Detaylı Kılavuzu
- **Responsive Design**: Tüm cihazlarda mükemmel görünüm
- **Smooth Animasyonlar**: Framer Motion ile akıcı geçişler

### 📊 Veri Senkronizasyonu
- **Güncel Veriler**: CarQuery API ile otomatik senkronizasyon
- **4 Seviyeli Veri Yapısı**: Brand → Model → Generation → Trim hiyerarşisi

## 🛠️ Teknoloji Yığını

### Backend
- **Go (Golang)**: Yüksek performanslı API sunucusu
- **SQLite + GORM**: Hafif ve hızlı veritabanı yönetimi
- **Standard Library HTTP**: Minimal bağımlılık ile native routing

### Frontend
- **React 18**: Modern UI geliştirme
- **TypeScript**: Type-safe kod yazımı
- **Vite**: Ultra hızlı build tool
- **Tailwind CSS**: Utility-first CSS framework
- **Framer Motion**: Smooth animasyonlar
- **Lucide React**: Modern icon kütüphanesi

### AI Chatbot Microservice
- **Python 3.x**: FastAPI ile yüksek performanslı API
- **LangChain**: AI agent orkestrayonu
- **Google Gemini 1.5 Flash**: Son teknoloji LLM
- **Tavily Search API**: Gerçek zamanlı web araması
- **FastAPI + Uvicorn**: Async API framework

## 🚀 Kurulum ve Çalıştırma

### Gereksinimler

- **Go** 1.22 veya üstü
- **Node.js** 18+ ve npm
- **Python** 3.8+ (Chatbot için)
- **API Keys** (Chatbot için):
  - Google AI API Key (Gemini)
  - Tavily API Key

### 📦 Hızlı Başlangıç (Önerilen)

Proje kök dizinindeki `start.bat` dosyasına çift tıklayın. Bu, bağımlılıkları kontrol edecek ve backend ile frontend servislerini otomatik olarak başlatacaktır.

```bash
# Windows için:
start.bat
```

### 🔧 Manuel Kurulum

#### 1. Backend API (Go)

```bash
cd backend
go mod tidy

# Sunucuyu çalıştır
go run cmd/api/main.go

# İlk kurulumda veri çekmek için:
go run cmd/api/main.go -scrape
```

**Backend çalışıyor:** `http://localhost:8080`

#### 2. Frontend (React)

```bash
cd frontend
npm install
npm run dev
```

**Frontend çalışıyor:** `http://localhost:5173`

#### 3. AI Chatbot Microservice (Python)

```bash
cd python-service

# Virtual environment oluştur (önerilen)
python -m venv venv
venv\Scripts\activate  # Windows
# source venv/bin/activate  # Linux/Mac

# Bağımlılıkları yükle
pip install -r requirements.txt

# Ortam değişkenlerini yapılandır
# .env dosyası oluşturun ve API anahtarlarınızı ekleyin:
# GOOGLE_API_KEY=your_google_api_key
# TAVILY_API_KEY=your_tavily_api_key

# Servisi çalıştır
python main.py
```

**Chatbot API çalışıyor:** `http://localhost:8000`

## 📁 Proje Yapısı

```
car-specs/
├── backend/                 # Go API sunucusu
│   ├── cmd/api/            # Ana uygulama
│   ├── internal/           # İç paketler
│   │   ├── database/       # Veritabanı yapılandırması
│   │   ├── handlers/       # HTTP handlers
│   │   ├── models/         # Veri modelleri
│   │   ├── repository/     # Veritabanı katmanı
│   │   └── service/        # İş mantığı katmanı
│   ├── migrations/         # Veritabanı migrasyonları
│   └── pkg/                # Harici paketler
│       └── carquery/       # CarQuery API client
│
├── frontend/               # React uygulaması
│   ├── src/
│   │   ├── components/     # Yeniden kullanılabilir componentler
│   │   ├── pages/          # Sayfa componentleri
│   │   ├── services/       # API servisleri
│   │   └── types/          # TypeScript tipleri
│   └── public/             # Statik dosyalar
│
├── python-service/         # AI Chatbot microservice
│   ├── main.py            # FastAPI uygulaması
│   ├── requirements.txt   # Python bağımlılıkları
│   └── .env               # API anahtarları (git'de değil)
│
└── data/                   # Paylaşılan veri kaynakları
```


**⭐ Projeyi beğendiyseniz yıldız vermeyi unutmayın! ⭐**

Made with ❤️ in Turkey 🇹🇷

</div>
