package locations

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type District struct {
	ID     int    `bson:"id"`
	CityID int    `bson:"city_id"`
	Name   string `bson:"name"`
}

type Location struct {
	ProvinceID int          `bson:"province_id" json:"province_id"`
	Name       string       `bson:"name" json:"name"`
	Cities     map[int]City `bson:"cities" json:"cities,omitempty"`
}

type City struct {
	ID         int              `bson:"id" json:"id"`
	ProvinceID int              `bson:"province_id" json:"province_id"`
	Name       string           `bson:"name" json:"name"`
	Districts  map[int]District `bson:"districts" json:"districts,omitempty"`
}

func PopulateData(db *mongo.Database) {
	districtCollection := db.Collection("locations")

	provinces := getProvinces()

	// insert per 500
	// page := len(districts) / 500
	// lastPageNoOfData := len(districts) % 500

	fmt.Println("starting indexing", len(provinces))
	for key, province := range provinces {
		fmt.Println(key)
		_, err := districtCollection.InsertOne(context.Background(), province, nil)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func getProvinces() map[int]interface{} {
	// URL to the CSV file
	url := "https://raw.githubusercontent.com/emsifa/api-wilayah-indonesia/refs/heads/master/data/provinces.csv"

	// Fetch the CSV data from the URL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch data, status code: %d", resp.StatusCode)
	}

	// Create a CSV reader
	reader := csv.NewReader(resp.Body)
	reader.Comma = ',' // default is ',' so this is optional
	provinces := make(map[int]Location)
	// Read all records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
		}

		// Check if the record has the expected number of columns
		if len(record) < 2 {
			log.Printf("Skipping incomplete record: %v", record)
			continue
		}

		// Create a District struct from the record
		provinceId, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			panic(err)
		}
		province := Location{
			ProvinceID: provinceId,
			Name:       strings.TrimSpace(record[1]),
		}
		fmt.Println(provinces)
		provinces[provinceId] = province
	}

	fmt.Println("done building provinces")

	url = "https://raw.githubusercontent.com/emsifa/api-wilayah-indonesia/refs/heads/master/data/regencies.csv"

	// Fetch the CSV data from the URL
	resp, err = http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch data, status code: %d", resp.StatusCode)
	}

	// Create a CSV reader
	reader = csv.NewReader(resp.Body)
	reader.Comma = ',' // default is ',' so this is optional

	cities := []City{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
		}

		// Check if the record has the expected number of columns
		if len(record) < 3 {
			log.Printf("Skipping incomplete record: %v", record)
			continue
		}

		// Create a District struct from the record
		provinceId, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			panic(err)
		}

		cityId, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			panic(err)
		}
		city := City{
			ID:         cityId,
			ProvinceID: provinceId,
			Name:       strings.TrimSpace(record[2]),
		}
		cities = append(cities, city)
	}

	fmt.Println("Looping through cities")

	cityMap := make(map[int]City)
	for i := 0; i < len(cities); i++ {
		city := cities[i]
		cityMap[city.ID] = city
	}

	url = "https://raw.githubusercontent.com/emsifa/api-wilayah-indonesia/refs/heads/master/data/districts.csv"

	// Fetch the CSV data from the URL
	resp, err = http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch data, status code: %d", resp.StatusCode)
	}

	// Create a CSV reader
	reader = csv.NewReader(resp.Body)
	reader.Comma = ',' // default is ',' so this is optional
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
		}

		// Check if the record has the expected number of columns
		if len(record) < 3 {
			log.Printf("Skipping incomplete record: %v", record)
			continue
		}

		// Create a District struct from the record
		districtId, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			panic(err)
		}

		cityId, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			panic(err)
		}
		district := District{
			ID:     districtId,
			CityID: cityId,
			Name:   strings.TrimSpace(record[2]),
		}
		city := cityMap[district.CityID]
		if city.Districts == nil {
			city.Districts = make(map[int]District)
		}

		fmt.Println(district)

		city.Districts[districtId] = district
		cityMap[cityId] = city
	}

	for _, city := range cityMap {
		province := provinces[city.ProvinceID]
		if province.Cities == nil {
			province.Cities = make(map[int]City)
		}

		province.Cities[city.ID] = city
		provinces[city.ProvinceID] = province
	}
	toWrite := map[int]interface{}{}

	for _, p := range provinces {
		toWrite[p.ProvinceID] = p
	}

	return toWrite
}
