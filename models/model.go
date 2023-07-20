package models

type ValCurs struct {
	Date      string `xml:"Date,attr"`
	ValuteArr []Valute
}

type Valute struct {
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}
