package shopify

import (
	"k-reviews-frontend-api/constant"
	"k-reviews-frontend-api/entity"
	"k-reviews-frontend-api/repository/mongodb"
	"k-reviews-frontend-api/usecase"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func GetReviewsButton(c *gin.Context) {

	buttonHTML := `<button data-toggle="modal" data-target="#static-modal" class="z-[2] border-0 rotate-90 fixed top-[50%] -right-[55px] text-center px-6 py-3 bg-[#252338] flex rounded-b-[10px] text-[18px] font-[500] text-[#FFFFFF]">
        <div class="h-[20px] w-[20px] mr-2.5">
            <img src="/Users/kspl-ragavedhra/Documents/Work/Projects/k-reviews-frontend-widget/templates/assets/images/star.svg" alt="start" width="100%" height="auto" />
        </div>
        <span>REVIEWS</span>
    </button>`

	// Send the HTML response
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, buttonHTML)

	// Sample review data
	// reviews := []entity.ProductReview{
	// 	{Title: "John Doe", AccountID: "1", Content: "Great product, highly recommended!"},
	// 	{Title: "Jane Smith", AccountID: "1", Content: "Very satisfied with my purchase."},
	// 	// Add more reviews as needed
	// }

	// // Render the review widget template with review data
	// c.HTML(http.StatusOK, "site_reviews_modal.html", gin.H{
	// 	"Reviews": reviews,
	// })
}

