package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nihi-lo/tempura/internal/templates"
	"github.com/nihi-lo/tempura/internal/tui"
	"github.com/spf13/cobra"
)

type TemplateData struct {
	ProjectName string
}

func copyTemplateFile(src string, dest string, data TemplateData) error {
	content, err := templates.Templates.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", src, err)
	}

	tmpl, err := template.New(filepath.Base(src)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", src, err)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", dest, err)
	}
	defer destFile.Close()

	if err := tmpl.Execute(destFile, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", src, err)
	}

	return nil
}

func copyFile(src string, dest string) error {
	if _, err := os.Stat(dest); err == nil {
		log.Printf("info: file %s already exists, overwriting...", dest)
	}

	srcFile, err := templates.Templates.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dest, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy data from %s to %s: %w", src, dest, err)
	}

	return nil
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Select a template to create a project",
	Example: "tempura create",
	RunE: func(cmd *cobra.Command, args []string) error {
		ml, err := tea.NewProgram(tui.InitialModel()).Run()
		if err != nil {
			return err
		}
		projectName := ml.(tui.ProjectNameInputModel).Input.Value()

		templateName := "vite-react-tw3-ts"

		err = fs.WalkDir(templates.Templates, templateName, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk dir %s: %w", path, err)
			}

			relPath, err := filepath.Rel(templateName, path)
			if err != nil {
				return fmt.Errorf("failed to calculate relative path %s: %w", path, err)
			}
			destPath := filepath.Join(projectName, relPath)

			if d.IsDir() {
				if err := os.MkdirAll(destPath, 0755); err != nil {
					log.Printf("warning: failed to create directory %s: %v", destPath, err)
				}
				return nil
			}

			if filepath.Ext(path) == ".tmpl" {
				data := TemplateData{
					ProjectName: projectName,
				}
				if err := copyTemplateFile(path, strings.TrimSuffix(destPath, ".tmpl"), data); err != nil {
					log.Printf("warning: failed to process template file %s: %v", path, err)
				}
			} else {
				if err := copyFile(path, destPath); err != nil {
					log.Printf("warning: failed to copy file %s: %v", path, err)
				}
			}

			return nil
		})
		if err != nil {
			log.Fatalf("error during file processing: %v\n", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
