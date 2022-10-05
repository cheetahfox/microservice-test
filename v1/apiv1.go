package v1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/cheetahfox/microservice-test/health"
	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {

	raw, err := getRawData()
	if err != nil {
		return err
	}

	s := string(raw)

	return c.SendString(s)
	//return c.SendString("v1")
}

type alphaVantageData struct {
	info     string            `json:"1. Information"`
	symbol   string            `json:"2. Symbol"`
	last     string            `json:"3. Last Refreshed"`
	timezone string            `json:"5. Time Zone"`
	days     []timeSeriesDaily `json:"Time Series (Daily)"`
}

type stockPrices struct {
	open   string `json:"1. open"`
	high   string `json:"2. high"`
	low    string `json:"3. low"`
	close  string `json:"4. close"`
	volume string `json:"5. volume"`
}

type timeSeriesDaily struct {
	day    string
	prices stockPrices
}

type runningValues struct {
	apiKey string
	ndays  int
	symbol string
}

// Verify and Fetch the required Env values
func fetchValues() (runningValues, error) {
	var values runningValues

	/*
		Do some basic checks on the env values that are present
		And if they aren't good; error and set not ready in the health check.
		Since I am double checking the apikey and symbol already I am going to to check ndays
		here as well.
	*/
	apiKey, present := os.LookupEnv("API_KEY")
	if !present {
		fmt.Println("The API key is missing somehow; it was present at startup...")
		health.ApiReady = false
		return values, errors.New("missing API Key")
	}
	values.apiKey = apiKey

	// Check if ndays is an int and if it's in the right range
	ndayString, present := os.LookupEnv("NDAYS")
	if !present {
		fmt.Println("Ndays Env is missing somehow; it was present at startup...")
		health.ApiReady = false
		return values, errors.New("missing ndays Env")
	}
	ndays, err := strconv.Atoi(ndayString)
	if err != nil {
		fmt.Println("Ndays env value isn't a number")
		health.ApiReady = false
		return values, errors.New("Ndays env not an int")
	}
	// we default to format compact in the api so we need to have NDAYS set to less than 100 and more than 0
	if ndays > 100 || ndays <= 0 {
		fmt.Println("Ndays outside the require range.")
		health.ApiReady = false
		return values, errors.New("Ndays outside of range")
	}
	values.ndays = ndays

	symbol, present := os.LookupEnv("SYMBOL")
	if !present {
		fmt.Println("Symbol Env is missing somehow; it was present at startup...")
		health.ApiReady = false
		return values, errors.New("missing SYMBOL Env")
	}
	values.symbol = symbol

	return values, nil
}

// Fetch the raw data from the Stock Ticker web service
func getRawData() ([]byte, error) {
	values, err := fetchValues()
	if err != nil {
		return nil, errors.New("Error while validating env values")
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", values.apiKey, values.symbol)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	// check the http status code
	if response.StatusCode != 200 {
		fmt.Printf("WARNING: Request Status Code %d %s\n", response.StatusCode, response.Status)
		return nil, errors.New("Not 200 return")
	}

	// Read the body response
	bodybytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bodybytes, nil
}

func parseData([]byte) (alphaVantageData, error) {
	var aData alphaVantageData

	return aData, nil

}
