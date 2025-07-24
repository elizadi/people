package enrichment

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"people/internal/types"
)

type Enrichment struct {
	AgeUrl         string
	GenderUrl      string
	NationalityUrl string
	Logger         *logrus.Logger
}

func New(ageUrl, genderUrl, nationalityUrl string, logger *logrus.Logger) (*Enrichment, error) {
	if ageUrl == "" || genderUrl == "" || nationalityUrl == "" {
		logrus.Errorf("Age URL or Gender URL or National URL are required")
		return &Enrichment{}, errors.New("Empty parameter")
	}
	return &Enrichment{
		AgeUrl:         ageUrl,
		GenderUrl:      genderUrl,
		NationalityUrl: nationalityUrl,
		Logger:         logger,
	}, nil
}

func (e *Enrichment) Age(name string) (uint8, error) {
	var ageData types.AgeData

	url := e.AgeUrl + types.NameParam + name
	body, err := e.httpGet(url)
	if err != nil {
		e.Logger.WithError(err).Errorln("Error getting age by request")
		return 0, err
	}

	err = json.Unmarshal(body, &ageData)
	if err != nil {
		e.Logger.WithError(err).Errorln("Error unmarshal Age response body")
		return 0, err
	}

	return ageData.Age, nil
}

func (e *Enrichment) Gender(name string) (string, error) {
	var genderData types.GenderData

	url := e.GenderUrl + types.NameParam + name
	body, err := e.httpGet(url)
	if err != nil {
		e.Logger.WithError(err).Errorln("Error getting gender by request")
		return "", err
	}

	err = json.Unmarshal(body, &genderData)
	if err != nil {
		e.Logger.WithError(err).Errorln("Error unmarshal gender response body")
		return "", err
	}

	return genderData.Gender, nil
}

func (e *Enrichment) Nationality(name string) (string, error) {
	var nationality types.NationalityData

	url := e.NationalityUrl + types.NameParam + name
	body, err := e.httpGet(url)
	if err != nil {
		e.Logger.WithError(err).Errorln("Error getting gender by request")
		return "", err
	}

	err = json.Unmarshal(body, &nationality)
	if err != nil {
		e.Logger.WithError(err).Errorln("Error unmarshal gender response body")
		return "", err
	}

	return nationality.Country[0].ID, nil
}

func (e *Enrichment) httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		e.Logger.WithError(err).Errorln("Request failed")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e.Logger.WithError(err).Errorln("Error get response body")
		return nil, err
	}

	return body, nil
}
