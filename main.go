package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/tidwall/gjson"
)

func main() {
	fmt.Println("Obteniendo info de comunas...")
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
	req, err := http.NewRequest(http.MethodGet, "https://e.infogram.com/81277d3a-5813-46f7-a270-79d1768a70b2", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	regex := regexp.MustCompile(`window.infographicData(.*)`)

	foundStr := regex.FindAllString(string(body), -1)[0]
	return parseData(foundStr)
}

func parseData(text string) (string, error) {
	str := strings.Trim(text, "window.infographicData=")
	str = strings.TrimRight(str, ";</script>")
	if !json.Valid([]byte(str)) {
		log.Println(str)
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
