package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Boostport/mjml-go"
)

func ConvertMJMLToHTML(inputDir, outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	files, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("failed to read input directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".mjml" {
			continue
		}

		inputPath := filepath.Join(inputDir, file.Name())
		outputPath := filepath.Join(outputDir, file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]+".html")

		mjmlContent, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read MJML file %s: %v", inputPath, err)
		}

		htmlContent, err := mjml.ToHTML(context.Background(), string(mjmlContent), mjml.WithMinify(true))
		if err != nil {
			var mjmlError mjml.Error
			if errors.As(err, &mjmlError) {
				fmt.Printf("MJML error: %s, details: %v\n", mjmlError.Message, mjmlError.Details)
			}
			return fmt.Errorf("failed to render MJML to HTML for file %s: %v", inputPath, err)
		}

		err = os.WriteFile(outputPath, []byte(htmlContent), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write HTML file %s: %v", outputPath, err)
		}
	}

	return nil
}
