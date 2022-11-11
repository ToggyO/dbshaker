package internal

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// TODO: add comment
func ParseSQLMigration(r io.Reader, direction bool) (statements []string, useTx bool, err error) {
	var markers sqlParseMarkers
	if !direction {
		markers = newMigrateDownParseMarkers()
	} else {
		markers = newMigrateUpParseMarkers()
	}

	var statementBuffer bytes.Buffer

	scanBuffer := bufferPool.Get().([]byte)
	defer bufferPool.Put(scanBuffer)

	scanner := bufio.NewScanner(r)
	scanner.Buffer(scanBuffer, scanBufferSize)

	state := parseState(startParse)

	useTx = true
	firstLineFound := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, SQLCommentPrefix) {
			marker := strings.TrimSpace(strings.TrimPrefix(line, SQLCommentPrefix))

			switch marker {
			case markers.parseStartMarker:
				switch state.Get() {
				case startParse:
					state.Set()
				}

			case markers.parseEndMarker:

			// TODO: add statements begin/end

			case markers.noTransactionMarker:
				useTx = false
				continue

			default:
				// Ignore line
				continue
			}
			}
		}
	}

	return
}
