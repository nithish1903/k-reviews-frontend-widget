package main

import (
	"fmt"
	"k-reviews-frontend-api/controllers"
	"k-reviews-frontend-api/usecase"

	"k-reviews-frontend-api/controllers/shopify"
	"k-reviews-frontend-api/repository/connectors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	connectors.LoadEnv()
	connectors.LoadLogger()
}

func ValidatePublicKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		publicKey := c.GetHeader("Public-Key")

		// Query MySQL database to validate publicKey and retrieve account_id
		account, err := usecase.GetAccountByKey(publicKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid public key"})
			c.Abort()
			return
		}

		// Store account_id in context
		c.Set("account_id", account.ID)
		c.Next()
	}
}

func main() {
	zap.L().Info("starting the main function")

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/v1/assets", "./templates/assets")

	r.GET("/health", controllers.Health)
	r.GET("/v1/get-site-review-button", shopify.GetReviewsButton)
	r.GET("/v1/get-site-reviews", shopify.GetSiteReviewsDataHandler)

	r.POST("/v1/save-product-reviews", shopify.SaveProductReviewsHandler)
	r.POST("/v1/save-vote", shopify.SaveVotesHandler)
	r.POST("/v1/save-site-reviews", shopify.SaveSiteReviewsHandler)
	r.POST("/v1/save-product-qa", shopify.SaveProductQuestionAnswersHandler)

	v1 := r.Group("/v1")
	v1.Use(ValidatePublicKey())
	{
		v1.GET("/get-product-reviews", shopify.GetProductReviewsDataHandler)
		v1.GET("/get-product-reviews-images", shopify.GetProductReviewImgagesHandler)
		v1.GET("/get-product-review-details", shopify.GetProductReviewDetailsHandler)
		v1.GET("/get-product-review-statistics", shopify.GetProductReviewStatisticsHandler)

		// v1.GET("/get-site-reviews", shopify.GetSiteReviewsDataHandler)
		v1.GET("/get-site-reviews-images", shopify.GetSiteReviewImgagesHandler)
		v1.GET("/get-site-review-details", shopify.GetSiteReviewDetailsHandler)

		v1.GET("/get-product-qa", shopify.GetProductQADataHandler)
		// v1.GET("/get-vote", shopify.GetVotesDataHandler)
	}

	// e.POST("/v1/sign-up", controllers.Signup)
	port := os.Getenv("PORT")
	log.Println("port-->", port)
	if port == "" {
		port = "9024"
	}
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Could not start server. Err: %s", err)
	}
}
