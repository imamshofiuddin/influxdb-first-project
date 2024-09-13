package bidcontroller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	models "bidding-project/models"

	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var bidderCount int = 0

func GetBidder(c *gin.Context) {
	query := `from(bucket:"tes-project")|> range(start: -24h) |> filter(fn: (r) => r._measurement == "bidding")`
	queryAPI := models.Client.QueryAPI("PENS")
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying InfluxDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get bidder data"})
		return
	}

	var bidders []gin.H
	for result.Next() {
		if result.Err() != nil {
			log.Printf("Query error: %v", result.Err())
			continue
		}
		bidders = append(bidders, gin.H{
			"hour":         result.Record().ValueByKey("hour"),
			"bidder_count": result.Record().Value(),
			"time":         result.Record().Time(),
		})
	}

	c.JSON(http.StatusOK, bidders)
}

func AddBidder(c *gin.Context) {
	bidderCount++
	c.JSON(http.StatusOK, gin.H{"message": "Bidder count recorded"})
}

func RecordBidderCount() error {
	for {
		time.Sleep(time.Minute)

		count := bidderCount
		bidderCount = 0

		now := time.Now()
		hour := now.Hour()

		// Buat point data
		tags := map[string]string{
			"hour": fmt.Sprintf("%d", hour),
		}
		fields := map[string]interface{}{
			"bidder_count": count,
		}

		// Tulis data ke InfluxDB
		point := write.NewPoint("bidding", tags, fields, time.Now())
		if err := models.WriteAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
	}
}
