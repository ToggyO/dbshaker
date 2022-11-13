package internal

type sqlParseMarkers struct {
	parseStartMarker    string
	statementBegin      string
	statementEnd        string
	parseEndMarker      string
	noTransactionMarker string
}

func newMigrateUpParseMarkers() sqlParseMarkers {
	return sqlParseMarkers{
		parseStartMarker:    markerMigrateUpStart,
		statementBegin:      markerStatementBegin,
		statementEnd:        markerStatementEnd,
		parseEndMarker:      markerMigrateUpEnd,
		noTransactionMarker: markerNoTransaction,
	}
}

func newMigrateDownParseMarkers() sqlParseMarkers {
	return sqlParseMarkers{
		parseStartMarker:    markerMigrateDownStart,
		statementBegin:      markerStatementBegin,
		statementEnd:        markerStatementEnd,
		parseEndMarker:      markerMigrateDownEnd,
		noTransactionMarker: markerNoTransaction,
	}
}