func OpenReviewModalHandler(c *gin.Context) {
	// Sample review data
	reviews := []entity.ProductReview{
		{Title: "John Doe", AccountID: "1", Content: "Great product, highly recommended!"},
		{Title: "Jane Smith", AccountID: "1", Content: "Very satisfied with my purchase."},
		// Add more reviews as needed
	}

	// Render the review widget template with review data
	c.HTML(http.StatusOK, "site_reviews_modal.html", gin.H{
		"Reviews": reviews,
	})

}
func SaveProductReviewsHandler(c *gin.Context) {
	zap.L().Info("SaveProductReviewsHandler called..!")

	// Parse JSON request body into a Person struct.
	mongoCon := mongodb.MongoConnect()

	var productReview entity.ProductReview
	if err := c.Bind(&productReview); err != nil {
		zap.L().Error("Could not insert product reviews data", zap.Any("error:", err))
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Bad request",
			Status:     true,
			StatusCode: 400,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusBadRequest, response)
	}

	account, err := usecase.GetAccountById(productReview.AccountID)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Account not found",
			Status:     true,
			StatusCode: 404,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusNotFound, response)
	}
	zap.L().Info("account details:", zap.Any("account:", account))

	account, err = usecase.GetAccountById(productReview.AccountID)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Account not found",
			Status:     true,
			StatusCode: 404,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusNotFound, response)
	}

	productInfo, err := usecase.FetchProductInfo(productReview.ProductID, *account)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Error fetching product info",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}

	bsonProductInfo, _ := bson.Marshal(productInfo)
	insert := mongodb.UpsertOne(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_COLLECTION, bson.M{"product_id": productInfo.ProductID}, bsonProductInfo)
	if !insert {
		zap.L().Info("Could not insert product data")
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Could not insert product data",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}

	productReview.Status = constant.REVIEW_STATUS_PENDING
	productReview.CreatedAt = time.Now().UTC()

	// Convert Product_reviews struct to BSON document.
	bsonProductReview, err := bson.Marshal(productReview)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Internal server error",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}

	insert = mongodb.InsertOne(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_REVIEWS_COLLECTION, bsonProductReview)
	if !insert {
		zap.L().Info("Could not insert product reviews data")
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Could not insert product reviews data",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}
	zap.L().Info("Product review saved successfully")

	// Return success response.
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       "Product review saved successfully",
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func SaveVotesHandler(c *gin.Context) {
	zap.L().Info("SaveVotes called..!")

	// Parse JSON request body into a Person struct.
	mongoCon := mongodb.MongoConnect()
	var votes entity.Votes
	if err := c.Bind(&votes); err != nil {
		zap.L().Error("Could not insert vote data", zap.Any("error:", err))
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Bad request",
			Status:     true,
			StatusCode: 400,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusBadRequest, response)
	}

	// Convert Product_reviews struct to BSON document.
	bsonVotes, err := bson.Marshal(votes)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Internal server error",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}
	zap.L().Info("bsonVotes:", zap.Any("bsonVotes:", bsonVotes))

	insert := mongodb.InsertOne(mongoCon.Connection, constant.K_REVIEWS_DB, constant.VOTES_COLLECTION, bsonVotes)
	if !insert {
		zap.L().Info("Could not insert votes data")
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Could not insert votes data",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}
	zap.L().Info("Votes saved successfully")

	// Return success response.
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       "Votes saved successfully",
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func SaveSiteReviewsHandler(c *gin.Context) {
	zap.L().Info("SaveSiteReviewsHandler called..!")

	// Parse JSON request body into a Person struct.
	mongoCon := mongodb.MongoConnect()
	var siteReview entity.SiteReview
	if err := c.Bind(&siteReview); err != nil {
		zap.L().Error("Could not insert site reviews data", zap.Any("error:", err))
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Bad request",
			Status:     true,
			StatusCode: 400,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusBadRequest, response)
	}

	siteReview.Status = constant.REVIEW_STATUS_PENDING
	siteReview.CreatedAt = time.Now().UTC()
	// Convert Product_reviews struct to BSON document.
	bsonSiteReview, err := bson.Marshal(siteReview)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Internal server error",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}
	zap.L().Info("siteReview:", zap.Any("siteReview:", siteReview))

	insert := mongodb.InsertOne(mongoCon.Connection, constant.K_REVIEWS_DB, constant.SITE_REVIEWS_COLLECTION, bsonSiteReview)
	if !insert {
		zap.L().Info("Could not insert site review data")
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Could not insert insert site review data",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)

	}
	zap.L().Info("site review saved successfully")

	// Return success response.
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       "Site review saved successfully",
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func SaveProductQuestionAnswersHandler(c *gin.Context) {
	zap.L().Info("SaveProductQuestionAnswersHandler called..!")

	// Parse JSON request body into a Person struct.
	mongoCon := mongodb.MongoConnect()
	var reviewQuestion entity.ProductQuestionAnswer
	if err := c.Bind(&reviewQuestion); err != nil {
		zap.L().Error("Could not insert product qa data", zap.Any("error:", err))
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Bad request",
			Status:     true,
			StatusCode: 400,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusBadRequest, response)
	}

	account, err := usecase.GetAccountById(reviewQuestion.AccountID)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Account not found",
			Status:     true,
			StatusCode: 404,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusNotFound, response)
	}
	zap.L().Info("account details:", zap.Any("account:", account))

	account, err = usecase.GetAccountById(reviewQuestion.AccountID)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Account not found",
			Status:     true,
			StatusCode: 404,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusNotFound, response)
	}

	productInfo, err := usecase.FetchProductInfo(reviewQuestion.ProductID, *account)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Error fetching product info",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)

	}

	bsonProductInfo, _ := bson.Marshal(productInfo)
	insert := mongodb.UpsertOne(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_COLLECTION, bson.M{"product_id": productInfo.ProductID}, bsonProductInfo)
	if !insert {
		zap.L().Info("Could not insert product data")
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Could not insert product data",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)

	}

	reviewQuestion.Status = constant.REVIEW_STATUS_PENDING
	reviewQuestion.CreatedAt = time.Now().UTC()
	// Convert Product_reviews struct to BSON document.
	bsonReviewQuestion, err := bson.Marshal(reviewQuestion)
	if err != nil {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Internal server error",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}
	zap.L().Info("ReviewQuestion:", zap.Any("ReviewQuestion:", reviewQuestion))

	insert = mongodb.InsertOne(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_QA_COLLECTION, bsonReviewQuestion)
	if !insert {
		zap.L().Info("Could not insert product question and answer data")
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Could not insert product question and answer data",
			Status:     true,
			StatusCode: 500,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusInternalServerError, response)
	}
	zap.L().Info("Product QA saved successfully")

	// Return success response.
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       "Product question and answer saved successfully",
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func GetProductReviewsDataHandler(c *gin.Context) {
	zap.L().Info("GetProductReviewsDataHandler called..!")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	accountId, ok := c.Get("account_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})
	}
	productIdStr := c.Query("productId")
	filterType := c.Query("filterType")

	mongoCon := mongodb.MongoConnect()

	limit := int64(10)
	skip := (int64(page) - 1) * int64(limit)
	var productId int64 = 0
	if productIdStr != "" {
		var err error
		productId, err = strconv.ParseInt(productIdStr, 10, 64)
		if err != nil {
			zap.L().Error("Error parsing productId:", zap.Any("error:", err))
		}
	}
	var filter bson.M
	if productId != 0 {
		filter = bson.M{"account_id": accountId, "product_id": productId}
	} else {
		filter = bson.M{"account_id": accountId}
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$rating"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
		{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "rating", Value: "$_id"}, {Key: "count", Value: 1}}}},
	}
	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(skip)
	switch filterType {
	case "oldest":
		filter["created_at"] = bson.M{"$exists": true}
		sortOption := bson.M{"created_at": -1}
		options.Sort = sortOption
	case "newest":
		filter["created_at"] = bson.M{"$exists": true}
		sortOption := bson.M{"created_at": 1}
		options.Sort = sortOption
	case "highestRating":
		sortOption := bson.M{"rating": -1}
		options.Sort = sortOption
	case "lowestRatings":
		sortOption := bson.M{"rating": 1}
		options.Sort = sortOption
	}

	results, total, err := mongodb.GetReviewsDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_REVIEWS_COLLECTION, constant.VOTES_COLLECTION, filter, options, int64(limit), int64(skip))
	if err != nil {
		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))
	}
	zap.L().Info("ratingCounts:", zap.Any("filter:", filter))

	// Call the GetAggregateDocuments function
	ratingCounts, err := mongodb.GetRatingCounts(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_REVIEWS_COLLECTION, pipeline)
	if err != nil {
		zap.L().Error("Error while fetching the ratings count:", zap.Any("error:", err))
	}

	zap.L().Info("ratingCounts:", zap.Any("ratingCounts:", ratingCounts))
	totalScore := 0
	for rating, count := range ratingCounts {
		totalScore += int(rating) * int(count)
	}

	// Calculate total votes
	totalVotes := 0
	for _, count := range ratingCounts {
		totalVotes += int(count)
	}
	averageRating := float64(0)
	if totalVotes != 0 {
		averageRating = float64(totalScore) / float64(totalVotes)
		averageRating = math.Round(averageRating*100) / 100
	}

	pagination := &entity.Pagination{
		Total:   total,
		PerPage: limit,
		Page:    int64(page),
	}
	zap.L().Info("sending response:")
	zap.L().Info("results:", zap.Any("response:", results))
	zap.L().Info("RatingCounts:", zap.Any("response:", ratingCounts))
	zap.L().Info("averageRating:", zap.Any("response:", averageRating))

	response := entity.ReviewResponse{
		Message: "Success",
		Data: struct {
			Results       interface{}     `json:"results,omitempty"`
			RatingCounts  map[int32]int32 `json:"ratings_counts,omitempty"` // Assuming ratingCounts is a map[string]int32
			AverageRating float64         `json:"avergae_ratings,omitempty"`
		}{results, ratingCounts, averageRating},
		Pagination: pagination,
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	zap.L().Info("sending sent:", zap.Any("response:", response))

	c.JSON(http.StatusOK, response)
}

