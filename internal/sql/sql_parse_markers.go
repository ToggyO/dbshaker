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
		parseStartMarker:    internal.markerMigrateUpStart,
		statementBegin:      internal.markerStatementBegin,
		statementEnd:        internal.markerStatementEnd,
		parseEndMarker:      internal.markerMigrateUpEnd,
		noTransactionMarker: internal.markerNoTransaction,
	}
}

func newMigrateDownParseMarkers() sqlParseMarkers {
	return sqlParseMarkers{
		parseStartMarker:    internal.markerMigrateDownStart,
		statementBegin:      internal.markerStatementBegin,
		statementEnd:        internal.markerStatementEnd,
		parseEndMarker:      internal.markerMigrateDownEnd,
		noTransactionMarker: internal.markerNoTransaction,
	}
}
