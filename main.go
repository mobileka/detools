package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"fmt"
	"sort"
)

type callback func(i int, s *goquery.Selection)

type Config struct {
	url string
	selector string
	callback
}

func main() {
	channel := make(chan string)

	gfesCallback := func(i int, s *goquery.Selection) {
		if i > 0 {
			tds := s.Children()
			infinitive := tds.First().Text()
			participle := tds.Filter("td:nth-child(4)").Text()
			toCheck := map[string]int{"hat ge" + infinitive: 1, "ist ge" + infinitive: 1}

			if toCheck[participle] == 1 {
				channel <- participle
			}
		}
	}

	coLanguageCallback := func(i int, s *goquery.Selection) {
		if i > 0 {
			tds := s.Children()
			infinitive := tds.Filter("td:nth-child(2)").Text()
			participle := tds.Filter("td:nth-child(5)").Text()
			toCheck := "ge" + infinitive

			if toCheck == participle {
				channel <- "hat " + participle
			}
		}
	}

	configs := []Config{
		{
			"https://www.colanguage.com/irregular-verbs-strong-verbs-german",
			"table tr",
			coLanguageCallback,
		},
		{
			"http://germanforenglishspeakers.com/reference/strong-verbs/",
			"table tr",
			gfesCallback,
		},
	}

	for _, config := range configs {
		go parse(config)
	}

	participles := map[string]int{}

	participles[<- channel] = 1

	for i := 0; i < len(channel); i++ {
		participles[<- channel] = 1
	}

	elements := len(participles)
	result := make([]string, elements)

	i := 0
	for key := range participles {
		result[i] = key
		i++
	}

	fmt.Printf("%d\n", elements)
	sort.Strings(result)
	for _, participle := range result {
		fmt.Printf("%s\n", participle)
	}
}

func parse(config Config)  {
	document, err := goquery.NewDocument(config.url)
	if err != nil {
		log.Fatal(err)
	}

	document.Find(config.selector).Each(config.callback)
}
