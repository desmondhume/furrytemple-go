package pageBuilder

import (
	"fmt"
	"html/template"
	"os"
)

var (
	templatesFolder = os.Getenv("FURRYTEMPLE_HTML_TEMPLATES_FOLDER")
	outputFolder    = os.Getenv("FURRYTEMPLE_HTML_OUTPUT_FOLDER")
)

func buildPage(templateName string, data interface{}, outputFile *os.File) error {
	templatePath := fmt.Sprintf("%s%s.html", templatesFolder, templateName)
	t, err := template.ParseFiles(templatePath)

	if err != nil {
		return err
	}

	t.Execute(outputFile, data)

	return err
}

type Video struct {
	Title string
}

func BuildHomepage() error {
	var err error
	outputFilePath := fmt.Sprintf("%s%s.html", outputFolder, "home")
	outputFile, err := os.Create(outputFilePath)

	if err != nil {
		return err
	}

	videos := []Video{
		Video{"Video#1"},
		Video{"Video#2"},
		Video{"Video#3"},
	}

	data := map[string]interface{}{
		"Title":  "Furrytemple",
		"Videos": videos,
	}

	buildPage("home", data, outputFile)
	return err
}
