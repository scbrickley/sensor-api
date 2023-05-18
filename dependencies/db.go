package dependencies

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	DB_USER = "postgres"
	DB_NAME = "sensor-db"
	DB_HOST = "localhost"
	DB_PORT = "5432"
)

type Sensor struct {
	Name string         `json:"name"`
	Lat  float64        `json:"latitude"`
	Lon  float64        `json:"longitude"`
	Tags pq.StringArray `json:"tags"`
}

type SensorDB struct {
	db *sql.DB
}

func initDB() (*SensorDB, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	_, err = db.db.Exec(`
	CREATE TABLE IF NOT EXISTS sensors (
		name VARCHAR NOT NULL UNIQUE,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		tags VARCHAR[] NOT NULL
	)
	`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB() (*SensorDB, error) {
	conString := fmt.Sprintf(
		"user=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_NAME, DB_HOST, DB_PORT,
	)

	db, err := sql.Open("postgres", conString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &SensorDB{db: db}, nil
}

func (db *SensorDB) ListSensors() ([]Sensor, error) {
	rows, err := db.db.Query("SELECT name, latitude, longitude, tags FROM sensors")
	if err != nil {
		return nil, err
	}

	var sensorList []Sensor
	for rows.Next() {
		var (
			name string
			lat  float64
			lon  float64
			tags pq.StringArray
		)
		err = rows.Scan(&name, &lat, &lon, &tags)
		if err != nil {
			return nil, err
		}
		sensorList = append(sensorList, Sensor{
			Name: name, Lat: lat, Lon: lon, Tags: tags,
		})
	}
	return sensorList, nil
}

func (db *SensorDB) InsertSensor(sensor Sensor) (*Sensor, error) {
	var (
		name string
		lat  float64
		lon  float64
		tags pq.StringArray
	)
	err := db.db.QueryRow(`
		INSERT INTO sensors (
			name, latitude, longitude, tags
		) VALUES ($1, $2, $3, $4)
		returning name, latitude, longitude, tags
		`, sensor.Name, sensor.Lat, sensor.Lon, sensor.Tags,
	).Scan(&name, &lat, &lon, &tags)
	if err != nil {
		return nil, err
	}

	// Return the sensor data received from postgres to confirm the data has been entered correctly
	recvSensor := &Sensor{
		Name: name, Lat: lat, Lon: lon, Tags: tags,
	}
	return recvSensor, nil
}

func (db *SensorDB) GetSensorByName(name string) (*Sensor, error) {
	// Names are unique, according to the schema established in `initDB()`, so
	// querying a single row should be fine.
	row := db.db.QueryRow(`
		SELECT name, latitude, longitude, tags FROM sensors
			WHERE name = $1
	`, name)

	var (
		sname string
		lat   float64
		lon   float64
		tags  pq.StringArray
	)
	err := row.Scan(&sname, &lat, &lon, &tags)
	if err != nil {
		return nil, err
	}
	sensor := &Sensor{
		Name: sname, Lat: lat, Lon: lon, Tags: tags,
	}
	return sensor, nil
}

func (db *SensorDB) UpdateSensor(name string, newSensor *Sensor) (*Sensor, error) {
	row := db.db.QueryRow(`
		UPDATE sensors
		SET name = $1, latitude = $2, longitude = $3, tags = $4
			WHERE name = $5
		returning name, latitude, longitude, tags
	`, newSensor.Name, newSensor.Lat, newSensor.Lon, newSensor.Tags, name)

	var (
		sname string
		lat   float64
		lon   float64
		tags  pq.StringArray
	)
	err := row.Scan(&sname, &lat, &lon, &tags)
	if err != nil {
		return nil, err
	}
	sensor := &Sensor{
		Name: sname, Lat: lat, Lon: lon, Tags: tags,
	}
	return sensor, nil
}

func (db *SensorDB) DeleteSensor(name string) (*Sensor, error) {
	row := db.db.QueryRow(`
		DELETE FROM sensors
			WHERE name = $1
		returning name, latitude, longitude, tags
	`, name)

	var (
		sname string
		lat   float64
		lon   float64
		tags  pq.StringArray
	)
	err := row.Scan(&sname, &lat, &lon, &tags)
	if err != nil {
		return nil, err
	}
	sensor := &Sensor{
		Name: sname, Lat: lat, Lon: lon, Tags: tags,
	}
	return sensor, nil
}

type Point struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

func distance(a, b Point) float64 {
	// Vertical and horizontal difference
	vDist := b.Lat - a.Lat
	hDist := b.Lon - a.Lon

	vSquare := math.Pow(vDist, 2)
	hSquare := math.Pow(hDist, 2)

	return math.Sqrt(vSquare + hSquare)
}

func nearestSensorToPoint(p Point, sensors []Sensor) *Sensor {
	if len(sensors) < 1 {
		return nil
	}
	nearestIdx := 0
	smallestDist := 0.0
	for i, s := range sensors {
		dist := distance(p, Point{Lat: s.Lat, Lon: s.Lon})
		if i == 0 {
			smallestDist = dist
		}

		if dist < smallestDist {
			smallestDist = dist
			nearestIdx = i
		}
	}
	return &sensors[nearestIdx]
}