func GetProductReviewDetailsHandler(c *gin.Context) {
	zap.L().Info("GetProductReviewDetailsHandler called..!")

	accountId, ok := c.Get("account_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})

	}
	productReviewIdStr := c.Query("id")

	mongoCon := mongodb.MongoConnect()
	productReviewId, err := primitive.ObjectIDFromHex(productReviewIdStr)
	if err != nil {
		zap.L().Error("Error converting string to objectId:", zap.Any("error:", err))
	}
	var filter bson.M
	if productReviewIdStr != "" {
		filter = bson.M{"account_id": accountId, "_id": productReviewId}
	} else {
		response := entity.ReviewResponse{
			Message:    "Validation error",
			Data:       "Please enter review id",
			Status:     true,
			StatusCode: 400,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusBadRequest, response)
	}

	results, err := mongodb.GetDocuments(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_REVIEWS_COLLECTION, filter)
	if err != nil {
		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))
	}
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       results,
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func GetProductReviewImgagesHandler(c *gin.Context) {
	zap.L().Info("GetProductReviewImgagesHandler called..!")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	accountId, ok := c.Get("account_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})

	}
	productIdStr := c.Query("productId")

	mongoCon := mongodb.MongoConnect()

	limit := int64(10)
	skip := (int64(page) - 1) * int64(limit)

	var productId int64 = 0
	if productIdStr != "" {
		var err error
		productId, err = strconv.ParseInt(productIdStr, 10, 64)
		if err != nil {
			zap.L().Error("Error parsing productId:", zap.Any("error:", err))
		}
	}
	projection := bson.M{"media_uploads": 1}
	var filter bson.M
	if productId != 0 {
		filter = bson.M{"account_id": accountId, "product_id": productId}
	} else {
		filter = bson.M{"account_id": accountId}
	}
	options := options.Find().SetLimit(limit).SetSkip(skip)

	results, total, err := mongodb.GetDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_REVIEWS_COLLECTION, filter, projection, options, int64(limit), int64(skip))
	if err != nil {
		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))
	}
	pagination := &entity.Pagination{
		Total:   total,
		PerPage: limit,
		Page:    int64(page),
	}
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       results,
		Pagination: pagination,
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func GetProductReviewStatisticsHandler(c *gin.Context) {
	zap.L().Info("GetProductReviewStatisticsHandler called..!")

	accountId, ok := c.Get("account_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})

	}
	productIdStr := c.Query("productId")

	mongoCon := mongodb.MongoConnect()

	var productId int64 = 0
	if productIdStr != "" {
		var err error
		productId, err = strconv.ParseInt(productIdStr, 10, 64)
		if err != nil {
			zap.L().Error("Error parsing productId:", zap.Any("error:", err))
		}
	}
	var filter bson.M
	if productIdStr != "" {
		filter = bson.M{"account_id": accountId, "product_id": productId}
	} else {
		response := entity.ReviewResponse{
			Message:    "Failed",
			Data:       "Please enter product id",
			Status:     true,
			StatusCode: 400,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusBadRequest, response)
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$rating"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
		{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "rating", Value: "$_id"}, {Key: "count", Value: 1}}}},
	}

	// Call the GetAggregateDocuments function
	ratingCounts, err := mongodb.GetRatingCounts(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_REVIEWS_COLLECTION, pipeline)
	if err != nil {
		zap.L().Error("Error while fetching the ratings count:", zap.Any("error:", err))
	}

	zap.L().Info("ratingCounts:", zap.Any("ratingCounts:", ratingCounts))
	totalScore := 0
	for rating, count := range ratingCounts {
		totalScore += int(rating) * int(count)
	}

	// Calculate total votes
	totalVotes := 0
	for _, count := range ratingCounts {
		totalVotes += int(count)
	}
	averageRating := float64(0)
	if totalVotes != 0 {
		averageRating = float64(totalScore) / float64(totalVotes)
		averageRating = math.Round(averageRating*100) / 100
	}

	zap.L().Info("sending response:")
	zap.L().Info("RatingCounts:", zap.Any("response:", ratingCounts))
	zap.L().Info("averageRating:", zap.Any("response:", averageRating))

	response := entity.ReviewResponse{
		Message: "Success",
		Data: struct {
			RatingCounts  map[int32]int32 `json:"ratings_counts,omitempty"`
			AverageRating float64         `json:"avergae_ratings,omitempty"`
		}{ratingCounts, averageRating},
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	zap.L().Info("sending sent:", zap.Any("response:", response))

	c.JSON(http.StatusOK, response)
}

