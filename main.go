package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	"github.com/tidwall/gjson"
)

func main() {
	fmt.Println("Fetching recent data...")
	data, err := scrapeData()
	if err != nil {
		log.Fatalln(err)
	}
	// data: [["Aisen", "4", "Apertura"],...]

	mappedData := buildMap(data)
	input, err := promptInput()
	if err != nil {
		log.Fatalln(err)
	}
	input = strings.Trim(input, "\n")

	results := make([]string, 0, len(mappedData))
	for commune := range mappedData {
		if strings.Contains(strings.ToLower(commune), input) {
			results = append(results, commune)
		}
	}

	selected, err := promptSelect(results)
	if err != nil {
		log.Fatalln(err)
	}
	selectedCommune := mappedData[selected]
	fmt.Println("Nombre:\t", selectedCommune.Name)
	fmt.Println("Fase:\t", selectedCommune.Phase)
	fmt.Println("Estado:\t", selectedCommune.Status)
}

func scrapeData() (string, error) {
	var (
		finalJSON string
		err       error
	)
	c := colly.NewCollector()

	c.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "window.infographicData") {
			finalJSON, err = parseData(e.Text)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://e.infogram.com/81277d3a-5813-46f7-a270-79d1768a70b2")
	return finalJSON, err
}

func parseData(text string) (string, error) {
	str := strings.Trim(text, "window.infographicData=")
	str = strings.TrimRight(str, ";")
	if !json.Valid([]byte(str)) {
		return "", fmt.Errorf("invalid json format")
	}

	return str, nil
}

type commune struct {
	Name   string
	Phase  string
	Status string
}

func buildMap(data string) map[string]*commune {
	mappedData := make(map[string]*commune)
	results := gjson.Get(data, "elements.content.content.entities")
	for _, result := range results.Map() {
		data := result.Get("props.chartData.data")
		if data.String() != "" {
			tableData := data.Array()[0]
			for _, datum := range tableData.Array()[1:] {
				// expected array: ["name", "phase", "status"]
				datumArray := datum.Array()
				communeName := strings.TrimSpace(datumArray[0].String())
				mappedData[communeName] = &commune{
					Name:   communeName,
					Phase:  datumArray[1].String(),
					Status: datumArray[2].String(),
				}
			}
		}
	}

	return mappedData
}

func promptInput() (string, error) {
	prompt := promptui.Prompt{
		Label: "Comuna",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("error while prompting input: %w", err)
	}

	return result, nil
}

func promptSelect(opts []string) (string, error) {
	prompt := promptui.Select{
		Label: "Resultados encontrados",
		Items: opts,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("error while prompting select: %w", err)
	}

	return result, nil
}
