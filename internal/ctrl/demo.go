package ctrl

import (
	"bytes"
	"errors"
	"scrable3/internal/dto"
	"text/template"

	"golang.org/x/exp/rand"
)

func GetRandomField() ([]byte, error) {
	return randomFieldHtmlContent()
}

func GetExampleError() ([]byte, error) {
	var g []byte
	return g, errors.New("example error message")
}

func randomFieldHtmlContent() ([]byte, error) {
	var htmlContent []byte
	tmpl, err := template.ParseFiles("views/game/field.html")
	if err != nil {
		return htmlContent, err
	}

	x := rand.Intn(15)
	y := rand.Intn(15)
	z := rand.Intn(15)

	letter := string(randomUppercaseLetter())

	data := dto.NewHtmlFieldData(
		letter,
		x,
		y,
		z,
	)
	var buff bytes.Buffer
	err = tmpl.Execute(&buff, data)
	if err != nil {
		return htmlContent, err
	}
	htmlContent = buff.Bytes()
	return htmlContent, nil
}

func randomUppercaseLetter() rune {
	// Generate a random number between 0 and 25
	randomIndex := rand.Intn(26)

	// Convert the random number to an uppercase letter
	return rune('A' + randomIndex)
}
