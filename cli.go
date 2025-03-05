// boot.dev project personal project 1, 2025

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TatoebaPage struct {
	Paging struct {
		Sentences struct {
			Finder           string `json:"finder"`
			Page             int    `json:"page"`
			Current          int    `json:"current"`
			Count            int    `json:"count"`
			PerPage          int    `json:"perPage"`
			Start            int    `json:"start"`
			End              int    `json:"end"`
			PrevPage         bool   `json:"prevPage"`
			NextPage         bool   `json:"nextPage"`
			PageCount        int    `json:"pageCount"`
			Sort             any    `json:"sort"`
			Direction        any    `json:"direction"`
			Limit            any    `json:"limit"`
			SortDefault      bool   `json:"sortDefault"`
			DirectionDefault bool   `json:"directionDefault"`
			Scope            any    `json:"scope"`
			CompleteSort     []any  `json:"completeSort"`
		} `json:"Sentences"`
	} `json:"paging"`
	Results []struct {
		ID           int    `json:"id"`
		Text         string `json:"text"`
		Lang         string `json:"lang"`
		Correctness  int    `json:"correctness"`
		Script       string `json:"script"`
		License      string `json:"license"`
		Translations [][]struct {
			ID             int    `json:"id"`
			Text           string `json:"text"`
			Lang           string `json:"lang"`
			Correctness    int    `json:"correctness"`
			Script         any    `json:"script"`
			Transcriptions []any  `json:"transcriptions"`
			Audios         []any  `json:"audios"`
			IsDirect       bool   `json:"isDirect"`
			LangName       string `json:"lang_name"`
			Dir            string `json:"dir"`
			LangTag        string `json:"lang_tag"`
		} `json:"translations"`
		Transcriptions []struct {
			ID          int       `json:"id"`
			SentenceID  int       `json:"sentence_id"`
			Script      string    `json:"script"`
			Text        string    `json:"text"`
			UserID      any       `json:"user_id"`
			NeedsReview bool      `json:"needsReview"`
			Modified    time.Time `json:"modified"`
			User        any       `json:"user"`
			Readonly    bool      `json:"readonly"`
			Type        string    `json:"type"`
			HTML        string    `json:"html"`
			Markup      any       `json:"markup"`
			InfoMessage string    `json:"info_message"`
		} `json:"transcriptions"`
		Audios []any `json:"audios"`
		User   struct {
			Username string `json:"username"`
		} `json:"user"`
		LangName               string `json:"lang_name"`
		Dir                    string `json:"dir"`
		LangTag                string `json:"lang_tag"`
		IsFavorite             any    `json:"is_favorite"`
		IsOwnedByCurrentUser   bool   `json:"is_owned_by_current_user"`
		Permissions            any    `json:"permissions"`
		MaxVisibleTranslations int    `json:"max_visible_translations"`
		CurrentUserReview      any    `json:"current_user_review"`
	} `json:"results"`
}

type Config struct {
	argument              string
	query                 string
	page                  TatoebaPage
	current_example       string
	current_example_index int
}

//const search_url = "https://tatoeba.org/en/sentences/search?query="
// 白天

func removeHTMLTags(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(input, "")
}

func main() {

	registered_commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the CLI",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Help with Tatoeba CLI",
			callback:    commandHelp,
		},
		"example": {
			name:        "example",
			description: "Get example sentences using arguments from Tatoeba.org",
			callback:    commandExample,
		},
		"query": {
			name:        "query",
			description: "Get current query",
			callback:    commandQuery,
		},
		"current": {
			name:        "current",
			description: "Get current example sentence",
			callback:    commandCurrent,
		},
		"next": {
			name:        "next",
			description: "Get next example",
			callback:    commandNext,
		},
		"back": {
			name:        "back",
			description: "Get previous example",
			callback:    commandBack,
		},
		"note": {
			name:        "note",
			description: "Create note based on current example",
			callback:    commandNote,
		},
	}

	config := &Config{
		argument:              "",
		page:                  TatoebaPage{},
		current_example:       "",
		query:                 "",
		current_example_index: -1,
	}

	m_scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Sentence Example Generator > ")
		m_scanner.Scan()
		t := m_scanner.Text()
		m_command := strings.Fields(t)
		value, exists := registered_commands[m_command[0]]
		if exists {
			if len(m_command) > 1 {
				config.argument = m_command[1]
			} else {
				config.argument = ""
			}
			value.callback(config)
		} else {
			fmt.Printf("Unknown command")
		}
	}
}

// save to .txt or db
func commandNote(c *Config) error {
	fmt.Println("Not yet implemented")
	return errors.New("")
}

func commandBack(c *Config) error {
	if c.current_example != "" {
		if c.current_example_index > 0 {
			c.current_example = c.page.Results[c.current_example_index-1].Text
			c.current_example_index = c.current_example_index - 1
			fmt.Println("Previous Example: " + c.page.Results[c.current_example_index].Text)
		}
	}
	return errors.New("")
}

func commandNext(c *Config) error {
	if c.current_example != "" {
		if c.page.Paging.Sentences.Count-1 > c.current_example_index {
			c.current_example = c.page.Results[c.current_example_index+1].Text
			c.current_example_index = c.current_example_index + 1
			fmt.Println("Next Example: " + c.page.Results[c.current_example_index].Text)
		}
	}
	return errors.New("")
}

func commandCurrent(c *Config) error {
	if c.current_example != "" {
		fmt.Println("Current example based on " + c.query + ": " + c.current_example)
	}
	return errors.New("")
}

func commandQuery(c *Config) error {
	if c.query != "" {
		fmt.Println("Current query from Tatoeba.org: " + c.query)
	}
	return errors.New("")
}

func commandExample(c *Config) error {
	resp, err := http.Get("https://tatoeba.org/en/api_v0/search?query=" + c.argument)
	if err != nil {
		fmt.Println("error: httpGET", err)
		return errors.New("")
	}
	defer resp.Body.Close()

	// Check the Content-Type header
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		fmt.Println("Unexpected Content-Type:", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error: io read ", err)
		return errors.New("")
	}

	//fmt.Println("Response Body:", string(body))
	cleanedBody := removeHTMLTags(string(body))

	sentences_page := TatoebaPage{}
	err = json.Unmarshal([]byte(cleanedBody), &sentences_page)
	if err != nil {
		fmt.Println("error: json unmarshal: ", err)
		return errors.New("")
	}

	c.page = sentences_page

	if c.page.Paging.Sentences.Count > 0 {
		fmt.Println("Number of examples for " + c.argument + ": " + strconv.Itoa(sentences_page.Paging.Sentences.Count))
		c.current_example = c.page.Results[0].Text
		c.current_example_index = 0
		c.query = c.argument
		fmt.Println("Current example: " + c.current_example)
	} else {
		fmt.Println("no examples for " + c.argument)
	}

	return errors.New("")
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Tatoeba API cli!\nUsage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the cli")
	fmt.Println("example: gets examples of argument from Tatoeba.org")
	fmt.Println("query: current query from Tatoeba.org")
	fmt.Println("current: Get current example sentence based on argument given to example")
	fmt.Println("next: gets next example sentence")
	fmt.Println("back: gets previous example sentence")
	return errors.New("")
}

func commandExit(c *Config) error {
	os.Exit(0)
	return errors.New("")
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config) error
}
