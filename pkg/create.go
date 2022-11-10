package dbshaker

import (
	"fmt"
	"github.com/ToggyO/dbshaker/internal"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	GoTemplate              MigrationTemplateType = "go"
	SqlTemplate             MigrationTemplateType = "sql"
	migrationFileNameFormat                       = ""
	timestampFormat                               = "20060102150405"
)

type MigrationTemplateType string

type templateData struct {
	MName string
}

// TODO: добавить логгирование
func CreateMigrationTemplate(name, directory string, templateType MigrationTemplateType) error {
	version := time.Now().Format(timestampFormat)
	filename := fmt.Sprintf("%v_%s.%s", version, strcase.ToSnake(name), templateType)

	path := filepath.Join(directory, filename)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// TODO: вынести
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		// TODO: вынести
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	defer file.Close()

	migrationFuncName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	tmplData := templateData{MName: strcase.ToCamel(migrationFuncName)}
	template := internal.GoMigrationTemplate
	if templateType == SqlTemplate {
		// TODO: add sql tempalte
		//template =
	}

	if err = template.Execute(file, tmplData); err != nil {
		return err
	}

	logger.Printf("Created new file: %s\n", file.Name())
	return nil
}
