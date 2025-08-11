package utils

import (
	"fmt"
	"time"
)

// CommonRequestFormat - Updated to preserve bestMatch and queries
type CommonRequestFormat struct {
	// Preserved from UO Current Format
	Personalized          bool                   `json:"personalized"`
	ContentfulEnvironment string                 `json:"contentfulEnvironment"`
	BestMatch             map[string]interface{} `json:"bestMatch"`
	Queries               map[string]interface{} `json:"queries"`

	// Abstracted from isEvent
	User      UserContext      `json:"user"`
	Session   SessionContext   `json:"session"`
	Event     EventContext     `json:"event"`
	Page      PageContext      `json:"page"`
	Products  []ProductContext `json:"products,omitempty"`
	Device    DeviceContext    `json:"device"`
	Timestamp string           `json:"timestamp"`
}

// UOCurrentRequestFormat - Current UO format structure
type UOCurrentRequestFormat struct {
	Personalized          bool                   `json:"personalized"`
	ContentfulEnvironment string                 `json:"contentfulEnvironment"`
	BestMatch             map[string]interface{} `json:"bestMatch"`
	Queries               map[string]interface{} `json:"queries"`
	IsEvent               IsEventContext         `json:"isEvent"`
}

// Context structures (from previous definitions)
type UserContext struct {
	ID                      string                 `json:"id"`
	Email                   string                 `json:"email,omitempty"`
	Type                    string                 `json:"type,omitempty"`
	Segments                []string               `json:"segments,omitempty"`
	Attributes              map[string]interface{} `json:"attributes,omitempty"`
}

type SessionContext struct {
	ID        string `json:"id"`
	IsNew     bool   `json:"isNew,omitempty"`
	StartTime string `json:"startTime,omitempty"`
}

type EventContext struct {
	Type       string `json:"type"`
	Action     string `json:"action"`
	ItemAction string `json:"itemAction,omitempty"`
	Source     string `json:"source,omitempty"`
}

type PageContext struct {
	Type     string `json:"type"`
	URL      string `json:"url"`
	Referrer string `json:"referrer,omitempty"`
	Title    string `json:"title,omitempty"`
	Language string `json:"language,omitempty"`
}

type ProductContext struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Brand      string                 `json:"brand,omitempty"`
	Price      float64                `json:"price,omitempty"`
	Currency   string                 `json:"currency,omitempty"`
	Quantity   int                    `json:"quantity,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type DeviceContext struct {
	Type      string `json:"type,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	IP        string `json:"ip,omitempty"`
	Platform  string `json:"platform,omitempty"`
}

type IsEventContext struct {
	Source     IsEventSource  `json:"source"`
	User       IsEventUser    `json:"user"`
	Flags      IsEventFlags   `json:"flags"`
	Action     string         `json:"action"`
	ItemAction string         `json:"itemAction,omitempty"`
	Catalog    IsEventCatalog `json:"catalog,omitempty"`
	Cart       interface{}    `json:"cart,omitempty"`
	Device     *IsEventDevice `json:"device,omitempty"`
	Timestamp  string         `json:"timestamp,omitempty"`
}

type IsEventSource struct {
	Locale      string `json:"locale"`
	Application string `json:"application"`
	URL         string `json:"url"`
	Channel     string `json:"channel"`
	PageType    string `json:"pageType"`
	Referrer    string `json:"referrer,omitempty"`
}

type IsEventUser struct {
	ID         string                `json:"id"`
	Attributes IsEventUserAttributes `json:"attributes"`
}

type IsEventUserAttributes struct {
	CustomerID                     string   `json:"customerId,omitempty"`
	CustomerAuthStatus             string   `json:"customer_auth_status"`
	CustomerIsEmployee             bool     `json:"customer_is_employee"`
	CustomerDeliveryPassMbr        bool     `json:"customer_delivery_pass_mbr"`
	CustomerNonConsent             bool     `json:"customer_non_consent"`
	Locale                         string   `json:"locale"`
	URBNIsLoyalty                  bool     `json:"urbn_is_loyalty"`
	TierStatus                     string   `json:"tier_status"`
	CustomerNotificationPermission string   `json:"customer_notification_permission,omitempty"`
	URBNMbrA                       bool     `json:"urbn_mbr_a"`
	URBNMbrB                       bool     `json:"urbn_mbr_b"`
	URBNMbrMarketA                 bool     `json:"urbn_mbr_market_a"`
	URBNMbrMarketB                 bool     `json:"urbn_mbr_market_b"`
	CountryCode                    string   `json:"countryCode"`
	RegionCode                     string   `json:"regionCode,omitempty"`
	Email                          string   `json:"email,omitempty"`
	LoyaltyTier                    string   `json:"loyaltyTier,omitempty"`
	Segments                       []string `json:"segments,omitempty"`
}