func GetSiteReviewsDataHandler(c *gin.Context) {
	zap.L().Info("GetSiteReviewsDataHandler called..!")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// accountId, ok := c.Get("account_id")
	// if !ok {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})

	// }
	accountId := "1"
	productIdStr := c.Query("productId")
	filterType := c.Query("filterType")

	mongoCon := mongodb.MongoConnect()

	limit := int64(8)
	skip := (int64(page) - 1) * int64(limit)
	var productId int64 = 0
	if productIdStr != "" {
		var err error
		productId, err = strconv.ParseInt(productIdStr, 10, 64)
		if err != nil {
			zap.L().Error("Error parsing productId:", zap.Any("error:", err))
		}
	}
	var filter bson.M
	if productId != 0 {
		filter = bson.M{"account_id": accountId, "product_id": productId, "status": constant.REVIEW_STATUS_APPROVED}
	} else {
		filter = bson.M{"account_id": accountId, "status": constant.REVIEW_STATUS_APPROVED}
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$rating"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
		{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "rating", Value: "$_id"}, {Key: "count", Value: 1}}}},
	}
	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(skip)
	switch filterType {
	case "oldest":
		sortOption := bson.M{"created_at": -1}
		options.Sort = sortOption
	case "newest":
		sortOption := bson.M{"created_at": 1}
		options.Sort = sortOption
	case "highestRating":
		sortOption := bson.M{"rating": -1}
		options.Sort = sortOption
	case "lowestRatings":
		sortOption := bson.M{"rating": 1}
		options.Sort = sortOption
	}

	results, total, err := mongodb.GetReviewsDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.SITE_REVIEWS_COLLECTION, constant.VOTES_COLLECTION, filter, options, int64(limit), int64(skip))
	if err != nil {
		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))
	}
	zap.L().Info("ratingCounts:", zap.Any("filter:", filter))

	// Call the GetAggregateDocuments function
	ratingCounts, err := mongodb.GetRatingCounts(mongoCon.Connection, constant.K_REVIEWS_DB, constant.SITE_REVIEWS_COLLECTION, pipeline)
	if err != nil {
		zap.L().Error("Error while fetching the ratings count:", zap.Any("error:", err))
	}

	zap.L().Info("ratingCounts:", zap.Any("ratingCounts:", ratingCounts))
	totalScore := 0
	for rating, count := range ratingCounts {
		totalScore += int(rating) * int(count)
	}

	// Calculate total votes
	totalVotes := 0
	for _, count := range ratingCounts {
		totalVotes += int(count)
	}
	averageRating := float64(0)
	if totalVotes != 0 {
		averageRating = float64(totalScore) / float64(totalVotes)
		averageRating = math.Round(averageRating*100) / 100
	}

	pagination := entity.Pagination{
		Total:   total,
		PerPage: limit,
		Page:    int64(page),
	}
	zap.L().Info("sending response:")
	zap.L().Info("results:", zap.Any("response:", results))
	zap.L().Info("RatingCounts:", zap.Any("response:", ratingCounts))
	zap.L().Info("averageRating:", zap.Any("response:", averageRating))
	zap.L().Info("pagination:", zap.Any("pagination:", pagination))

	// response := entity.ReviewResponse{
	// 	Message: "Success",
	// 	Data: struct {
	// 		Results       interface{}     `json:"results,omitempty"`
	// 		RatingCounts  map[int32]int32 `json:"ratings_counts,omitempty"` // Assuming ratingCounts is a map[string]int32
	// 		AverageRating float64         `json:"avergae_ratings,omitempty"`
	// 	}{results, ratingCounts, averageRating},
	// 	Pagination: pagination,
	// 	Status:     true,
	// 	StatusCode: 200,
	// 	Error:      nil,
	// 	Timestamp:  time.Now().UTC(),
	// }
	// c.JSON(http.StatusOK, response)

	if page == 1 {
		// Render the initial page with the complete HTML
		c.HTML(http.StatusOK, "site_reviews_modal.html", gin.H{
			"Reviews":       results,
			"RatingCounts":  ratingCounts,
			"AverageRating": averageRating,
			"Pagination":    pagination,
		})
	} else {
		// Render only the reviews div for subsequent pages
		c.HTML(http.StatusOK, "reviews_div.html", gin.H{
			"Reviews": results,
		})
	}
}

