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

type TemplateInputData struct {
	ProjectName string
}

func copyTemplateFile(src string, dest string, data TemplateInputData) error {
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
	Aliases: []string{"c", "init"},
	Short:   "Select a template to create a project",
	Example: "tempura create",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 事前設定: エラー発生時にコマンドの使い方を表示させない
		cmd.SilenceUsage = true

		// ユーザーにプロジェクト名の入力を求める
		ml, err := tea.NewProgram(tui.InitialModel()).Run()
		if err != nil {
			return err
		}

		projectName := ml.(tui.ProjectNameInputModel).Input.Value()
		if projectName == "" {
			return fmt.Errorf("project name is empty")
		}

		// 現在の実行ディレクトリにすでに同名プロジェクトがないか確認する
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		projectPath := filepath.Join(wd, projectName)
		if fi, err := os.Stat(projectPath); err == nil {
			if !fi.IsDir() {
				return fmt.Errorf("not a directory. scandir '%s'", filepath.Base(projectPath))
			}

			entries, err := os.ReadDir(projectPath)
			if err != nil {
				return fmt.Errorf("failed to read directory: %w", err)
			}
			if len(entries) != 0 {
				return fmt.Errorf("\"%s\" already exists and isn't empty", filepath.Base(projectPath))
			}
		}

		// ユーザーにテンプレートの選択を求める
		templateName := "vite-react-tw3-ts"

		// テンプレートからプロジェクトを作成する
		data := TemplateInputData{
			ProjectName: projectName,
		}

		err = fs.WalkDir(templates.Templates, templateName, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk dir %s: %w", path, err)
			}

			relPath, err := filepath.Rel(templateName, path)
			if err != nil {
				return fmt.Errorf("failed to calculate relative path %s: %w", path, err)
			}
			destPath := filepath.Join(projectPath, relPath)

			if d.IsDir() {
				if err := os.MkdirAll(destPath, 0755); err != nil {
					log.Printf("warning: failed to create directory %s: %v", destPath, err)
				}
				return nil
			}

			if filepath.Ext(path) == ".tmpl" {
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
