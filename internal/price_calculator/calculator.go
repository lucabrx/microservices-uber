package pricecalculator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func CalculatePrice(startLat, startLon, endLat, endLon float64) (float64, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/car/%.6f,%.6f;%.6f,%.6f?overview=false",
		startLon, startLat, endLon, endLat,
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("OSRM request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	type osrmRoute struct {
		Distance float64 `json:"distance"`
	}
	type osrmResponse struct {
		Routes []osrmRoute `json:"routes"`
	}

	var osrmResp osrmResponse
	if err := json.Unmarshal(body, &osrmResp); err != nil {
		return 0, err
	}
	if len(osrmResp.Routes) == 0 {
		return 0, errors.New("no routes found")
	}

	distanceInKm := osrmResp.Routes[0].Distance / 1000
	baseFare := 2.50
	perKmRate := 1.50
	price := baseFare + (distanceInKm * perKmRate)

	return float64(int(price*100)) / 100, nil
}
