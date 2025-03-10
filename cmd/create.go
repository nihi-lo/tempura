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

func templateExists(folderName string) (bool, error) {
	// ディレクトリ内のファイルとサブディレクトリを取得
	entries, err := fs.ReadDir(templates.Templates, ".")
	if err != nil {
		return false, fmt.Errorf("failed to read template dir: %v", err)
	}

	// 各エントリをチェック
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == folderName {
			// フォルダが見つかった場合
			return true, nil
		}
	}

	// フォルダが見つからなかった場合
	return false, nil
}

func createProject(projectPath string, templateName string, data TemplateInputData) error {
	return fs.WalkDir(templates.Templates, templateName, func(path string, d fs.DirEntry, err error) error {
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

		// コマンドライン引数および、フラグを取得する
		projectName := ""
		if len(args) >= 1 {
			projectName = args[0]
		}

		templateName, err := cmd.Flags().GetString("template")
		if err != nil {
			return fmt.Errorf("failed to read template flag: %w", err)
		}

		// プロジェクト名が未指定の場合、ユーザーにプロジェクト名の入力を求める
		// その後、現在の実行ディレクトリにすでに同名プロジェクトがないか確認する
		if projectName == "" {
			ml, err := tea.NewProgram(tui.InitialProjectNameInputModel()).Run()
			if err != nil {
				return err
			}

			projectName = ml.(tui.ProjectNameInputModel).Input.Value()
			if projectName == "" {
				return fmt.Errorf("project name is empty")
			}
		}

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

		// テンプレート名が未指定の場合、ユーザーにテンプレートの選択を求める
		// その後、テンプレート名が有効かどうか確認する
		if templateName == "" {
			tsm, err := tea.NewProgram(tui.InitialTemplateSelectModel()).Run()
			if err != nil {
				return err
			}

			templateName = tsm.(tui.TemplateSelectModel).Choice
			if templateName == "" {
				return nil
			}
		}

		exists, err := templateExists(templateName)
		if err != nil {
			return fmt.Errorf("failed to template name validity check: %w", err)
		}
		if !exists {
			return fmt.Errorf("invalid template name")
		}

		// テンプレートからプロジェクトを作成する
		err = createProject(projectPath, templateName, TemplateInputData{
			ProjectName: projectName,
		})
		if err != nil {
			log.Fatalf("error during file processing: %v\n", err)
		}

		fmt.Printf("\nDone. Let's start developing!\n\n")
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("template", "t", "", "template name")
	rootCmd.AddCommand(createCmd)
}
