package usecase

import (
	"encoding/json"
	"fmt"
	"k-reviews-frontend-api/entity"
	"os"

	goshopify "github.com/bold-commerce/go-shopify"
)

func FetchProductInfo(productId int64, account entity.Account) (entity.ProductInfo, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	app := goshopify.App{
		ApiKey:      clientID,
		ApiSecret:   clientSecret,
		RedirectUrl: "http://localhost:9024/v1/shopify/redirect",
		Scope:       "read_content,write_content,read_themes,write_themes,read_products,read_customers,read_orders,write_script_tags",
	}

	settings := entity.Settings{}
	err := json.Unmarshal([]byte(account.Settings), &settings)
	if err != nil {
		return entity.ProductInfo{}, err
	}

	clientInfo := goshopify.NewClient(app, account.WebUrl, settings.AccessToken)

	product, err := clientInfo.Product.Get(productId, nil)
	if err != nil {
		return entity.ProductInfo{}, err
	}
	var productURL = ""
	if product.Handle != "" {
		productURL = fmt.Sprintf("%s/products/%s", account.WebUrl, product.Handle)
	}

	if product != nil {
		productInfo := entity.ProductInfo{
			AccountID: account.ID,
			ProductID: product.ID,
			Name:      product.Title,
			PageURL:   productURL,
			Status:    "1",
		}
		return productInfo, nil
	}
	return entity.ProductInfo{}, nil
}
