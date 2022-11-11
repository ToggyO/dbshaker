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

const (
	startParse parsingState = iota
	migrateUp
	statementBegin
	statementEnd
	migrateDown
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
func ParseSQLMigration2(r io.Reader, direction bool) (statements []string, useTx bool, err error) {
	target := markerMigrateUp
	if !direction {
		target = markerMigrateDown
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
			case markerMigrateUp:
				firstLineFound = true

				if target != markerMigrateUp {
					continue
				}

				switch state.Get() {
				case startParse:
					state.Set(migrateUp)
				case migrateDown:
					return nil, false, fmt.Errorf("sql migration file must start from `%s %s`, state=%v", SQLCommentPrefix, markerMigrateUp, state)
				default:
					return nil, false, fmt.Errorf("duplicate statement `%s %s`, state=%v", SQLCommentPrefix, markerMigrateUp, state)
				}
				continue

			case markerMigrateDown:
				firstLineFound = true

				if target != markerMigrateDown {
					return
				}

				switch state.Get() {
				case startParse, migrateUp:
					state.Set(migrateDown)
				case migrateDown:
					return nil, false, fmt.Errorf("duplicate statement `%s %s`, state=%v", SQLCommentPrefix, markerMigrateDown, state)
				default:
					return nil, false, fmt.Errorf("sql migration file must start from `%s %s`, state=%v", SQLCommentPrefix, markerMigrateUp, state)
				}
				continue

			// TODO: add statements begin/end

			case markerNoTransaction:
				useTx = false
				continue

			default:
				// Ignore line
				continue
			}
		}

		emptyLineRegex.FindAllStringIndex(line, -1)

		if !firstLineFound && emptyLineRegex.MatchString(line) {
			continue
		}

		if _, err = statementBuffer.WriteString(line + "\n"); err != nil {
			return nil, false, fmt.Errorf("failed to scan migration: %w", err)
		}

		// TODO: проверить окончание с помощью stametentBegin/End
		if checkOnStatementEnds(line) {
			statements = append(statements, statementBuffer.String())
			statementBuffer.Reset()
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, false, fmt.Errorf("failed to scan migration: %w", err)
	}

	if state.Get() == startParse {
		return nil, false, fmt.Errorf("failed to parse migration: must start with `%s %s`", SQLCommentPrefix, markerMigrateUp)
	}

	if bufferRemaining := strings.TrimSpace(statementBuffer.String()); len(bufferRemaining) > 0 {
		return nil, false, fmt.Errorf("failed to parse migration: state %q, direction: %v: unexpected unfinished SQL query: %q: missing semicolon", state, direction, bufferRemaining)
	}

	return
}
