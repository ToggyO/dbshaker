package postgres

import (
	"fmt"
	"strings"

	"github.com/ToggyO/dbshaker/tests/suites"
)

func CreatePgConnectionString(conf suites.DBConf) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("host=%s port=%d", conf.Host, conf.Port))

	if conf.User != "" {
		sb.WriteString(fmt.Sprintf(" user=%s", conf.User))
	}
	if conf.Password != "" {
		sb.WriteString(fmt.Sprintf(" password=%s", conf.Password))
	}
	if conf.Name != "" {
		sb.WriteString(fmt.Sprintf(" dbname=%s", conf.Name))
	}

	sb.WriteString(" sslmode=disable")
	return sb.String()
}
