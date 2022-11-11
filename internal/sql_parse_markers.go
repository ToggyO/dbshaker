package internal

type sqlParseMarkers struct {
	parseStartMarker    string
	parseEndMarker      string
	noTransactionMarker string
}

func newMigrateUpParseMarkers() sqlParseMarkers {
	return sqlParseMarkers{
		parseStartMarker:    markerMigrateUpStart,
		parseEndMarker:      markerMigrateUpEnd,
		noTransactionMarker: markerNoTransaction,
	}
}

func newMigrateDownParseMarkers() sqlParseMarkers {
	return sqlParseMarkers{
		parseStartMarker:    markerMigrateDownStart,
		parseEndMarker:      markerMigrateDownEnd,
		noTransactionMarker: markerNoTransaction,
	}
}
