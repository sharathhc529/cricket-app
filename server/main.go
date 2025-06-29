package main
import (
	"context"
	"net/http"
	"os" // 👈 This line is what was missing
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"golang.org/x/oauth2/google"
	"github.com/gin-contrib/cors"
)

func main() {
	r := gin.Default()

	// 🔐 CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://localhost:5173",
			"https://sharathhc529.github.io",
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}
	log.Printf("Server starting on port %s", port)
	http.ListenAndServe(":"+port, nil)

	r.GET("/entries", func(c *gin.Context) {
		values, err := readSheetData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": values})
	})

	r.POST("/submit", func(c *gin.Context) {
		var body struct {
			Value string `json:"value"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		err := appendToSheet(body.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)

}

func readSheetData() ([][]interface{}, error) {
	ctx := context.Background()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("error reading credentials: %w", err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing credentials: %w", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating Sheets service: %w", err)
	}

	spreadsheetId := "1isWFS031E7-DDv3i1I7xi73okYihWtuViMpxzeISeHM"
	readRange := "Sheet1"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("error reading sheet data: %w", err)
	}

	return resp.Values, nil
}

func appendToSheet(value string) error {
	ctx := context.Background()
	
		fmt.Println("Incoming value:",value)

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return err
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return err
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	spreadsheetId := "1isWFS031E7-DDv3i1I7xi73okYihWtuViMpxzeISeHM"
	writeRange := "Sheet1"

	_, err = srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, &sheets.ValueRange{
		Values: [][]interface{}{{value}},
	}).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()

	return err
}