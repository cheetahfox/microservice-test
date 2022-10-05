package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cheetahfox/microservice-test/health"
	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	s, err := v1Get()
	if err != nil {
		return err
	}
	c.Set("Content-Type", "application/json")

	return c.SendString(string(s))
}

type AlphaVantageData struct {
	Meta MetaData                   `json:"Meta Data"`
	Tds  map[string]TimeSeriesDaily `json:"Time Series (Daily)"`
}

type MetaData struct {
	Info       string `json:"1. Information"`
	Symbol     string `json:"2. Symbol"`
	Last       string `json:"3. Last Refreshed"`
	Outputsize string `json:"4. Output Size"`
	Timezone   string `json:"5. Time Zone"`
}

type TimeSeriesDaily struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type runningValues struct {
	apiKey string
	ndays  int
	symbol string
}

type Formated struct {
	Symbol       string                     `json:"Stock Symbol"`
	AverageClose string                     `json:"Average Closing Price"`
	Tds          map[string]TimeSeriesDaily `json:"Time Series (Daily)"`
}

/*
Main v1 Get func: most of everthing is done here.
We call the helpper func to fetch the data, we compute average data
*/
func v1Get() ([]byte, error) {
	var output Formated
	output.Tds = make(map[string]TimeSeriesDaily)

	// Get the data from the AlaphVantage API
	raw, err := getRawData()
	if err != nil {
		return nil, err
	}
	data, err := parseData(raw)
	if err != nil {
		return nil, err
	}

	// Get a list of the last ndays
	days, err := listOfDays(data)
	if err != nil {
		return nil, err
	}

	// Start formating the output Struct
	output.Symbol = data.Meta.Symbol

	var closedtotal float64
	for index := range days {
		if val, key := data.Tds[days[index]]; key {
			closed, err := strconv.ParseFloat(val.Close, 64)
			if err == nil {
				closedtotal = closedtotal + closed
			}
			output.Tds[days[index]] = val
		}
	}
	// only compute the average if we actually have data
	if len(output.Tds) != 0 {
		closedAverage := closedtotal / float64(len(output.Tds))
		output.AverageClose = fmt.Sprintf("%.2f", closedAverage)
	}

	return json.Marshal(output)

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
		return values, errors.New("ndays env not an int")
	}
	// we default to format compact in the api so we need to have NDAYS set to less than 100 and more than 0
	if ndays > 100 || ndays <= 0 {
		fmt.Println("Ndays outside the require range.")
		health.ApiReady = false
		return values, errors.New("ndays outside of range")
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
		return nil, errors.New("not 200 return")
	}

	// Read the body response
	bodybytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bodybytes, nil
}

// Parse raw data from the url fetch : this was way easier than I remember doing with XML.
func parseData(raw []byte) (AlphaVantageData, error) {
	var aData AlphaVantageData

	err := json.Unmarshal(raw, &aData)
	return aData, err

}

// Looks at current time and gets an array of dates ndays into the past
func listOfDays(data AlphaVantageData) ([]string, error) {
	var days []string
	values, err := fetchValues()
	if err != nil {
		return nil, errors.New("error while validating env values")
	}

	// Make sure the current timezone matches the data
	location, err := time.LoadLocation(data.Meta.Timezone)
	if err != nil {
		return days, err
	}
	now := time.Now().In(location)
	days = append(days, now.Format("2006-01-02"))

	pastDays := 1
	for pastDays < values.ndays {
		past := now.AddDate(0, 0, -pastDays)
		days = append(days, past.Format("2006-01-02"))
		pastDays++

	}

	return days, nil
}
