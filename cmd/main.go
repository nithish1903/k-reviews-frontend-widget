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

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func init() {
	connectors.LoadEnv()
	connectors.LoadLogger()
}

func ValidatePublicKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		publicKey := c.Request().Header.Get("Public-Key")

		// Query MySQL database to validate publicKey and retrieve account_id
		account, err := usecase.GetAccountByKey(publicKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid public key"})
		}

		// Store account_id in context
		c.Set("account_id", account.ID)

		return next(c)
	}
}
func main() {
	zap.L().Info("starting the main function")

	e := echo.New()

	e.GET("/health", controllers.Health)
	e.POST("/v1/save-product-reviews", shopify.SaveProductReviewsHandler)
	e.POST("/v1/save-vote", shopify.SaveVotesHandler)
	e.POST("/v1/save-site-reviews", shopify.SaveSiteReviewsHandler)
	e.POST("/v1/save-product-qa", shopify.SaveProductQuestionAnswersHandler)

	e.GET("/v1/get-product-reviews", shopify.GetProductReviewsDataHandler, ValidatePublicKey)
	e.GET("/v1/get-product-reviews-images", shopify.GetProductReviewImgagesHandler, ValidatePublicKey)
	e.GET("/v1/get-product-review-details", shopify.GetProductReviewDetailsHandler, ValidatePublicKey)
	e.GET("/v1/get-product-review-statistics", shopify.GetProductReviewStatisticsHandler, ValidatePublicKey)

	e.GET("/v1/get-site-reviews", shopify.GetSiteReviewsDataHandler, ValidatePublicKey)
	e.GET("/v1/get-site-reviews-images", shopify.GetSiteReviewImgagesHandler, ValidatePublicKey)
	e.GET("/v1/get-site-review-details", shopify.GetSiteReviewDetailsHandler, ValidatePublicKey)

	e.GET("/v1/get-product-qa", shopify.GetProductQADataHandler, ValidatePublicKey)
	// e.GET("/v1/get-vote", shopify.GetVotesDataHandler)

	// e.POST("/v1/sign-up", controllers.Signup)
	port := os.Getenv("PORT")
	log.Println("port-->", port)
	if port == "" {
		port = "9024"
	}
	if err := e.Start(fmt.Sprintf(":%s", port)); err != http.ErrServerClosed {
		log.Fatalf("Could not start server. Err: %s", err)
	}
}
