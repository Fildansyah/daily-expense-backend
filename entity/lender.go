package entity

type LenderDTO struct {
	CommonDTOFields

	UserID string `json:"-"`
	Name   string `json:"name"`
}

type LenderPayload struct {
	Name string `json:"name"`
}
