package types

type SalesforceCodeRequest struct {
	Code        string `json:"code" validate:"numeric"`
	PhoneNumber string `json:"phoneNumber" validate:"min=1"`
	CaseNumber  string `json:"caseNumber" validate:"numeric"`
	Country     string `json:"country" validate:"min=1"`
}
