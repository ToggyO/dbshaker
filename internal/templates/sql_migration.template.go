package templates

import "text/template"

var SQLMigrationTemplate = template.Must(template.New("dbshaker.sql-migration").Parse(`-- +dbshaker UpStart

SELECT 'up SQL query';

-- +dbshaker UpEnd

-- +dbshaker DownStart

SELECT 'down SQL query';

-- +dbshaker DownEnd`))
