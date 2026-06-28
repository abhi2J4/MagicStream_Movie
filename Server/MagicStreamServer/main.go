package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/database"
	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	// This is the main function

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMovies!")
	})

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}

	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	//config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	var client *mongo.Client = database.Connect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to reach server: %v", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}

	}()

	routes.SetupUnProtectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/database"
// 	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/routes"
// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

// func main() {

// 	// ✅ STEP 1: Load .env FIRST before anything else
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Println("Warning: unable to find .env file")
// 		// ✅ Fallback: set manually if .env fails
// 		os.Setenv("MONGODB_URI", "mongodb://localhost:27017/magicstream")
// 		os.Setenv("DATABASE_NAME", "magicstream")
// 		os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
// 	}

// 	// ✅ STEP 2: Now start gin
// 	router := gin.Default()

// 	router.GET("/hello", func(c *gin.Context) {
// 		c.String(200, "Hello, MagicStreamMovies!")
// 	})

// 	// ✅ STEP 3: Read env vars (now they are loaded)
// 	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

// 	var origins []string
// 	if allowedOrigins != "" {
// 		origins = strings.Split(allowedOrigins, ",")
// 		for i := range origins {
// 			origins[i] = strings.TrimSpace(origins[i])
// 			log.Println("Allowed Origin:", origins[i])
// 		}
// 	} else {
// 		origins = []string{"http://localhost:5173"}
// 		log.Println("Allowed Origin: http://localhost:5173")
// 	}

// 	config := cors.Config{}
// 	config.AllowOrigins = origins
// 	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
// 	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
// 	config.ExposeHeaders = []string{"Content-Length"}
// 	config.AllowCredentials = true
// 	config.MaxAge = 12 * time.Hour

// 	router.Use(cors.New(config))
// 	router.Use(gin.Logger())

// 	// ✅ STEP 4: Connect to MongoDB (env vars are ready now)
// 	var client *mongo.Client = database.Connect()

// 	if err := client.Ping(context.Background(), nil); err != nil {
// 		log.Fatalf("Failed to reach server: %v", err)
// 	}
// 	defer func() {
// 		err := client.Disconnect(context.Background())
// 		if err != nil {
// 			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
// 		}
// 	}()

// 	routes.SetupUnProtectedRoutes(router, client)
// 	routes.SetupProtectedRoutes(router, client)

// 	if err := router.Run(":8080"); err != nil {
// 		fmt.Println("Failed to start server", err)
// 	}
// }