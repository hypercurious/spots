package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"test/db"
)

type Spots struct {
	ID          string
	Name        string
	Website     string
	Coordinates string
	Description string
	Rating      float64
	Distance    float64
}

//Variadic error handler function
func CheckError(err ...error) {
	for _, errx := range err {
		if errx != nil {
			log.Println("Error:", err)
			panic(err)
		}
	}
}

func Circle(db *sql.DB, long float64, lat float64, radius float64) *sql.Rows {
	rows, err := db.Query(`	SELECT id, name, COALESCE(website, '') AS "website", ST_AsText(coordinates) AS "coordinates", COALESCE(description, '') AS "description", rating,
							ST_DISTANCE(coordinates, ST_Point($1, $2, 4326)::geography) AS "distance" FROM "MY_TABLE"
							WHERE ST_DWithin(coordinates, ST_Point($1, $2, 4326)::geography, $3)
							ORDER BY ST_Distance(coordinates,ST_Point($1, $2, 4326)::geography)`, long, lat, radius)
	CheckError(err)
	return rows
}

func Square(db *sql.DB, long float64, lat float64, radius float64) *sql.Rows {
	rows, err := db.Query(`	SELECT id, name, COALESCE(website, '') AS "website", ST_AsText(coordinates) AS "coordinates", COALESCE(description, '') AS "description", rating,
							ST_DISTANCE(coordinates, ST_Point($1, $2, 4326)::geography) AS "distance" FROM "MY_TABLE"
							WHERE ST_DWithin(coordinates, ST_Buffer(ST_Point($1, $2, 4326)::geography, $3,'endcap=square'), $3)
							ORDER BY ST_Distance(coordinates, ST_Point($1, $2, 4326)::geography)`, long, lat, radius)
	CheckError(err)
	return rows
}

func GetAllSpots(db *sql.DB, long float64, lat float64, radius float64, shape string) []Spots {
	var rows *sql.Rows

	switch shape {
	case "circle":
		rows = Circle(db, long, lat, radius)
	case "square":
		rows = Square(db, long, lat, radius)
	default:
		CheckError(errors.New(" Invalid shape! "))
	}

	var spots []Spots
	for rows.Next() {
		spot := Spots{}
		err := rows.Scan(&spot.ID, &spot.Name, &spot.Website, &spot.Coordinates, &spot.Description, &spot.Rating, &spot.Distance)
		CheckError(err)
		spots = append(spots, spot)
	}

	return spots
}

func SortByRating(spots []Spots) {
	sort.Slice(spots, func(i, j int) bool {
		return spots[i].Rating > spots[j].Rating
	})
}

func FindDistance(db *sql.DB, coordinate0 string, coordinate1 string) float64 {
	rows, err := db.Query(`SELECT ST_Distance(ST_Point(ST_X($1),ST_Y($1),4326)::geography, ST_Point(ST_X($2),ST_Y($2),4326)::geography)`, coordinate0, coordinate1)
	CheckError(err)
	var distance float64
	for rows.Next() {
		err = rows.Scan(&distance)
		CheckError(err)
	}
	return distance
}

func GetSortedSpots(db *sql.DB, spots []Spots) []Spots {
	var nearSpots []Spots
	var endpoints []Spots
	start := 0
	for i := range spots[:len(spots)-1] {
		distance := FindDistance(db, spots[start].Coordinates, spots[i+1].Coordinates)
		if distance <= 50 {
			if start == i {
				nearSpots = append(nearSpots, spots[start])
			}
			nearSpots = append(nearSpots, spots[i+1])
			continue
		}
		if nearSpots != nil {
			SortByRating(nearSpots)
			endpoints = append(endpoints, nearSpots...)
			nearSpots = nil
			continue
		}
		start = i + 1
		endpoints = append(endpoints, spots[i])
	}
	return endpoints
}

func main() {

	db := db.SetupDB()
	defer db.Close()

	longitude := 23.79847
	latitude := 37.97839
	radius := 10000.0
	shape := "square"

	spots := GetAllSpots(db, longitude, latitude, radius, shape)
	if spots == nil {
		CheckError(errors.New(" The spot list is empty! "))
	}
	endpoints := GetSortedSpots(db, spots)
	fmt.Println()
	for _, x := range endpoints {
		fmt.Printf("Name: %s, Distance: %f, Rating: %f\n", x.Name, x.Distance, x.Rating)
	}
	fmt.Println()
}
