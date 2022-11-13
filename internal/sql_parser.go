package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

type parsingState int

const (
	startParse parsingState = iota
	onParseTarget
	statementBegin
	statementEnd
	// TODO: check
	endParse
)

type parseState parsingState

func (s *parseState) Get() parsingState {
	return parsingState(*s)
}

func (s *parseState) Set(val parsingState) {
	*s = parseState(val)
}

const scanBufferSize = 4 * 1024 * 1024

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, scanBufferSize)
	},
}

var emptyLineRegex = regexp.MustCompile(`^\s*$`)

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
				firstLineFound = true
				switch state.Get() {
				case startParse:
					state.Set(onParseTarget)
				default:
					return nil, false, fmt.
						Errorf("duplicate statement `-- %s`, state=%v", markers.parseStartMarker, state)
				}
				continue

			case markers.parseEndMarker:
				firstLineFound = true
				switch state.Get() {
				case onParseTarget, statementEnd:
					// TODO: check
					//state.Set(endParse)
					if bufferRemaining := strings.TrimSpace(statementBuffer.String()); len(bufferRemaining) > 0 {
						return nil, false, errUnfinishedSQLQuery(int(state), direction, bufferRemaining)
					}
					return
				case statementBegin:
					return nil, false, errMissingSQLParsingAnnotation(markers.statementEnd)
				default:
					return nil, false, fmt.
						Errorf("sql migration file must start from `-- %s`, state=%v", markers.parseStartMarker, state)
				}

			case markers.statementBegin:
				firstLineFound = true
				switch state.Get() {
				case startParse:
					break
				case onParseTarget, statementEnd:
					state.Set(statementBegin)
				default:
					return nil, false, fmt.
						Errorf("`-- %s` must be defined after `-- %s` or `-- %s` annotation,"+
							" state=%v", markerStatementBegin, markerMigrateUpStart, markerMigrateDownStart, state)
				}
				continue

			case markers.statementEnd:
				firstLineFound = true
				switch state.Get() {
				case startParse:
					// TODO: check
					continue
				case statementBegin:
					state.Set(statementEnd)
				default:
					return nil, false, fmt.
						Errorf("`-- %s` must be defined after `-- %s`", markerStatementEnd, markerStatementBegin)
				}

			case markers.noTransactionMarker:
				useTx = false
				continue

			default:
				// Ignore line
				continue
			}
		}

		if !firstLineFound && emptyLineRegex.MatchString(line) {
			continue
		}

		switch state.Get() {
		case onParseTarget, statementBegin:
			if _, err = statementBuffer.WriteString(line + "\n"); err != nil {
				return nil, false, err
			}
		}

		switch state.Get() {
		case onParseTarget:
			if checkOnStatementEnds(line) {
				statements = append(statements, statementBuffer.String())
				statementBuffer.Reset()
			}
		case statementEnd:
			statements = append(statements, statementBuffer.String())
			statementBuffer.Reset()
			state.Set(onParseTarget)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, false, fmt.Errorf("failed to scan migration: %w", err)
	}

	switch state.Get() {
	case startParse:
		return nil, false, errMissingSQLParsingAnnotation(markers.parseStartMarker)
	case onParseTarget, statementBegin, statementEnd:
		return nil, false, errMissingSQLParsingAnnotation(markers.parseEndMarker)
	}

	// TODO: check
	//if bufferRemaining := strings.TrimSpace(statementBuffer.String()); len(bufferRemaining) > 0 {
	//	return nil, false, fmt.Errorf("failed to parse migration: state %q, direction: %v: unexpected unfinished SQL query: %q: missing semicolon", state, direction, bufferRemaining)
	//}

	return
}

func checkOnStatementEnds(line string) bool {
	scannerBuffer := bufferPool.Get().([]byte)
	defer bufferPool.Put(scannerBuffer)

	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Buffer(scannerBuffer, scanBufferSize)
	scanner.Split(bufio.ScanWords)

	prevWord := ""
	for scanner.Scan() {
		currentWord := scanner.Text()
		if strings.HasSuffix(currentWord, SQLCommentPrefix) {
			break
		}
		prevWord = currentWord
	}

	return strings.HasSuffix(prevWord, SQLSemicolon)
}