func GetSiteReviewImgagesHandler(c *gin.Context) {
	zap.L().Info("GetProductReviewImgagesHandler called..!")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	accountId, ok := c.Get("account_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})

	}
	productIdStr := c.Query("productId")

	mongoCon := mongodb.MongoConnect()

	limit := int64(10)
	skip := (int64(page) - 1) * int64(limit)

	var productId int64 = 0
	if productIdStr != "" {
		var err error
		productId, err = strconv.ParseInt(productIdStr, 10, 64)
		if err != nil {
			zap.L().Error("Error parsing productId:", zap.Any("error:", err))

		}
	}
	projection := bson.M{"media_uploads": 1}
	var filter bson.M
	if productId != 0 {
		filter = bson.M{"account_id": accountId, "product_id": productId}
	} else {
		filter = bson.M{"account_id": accountId}
	}
	options := options.Find().SetLimit(limit).SetSkip(skip)

	results, total, err := mongodb.GetDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_REVIEWS_COLLECTION, filter, projection, options, int64(limit), int64(skip))
	if err != nil {
		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))

	}
	pagination := &entity.Pagination{
		Total:   total,
		PerPage: limit,
		Page:    int64(page),
	}
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       results,
		Pagination: pagination,
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func GetSiteReviewDetailsHandler(c *gin.Context) {
	zap.L().Info("GetSiteReviewDetailsHandler called..!")

	accountId, ok := c.Get("account_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})

	}
	siteReviewIdStr := c.Query("id")

	mongoCon := mongodb.MongoConnect()
	siteReviewId, err := primitive.ObjectIDFromHex(siteReviewIdStr)
	if err != nil {
		zap.L().Error("Error converting string to objectId:", zap.Any("error:", err))
	}
	var filter bson.M
	if siteReviewIdStr != "" {
		filter = bson.M{"account_id": accountId, "_id": siteReviewId}
	} else {
		response := entity.ReviewResponse{
			Message:    "Validation error",
			Data:       "Please enter review id",
			Status:     true,
			StatusCode: 400,
			Error:      nil,
			Timestamp:  time.Now().UTC(),
		}
		c.JSON(http.StatusBadRequest, response)
	}

	results, err := mongodb.GetDocuments(mongoCon.Connection, constant.K_REVIEWS_DB, constant.SITE_REVIEWS_COLLECTION, filter)
	if err != nil {
		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))

	}
	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       results,
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

