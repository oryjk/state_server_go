package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"state_server_go/db"
	"state_server_go/routes"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// Initialize the database
	dsn := "root:beifa888@tcp(49.234.55.170:3306)/state_server"
	dbConnection := db.InitDB(dsn)

	// Start queue processor
	go routes.ProcessQueue(dbConnection)

	// Set up Gin router
	r := gin.Default()

	// Define routes
	r.POST("/client/status", routes.ReceiveStatus)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
