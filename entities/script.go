package entities

type Script struct {
	Name           string         `json:"companyName"`
	NSECode        string         `json:"symbol"`
	StatementsList StatementsList `json:"statements"`
}