func GetProductQADataHandler(c *gin.Context) {
	zap.L().Info("GetProductQADataHandler called..!")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	accountId, ok := c.Get("account_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account ID"})

	}
	productIdStr := c.Query("productId")

	mongoCon := mongodb.MongoConnect()

	limit := int64(10)
	skip := (int64(page) - 1) * int64(limit)
	var productId int64 = 0
	if productIdStr != "" {
		var err error
		productId, err = strconv.ParseInt(productIdStr, 10, 64)
		if err != nil {
			zap.L().Error("Error parsing productId:", zap.Any("error:", err))

		}
	}
	var filter bson.M
	if productId != 0 {
		filter = bson.M{"account_id": accountId, "product_id": productId}
	} else {
		filter = bson.M{"account_id": accountId}
	}

	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(skip)

	results, total, err := mongodb.GetProductReviewsQaDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_QA_COLLECTION, constant.VOTES_COLLECTION, filter, options, int64(limit), int64(skip))
	if err != nil {
		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))

	}
	zap.L().Info("ratingCounts:", zap.Any("filter:", filter))

	pagination := &entity.Pagination{
		Total:   total,
		PerPage: limit,
		Page:    int64(page),
	}
	zap.L().Info("sending response:")
	zap.L().Info("results:", zap.Any("response:", results))

	response := entity.ReviewResponse{
		Message:    "Success",
		Data:       results,
		Pagination: pagination,
		Status:     true,
		StatusCode: 200,
		Error:      nil,
		Timestamp:  time.Now().UTC(),
	}
	zap.L().Info("sending sent:", zap.Any("response:", response))

	c.JSON(http.StatusOK, response)
}

// func GetSiteReviewsDataHandler(c *gin.Context){
// 	zap.L().Info("GetSiteReviewsDataHandler called..!")

