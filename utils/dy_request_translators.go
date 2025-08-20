package utils

import "time"

// DYChooseRequest represents the request payload for the Dynamic Yield choose API
type DYChooseRequest struct {
	User    DYUser    `json:"user"`
	Session DYSession `json:"session"`
	Context DYContext `json:"context"`
	Selector DYSelector `json:"selector"`
	Options  DYOptions  `json:"options"`
}

// DYUser represents the user object in the DY request
type DYUser struct {
	ActiveConsentAccepted bool   `json:"active_consent_accepted"`
	DyidServer            string `json:"dyid_server"`
	Dyid                  string `json:"dyid"`
}

// DYSession represents the session object in the DY request
type DYSession struct {
	Dy string `json:"dy"`
}

// DYContext represents the context object in the DY request
type DYContext struct {
	Page   DYPage   `json:"page"`
	Device DYDevice `json:"device"`
}

// DYPage represents the page object in the DY request
type DYPage struct {
	Type     string   `json:"type"`
	Data     []string `json:"data"`
	Location string   `json:"location"`
}

// DYDevice represents the device object in the DY request
type DYDevice struct {
	UserAgent string `json:"userAgent"`
	Type      string `json:"type"`
	Browser   string `json:"browser"`
	Ip        string `json:"ip"`
}

// DYSelector represents the selector object in the DY request
type DYSelector struct {
	Names []string `json:"names"`
}

// DYOptions represents the options object in the DY request
type DYOptions struct {
	IsImplicitPageview       bool              `json:"isImplicitPageview"`
	ReturnAnalyticsMetadata  bool              `json:"returnAnalyticsMetadata"`
	IsImplicitImpressionMode bool              `json:"isImplicitImpressionMode"`
	IsImplicitClientData     bool              `json:"isImplicitClientData"`
	RecsProductData          *DYRecsProductData `json:"recsProductData,omitempty"`
}

// DYRecsProductData represents the recsProductData object in the DY request
type DYRecsProductData struct {
	FieldFilter []string `json:"fieldFilter"`
}

// CommonToDYRequestTranslator translates from the common format to the DY format
type CommonToDYRequestTranslator struct{}

// Translate performs the translation
func (t *CommonToDYRequestTranslator) Translate(commonRequest *CommonRequestFormat) (*DYChooseRequest, error) {
	user := DYUser{
		Dyid: commonRequest.User.ID,
	}
	if val, ok := commonRequest.User.Attributes["active_consent_accepted"].(bool); ok {
		user.ActiveConsentAccepted = val
	}
	if val, ok := commonRequest.User.Attributes["dyid_server"].(string); ok {
		user.DyidServer = val
	}

	session := DYSession{
		Dy: commonRequest.Session.ID,
	}

	pageType := "OTHER"
	switch commonRequest.Page.Type {
	case "homepage":
		pageType = "HOMEPAGE"
	case "product":
		pageType = "PRODUCT"
	}

	var productData []string
	for _, product := range commonRequest.Products {
		productData = append(productData, product.ID)
	}

	page := DYPage{
		Type:     pageType,
		Location: commonRequest.Page.URL,
		Data:     productData,
	}

	device := DYDevice{
		UserAgent: commonRequest.Device.UserAgent,
		Type:      commonRequest.Device.Type,
		Browser:   commonRequest.Device.Platform, // Assuming platform is the browser
		Ip:        commonRequest.Device.IP,
	}

	context := DYContext{
		Page:   page,
		Device: device,
	}

	selector := DYSelector{}
	if val, ok := commonRequest.Queries["selector"].(map[string]interface{}); ok {
		if names, ok := val["names"].([]interface{}); ok {
			for _, name := range names {
				if strName, ok := name.(string); ok {
					selector.Names = append(selector.Names, strName)
				}
			}
		}
	}

	options := DYOptions{}
	if val, ok := commonRequest.Queries["options"].(map[string]interface{}); ok {
		if v, ok := val["isImplicitPageview"].(bool); ok {
			options.IsImplicitPageview = v
		}
		if v, ok := val["returnAnalyticsMetadata"].(bool); ok {
			options.ReturnAnalyticsMetadata = v
		}
		if v, ok := val["isImplicitImpressionMode"].(bool); ok {
			options.IsImplicitImpressionMode = v
		}
		if v, ok := val["isImplicitClientData"].(bool); ok {
			options.IsImplicitClientData = v
		}
		if recsProductData, ok := val["recsProductData"].(map[string]interface{}); ok {
			if fieldFilter, ok := recsProductData["fieldFilter"].([]interface{}); ok {
				var filters []string
				for _, filter := range fieldFilter {
					if strFilter, ok := filter.(string); ok {
						filters = append(filters, strFilter)
					}
				}
				options.RecsProductData = &DYRecsProductData{FieldFilter: filters}
			}
		}
	}

	dyRequest := &DYChooseRequest{
		User:     user,
		Session:  session,
		Context:  context,
		Selector: selector,
		Options:  options,
	}

	return dyRequest, nil
}

// DYToCommonRequestTranslator translates from the DY format to the common format
type DYToCommonRequestTranslator struct{}

// Translate performs the translation
func (t *DYToCommonRequestTranslator) Translate(dyRequest *DYChooseRequest) (*CommonRequestFormat, error) {
	user := UserContext{
		ID: dyRequest.User.Dyid,
		Attributes: map[string]interface{}{
			"active_consent_accepted": dyRequest.User.ActiveConsentAccepted,
			"dyid_server":             dyRequest.User.DyidServer,
		},
	}

	session := SessionContext{
		ID: dyRequest.Session.Dy,
	}

	eventType := "page_view"
	if dyRequest.Context.Page.Type == "PRODUCT" {
		eventType = "product_view"
	}

	event := EventContext{
		Type:   eventType,
		Action: dyRequest.Context.Page.Type,
		Source: "Dynamic Yield",
	}

	pageType := "other"
	switch dyRequest.Context.Page.Type {
	case "HOMEPAGE":
		pageType = "homepage"
	case "PRODUCT":
		pageType = "product"
	}

	page := PageContext{
		Type: pageType,
		URL:  dyRequest.Context.Page.Location,
	}

	var products []ProductContext
	for _, productID := range dyRequest.Context.Page.Data {
		products = append(products, ProductContext{ID: productID})
	}

	device := DeviceContext{
		UserAgent: dyRequest.Context.Device.UserAgent,
		Type:      dyRequest.Context.Device.Type,
		Platform:  dyRequest.Context.Device.Browser, // Assuming browser is the platform
		IP:        dyRequest.Context.Device.Ip,
	}

	commonRequest := &CommonRequestFormat{
		Personalized: true,
		User:         user,
		Session:      session,
		Event:        event,
		Page:         page,
		Products:     products,
		Device:       device,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Queries: map[string]interface{}{
			"selector": dyRequest.Selector,
			"options":  dyRequest.Options,
		},
	}

	return commonRequest, nil
}
