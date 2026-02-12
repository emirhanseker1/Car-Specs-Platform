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

	log.Println("=== VW Group 1.0L Engines Data Insertion ===")
	log.Printf("Database: %s\n", dbPath)

	engines := []EngineData{
		// 1. 1.0 MPI (60/65 PS)
		{
			Code:              "CHYA",
			Name:              "1.0 MPI 60",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Naturally Aspirated",
			Description:       "Turbo beslemesi olmayan, şehir içi kullanıma yönelik üretilmiş, basit yapılı ve uzun ömürlü giriş seviyesi motor. 60 PS güç üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "MPI (Multi-Point Injection)", Description: "Çoklu nokta enjeksiyon sistemi"},
				{Name: "4 Valf/Silindir", Description: "Modern supap tasarımı"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   nil,
			CommonProblems: []EngineProblem{
				{
					Title:       "Rölanti Titreşimi",
					Description: "3 silindirli yapısından dolayı rölantide direksiyona yansıyan titreşim mevcuttur",
					Severity:    "low",
				},
				{
					Title:       "Debriyaj Bilyası Sesi",
					Description: "Manuel şanzımanla eşleştiğinde prizdirek (debriyaj bilyası) sesi şikayetleri yaygındır",
					Severity:    "low",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Rölanti Titreşimi",
					Solution:      "Motor takozu kontrolü ve gerekirse değişimi. Normal karakteristik davranış olduğu için tam çözüm yoktur.",
					EstimatedCost: "500-1500 TL",
				},
				{
					ProblemTitle:  "Debriyaj Bilyası Sesi",
					Solution:      "Debriyaj bilyası (prizdirek rulmanı) değişimi",
					EstimatedCost: "800-2000 TL",
				},
			},
			MaintenanceNotes:  "Mekanik olarak son derece basittir ve kronik arızası neredeyse yoktur. Düzenli yağ ve filtre bakımı yeterlidir.",
			ReliabilityRating: 9,
		},
		{
			Code:              "CHYC",
			Name:              "1.0 MPI 65",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Naturally Aspirated",
			Description:       "Turbo beslemesi olmayan, şehir içi kullanıma yönelik üretilmiş, basit yapılı ve uzun ömürlü giriş seviyesi motor. 65 PS güç üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "MPI (Multi-Point Injection)", Description: "Çoklu nokta enjeksiyon sistemi"},
				{Name: "4 Valf/Silindir", Description: "Modern supap tasarımı"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   nil,
			CommonProblems: []EngineProblem{
				{
					Title:       "Rölanti Titreşimi",
					Description: "3 silindirli yapısından dolayı rölantide direksiyona yansıyan titreşim mevcuttur",
					Severity:    "low",
				},
				{
					Title:       "Debriyaj Bilyası Sesi",
					Description: "Manuel şanzımanla eşleştiğinde prizdirek (debriyaj bilyası) sesi şikayetleri yaygındır",
					Severity:    "low",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Rölanti Titreşimi",
					Solution:      "Motor takozu kontrolü ve gerekirse değişimi. Normal karakteristik davranış olduğu için tam çözüm yoktur.",
					EstimatedCost: "500-1500 TL",
				},
				{
					ProblemTitle:  "Debriyaj Bilyası Sesi",
					Solution:      "Debriyaj bilyası (prizdirek rulmanı) değişimi",
					EstimatedCost: "800-2000 TL",
				},
			},
			MaintenanceNotes:  "CHYA ile aynı blok, yazılımsal farklılık. Mekanik olarak son derece basit ve güvenilir.",
			ReliabilityRating: 9,
		},
		// 2. 1.0 MPI (75/80 PS)
		{
			Code:              "CHYB",
			Name:              "1.0 MPI 75",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Naturally Aspirated",
			Description:       "CHYA ile tamamen aynı blok, sadece yazılımsal olarak üst devirleri açılmıştır. 75 PS güç üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "MPI (Multi-Point Injection)", Description: "Çoklu nokta enjeksiyon sistemi"},
				{Name: "Yazılımsal Optimizasyon", Description: "Üst devir açılımı"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   nil,
			CommonProblems: []EngineProblem{
				{
					Title:       "Uzun Yol ve Yokuş Çekiş Düşüklüğü",
					Description: "Uzun yolda ve yokuşlarda çekiş düşüklüğü en büyük şikayettir (kronik arıza değil, karakteristik)",
					Severity:    "medium",
				},
				{
					Title:       "Rölanti Titreşimi",
					Description: "3 silindirli yapısından dolayı rölantide titreşim",
					Severity:    "low",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Uzun Yol ve Yokuş Çekiş Düşüklüğü",
					Solution:      "Motor karakteristiği nedeniyle çözüm yok. Aracı yük altında kullanmamak en iyi yaklaşım.",
					EstimatedCost: "-",
				},
			},
			MaintenanceNotes:  "CHYA/CHYC ile aynı güvenilirlik. Şehir içi kullanıma optimize edilmiş, uzun yolda performans beklentisi düşük tutulmalı.",
			ReliabilityRating: 8,
		},
		{
			Code:              "DSGB",
			Name:              "1.0 MPI 80",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Naturally Aspirated",
			Description:       "1.0 MPI serisinin 80 PS versiyonu. Yazılımsal olarak optimize edilmiş.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "MPI (Multi-Point Injection)"},
				{Name: "Yazılımsal Optimizasyon"},
			},
			ProductionStartYear: 2014,
			ProductionEndYear:   nil,
			CommonProblems: []EngineProblem{
				{
					Title:       "Uzun Yol Çekiş Düşüklüğü",
					Description: "Yokuş ve uzun yollarda çekiş performansı düşüktür",
					Severity:    "medium",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "Şehir içi kullanıma uygun motor. Düzenli bakım yeterlidir.",
			ReliabilityRating: 8,
		},
		// 3. 1.0 TSI/TFSI (90/95 PS)
		{
			Code:              "CHZB",
			Name:              "1.0 TSI 95",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "B segmenti araçlarda kullanılan, yakıt tasarrufu odaklı ilk seviye turbo motor. 95 PS güç, 160-175 Nm tork üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Turbocharger", Description: "Turbo besleme sistemi"},
				{Name: "Direct Injection", Description: "Direkt enjeksiyon"},
				{Name: "Variable Valve Timing"},
			},
			ProductionStartYear: 2015,
			ProductionEndYear:   nil,
			CommonProblems: []EngineProblem{
				{
					Title:       "Wastegate Zırıltısı",
					Description: "Turbonun atık kapağından gelen metalik zırıltı sesi (özellikle ivmelenme sırasında duyulur, arıza değildir ancak şikayet edilir)",
					Severity:    "low",
				},
				{
					Title:       "DSG Kavrama Titremesi",
					Description: "95 PS versiyonlar genelde DSG (DQ200) ile eşleştirilir, 2. viteste düşük devirde kavrama titremesi görülür",
					Severity:    "medium",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Wastegate Zırıltısı",
					Solution:      "Normal çalışma sesi. Wastegate kapağı değişimi genelde gerekli değil, ses karakteristiktir.",
					EstimatedCost: "-",
				},
				{
					ProblemTitle:  "DSG Kavrama Titremesi",
					Solution:      "DSG adaptasyon sıfırlama, gerekirse DSG yazılım güncellemesi",
					EstimatedCost: "500-1000 TL (yazılım) / 8000-15000 TL (DSG onarımı)",
				},
			},
			MaintenanceNotes:  "Turbo motorlar için düzenli yağ bakımı kritiktir. Her 15.000 km'de sentetik yağ kullanımı önerilir.",
			ReliabilityRating: 7,
		},
		{
			Code:              "DKLA",
			Name:              "1.0 TSI 95",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "1.0 TSI serisinin güncel versiyonu. 95 PS güç üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
				{Name: "Variable Valve Timing"},
			},
			ProductionStartYear: 2015,
			ProductionEndYear:   nil,
			CommonProblems: []EngineProblem{
				{
					Title:       "Wastegate Zırıltısı",
					Description: "Turbo atık kapağından metalik ses",
					Severity:    "low",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "Düzenli turbo bakımı ve sentetik yağ kullanımı önemlidir.",
			ReliabilityRating: 7,
		},
		// 4. 1.0 TSI/TFSI (110/115/116 PS)
		{
			Code:              "CHZC",
			Name:              "1.0 TSI 110",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "C segmenti araçların ağırlığını çekebilmesi için turbosu ve yazılımı güçlendirilmiş versiyon. 110 PS güç, 200 Nm tork üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Upgraded Turbocharger", Description: "Güçlendirilmiş turbo"},
				{Name: "Direct Injection"},
				{Name: "Variable Valve Timing"},
			},
			ProductionStartYear: 2015,
			ProductionEndYear:   intPtr(2020),
			CommonProblems: []EngineProblem{
				{
					Title:       "Devirdaim / Su Pompası Sızıntısı",
					Description: "Plastik termostat gövdesi ve devirdaim pompasından antifriz sızdırması (EA211 serisinin genel kroniği)",
					Severity:    "high",
				},
				{
					Title:       "Çift Kütleli Volan Sesi",
					Description: "S tronic / DSG şanzımanla eşleştiğinde alt devirlerde (1200-1400 d/d arası) volandan gelen şıngırtı sesi",
					Severity:    "medium",
				},
				{
					Title:       "Turbo Gecikmesi (Lag)",
					Description: "Alt devirlerde torkun gelmesi 2000 devri bulduğu için şehir içi dur-kalk trafikte hantallık şikayeti",
					Severity:    "low",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Devirdaim / Su Pompası Sızıntısı",
					Solution:      "Orijinal termostat gövdesi ve su pompası değişimi. Metal gövdeli aftermarket parçalar daha dayanıklıdır.",
					EstimatedCost: "3000-5000 TL",
				},
				{
					ProblemTitle:  "Çift Kütleli Volan Sesi",
					Solution:      "Çift kütleli volan (DMF) değişimi. Tek kütleli volana dönüşüm de yapılabilir.",
					EstimatedCost: "8000-12000 TL (çift kütleli) / 4000-6000 TL (tek kütleli)",
				},
			},
			MaintenanceNotes:  "Su pompası ve termostat gövdesini 80.000-100.000 km'de proaktif değiştirmek önerilir. Turbo motorlar için sentetik yağ şarttır.",
			ReliabilityRating: 6,
		},
		{
			Code:              "CHZD",
			Name:              "1.0 TSI 115",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "1.0 TSI yüksek güç serisinin 115 PS versiyonu. 200 Nm tork üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Upgraded Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2015,
			ProductionEndYear:   intPtr(2020),
			CommonProblems: []EngineProblem{
				{
					Title:       "Devirdaim Pompası Sızıntısı",
					Description: "Plastik gövdeden antifriz sızıntısı",
					Severity:    "high",
				},
				{
					Title:       "Çift Kütleli Volan Sesi",
					Description: "Alt devirlerde volan şıngırtısı",
					Severity:    "medium",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Devirdaim Pompası Sızıntısı",
					Solution:      "Su pompası ve termostat gövdesi değişimi",
					EstimatedCost: "3000-5000 TL",
				},
			},
			MaintenanceNotes:  "Proaktif bakım kritiktir. Su pompası 100.000 km öncesi değiştirilmeli.",
			ReliabilityRating: 6,
		},
		{
			Code:              "DKRF",
			Name:              "1.0 TSI 115",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged",
			Description:       "1.0 TSI serisinin güncel 115 PS versiyonu.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
			},
			ProductionStartYear: 2015,
			ProductionEndYear:   intPtr(2020),
			CommonProblems: []EngineProblem{
				{
					Title:       "Su Pompası Kroniği",
					Description: "EA211 serisinin genel sorunu - su pompası sızıntısı",
					Severity:    "high",
				},
			},
			Solutions:         []EngineSolution{},
			MaintenanceNotes:  "Düzenli bakım ve su pompası takibi önemlidir.",
			ReliabilityRating: 6,
		},
		// 5. 1.0 eTSI (Hafif Hibrit)
		{
			Code:              "DLAA",
			Name:              "1.0 eTSI 110 (MHEV)",
			Manufacturer:      "VW Group",
			EngineType:        "I3",
			FuelType:          "Petrol",
			DisplacementCC:    999,
			Cylinders:         3,
			CylinderLayout:    "Inline",
			ValvesPerCylinder: 4,
			Aspiration:        "Turbocharged + 48V MHEV",
			Description:       "1.0 TSI motorun yanına 48 Volt'luk lityum-iyon batarya ve kayışlı marş jeneratörü (BSG) eklenmiş en güncel versiyon. 110 PS güç, 200 Nm tork üretir.",
			TechnologyFeatures: []TechnologyFeature{
				{Name: "48V Mild Hybrid System", Description: "Hafif hibrit sistemi"},
				{Name: "Belt Starter Generator (BSG)", Description: "Kayışlı marş jeneratörü"},
				{Name: "Turbocharger"},
				{Name: "Direct Injection"},
				{Name: "Energy Recovery", Description: "Frenleme enerjisi geri kazanımı"},
			},
			ProductionStartYear: 2020,
			ProductionEndYear:   nil,
			CommonProblems: []EngineProblem{
				{
					Title:       "Yazılımsal Hatalar (İlk Üretim)",
					Description: "İlk üretim yıllarında (2020-2021) 48V sistemiyle ilgili gösterge panelinde 'Elektrik sistemi arızası' uyarısı veren yazılımsal bug'lar mevcuttu (Güncellemelerle çözüldü)",
					Severity:    "medium",
				},
				{
					Title:       "Start-Stop Hassasiyeti",
					Description: "Araç henüz tam durmadan (20 km/s hızın altına inildiğinde) motoru kapattığı için direksiyon sertleşmesi veya klima performansında anlık düşüş yaşanması",
					Severity:    "low",
				},
			},
			Solutions: []EngineSolution{
				{
					ProblemTitle:  "Yazılımsal Hatalar",
					Solution:      "Yetkili serviste ECU yazılım güncellemesi. 2021 sonrası üretim araçlarda sorun giderilmiştir.",
					EstimatedCost: "Garanti kapsamında ücretsiz / 500-1500 TL",
				},
				{
					ProblemTitle:  "Start-Stop Hassasiyeti",
					Solution:      "Start-stop sistemini kapatmak veya yazılım güncellemesi ile davranış iyileştirmesi yapılabilir.",
					EstimatedCost: "-",
				},
			},
			MaintenanceNotes:  "Mekanik bloğu (DLAA) sorunsuzdur. 48V sistem bakımı özel ekipman gerektirir, yetkili serviste yaptırılmalıdır. Akü ömrü 6-8 yıldır.",
			ReliabilityRating: 7,
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

	log.Println("\n✅ All VW 1.0L engines have been added to the database!")
	log.Println("\n📊 Summary:")
	log.Printf("   - Total engines added: %d", len(engines))
	log.Println("   - Engine codes: CHYA, CHYC, CHYB, DSGB, CHZB, DKLA, CHZC, CHZD, DKRF, DLAA")
	log.Println("\nYou can now view these engines at: http://localhost:5173/guides/engines")
}

func intPtr(i int) *int {
	return &i
}