type IsEventFlags map[string]interface{}

type IsEventCatalog struct {
	Product  *IsEventProduct  `json:"Product,omitempty"`
	Category *IsEventCategory `json:"Category,omitempty"`
}

type IsEventProduct struct {
	ID         string                 `json:"_id"`
	Name       string                 `json:"name,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Brand      string                 `json:"brand,omitempty"`
	Price      float64                `json:"price,omitempty"`
	Currency   string                 `json:"currency,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type IsEventCategory struct {
	ID string `json:"_id"`
}

type IsEventDevice struct {
	Type      string `json:"type,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	IP        string `json:"ip,omitempty"`
	Platform  string `json:"platform,omitempty"`
}

// UOToCommonTranslator - Translates UO Current Format to Common Request Format
type UOToCommonTranslator struct{}

func (t *UOToCommonTranslator) Translate(uoRequest *UOCurrentRequestFormat) (*CommonRequestFormat, error) {
	// Abstract user from isEvent.user
	user := t.extractUser(&uoRequest.IsEvent.User)

	// Preserve flags, source channel, itemAction, catalog, and cart in user attributes
	user.Attributes["flags"] = uoRequest.IsEvent.Flags
	user.Attributes["source_channel"] = uoRequest.IsEvent.Source.Channel
	user.Attributes["item_action"] = uoRequest.IsEvent.ItemAction
	user.Attributes["catalog"] = uoRequest.IsEvent.Catalog
	user.Attributes["cart"] = uoRequest.IsEvent.Cart

	// Generate session (UO doesn't explicitly track this)
	session := SessionContext{
		ID:    t.generateSessionID(),
		IsNew: true, // Default assumption
	}

	// Abstract event from isEvent
	event := t.extractEvent(&uoRequest.IsEvent)

	// Abstract page from isEvent.source
	page := t.extractPage(&uoRequest.IsEvent.Source)

	// Abstract products from isEvent.catalog
	products := t.extractProducts(&uoRequest.IsEvent.Catalog)

	// Abstract device from isEvent.device
	device := t.extractDevice(uoRequest.IsEvent.Device)

	// Extract timestamp
	timestamp := uoRequest.IsEvent.Timestamp
	if timestamp == "" {
		timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	// Build common format - preserving bestMatch and queries exactly
	commonRequest := &CommonRequestFormat{
		// Preserved sections
		Personalized:          uoRequest.Personalized,
		ContentfulEnvironment: uoRequest.ContentfulEnvironment,
		BestMatch:             uoRequest.BestMatch,
		Queries:               uoRequest.Queries,

		// Abstracted sections
		User:      user,
		Session:   session,
		Event:     event,
		Page:      page,
		Products:  products,
		Device:    device,
		Timestamp: timestamp,
	}

	return commonRequest, nil
}

func (t *UOToCommonTranslator) extractUser(isEventUser *IsEventUser) UserContext {
	user := UserContext{
		ID: isEventUser.ID,
		Attributes: make(map[string]interface{}),
	}

	// Extract email
	if isEventUser.Attributes.Email != "" {
		user.Email = isEventUser.Attributes.Email
	}

	// Map auth status to user type
	switch isEventUser.Attributes.CustomerAuthStatus {
	case "AUTHORIZED":
		user.Type = "member"
	case "GUEST":
		user.Type = "guest"
	default:
		user.Type = "guest"
	}

	// Extract segments
	if len(isEventUser.Attributes.Segments) > 0 {
		user.Segments = isEventUser.Attributes.Segments
	}

	// Preserve all additional user attributes
	user.Attributes["customer_auth_status"] = isEventUser.Attributes.CustomerAuthStatus
	user.Attributes["customer_delivery_pass_mbr"] = isEventUser.Attributes.CustomerDeliveryPassMbr
	user.Attributes["customer_is_employee"] = isEventUser.Attributes.CustomerIsEmployee
	user.Attributes["customer_non_consent"] = isEventUser.Attributes.CustomerNonConsent
	user.Attributes["locale"] = isEventUser.Attributes.Locale
	user.Attributes["urbn_is_loyalty"] = isEventUser.Attributes.URBNIsLoyalty
	user.Attributes["tier_status"] = isEventUser.Attributes.TierStatus
	user.Attributes["countryCode"] = isEventUser.Attributes.CountryCode
	user.Attributes["regionCode"] = isEventUser.Attributes.RegionCode

	return user
}

func (t *UOToCommonTranslator) extractEvent(isEvent *IsEventContext) EventContext {
	// Map UO actions to common event types
	actionToEventType := map[string]string{
		"Page View":      "page_view",
		"Product Detail": "product_view",
		"Add to Cart":    "add_to_cart",
		"Purchase":       "purchase",
		"CategoryView":   "category_view",
		"Cart":           "cart_view",
		"Search":         "search",
		"Login":          "login",
		"Signup":         "signup",
		"ContentView":    "page_view",
	}

	eventType, exists := actionToEventType[isEvent.Action]
	if !exists {
		eventType = "page_view"
	}

	return EventContext{
		Type:       eventType,
		Action:     isEvent.Action,
		ItemAction: isEvent.ItemAction,
		Source:     isEvent.Source.Application,
	}
}

func (t *UOToCommonTranslator) extractPage(source *IsEventSource) PageContext {
	// Map UO page types to common page types
	pageTypeMapping := map[string]string{
		"homepage": "homepage",
		"home":     "homepage", // Handle both home and homepage
		"product":  "product",
		"category": "category",
		"cart":     "cart",
		"Cart":     "cart", // Handle capital C
		"checkout": "checkout",
		"search":   "search",
		"content":  "other",
	}

	pageType, exists := pageTypeMapping[source.PageType]
	if !exists {
		pageType = "other"
	}

	return PageContext{
		Type:     pageType,
		URL:      source.URL,
		Referrer: source.Referrer,
		Language: "en", // Default
	}
}

func (t *UOToCommonTranslator) extractProducts(catalog *IsEventCatalog) []ProductContext {
	if catalog.Product == nil {
		return nil
	}

	product := ProductContext{
		ID:         catalog.Product.ID,
		Name:       catalog.Product.Name,
		Category:   catalog.Product.Category,
		Brand:      catalog.Product.Brand,
		Price:      catalog.Product.Price,
		Currency:   catalog.Product.Currency,
		Attributes: catalog.Product.Attributes,
	}

	return []ProductContext{product}
}

func (t *UOToCommonTranslator) extractDevice(isEventDevice *IsEventDevice) DeviceContext {
	device := DeviceContext{
		Platform: "web", // Default
	}

	if isEventDevice != nil {
		device.Type = isEventDevice.Type
		device.UserAgent = isEventDevice.UserAgent
		device.IP = isEventDevice.IP
		if isEventDevice.Platform != "" {
			device.Platform = isEventDevice.Platform
		}
	}

	return device
}

func (t *UOToCommonTranslator) generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

// CommonToUOTranslator - Translates Common Request Format to UO Current Format
type CommonToUOTranslator struct{}

func (t *CommonToUOTranslator) Translate(commonRequest *CommonRequestFormat) (*UOCurrentRequestFormat, error) {
	// Reconstruct isEvent from abstracted data
	isEvent := t.buildIsEvent(commonRequest)

	// Build UO format - preserving bestMatch and queries exactly
	uoRequest := &UOCurrentRequestFormat{
		// Preserved sections
		Personalized:          commonRequest.Personalized,
		ContentfulEnvironment: commonRequest.ContentfulEnvironment,
		BestMatch:             commonRequest.BestMatch,
		Queries:               commonRequest.Queries,

		// Reconstructed section
		IsEvent: isEvent,
	}

	return uoRequest, nil
}

func (t *CommonToUOTranslator) buildIsEvent(commonRequest *CommonRequestFormat) IsEventContext {
	// Build isEvent.source from page
	source := IsEventSource{
		Locale:      "en_US", // Default
		Application: commonRequest.Event.Source,
		URL:         commonRequest.Page.URL,
		Channel:     "Server", // UO default
		PageType:    t.mapPageType(commonRequest.Page.Type),
		Referrer:    commonRequest.Page.Referrer,
	}
	
	// Restore channel from user attributes if available
	if channelData, exists := commonRequest.User.Attributes["source_channel"]; exists {
		if channel, ok := channelData.(string); ok {
			source.Channel = channel
		}
	}

	// Build isEvent.user from user
	user := IsEventUser{
		ID:         commonRequest.User.ID,
		Attributes: t.buildUserAttributes(commonRequest.User),
	}

	// Build isEvent.flags from preserved data
	flags := IsEventFlags{
		"pageView":    true, // Default
		"noCampaigns": false, // Default
	}
	
	// Restore flags from user attributes if available
	if flagsData, exists := commonRequest.User.Attributes["flags"]; exists {
		if flagsMap, ok := flagsData.(map[string]interface{}); ok {
			flags = flagsMap
		}
	}

	// Build isEvent.catalog from products
	catalog := t.buildCatalog(commonRequest.Products)
	
	// Restore catalog from user attributes if available
	if catalogData, exists := commonRequest.User.Attributes["catalog"]; exists {
		if catalogMap, ok := catalogData.(IsEventCatalog); ok {
			catalog = catalogMap
		} else if catalogInterface, ok := catalogData.(map[string]interface{}); ok {
			// Handle case where catalog comes as generic map
			if categoryData, exists := catalogInterface["Category"]; exists {
				if categoryMap, ok := categoryData.(map[string]interface{}); ok {
					if id, exists := categoryMap["_id"]; exists {
						if idStr, ok := id.(string); ok {
							catalog.Category = &IsEventCategory{ID: idStr}
						}
					}
				}
			}
		}
	}
	
	// Restore cart from user attributes if available
	var cart interface{}
	if cartData, exists := commonRequest.User.Attributes["cart"]; exists {
		cart = cartData
	}

	// Build isEvent.device from device
	var device *IsEventDevice
	if commonRequest.Device.Type != "" || commonRequest.Device.UserAgent != "" || commonRequest.Device.IP != "" {
		device = &IsEventDevice{
			Type:      commonRequest.Device.Type,
			UserAgent: commonRequest.Device.UserAgent,
			IP:        commonRequest.Device.IP,
			Platform:  commonRequest.Device.Platform,
		}
	}

	// Use the preserved action directly if available, otherwise map from event type
	action := commonRequest.Event.Action
	if action == "" {
		action = t.mapEventTypeToAction(commonRequest.Event.Type)
	}
	
	itemAction := commonRequest.Event.ItemAction
	if itemAction == "" {
		itemAction = t.mapEventTypeToItemAction(commonRequest.Event.Type)
		// Restore itemAction from user attributes if available
		if itemActionData, exists := commonRequest.User.Attributes["item_action"]; exists {
			if ia, ok := itemActionData.(string); ok {
				itemAction = ia
			}
		}
	}

	return IsEventContext{
		Source:     source,
		User:       user,
		Flags:      flags,
		Action:     action,
		ItemAction: itemAction,
		Catalog:    catalog,
		Cart:       cart,
		Device:     device,
		Timestamp:  commonRequest.Timestamp,
	}
}

func (t *CommonToUOTranslator) buildUserAttributes(user UserContext) IsEventUserAttributes {
	// Map user type to auth status
	authStatus := "GUEST"
	if user.Type == "member" {
		authStatus = "AUTHORIZED"
	}

	attributes := IsEventUserAttributes{
		CustomerAuthStatus:             authStatus,
		CustomerIsEmployee:             false,   // Default
		CustomerDeliveryPassMbr:        false,   // Default
		CustomerNonConsent:             false,   // Default
		Locale:                         "en_US", // Default
		URBNIsLoyalty:                  user.Type == "member",
		TierStatus:                     "",      // Default
		CustomerNotificationPermission: "default",
		URBNMbrA:                       false, // Default
		URBNMbrB:                       false, // Default
		URBNMbrMarketA:                 false, // Default
		URBNMbrMarketB:                 false, // Default
		CountryCode:                    "US",  // Default
		Email:                          user.Email,
		Segments:                       user.Segments,
	}

	// Restore preserved attributes from user.Attributes
	if user.Attributes != nil {
		if val, exists := user.Attributes["customer_auth_status"]; exists {
			if str, ok := val.(string); ok {
				attributes.CustomerAuthStatus = str
			}
		}
		if val, exists := user.Attributes["customer_delivery_pass_mbr"]; exists {
			if b, ok := val.(bool); ok {
				attributes.CustomerDeliveryPassMbr = b
			}
		}
		if val, exists := user.Attributes["customer_is_employee"]; exists {
			if b, ok := val.(bool); ok {
				attributes.CustomerIsEmployee = b
			}
		}
		if val, exists := user.Attributes["customer_non_consent"]; exists {
			if b, ok := val.(bool); ok {
				attributes.CustomerNonConsent = b
			}
		}
		if val, exists := user.Attributes["locale"]; exists {
			if str, ok := val.(string); ok {
				attributes.Locale = str
			}
		}
		if val, exists := user.Attributes["urbn_is_loyalty"]; exists {
			if b, ok := val.(bool); ok {
				attributes.URBNIsLoyalty = b
			}
		}
		if val, exists := user.Attributes["tier_status"]; exists {
			if str, ok := val.(string); ok {
				attributes.TierStatus = str
			}
		}
		if val, exists := user.Attributes["countryCode"]; exists {
			if str, ok := val.(string); ok {
				attributes.CountryCode = str
			}
		}
		if val, exists := user.Attributes["regionCode"]; exists {
			if str, ok := val.(string); ok {
				attributes.RegionCode = str
			}
		}
	}

	return attributes
}

func (t *CommonToUOTranslator) buildCatalog(products []ProductContext) IsEventCatalog {
	catalog := IsEventCatalog{}

	if len(products) > 0 {
		product := products[0] // UO typically handles single product
		catalog.Product = &IsEventProduct{
			ID:         product.ID,
			Name:       product.Name,
			Category:   product.Category,
			Brand:      product.Brand,
			Price:      product.Price,
			Currency:   product.Currency,
			Attributes: product.Attributes,
		}
	}

	return catalog
}

func (t *CommonToUOTranslator) mapPageType(commonPageType string) string {
	// Map common page types to UO page types
	pageTypeMapping := map[string]string{
		"homepage": "home", // UO uses "home" for homepage
		"product":  "product",
		"category": "category",
		"cart":     "Cart", // UO uses capital C
		"checkout": "checkout",
		"search":   "search",
		"other":    "content",
	}

	if uoPageType, exists := pageTypeMapping[commonPageType]; exists {
		return uoPageType
	}
	return "content"
}

func (t *CommonToUOTranslator) mapEventTypeToAction(eventType string) string {
	// Map common event types to UO actions
	eventTypeToAction := map[string]string{
		"page_view":     "Page View",
		"product_view":  "Product Detail",
		"add_to_cart":   "Add to Cart",
		"purchase":      "Purchase",
		"category_view": "CategoryView",
		"cart_view":     "Cart",
		"search":        "Search",
		"login":         "Login",
		"signup":        "Signup",
	}

	if action, exists := eventTypeToAction[eventType]; exists {
		return action
	}
	return "Page View"
}

func (t *CommonToUOTranslator) mapEventTypeToItemAction(eventType string) string {
	// Map common event types to UO item actions
	eventTypeToItemAction := map[string]string{
		"page_view":     "View Category",
		"product_view":  "View Product",
		"add_to_cart":   "Add to Cart",
		"purchase":      "Purchase",
		"category_view": "View Category",
		"cart_view":     "View Cart",
		"search":        "Search",
		"login":         "Login",
		"signup":        "Signup",
	}

	if itemAction, exists := eventTypeToItemAction[eventType]; exists {
		return itemAction
	}
	return "View Category"
}

// Helper function to compare maps (simplified)
func (t *CommonToUOTranslator) CompareMaps(map1, map2 map[string]interface{}) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value1 := range map1 {
		if value2, exists := map2[key]; !exists || fmt.Sprintf("%v", value1) != fmt.Sprintf("%v", value2) {
			return false
		}
	}

	return true
}
