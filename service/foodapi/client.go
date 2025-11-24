package foodapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lorenzougolini/wimf-app/service/models"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://world.openfoodfacts.org/api/v0/product",
	}
}

func (c *Client) GetProductByBarcode(barcode string) (models.ProductInfo, error) {
	url := fmt.Sprintf("%s/%s.json", c.baseURL, barcode)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return models.ProductInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ProductInfo{}, fmt.Errorf("external api returned status: %d", resp.StatusCode)
	}

	// Decode response
	var result struct {
		Product models.ProductInfo `json:"product"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.ProductInfo{}, err
	}

	// Pick the name
	finalName := result.Product.Name
	if finalName == "" {
		finalName = result.Product.NameIT
	}
	if finalName == "" {
		finalName = result.Product.NameEN
	}
	if finalName == "" {
		finalName = "Unkown product"
	}
	result.Product.Name = finalName

	// fallback for missing barcode
	if result.Product.Barcode == "" {
		result.Product.Barcode = barcode
	}

	return result.Product, nil
}
