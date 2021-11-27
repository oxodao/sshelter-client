package cmd

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/oxodao/sshelter_client/services"
)

//go:embed template.yml
var templateFile string

func GenerateSkeleton(prv *services.Provider, filename string) {
	if !strings.HasSuffix(filename, ".yml") || !strings.HasSuffix(filename, ".yaml") {
		filename += ".yml"
	}

	fmt.Println("Generating the template to create a machine...")
	fmt.Println("Edit the file " + filename + " and then run `sshelter --create " + filename + "`")

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file: " + err.Error())
		os.Exit(1)
	}

	f.WriteString(templateFile)

	f.Close()

	os.Exit(0)
}
