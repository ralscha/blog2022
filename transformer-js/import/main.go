package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type country struct {
	name             string
	area             int
	areaLand         int
	areaWater        int
	population       int
	populationGrowth float64
	birthRate        float64
	deathRate        float64
	migrationRate    float64
	flagDescription  string
}

func main() {
	conn, err := sqlite.OpenConn("./countries.sqlite3", sqlite.OpenReadWrite, sqlite.OpenCreate)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = sqlitex.ExecuteTransient(conn, "DROP TABLE IF EXISTS countries", nil)
	if err != nil {
		panic(err)
	}

	err = sqlitex.ExecuteTransient(conn, `CREATE TABLE countries (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,

    area       INTEGER,
    area_land  INTEGER,
    area_water INTEGER,

    population        INTEGER,
    population_growth REAL,
    birth_rate        REAL,
    death_rate        REAL,
    migration_rate    REAL,

    flag_description TEXT
  )`, nil)
	if err != nil {
		panic(err)
	}

	inputFolder := "../factbook.json"
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			if file.Name() == "world" || file.Name() == "oceans" || file.Name() == "meta" || file.Name() == "antarctica" {
				continue
			}
			subFolder := inputFolder + "/" + file.Name()
			subFiles, err := os.ReadDir(subFolder)
			if err != nil {
				panic(err)
			}
			for _, subFile := range subFiles {
				if strings.HasSuffix(subFile.Name(), ".json") {
					country := extractCountry(subFolder + "/" + subFile.Name())
					if country != nil {
						insertCountry(conn, country)
					}
				}
			}
		}
	}
}

func insertCountry(conn *sqlite.Conn, country *country) {
	err := sqlitex.Execute(conn, `INSERT INTO countries (
    name,
    area,
    area_land,
    area_water,
    population,
    population_growth,
    birth_rate,
    death_rate,
    migration_rate,
    flag_description
  ) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
  )`, &sqlitex.ExecOptions{
		Args: []any{country.name, country.area, country.areaLand, country.areaWater, country.population, country.populationGrowth, country.birthRate, country.deathRate, country.migrationRate, country.flagDescription},
	})
	if err != nil {
		panic(err)
	}
}

func extractCountry(fileName string) *country {
	content, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	in := map[string]any{}
	err = json.NewDecoder(strings.NewReader(string(content))).Decode(&in)
	if err != nil {
		panic(err)

	}

	c := &country{}

	if government, ok := in["Government"].(map[string]any); ok {
		if countryName, ok := government["Country name"].(map[string]any); ok {
			if shortForm, ok := countryName["conventional short form"].(map[string]any); ok {
				text := shortForm["text"].(string)
				c.name = text
			}
		}
		if flag, ok := government["Flag description"].(map[string]any); ok {
			c.flagDescription = flag["text"].(string)
		}
	}

	if people, ok := in["People and Society"].(map[string]any); ok {
		if population, ok := people["Population"].(map[string]any); ok {
			if total, ok := population["total"].(map[string]any); ok {
				totalText := total["text"].(string)
				c.population = num(totalText)
			}
		}
		if populationGrowth, ok := people["Population growth rate"].(map[string]any); ok {
			text := populationGrowth["text"].(string)
			c.populationGrowth = percent(text)
		}
		if birthRate, ok := people["Birth rate"].(map[string]any); ok {
			text := birthRate["text"].(string)
			c.birthRate = ratePerThousand(text)
		}
		if deathRate, ok := people["Death rate"].(map[string]any); ok {
			text := deathRate["text"].(string)
			c.deathRate = ratePerThousand(text)
		}
		if migrationRate, ok := people["Net migration rate"].(map[string]any); ok {
			text := migrationRate["text"].(string)
			c.migrationRate = ratePerThousand(text)
		}
	}

	if geography, ok := in["Geography"].(map[string]any); ok {
		if area, ok := geography["Area"].(map[string]any); ok {
			if total, ok := area["total "].(map[string]any); ok {
				text := total["text"].(string)
				c.area = sqKm(text)
			}
			if land, ok := area["land"].(map[string]any); ok {
				text := land["text"].(string)
				c.areaLand = sqKm(text)
			}
			if water, ok := area["water"].(map[string]any); ok {
				text := water["text"].(string)
				c.areaWater = sqKm(text)
			}
		}
	}

	if c.name == "" || c.area == 0 || c.population == 0 || c.flagDescription == "" {
		return nil
	}
	return c
}

var perThousandRegex = regexp.MustCompile("([0-9.]+) [a-z()]+/1,000")

func ratePerThousand(text string) float64 {
	match := perThousandRegex.FindStringSubmatch(text)
	if len(match) > 1 {
		return parseFloat(match[1])
	}

	fmt.Println("*** warn: unknown rate <name>/1,000 format (no match): >" + text + "<")
	return 0.0
}

var sqKmRegex = regexp.MustCompile("([0-9,.]+) sq km")

func sqKm(text string) int {
	match := sqKmRegex.FindStringSubmatch(text)
	if len(match) > 1 {
		return num(match[1])
	}

	fmt.Println("*** warn: unknown sq km format (no match): >" + text + "<")
	return 0.0
}

var numRegex = regexp.MustCompile("([0-9,.]+)")

func num(text string) int {
	match := numRegex.FindStringSubmatch(text)
	if len(match) > 1 {
		s := strings.ReplaceAll(match[1], ",", "")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return i
	}

	fmt.Println("*** warn: unknown number format (no match): >" + text + "<")
	return 0
}

var percentRegex = regexp.MustCompile("([0-9.]+)%")

func percent(text string) float64 {
	match := percentRegex.FindStringSubmatch(text)
	if len(match) > 1 {
		return parseFloat(match[1])
	}

	fmt.Println("*** warn: unknown percent format (no match): >" + text + "<")
	return 0.0
}

func parseFloat(text string) float64 {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		fmt.Println(err)
		return 0.0
	}
	return f
}
