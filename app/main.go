package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testfromQiwi/models"
	"time"
)

// ссылка ресурса
const cbrAPIURL = "https://www.cbr.ru/scripts/XML_daily.asp"

func getCurrencyRates(code, date string) (models.ValCurs, error) {
	// Преобразуем дату в формат "dd/mm/yyyy"
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return models.ValCurs{}, fmt.Errorf("неверный формат даты, нужен: yyyy-mm-dd")
	}
	formattedDate := fmt.Sprintf("%s/%s/%s", parts[2], parts[1], parts[0])

	url := fmt.Sprintf("%s?date_req=%s", cbrAPIURL, formattedDate)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return models.ValCurs{}, err
	}
	defer resp.Body.Close()

	fmt.Println(resp)

	if resp.StatusCode != http.StatusOK {
		return models.ValCurs{}, fmt.Errorf("Ошибка получения данных. Status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.ValCurs{}, err
	}

	var valCurs models.ValCurs
	err = xml.Unmarshal(body, &valCurs)
	if err != nil {
		return models.ValCurs{}, err
	}

	return valCurs, nil
}
func isValidDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func main() {
	codePtr := flag.String("code", "", "Currency code in ISO 4217 format")
	datePtr := flag.String("date", "", "Date in YYYY-MM-DD format")
	flag.Parse()

	if *codePtr == "" || *datePtr == "" || !isValidDateFormat(*datePtr) {
		fmt.Println("Usage: currency_rates --code=USD --date=2022-10-08")
		return
	}

	valCurs, err := getCurrencyRates(*codePtr, *datePtr)
	if err != nil {
		fmt.Printf("Error getting currency rates: %s\n", err)
		return
	}

	fmt.Printf("Currency rates for %s on %s:\n", *codePtr, valCurs.Date)
	found := false
	for _, valute := range valCurs.ValuteArr {
		if strings.ToUpper(valute.CharCode) == strings.ToUpper(*codePtr) {
			fmt.Printf("%s (%s): %s \n", valute.CharCode, valute.Name, valute.Value)
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Currency with code %s not found for the specified date.\n", *codePtr)
	}
}
