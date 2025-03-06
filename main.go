package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	internal "github.com/fernandofreamunde/pokedexcli/internal/cache"
)

const prompt = "Pokedex > "

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	next     string
	previous string
	cache    *internal.PokeCache
}

var supportedCommands map[string]cliCommand

func main() {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	config := config{}
	config.cache = internal.NewCache(30)

	supportedCommands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the names of location areas in the Pokemon World 20 items per page use 'map' for each page.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Displays previous 20 the names of location areas.",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Displays pokemon names of provided location areas.",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	var userInput []string
	for scanner.Scan() {
		userInput = cleanInput(scanner.Text())
		parameters := userInput[1:]
		command, ok := supportedCommands[userInput[0]]
		if !ok {
			fmt.Println("Unknown command")
			fmt.Print(prompt)
			continue
		}

		command.callback(&config, parameters)
		fmt.Print(prompt)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, " shouldnt see an error scanning a string")
	}
}

func cleanInput(text string) []string {
	t := strings.Fields(strings.ToLower(text))

	return t
}

func commandExit(config *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, params []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, v := range supportedCommands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(config *config, params []string) error {

	url := "https://pokeapi.co/api/v2/location-area"
	if len(config.next) > 0 {
		url = config.next
	}

	locations, err := queryLocationAreas(config, url)
	if err != nil {
		return err
	}

	for _, l := range locations.Results {
		fmt.Println(l.Name)
	}

	return nil
}

func commandMapBack(config *config, params []string) error {

	url := "https://pokeapi.co/api/v2/location-area"
	if len(config.previous) > 0 {
		url = config.previous
	}

	locations, err := queryLocationAreas(config, url)
	if err != nil {
		return err
	}

	for _, l := range locations.Results {
		fmt.Println(l.Name)
	}

	return nil
}

func commandExplore(config *config, params []string) error {
	city := params[0]
	url := "https://pokeapi.co/api/v2/location-area/" + city

	location, err := queryLocationArea(config, url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("exploring: %s...\n", location.Name)
	fmt.Println("Found Pokemon:")
	for _, e := range location.PokemonEncounters {
		fmt.Println(" - " + e.Pokemon.Name)
	}

	return nil
}

type PokemonVersion struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationPokeApiResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int            `json:"rate"`
			Version PokemonVersion `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	Id        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int `json:"chance"`
				ConditionValues any `json:"condition_values"`
				MaxLevel        int `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int            `json:"max_chance"`
			Version   PokemonVersion `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func queryLocationArea(config *config, url string) (LocationPokeApiResponse, error) {

	body, err := queryApi(config, url)
	if err != nil {
		return LocationPokeApiResponse{}, fmt.Errorf("error: %s", err)
	}

	location := LocationPokeApiResponse{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return LocationPokeApiResponse{}, fmt.Errorf("error: %s", err)
	}

	return location, nil
}

func queryApi(config *config, url string) ([]byte, error) {
	body, hit := config.cache.Get(url)

	if !hit {
		// time.Sleep(1 * time.Second) // this does work
		res, err := http.Get(url)
		if err != nil {
			return []byte{}, fmt.Errorf("something whent wrong when requesting location data")
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return []byte{}, fmt.Errorf("error: %s", err)
		}
		if res.StatusCode > 299 {
			return []byte{}, fmt.Errorf("NetworkError: Response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		}

		config.cache.Add(url, body)
	}

	return body, nil
}

type PokeApiResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func queryLocationAreas(config *config, url string) (PokeApiResponse, error) {

	body, err := queryApi(config, url)
	if err != nil {
		return PokeApiResponse{}, fmt.Errorf("error: %s", err)
	}

	locations := PokeApiResponse{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return PokeApiResponse{}, fmt.Errorf("error: %s", err)
	}

	config.next = locations.Next
	if locations.Previous != nil {
		if prevURL, ok := locations.Previous.(string); ok {
			config.previous = prevURL
		} else {
			fmt.Println("Previous is not a string")
		}
	} else {
		config.previous = ""
	}

	return locations, nil
}
