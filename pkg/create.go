package dbshaker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/internal/templates"
	"github.com/iancoleman/strcase"
)

const (
	GoTemplate              MigrationTemplateType = "go"
	SQLTemplate             MigrationTemplateType = "sql"
	migrationFileNameFormat                       = "%v_%s.%s"
	timestampFormat                               = "20060102150405"
)

type MigrationTemplateType string

type templateData struct {
	MName string
}

func CreateMigrationTemplate(name, directory string, templateType MigrationTemplateType) error {
	version := time.Now().Format(timestampFormat)
	filename := fmt.Sprintf(migrationFileNameFormat, version, strcase.ToSnake(name), templateType)

	path := filepath.Join(directory, filename)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return internal.ErrFailedToCreateMigration(err)
	}

	file, err := os.Create(path)
	if err != nil {
		return internal.ErrFailedToCreateMigration(err)
	}

	defer file.Close()

	migrationFuncName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	tmplData := templateData{MName: strcase.ToCamel(migrationFuncName)}
	template := templates.GoMigrationTemplate
	if templateType == SQLTemplate {
		template = templates.SQLMigrationTemplate
	}

	if err = template.Execute(file, tmplData); err != nil {
		return err
	}

	logger.Printf("created new file: %s\n", file.Name())
	return nil
}
