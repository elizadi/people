package types

const NameParam = "/?name="

type AgeData struct {
	Count uint64
	Name  string
	Age   uint8
}

type GenderData struct {
	Count       uint64  `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

type NationalityData struct {
	Count   uint64        `json:"count"`
	Name    string        `json:"name"`
	Country []CountryData `json:"country"`
}

type CountryData struct {
	ID          string  `json:"country_id"`
	Probability float32 `json:"probability"`
}
