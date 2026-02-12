package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type EngineProblem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type EngineSolution struct {
	ProblemTitle  string `json:"problemTitle"`
	Solution      string `json:"solution"`
	EstimatedCost string `json:"estimatedCost,omitempty"`
}

type TechnologyFeature struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type EngineData struct {
	Code                string
	Name                string
	Manufacturer        string
	EngineType          string
	FuelType            string
	DisplacementCC      int
	Cylinders           int
	CylinderLayout      string
	ValvesPerCylinder   int
	Aspiration          string
	Description         string
	TechnologyFeatures  []TechnologyFeature
	ProductionStartYear int
	ProductionEndYear   *int
	CommonProblems      []EngineProblem
	Solutions           []EngineSolution
	MaintenanceNotes    string
	ReliabilityRating   int
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./vehicles.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	log.Println("=== VW Group 1.2L & 1.4L Engines Data Insertion ===")
	log.Printf("Database: %s\n", dbPath)

	engines := []EngineData{
		// === 1.2 TSI EA111 (Zincirli) ===
		{
			Code:              "CBZA",
			Name:              "1.2 TSI 86",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1197,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "Küçük hacimli turbo devrimini başlatan EA111 serisi 1.2 TSI motor. Eksantrik tahriki zincir ile sağlanır. 86 PS güç, 160 Nm tork üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Chain", Description: "Zincirli triger sistemi"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2010,
			ProductionEndYear:   intPtr(2014),
			CommonProblems: []EngineProblem{
				{
					Title:       "Triger Zinciri Uzaması (EN CİDDİ)",
					Description: "İlk çalıştırmada (özellikle soğukken) 2-3 saniye süren metalik 'şakırtı' sesi zincirin uzadığının belirtisidir. Değiştirilmezse sente atlar.",
					Severity:    "high",
				},
				{
					Title:       "Turbo Wastegate Zırıltısı",
					Description: "Turbonun atık kapakçığı (wastegate) zamanla boşluk yapar ve gaza basıp bırakmalarda zırıltı sesi çıkarır",
					Severity:    "medium",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Triger Zinciri Uzaması",
					Solution:      "Zincir seti, gerici ve yağ pompası komple değişimi. Erken müdahale kritiktir, sente atlatırsa motor inmesi gerekir.",
					EstimatedCost: "8000-15000 TL (zincir seti) / 30000+ TL (motor inme)",
				},
				{
					ProblemTitle:  "Turbo Wastegate Zırıltısı",
					Solution:      "Wastegate kapağı araya pul takılarak ayarlanır, gerekirse turbo revizyonu",
					EstimatedCost: "500-1000 TL (ayar) / 8000-12000 TL (turbo)",
				},
			},
			MaintenanceNotes:  "Zincir sesi başlar başlamaz müdahale edilmeli. Her 7500 km yağ değişimi şart. Bu motorlardan kaçınılması önerilir.",
			ReliabilityRating: 3,
		},
		{
			Code:              "CBZB",
			Name:              "1.2 TSI 105",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1197,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "CBZA ile aynı bloğu kullanan 105 PS versiyon. Zincir uzaması ve wastegate sorunu birebir aynıdır. 175 Nm tork üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Chain"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2009,
			ProductionEndYear:   intPtr(2014),
			CommonProblems: []EngineProblem{
				{
					Title:       "Triger Zinciri Uzaması",
					Description: "CBZA ile aynı kronik sorun",
					Severity:    "high",
				},
				{
					Title:       "Turbo Wastegate Zırıltısı",
					Description: "Wastegate boşluğu ve zırıltı",
					Severity:    "medium",
				},
				{
					Title:       "Buji Kablosu ve Ateşleme Bobini Arızası",
					Description: "Daha yüksek basınca maruz kaldığı için buji kabloları ve ateşleme bobini arızaları sık görülür",
					Severity:    "medium",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Triger Zinciri Uzaması",
					Solution:      "Zincir seti ve gerici komple değişimi",
					EstimatedCost: "8000-15000 TL",
				},
				{
					ProblemTitle:  "Buji Kablosu ve Ateşleme Bobini Arızası",
					Solution:      "Orijinal buji kablosu ve bobin değişimi",
					EstimatedCost: "2000-4000 TL",
				},
			},
			MaintenanceNotes:  "CBZA ile aynı sorunlara sahip. Yüksek riskli motor, satın alınması önerilmez.",
			ReliabilityRating: 3,
		},

		// === 1.2 TSI EA211 (Kayışlı - Yeni Nesil) ===
		{
			Code:              "CJZC",
			Name:              "1.2 TSI 90",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1197,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "EA111'deki zincir sorunlarından ders çıkarılarak üretilen tamamen yeni nesil, triger kayışlı ve alüminyum bloklu 1.2 TSI. 90 PS güç, 160 Nm tork.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt", Description: "Kayışlı triger sistemi (zincir sorunu yok)"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
				{Name: "Aluminum Block", Description: "Hafif alüminyum blok"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   intPtr(2017),
			CommonProblems: []EngineProblem{
				{
					Title:       "Turbo Aktüatör Gıcırtısı",
					Description: "Nadiren turbo aktüatör çubuğundan gıcırtı sesi şikayeti gelir (Yağlama ile çözülür)",
					Severity:    "low",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Turbo Aktüatör Gıcırtısı",
					Solution:      "Aktüatör çubuğu yağlaması",
					EstimatedCost: "200-500 TL",
				},
			},
			MaintenanceNotes:  "Mekanik olarak son derece sağlam. Zincir sorunu yok. Günümüz standardlarında güvenilir motor.",
			ReliabilityRating: 8,
		},
		{
			Code:              "CYVA",
			Name:              "1.2 TSI 90",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1197,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "CJZC'nin güncellenmiş versiyonu. Aynı motor bloku, yazılımsal optimizasyonlar.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   intPtr(2017),
			CommonProblems:      []EngineProblem{},
			Solutions:           []EngineSolution{},
			MaintenanceNotes:    "CJZC ile aynı güvenilirlik seviyesinde.",
			ReliabilityRating:   8,
		},
		{
			Code:              "CJZA",
			Name:              "1.2 TSI 105",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1197,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "EA211 serisi kayışlı 1.2 TSI motorun 105 PS versiyonu. 175 Nm tork.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2012,
			ProductionEndYear:   intPtr(2017),
			CommonProblems: []EngineProblem{
				{
					Title:       "Termostat / Devirdaim Pompası Su Kaçağı",
					Description: "EA211 serisinin en meşhur kroniği. Plastik gövde zamanla kılcal çatlak yapar ve antifriz eksiltir",
					Severity:    "high",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Termostat / Devirdaim Pompası Su Kaçağı",
					Solution:      "Termostat gövdesi komple değişimi. Metal aftermarket parçalar daha dayanıklı.",
					EstimatedCost: "2500-4000 TL",
				},
			},
			MaintenanceNotes:  "Zincir uzama sorunu bitmiştir. 80.000-100.000 km'de termostat gövdesi proaktif değiştirilmeli.",
			ReliabilityRating: 7,
		},
		{
			Code:              "CYVB",
			Name:              "1.2 TSI 110",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1197,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "1.2 TSI EA211 serisinin en yüksek güçlü versiyonu. 110 PS, 175 Nm.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2012,
			ProductionEndYear:   intPtr(2017),
			CommonProblems: []EngineProblem{
				{
					Title:       "Termostat Gövdesi Çatlakları",
					Description: "EA211 kroniği - plastik termostat gövdesi su kaçağı",
					Severity:    "high",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Termostat Gövdesi Çatlakları",
					Solution:      "Gövde komple değişimi",
					EstimatedCost: "2500-4000 TL",
				},
			},
			MaintenanceNotes:  "Güvenilir motor. Termostat bakımı önemli.",
			ReliabilityRating: 7,
		},

		// === 1.4 TSI EA111 (Sadece Turbo) ===
		{
			Code:              "CAXA",
			Name:              "1.4 TSI 122",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "Eski nesil zincirli motorun tek aşırı beslemeli (sadece turbo) standart versiyonu. 122 PS güç, 200 Nm tork.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Chain"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2007,
			ProductionEndYear:   intPtr(2014),
			CommonProblems: []EngineProblem{
				{
					Title:       "Sente Atlama Riski (KRİTİK)",
					Description: "Zincir gergisi yağ basıncıyla çalışır. Araç viteste yokuşta park edilirse ve zincir uzamışsa, ilk marşta sente atlatıp supap eğme riski vardır",
					Severity:    "high",
				},
				{
					Title:       "N75 Turbo Basınç Valfi Arızası",
					Description: "Turbo boost kontrol valfi arızalanır, güç kaybı ve MIL lambası yanar",
					Severity:    "medium",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Sente Atlama Riski",
					Solution:      "Zincir uzaması başladığında acilen zincir seti değiştirilmeli. Viteste yokuşta park etmekten kaçınılmalı.",
					EstimatedCost: "8000-15000 TL (zincir) / 30000+ TL (motor inme)",
				},
				{
					ProblemTitle:  "N75 Turbo Basınç Valfi Arızası",
					Solution:      "N75 valfi değişimi",
					EstimatedCost: "500-1200 TL",
				},
			},
			MaintenanceNotes:  "Yüksek riskli motor. Zincir sesi başladığında acil müdahale gerekir. Satın alırken çok dikkatli olunmalı.",
			ReliabilityRating: 4,
		},
		{
			Code:              "CMSB",
			Name:              "1.4 TSI 122",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "CAXA'nın güncellenmiş versiyonu. Kronik sorunlar aynı.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Chain"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2007,
			ProductionEndYear:   intPtr(2014),
			CommonProblems: []EngineProblem{
				{
					Title:       "Zincir Uzaması ve Sente Atlama",
					Description: "CAXA ile aynı risk",
					Severity:    "high",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "CAXA ile aynı riskler. Dikkatli olunmalı.",
			ReliabilityRating: 4,
		},
		{
			Code:              "CAXC",
			Name:              "1.4 TFSI 125",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "CAXA motorun Audi/Seat/Skoda grubu için yazılımla 3 beygir güçlendirilmiş hali. 125 PS, 200 Nm.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Chain"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2007,
			ProductionEndYear:   intPtr(2014),
			CommonProblems: []EngineProblem{
				{
					Title:       "Zincir Uzaması",
					Description: "CAXA ile aynı kronik",
					Severity:    "high",
				},
				{
					Title:       "Wastegate Zırıltısı",
					Description: "Turbo wastegate boşluğu",
					Severity:    "medium",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "CAXA kronikleri birebir aynı.",
			ReliabilityRating: 4,
		},

		// === 1.4 TSI Twincharger (Turbo + Kompresör) - EN SORUNLU SERİ ===
		{
			Code:              "BMY",
			Name:              "1.4 TSI Twincharger 140",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Twincharged (Turbo + Supercharger)",
			Description:       "Alt devirlerde kompresör, üst devirlerde turbo devreye girer. Mühendislik harikası ancak kronik sorunlarıyla çok ünlü. 140 PS, 220 Nm.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Twincharger", Description: "Turbo + Süperşarj (Kompresör)"},
				{Name: "Timing Chain"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2006,
			ProductionEndYear:   intPtr(2010),
			CommonProblems: []EngineProblem{
				{
					Title:       "Manyetik Kavrama Arızası",
					Description: "Kompresörü devreye sokan manyetik kavrama (su pompası ile bütün) sık arızalanır ve araç 3000 devire kadar hantal kalır (Ciyaklama sesi yapar)",
					Severity:    "high",
				},
				{
					Title:       "Yağ Yakma Sorunu",
					Description: "Motor yağ tüketimi başlar, 1000 km'de 0.5-1 litre yağ eksiltme normal kabul edilir",
					Severity:    "high",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Manyetik Kavrama Arızası",
					Solution:      "Su pompası + manyetik kavrama komple değişimi",
					EstimatedCost: "5000-8000 TL",
				},
			},
			MaintenanceNotes:  "Satın alınması ÖNERİLMEZ. Bu motordan uzak durulmalı.",
			ReliabilityRating: 2,
		},
		{
			Code:              "BWK",
			Name:              "1.4 TSI Twincharger 150",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Twincharged (Turbo + Supercharger)",
			Description:       "Twincharger serisinin 150 PS versiyonu. 240 Nm tork.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Twincharger"},
				{Name: "Timing Chain"},
			},
			ProductionStartYear: 2006,
			ProductionEndYear:   intPtr(2010),
			CommonProblems: []EngineProblem{
				{
					Title:       "Kompresör Kavraması",
					Description: "Manyetik kavrama arızası",
					Severity:    "high",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "BMY ile aynı sorunlar. Satın alınmamalı.",
			ReliabilityRating: 2,
		},
		{
			Code:              "CAVA",
			Name:              "1.4 TSI Twincharger 150",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Twincharged (Turbo + Supercharger)",
			Description:       "Twincharger serisinin 2008+ versiyonu.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Twincharger"},
			},
			ProductionStartYear: 2006,
			ProductionEndYear:   intPtr(2010),
			CommonProblems:      []EngineProblem{},
			Solutions:           []EngineSolution{},
			MaintenanceNotes:    "Twincharger serisi - yüksek riskli.",
			ReliabilityRating:   2,
		},
		{
			Code:              "CAVD",
			Name:              "1.4 TSI Twincharger 160 (SORUNLU)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Twincharged (Turbo + Supercharger)",
			Description:       "Twincharger serisinin EN SORUNLU motoru. Piston kırma (Ringland Failure) kroniği ile meşhur. 160 PS, 240 Nm.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Twincharger"},
				{Name: "Timing Chain"},
			},
			ProductionStartYear: 2008,
			ProductionEndYear:   intPtr(2012),
			CommonProblems: []EngineProblem{
				{
					Title:       "Piston Kırma (Ringland Failure) - KRİTİK",
					Description: "Piston segman kanalları çok ince. Kötü yakıt (vuruntu) veya zorlama sonucu 2. veya 3. piston segman kanalından kırılır. Kompresyon kaybı ve yoğun yağ yakma başlar. Motor inmesi gerekir",
					Severity:    "high",
				},
				{
					Title:       "Kompresör Kavraması",
					Description: "Manyetik kavrama arızası",
					Severity:    "high",
				},
				{
					Title:       "Yoğun Yağ Yakma",
					Description: "1000 km'de 1+ litre yağ eksiltme",
					Severity:    "high",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Piston Kırma",
					Solution:      "Motor inmesi, piston değişimi, honing, yeni segman takımı",
					EstimatedCost: "25000-40000+ TL",
				},
			},
			MaintenanceNotes:  "KESINLIKLE ALINMAMALI! VW'nin kara leke motoru.",
			ReliabilityRating: 1,
		},
		{
			Code:              "CTHD",
			Name:              "1.4 TSI Twincharger 160 (Revizyonlu)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Twincharged (Turbo + Supercharger)",
			Description:       "VW, 2013'te piston yapısını kalınlaştırarak ve yağ jetlerini büyüterek bu koda geçti. Sorunlar büyük oranda çözülse de tamamen bitmemiştir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Twincharger"},
				{Name: "Reinforced Pistons", Description: "Güçlendirilmiş pistonlar"},
			},
			ProductionStartYear: 2013,
			ProductionEndYear:   intPtr(2014),
			CommonProblems: []EngineProblem{
				{
					Title:       "Yağ Yakma (Azalmış)",
					Description: "CAVD'den daha az ama hala mevcut",
					Severity:    "medium",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "CAVD'den daha iyi ama yine de riskli. Dikkatli olunmalı.",
			ReliabilityRating: 3,
		},
		{
			Code:              "CAVE",
			Name:              "1.4 TSI Twincharger 180 (GTI/Cupra)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Twincharged (Turbo + Supercharger)",
			Description:       "Twincharger serisinin performans versiyonu. 180 PS, 250 Nm. Yağ eksiltme bu motorлдa fabrikasyon kabul edilir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Twincharger"},
				{Name: "Performance Tuning"},
			},
			ProductionStartYear: 2010,
			ProductionEndYear:   intPtr(2012),
			CommonProblems: []EngineProblem{
				{
					Title:       "Yağ Eksiltme (Fabrikasyon Seviye)",
					Description: "1000 km'de 1 litreye kadar yağ yakma normal kabul edilir",
					Severity:    "high",
				},
				{
					Title:       "Buji Erimesi",
					Description: "Yüksek sıcaklık ve basınç nedeniyle buji erimesi",
					Severity:    "medium",
				},
				{
					Title:       "Piston Kırma Riski",
					Description: "CAVD'deki gibi piston kırma riskleri en üst seviyede",
					Severity:    "high",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "Performans tutkunları için, ama yüksek bakım maliyeti ve arıza riski. Her 5000 km yağ kontrolü şart.",
			ReliabilityRating: 2,
		},
		{
			Code:              "CTHE",
			Name:              "1.4 TSI Twincharger 180 (Revizyonlu)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1390,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Twincharged (Turbo + Supercharger)",
			Description:       "CAVE'nin 2013+ revizyonlu versiyonu. Güçlendirilmiş pistonlar.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Twincharger"},
				{Name: "Reinforced Pistons"},
			},
			ProductionStartYear: 2013,
			ProductionEndYear:   intPtr(2014),
			CommonProblems:      []EngineProblem{},
			Solutions:           []EngineSolution{},
			MaintenanceNotes:    "CAVE'den daha iyi ama yine de performans motoru riskleri var.",
			ReliabilityRating:   3,
		},

		// === 1.4 TSI EA211 (Kayışlı - Yeni Nesil - GÜVENİLİR) ===
		{
			Code:              "CMBA",
			Name:              "1.4 TSI 122 (EA211)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1395,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "Twincharger fiyaskosundan sonra tamamen yeni tasarlanan, kayışlı, alüminyum bloklu, güvenilir 1.4 TSI. 122 PS, 200 Nm.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
				{Name: "Integrated Exhaust Manifold", Description: "Egzoz manifoldu silindir kapağına entegre"},
				{Name: "Aluminum Block"},
			},
			ProductionStartYear: 2012,
			ProductionEndYear:   intPtr(2018),
			CommonProblems: []EngineProblem{
				{
					Title:       "Termostat Gövdesi Çatlaması",
					Description: "EA211 kroniği - plastik termostat gövdesi su sızıntısı",
					Severity:    "medium",
				},
				{
					Title:       "Turbo Aktüatör Paslanması",
					Description: "Nadiren turbo aktüatör çubuğu paslanıp takılı kalabilir",
					Severity:    "low",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Termostat Gövdesi Çatlaması",
					Solution:      "Termostat gövdesi değişimi",
					EstimatedCost: "2000-3500 TL",
				},
			},
			MaintenanceNotes:  "Mekanik zincir uzaması sorunları yok. İnanılmaz sessiz ve dayanıklı. Güvenilir motor.",
			ReliabilityRating: 8,
		},
		{
			Code:              "CXSA",
			Name:              "1.4 TSI 122 (EA211)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1395,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "CMBA'nın güncellenmiş versiyonu.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2012,
			ProductionEndYear:   intPtr(2018),
			CommonProblems:      []EngineProblem{},
			Solutions:           []EngineSolution{},
			MaintenanceNotes:    "Güvenilir EA211 motoru.",
			ReliabilityRating:   8,
		},
		{
			Code:              "CZCA",
			Name:              "1.4 TFSI 125 (EA211)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1395,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "EA211 1.4 TSI'ın 125 PS Audi versiyonu.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2012,
			ProductionEndYear:   intPtr(2018),
			CommonProblems:      []EngineProblem{},
			Solutions:           []EngineSolution{},
			MaintenanceNotes:    "Güvenilir motor.",
			ReliabilityRating:   8,
		},
		{
			Code:              "CHPA",
			Name:              "1.4 TSI 140 (EA211)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1395,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "EA211 1.4 TSI 140 PS normal versiyon. 250 Nm tork.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2012,
			ProductionEndYear:   intPtr(2014),
			CommonProblems:      []EngineProblem{},
			Solutions:           []EngineSolution{},
			MaintenanceNotes:    "Güvenilir motor.",
			ReliabilityRating:   8,
		},
		{
			Code:              "CPTA",
			Name:              "1.4 TSI 140 ACT/COD",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1395,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged + Cylinder Deactivation",
			Description:       "Silindir kapatmalı (ACT/COD) teknolojili versiyon. Ortadaki 2 silindir kapatılabilir. 140 PS, 250 Nm.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "ACT - Active Cylinder Technology", Description: "Silindir kapatma sistemi"},
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2012,
			ProductionEndYear:   intPtr(2014),
			CommonProblems: []EngineProblem{
				{
					Title:       "ACT Sistem Titremesi",
					Description: "Ortadaki 2 silindir kapatıldığında ve tekrar devreye girdiğinde hafif bir 'titreme' veya 'hırıltı' hissedilir",
					Severity:    "low",
				},
				{
					Title:       "Aktüatör Yağ Terlemesi",
					Description: "Eksantrik mili üzerindeki silindir kapatma aktüatörlerinde (mıknatıslarda) zamanla yağ terlemesi görülebilir",
					Severity:    "low",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "ACT sistemi karakteristik titreme yapar, arıza değil. Genel olarak güvenilir.",
			ReliabilityRating: 7,
		},
		{
			Code:              "CZDA",
			Name:              "1.4 TSI 150 (EA211)",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1395,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "EA211 1.4 TSI'ın en yüksek güçlü normal versiyonu. 150 PS, 250 Nm.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   intPtr(2019),
			CommonProblems: []EngineProblem{
				{
					Title:       "Su Pompası Terlemesi",
					Description: "Klasik EA211 su pompası kroniği",
					Severity:    "low",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "Grubun en verimli ve sağlam motorlarından biri. Arıza kaydı oldukça nadir.",
			ReliabilityRating: 9,
		},
		{
			Code:              "CZEA",
			Name:              "1.4 TSI 150 ACT/COD",
			Manufacturer:      "VW Group",
			EngineType:        "I4",
			FuelType:          "Petrol",
			DisplacementCC:    1395,
			Cylinders:         4,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged + Cylinder Deactivation",
			Description:       "Silindir kapama teknolojili 150 PS versiyon. Grubun en verimli motorlarından biri. 250 Nm tork.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "ACT - Active Cylinder Technology"},
				{Name: "Timing Belt"},
				{Name: "Turbocharger"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   intPtr(2019),
			CommonProblems: []EngineProblem{
				{
					Title:       "ACT Titreşimi",
					Description: "Silindir kapatma sırasında hafif titreme (karakteristik)",
					Severity:    "low",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "Klasik EA211 su pompası terlemesi ve ACT titreşimi dışında arıza kaydı oldukça nadir. Çok güvenilir.",
			ReliabilityRating: 9,
		},
	}

	log.Printf("\n📝 Adding %d engines to database...\n", len(engines))

	for i, engine := range engines {
		log.Printf("\n[%d/%d] Adding engine: %s (%s)", i+1, len(engines), engine.Code, engine.Name)

		techFeaturesJSON, _ := json.Marshal(engine.TechnologyFeatures)
		problemsJSON, _ := json.Marshal(engine.CommonProblems)
		solutionsJSON, _ := json.Marshal(engine.Solutions)

		query := `
			INSERT INTO engines (
				code, name, manufacturer, engine_type, fuel_type,
				displacement_cc, cylinders, cylinder_layout, valves_per_cylinder,
				aspiration, description, technology_features,
				production_start_year, production_end_year,
				common_problems, solutions, maintenance_notes, reliability_rating
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err := db.Exec(query,
			engine.Code, engine.Name, engine.Manufacturer, engine.EngineType, engine.FuelType,
			engine.DisplacementCC, engine.Cylinders, engine.CylinderLayout, engine.ValvesPerCylinder,
			engine.Aspiration, engine.Description, string(techFeaturesJSON),
			engine.ProductionStartYear, engine.ProductionEndYear,
			string(problemsJSON), string(solutionsJSON), engine.MaintenanceNotes, engine.ReliabilityRating,
		)

		if err != nil {
			log.Printf("  ❌ Error: %v", err)
			continue
		}

		log.Printf("  ✅ Successfully added")
	}

	log.Println("\n✅ All VW 1.2L & 1.4L engines have been added to the database!")
	log.Println("\n📊 Summary:")
	log.Printf("   - Total engines added: %d", len(engines))
	log.Println("\n🔧 Engine Categories:")
	log.Println("   - 1.2 TSI EA111 (Zincirli - SORUNLU): CBZA, CBZB")
	log.Println("   - 1.2 TSI EA211 (Kayışlı - GÜVENİLİR): CJZC, CYVA, CJZA, CYVB")
	log.Println("   - 1.4 TSI EA111 (Sadece Turbo - RİSKLİ): CAXA, CMSB, CAXC")
	log.Println("   - 1.4 TSI Twincharger (KESINLIKLE ALINMAMALI): BMY, BWK, CAVA, CAVD, CTHD, CAVE, CTHE")
	log.Println("   - 1.4 TSI EA211 (Kayışlı - ÇOK GÜVENİLİR): CMBA, CXSA, CZCA, CHPA, CPTA, CZDA, CZEA")
	log.Println("\nYou can now view these engines at: http://localhost:5173/guides/engines")
}

func intPtr(i int) *int {
	return &i
}
