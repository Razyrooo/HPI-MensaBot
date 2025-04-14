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

// find the correct content element containing the meals of the current day
func findContent(t time.Time, response Mensa_API_Response) int {
	for i, content := range response.Content {
		mealCounter := 0
		for _, meal := range content.SpeiseplanGerichtData {
			if meal.SpeiseplanAdvancedGericht.Datum.Day() == t.Day() {
				mealCounter++
			}
		}
		if mealCounter >= 5 {
			return i
		}
	}
	log.Panic("No content found.")
	return -1
}

// get mealdata from the mensa-website-api und return sorted meals and prices
func getMeals(t time.Time) ([10]string, [10]float32) {
	//creates the mensa-api URL and makes a http request to it
	//var t time.Time = time.Now()
	mensa_api_url := fmt.Sprintf("https://stwwb.webspeiseplan.de/index.php?token=55ed21609e26bbf68ba2b19390bf7961&model=menu&location=9601&languagetype=1&_=%s", strconv.FormatInt(t.Unix(), 10))

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

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("failed to clode response body: %w", err)
		}
	}()

	//reads out the response body into an array
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("failed to clode response body: %w", err)
		}
	}()

	//unmarshalls the array into readable json
	var response Mensa_API_Response
	err_j := json.Unmarshal(body, &response)
	if err_j != nil {
		log.Fatal(err)
	}

	// searches response data for the meals and prices of the day
	/*
		the content from the json response always begins with the relevo-schale(id 112) and ends with either heiße Theke (id 294) or nachtisch (119/120)
	*/
	var meals [10]string
	var prices [10]float32
	{
		i := 0
		for _, meal := range response.Content[findContent(t, response)].SpeiseplanGerichtData {

			// meal filter
			if meal.SpeiseplanAdvancedGericht.Datum.Day() != t.Day() {
				continue

			} else if meal.SpeiseplanAdvancedGericht.GerichtkategorieID == 112 { // Relevo-Schale
				continue

			} else if meal.SpeiseplanAdvancedGericht.GerichtkategorieID == 294 { // Heiße Theke
				break

			} else if meal.SpeiseplanAdvancedGericht.GerichtkategorieID == 119 { // Nachtisch
				break

			} else if meal.SpeiseplanAdvancedGericht.GerichtkategorieID == 120 { //Nachtisch
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
