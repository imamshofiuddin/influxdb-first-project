package models

import (
	"log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/joho/godotenv"
)

var Client influxdb2.Client
var WriteAPI api.WriteAPIBlocking

func ConnectInfluxDb() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("INFLUXDB_TOKEN")

	url := "http://localhost:8086"
	Client = influxdb2.NewClient(url, token)

	org := "PENS"
	bucket := "tes-project"
	WriteAPI = Client.WriteAPIBlocking(org, bucket)
}
