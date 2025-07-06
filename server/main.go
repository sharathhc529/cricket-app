package main
import (
   "context"
   "net/http"
   "os"
   "fmt"
   "log"
   "time"
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
	   AllowAllOrigins: true,
	   AllowMethods:     []string{"GET", "POST", "OPTIONS"},
	   AllowHeaders:     []string{"*"},
	   ExposeHeaders:    []string{"*"},
	   AllowCredentials: true,
	   MaxAge:           12 * time.Hour,
   }))

	// Add this OPTIONS handler for preflight requests
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

   port := os.Getenv("PORT")
   if port == "" {
	   port = "8080" // fallback for local dev
   }
   log.Printf("Server starting on port %s", port)

   r.GET("/entries", func(c *gin.Context) {
		values, err := readSheetData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": values,
		})

	})

   r.POST("/submit", func(c *gin.Context) {
	   // Try multi-input first
	   var multi struct {
		   Player string `json:"player"`
		   Score  float64 `json:"score"`
	   }
	   var single struct {
		   Value string `json:"value"`
	   }

	   if err := c.ShouldBindJSON(&multi); err == nil && multi.Player != "" {
		   err := appendToSheet(multi.Player, fmt.Sprintf("%v", multi.Score))
		   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
		   }
		   c.JSON(http.StatusOK, gin.H{"status": "success"})
		   return
	   }

	   if err := c.ShouldBindJSON(&single); err == nil && single.Value != "" {
		   err := appendToSheet(single.Value, "")
		   if err != nil {
			   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   return
		   }
		   c.JSON(http.StatusOK, gin.H{"status": "success"})
		   return
	   }

	   c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
   })

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

func appendToSheet(player string, score string) error {

	ctx := context.Background()
	
	fmt.Println("Incoming data:", player, score)

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

	var values [][]interface{}
	if score == "" {
		values = [][]interface{}{{player}}
	} else {
		values = [][]interface{}{{player, score}}
	}
	_, err = srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
	return err
}