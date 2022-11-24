package sql

import "github.com/ToggyO/dbshaker/internal"

type sqlParseMarkers struct {
	parseStartMarker    string
	statementBegin      string
	statementEnd        string
	parseEndMarker      string
	noTransactionMarker string
}

func newMigrateUpParseMarkers() sqlParseMarkers {
	return sqlParseMarkers{
		parseStartMarker:    internal.MarkerMigrateUpStart,
		statementBegin:      internal.MarkerStatementBegin,
		statementEnd:        internal.MarkerStatementEnd,
		parseEndMarker:      internal.MarkerMigrateUpEnd,
		noTransactionMarker: internal.MarkerNoTransaction,
	}
}

func newMigrateDownParseMarkers() sqlParseMarkers {
	return sqlParseMarkers{
		parseStartMarker:    internal.MarkerMigrateDownStart,
		statementBegin:      internal.MarkerStatementBegin,
		statementEnd:        internal.MarkerStatementEnd,
		parseEndMarker:      internal.MarkerMigrateDownEnd,
		noTransactionMarker: internal.MarkerNoTransaction,
	}
}