// 	page, err := strconv.Atoi(c.Query("page"))
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	accountId, _ := strconv.Atoi(c.Get("account_id").(string))
// 	mongoCon := mongodb.MongoConnect()

// 	limit := 10
// 	skip := (page - 1) * limit

// 	zap.L().Info("accountId:", zap.Any("results:", accountId))
// 	zap.L().Info("limit:", zap.Any("results:", limit))
// 	zap.L().Info("skip:", zap.Any("results:", skip))

// 	// Perform aggregation
// 	results, total, err := mongodb.GetDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.SITE_REVIEWS_COLLECTION, bson.M{"account_id": fmt.Sprint(accountId)}, int64(limit), int64(skip))
// 	if err != nil {
// 		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))
//
// 	}
// 	pagination := &entity.Pagination{
// 		Total:   total,
// 		PerPage: limit,
// 		Page:    page,
// 	}
// 	response := entity.ReviewResponse{
// 		Message:    "Success",
// 		Data:       results,
// 		Pagination: pagination,
// 		Status:     true,
// 		StatusCode: 200,
// 		Error:      nil,
// 		Timestamp:  time.Now().UTC(),
// 	}
// 	c.JSON(http.StatusOK, response)
// }

// func GetProductQuestionAndAnswersDataHandler(c *gin.Context){
// 	zap.L().Info("GetProductQuestionAndAnswersDataHandler called..!")

// 	page, err := strconv.Atoi(c.Query("page"))
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	accountId, _ := strconv.Atoi(c.Get("account_id").(string))
// 	mongoCon := mongodb.MongoConnect()

// 	limit := 10
// 	skip := (page - 1) * limit

// 	zap.L().Info("accountId:", zap.Any("results:", accountId))
// 	zap.L().Info("limit:", zap.Any("results:", limit))
// 	zap.L().Info("skip:", zap.Any("results:", skip))

// 	// Perform aggregation
// 	results, total, err := mongodb.GetDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.PRODUCT_QA_COLLECTION, bson.M{"account_id": fmt.Sprint(accountId)}, int64(limit), int64(skip))
// 	if err != nil {
// 		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))
//
// 	}
// 	pagination := &entity.Pagination{
// 		Total:   total,
// 		PerPage: limit,
// 		Page:    page,
// 	}

// 	response := entity.ReviewResponse{
// 		Message:    "Success",
// 		Data:       results,
// 		Pagination: pagination,
// 		Status:     true,
// 		StatusCode: 200,
// 		Error:      nil,
// 		Timestamp:  time.Now().UTC(),
// 	}
// 	c.JSON(http.StatusOK, response)
// }

// func GetVotesDataHandler(c *gin.Context){
// 	zap.L().Info("GetVotesDataHandler called..!")

// 	page, err := strconv.Atoi(c.Query("page"))
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	accountId, _ := strconv.Atoi(c.Get("account_id").(string))
// 	mongoCon := mongodb.MongoConnect()

// 	limit := 10
// 	skip := (page - 1) * limit

// 	zap.L().Info("accountId:", zap.Any("results:", accountId))
// 	zap.L().Info("limit:", zap.Any("results:", limit))
// 	zap.L().Info("skip:", zap.Any("results:", skip))

// 	// Perform aggregation
// 	results, total, err := mongodb.GetDocumentsOnPagination(mongoCon.Connection, constant.K_REVIEWS_DB, constant.VOTES_COLLECTION, bson.M{"account_id": fmt.Sprint(accountId)}, int64(limit), int64(skip))
// 	if err != nil {
// 		zap.L().Error("Error while fetching the aggregate results:", zap.Any("error:", err))
//
// 	}
// 	pagination := &entity.Pagination{
// 		Total:   total,
// 		PerPage: limit,
// 		Page:    page,
// 	}

// 	response := entity.ReviewResponse{
// 		Message:    "Success",
// 		Data:       results,
// 		Pagination: pagination,
// 		Status:     true,
// 		StatusCode: 200,
// 		Error:      nil,
// 		Timestamp:  time.Now().UTC(),
// 	}
// 	c.JSON(http.StatusOK, response)
// }
