package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//Check ENV variables.
	envChecks()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Currency Converter",
		})
	})
	router.GET("/converter", func(c *gin.Context) {
		amount := c.Query("amount")
		from := c.Query("from")
		to := c.Query("to")

		data, err := GetAllCurrency()
		if err != nil {
			log.Fatal("NewRequest: ", err)
		}

		var fromValue float64
		fromValue = getCurrencyValue(from, data)

		var toValue float64
		toValue = getCurrencyValue(to, data)

		var convertedAmount float64

		floatAmount, _ := strconv.ParseFloat(amount, 8)
		convertedAmount = (floatAmount * toValue) / fromValue

		var response struct {
			From   string  `json:"from"`
			To     string  `json:"to"`
			Amount float64 `json:"amount"`
		}
		response.From = from
		response.To = to
		response.Amount = convertedAmount

		c.JSON(http.StatusOK, response)
	})
	port := os.Getenv("PORT")
	router.Run(":" + port)
}

func envChecks() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	port, portExist := os.LookupEnv("PORT")

	if !portExist || port == "" {
		log.Fatal("PORT must be set in .env and not empty")
	}

	apiToken, apiExist := os.LookupEnv("FIXER_API")

	if !apiExist || apiToken == "" {
		log.Fatal("FIXER_API must be set in .env and not empty")
	}
}

func getCurrencyValue(currency string, data Currency) float64 {
	var value float64
	if currency == "EUR" {
		value = data.Currencies.EUR
	} else if currency == "USD" {
		value = data.Currencies.USD
	} else if currency == "NGN" {
		value = data.Currencies.NGN
	} else if currency == "CAD" {
		value = data.Currencies.CAD
	} else if currency == "BTC" {
		value = data.Currencies.BTC
	}
	return value
}

type Currency struct {
	BaseCurrency string `json:"base"`
	Currencies   struct {
		USD float64 `json:"USD"`
		NGN float64 `json:"NGN"`
		CAD float64 `json:"CAD"`
		EUR float64 `json:"EUR"`
		BTC float64 `json:"BTC"`
	} `json:"rates"`
}

func GetAllCurrency() (Currency, error) {
	apiToken := os.Getenv("FIXER_API")
	url := fmt.Sprintf("http://data.fixer.io/api/latest?access_key=%s&format=1", apiToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return Currency{}, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return Currency{}, err
	}

	defer resp.Body.Close()

	var record Currency

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	return record, nil
}
