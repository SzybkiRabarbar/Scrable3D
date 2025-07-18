package ctrl

import (
	"bufio"
	"fmt"
	"os"
	"scrable3/internal/cfg"
	"strings"
)

type WordsController interface {
	WordsNumber() int
	CheckWord(word string) error
}

type wordsController struct {
	words *map[string]bool
}

func NewWordsController(filePath string) (WordsController, error) {
	wordsController := wordsController{}
	wordMap := make(map[string]bool)
	wordsController.words = &wordMap

	fl, err := os.Open(filePath)
	if err != nil {
		return &wordsController, err
	}
	defer fl.Close()

	scanner := bufio.NewScanner(fl)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < cfg.MIN_WORD_LEN {
			continue
		}
		(*wordsController.words)[line] = true
	}
	return &wordsController, nil
}

func (wc *wordsController) WordsNumber() int {
	return len(*wc.words)
}

func (wc *wordsController) CheckWord(word string) error {
	if len(word) < cfg.MIN_WORD_LEN {
		return fmt.Errorf("word %v is too short", word)
	}
	_, ok := (*wc.words)[strings.ToLower(word)]
	if !ok {
		return fmt.Errorf("word %v is not on words list", word)
	}
	return nil
}
