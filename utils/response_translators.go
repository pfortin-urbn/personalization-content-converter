package utils

// CommonResponseFormat - Common response format
type CommonResponseFormat struct {
	RequestID string         `json:"requestId"`
	UserID    string         `json:"userId"`
	AccountID string         `json:"accountId"`
	EntityID  string         `json:"entityId"`
	ErrorCode int            `json:"errorCode"`
	Campaigns []CommonCampaign `json:"campaigns"`
}

// ISResponseFormat - IS response format
type ISResponseFormat struct {
	ID               string               `json:"id"`
	ResolvedUserID   string               `json:"resolvedUserId"`
	PersistedUserID  ISPersistedUserID    `json:"persistedUserId"`
	ErrorCode        int                  `json:"errorCode"`
	CampaignResponses []ISCampaignResponse `json:"campaignResponses"`
}

// CommonCampaign - Campaign structure in Common format
type CommonCampaign struct {
	CampaignID                  string      `json:"campaignId"`
	CampaignName                string      `json:"campaignName"`
	CampaignType                string      `json:"campaignType"`
	CampaignJavascriptContent   interface{} `json:"campaignJavascriptContent"`
	ExperienceID                string      `json:"experienceId"`
	ExperienceName              string      `json:"experienceName"`
	ExperienceSourceCode        string      `json:"experienceSourceCode"`
	State                       string      `json:"state"`
	Type                        string      `json:"type"`
	UserGroup                   string      `json:"userGroup"`
	TemplateNames               []string    `json:"templateNames"`
	Payload                     interface{} `json:"payload"`
}

// ISPersistedUserID - Persisted user ID structure in IS format
type ISPersistedUserID struct {
	AccountID string `json:"accountId"`
	EntityID  string `json:"entityId"`
}

// ISCampaignResponse - Campaign response structure in IS format
type ISCampaignResponse struct {
	CampaignID                  string      `json:"campaignId"`
	CampaignName                string      `json:"campaignName"`
	CampaignType                string      `json:"campaignType"`
	CampaignJavascriptContent   interface{} `json:"campaignJavascriptContent"`
	ExperienceID                string      `json:"experienceId"`
	ExperienceName              string      `json:"experienceName"`
	ExperienceSourceCode        string      `json:"experienceSourceCode"`
	State                       string      `json:"state"`
	Type                        string      `json:"type"`
	UserGroup                   string      `json:"userGroup"`
	TemplateNames               []string    `json:"templateNames"`
	Payload                     interface{} `json:"payload"`
}

// CommonToISResponseTranslator - Translates Common Response Format to IS Response Format
type CommonToISResponseTranslator struct{}

func (t *CommonToISResponseTranslator) Translate(commonResponse *CommonResponseFormat) (*ISResponseFormat, error) {
	// Convert campaigns
	campaignResponses := make([]ISCampaignResponse, len(commonResponse.Campaigns))
	for i, campaign := range commonResponse.Campaigns {
		campaignResponses[i] = ISCampaignResponse{
			CampaignID:                campaign.CampaignID,
			CampaignName:              campaign.CampaignName,
			CampaignType:              campaign.CampaignType,
			CampaignJavascriptContent: campaign.CampaignJavascriptContent,
			ExperienceID:              campaign.ExperienceID,
			ExperienceName:            campaign.ExperienceName,
			ExperienceSourceCode:      campaign.ExperienceSourceCode,
			State:                     campaign.State,
			Type:                      campaign.Type,
			UserGroup:                 campaign.UserGroup,
			TemplateNames:             campaign.TemplateNames,
			Payload:                   campaign.Payload,
		}
	}

	// Build IS response
	isResponse := &ISResponseFormat{
		ID:             commonResponse.RequestID,
		ResolvedUserID: commonResponse.UserID,
		PersistedUserID: ISPersistedUserID{
			AccountID: commonResponse.AccountID,
			EntityID:  commonResponse.EntityID,
		},
		ErrorCode:         commonResponse.ErrorCode,
		CampaignResponses: campaignResponses,
	}

	return isResponse, nil
}

// ISToCommonResponseTranslator - Translates IS Response Format to Common Response Format
type ISToCommonResponseTranslator struct{}

func (t *ISToCommonResponseTranslator) Translate(isResponse *ISResponseFormat) (*CommonResponseFormat, error) {
	// Convert campaign responses
	campaigns := make([]CommonCampaign, len(isResponse.CampaignResponses))
	for i, campaignResponse := range isResponse.CampaignResponses {
		campaigns[i] = CommonCampaign{
			CampaignID:                campaignResponse.CampaignID,
			CampaignName:              campaignResponse.CampaignName,
			CampaignType:              campaignResponse.CampaignType,
			CampaignJavascriptContent: campaignResponse.CampaignJavascriptContent,
			ExperienceID:              campaignResponse.ExperienceID,
			ExperienceName:            campaignResponse.ExperienceName,
			ExperienceSourceCode:      campaignResponse.ExperienceSourceCode,
			State:                     campaignResponse.State,
			Type:                      campaignResponse.Type,
			UserGroup:                 campaignResponse.UserGroup,
			TemplateNames:             campaignResponse.TemplateNames,
			Payload:                   campaignResponse.Payload,
		}
	}

	// Build Common response
	commonResponse := &CommonResponseFormat{
		RequestID: isResponse.ID,
		UserID:    isResponse.ResolvedUserID,
		AccountID: isResponse.PersistedUserID.AccountID,
		EntityID:  isResponse.PersistedUserID.EntityID,
		ErrorCode: isResponse.ErrorCode,
		Campaigns: campaigns,
	}

	return commonResponse, nil
}