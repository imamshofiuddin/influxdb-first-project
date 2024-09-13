package main

import (
	"bidding-project/models"

	BidController "bidding-project/controllers/BidController"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.ConnectInfluxDb()

	go BidController.RecordBidderCount()

	r.POST("/api/add-bidder-count", BidController.AddBidder)
	r.GET("/api/get-bidder-count", BidController.GetBidder)

	r.Run()
}
