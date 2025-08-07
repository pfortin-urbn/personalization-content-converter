package main

import (
	"encoding/json"
	"fmt"
	"personalization-content-converter/utils"
)

// Example usage demonstrating bidirectional translation
func main() {
	// Create original UO format request
	originalUO := &utils.UOCurrentRequestFormat{
		Personalized:          true,
		ContentfulEnvironment: "master",
		BestMatch: map[string]interface{}{
			"tokenScope":      "AUTHORIZED",
			"country":         "US",
			"region":          "PA",
			"city":            "PHILADELPHIA",
			"url":             "/products/uo-cropped-cardigan-womens",
			"homepage":        false,
			"sort":            false,
			"userType":        "member",
			"loyaltyTier":     "Gold",
			"promoCode":       "WINTER25",
			"experimentGroup": "test-group-a",
		},
		Queries: map[string]interface{}{
			"productRecommendations": map[string]interface{}{
				"include":         5,
				"content_type":    "productRecommendation",
				"fields.category": "womens-sweaters",
				"fields.active":   true,
				"select":          "fields.title,fields.products,fields.dyIdentifier,sys",
			},
			"globalPromo": map[string]interface{}{
				"include":      3,
				"content_type": "globalPromoContainer",
			},
		},
		IsEvent: utils.IsEventContext{
			Source: utils.IsEventSource{
				Locale:      "en_US",
				Application: "web",
				URL:         "https://urbanoutfitters.com/products/uo-cropped-cardigan-womens",
				Channel:     "Server",
				PageType:    "product",
				Referrer:    "https://urbanoutfitters.com/womens/sweaters",
			},
			User: utils.IsEventUser{
				ID: "4a095131-2682-4503-b1bf-d3bd8bfec1f7",
				Attributes: utils.IsEventUserAttributes{
					CustomerAuthStatus: "AUTHORIZED",
					CustomerIsEmployee: false,
					Locale:             "en_US",
					URBNIsLoyalty:      true,
					CountryCode:        "US",
					RegionCode:         "PA",
					Email:              "sarah.chen@gmail.com",
					Segments:           []string{"returning-customers", "womens-clothing-buyers"},
				},
			},
			Flags: utils.IsEventFlags{
				PageView: true,
			},
			Action:     "Product Detail",
			ItemAction: "View Product",
			Catalog: utils.IsEventCatalog{
				Product: &utils.IsEventProduct{
					ID:       "uo-cardigan-w-2025-001",
					Name:     "UO Cropped Cardigan",
					Category: "Womens-Sweaters",
					Brand:    "Urban Outfitters",
					Price:    69.99,
					Currency: "USD",
					Attributes: map[string]interface{}{
						"color":    "Sage Green",
						"size":     "Medium",
						"material": "Cotton Blend",
					},
				},
			},
			Timestamp: "2025-01-15T10:30:45Z",
		},
	}

	// 1. UO → Common
	uoToCommon := &utils.UOToCommonTranslator{}
	commonRequest, err := uoToCommon.Translate(originalUO)
	if err != nil {
		fmt.Printf("Error translating UO to Common: %v\n", err)
		return
	}

	fmt.Println("=== UO → Common Translation ===")
	commonJSON, _ := json.MarshalIndent(commonRequest, "", "  ")
	fmt.Printf("Common Format:\n%s\n\n", string(commonJSON))

	// 2. Common → UO (round-trip test)
	commonToUO := &utils.CommonToUOTranslator{}
	reconstructedUO, err := commonToUO.Translate(commonRequest)
	if err != nil {
		fmt.Printf("Error translating Common to UO: %v\n", err)
		return
	}

	fmt.Println("=== Common → UO Translation ===")
	reconstructedJSON, _ := json.MarshalIndent(reconstructedUO, "", "  ")
	fmt.Printf("Reconstructed UO Format:\n%s\n\n", string(reconstructedJSON))

	// 3. Validate round-trip accuracy
	fmt.Println("=== Round-trip Validation ===")
	fmt.Printf("bestMatch preserved: %v\n", commonToUO.CompareMaps(originalUO.BestMatch, reconstructedUO.BestMatch))
	fmt.Printf("queries preserved: %v\n", commonToUO.CompareMaps(originalUO.Queries, reconstructedUO.Queries))
	fmt.Printf("user ID preserved: %v\n", originalUO.IsEvent.User.ID == reconstructedUO.IsEvent.User.ID)
	fmt.Printf("action preserved: %v\n", originalUO.IsEvent.Action == reconstructedUO.IsEvent.Action)
	fmt.Printf("product ID preserved: %v\n", originalUO.IsEvent.Catalog.Product.ID == reconstructedUO.IsEvent.Catalog.Product.ID)
}
