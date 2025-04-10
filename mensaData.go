package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Mensa_API_Response struct {
	Content []struct {
		SpeiseplanAdvanced struct {
			Titel                  string `json:"titel"`
			Anzeigename            string `json:"anzeigename"`
			SpeiseplanLayoutTypeID any    `json:"speiseplanLayoutTypeID"`
		} `json:"speiseplanAdvanced"`

		SpeiseplanGerichtData []struct {
			SpeiseplanAdvancedGericht struct {
				Datum              time.Time `json:"datum"`
				GerichtkategorieID int       `json:"gerichtkategorieID"`
				Gerichtname        string    `json:"gerichtname"`
			} `json:"speiseplanAdvancedGericht"`

			Zusatzinformationen struct {
				GerichtnameAlternative   string  `json:"gerichtnameAlternative"`
				MitarbeiterpreisDecimal2 float32 `json:"mitarbeiterpreisDecimal2"`
				GaestepreisDecimal2      float32 `json:"gaestepreisDecimal2"`
			} `json:"zusatzinformationen"`
		} `json:"speiseplanGerichtData"`
	} `json:"content"`
}

func getMeals(t time.Time) ([10]string, [10]float32) {
	//creates the mensa-api URL and makes a http request to it
	//var t time.Time = time.Now()
	var mensa_api_url string = fmt.Sprintf("https://stwwb.webspeiseplan.de/index.php?token=55ed21609e26bbf68ba2b19390bf7961&model=menu&location=9601&languagetype=1&_=%s", strconv.FormatInt(t.Unix(), 10))

	req, err := http.NewRequest("GET", mensa_api_url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "de-DE,de;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://stwwb.webspeiseplan.de/Menu")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//reads out the response body into an array
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//unmarshalls the array into readable json
	var response Mensa_API_Response
	json.Unmarshal(body, &response)

	var meals [10]string
	var prices [10]float32
	{
		var i int = 0
		for _, meal := range response.Content[0].SpeiseplanGerichtData {
			if meal.SpeiseplanAdvancedGericht.Datum.Day() != t.Day() {
				continue

			} else if meal.SpeiseplanAdvancedGericht.GerichtkategorieID == 112 {
				continue

			} else if meal.SpeiseplanAdvancedGericht.GerichtkategorieID == 294 {
				break

			} else {
				meals[i] = meal.SpeiseplanAdvancedGericht.Gerichtname
				prices[i] = meal.Zusatzinformationen.MitarbeiterpreisDecimal2
				i++
			}

		}
	}
	return meals, prices
}
