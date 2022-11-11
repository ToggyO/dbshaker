package internal

import (
	"bufio"
	"strings"
)

// import (
//
//	"bufio"
//	"bytes"
//	"fmt"
//	"io"
//	"strings"
//
// )
type parsingState int

// const (
//
//	start1 parsingState = iota
//	//migrateUp
//	statementBeginUp
//	statementEndUp
//	//migrateDown
//	statementBeginDown
//	statementEndDown
//
// )
//
// type parseState parsingState
//
//	func (s *parseState) Get() parsingState {
//		return parsingState(*s)
//	}
//
//	func (s *parseState) Set(val parsingState) {
//		*s = parseState(val)
//	}
//
// //const scanBufferSize = 4 * 1024 * 1024
// //
// //var bufferPool = sync.Pool{
// //	New: func() interface{} {
// //		return make([]byte, scanBufferSize)
// //	},
// //}
// //
// //var emptyLineRegex = regexp.MustCompile(`^\s*$`)
//
// // TODO: add comment
//
//	func ParseSQLMigration1(r io.Reader, direction bool) (statements []string, useTx bool, err error) {
//		var statementBuffer bytes.Buffer
//
//		scanBuffer := bufferPool.Get().([]byte)
//		defer bufferPool.Put(scanBuffer)
//
//		scanner := bufio.NewScanner(r)
//		scanner.Buffer(scanBuffer, scanBufferSize)
//
//		state := parseState(start)
//
//		useTx = true
//		firstLineFound := false
//		for scanner.Scan() {
//			line := scanner.Text()
//
//			if strings.HasPrefix(line, SQLCommentPrefix) {
//				marker := strings.TrimSpace(strings.TrimPrefix(line, SQLCommentPrefix))
//
//				switch marker {
//				case statementMarkerUp:
//					firstLineFound = true
//
//					if !direction {
//						continue
//					}
//
//					switch state.Get() {
//					case start:
//						state.Set(migrateUp)
//					default:
//						return nil, false, fmt.Errorf("duplicate statement `%s %s`, state=%v", SQLCommentPrefix, statementMarkerUp, state)
//					}
//					continue
//
//				case statementMarkerDown:
//					firstLineFound = true
//
//					if direction {
//						continue
//					}
//
//					switch state.Get() {
//					case migrateUp, statementBeginUp, statementEndUp:
//						state.Set(migrateDown)
//					default:
//						return nil, false, fmt.Errorf("sql migration file must start from `%s %s`, state=%v", SQLCommentPrefix, statementMarkerUp, state)
//					}
//
//				// TODO: add statements begin/end
//
//				case statementMarkerNoTransaction:
//					firstLineFound = true
//					useTx = false
//					continue
//
//				default:
//					// Ignore line
//					continue
//				}
//
//			}
//
//			if !firstLineFound && emptyLineRegex.MatchString(line) {
//				continue
//			}
//
//			if _, err = statementBuffer.WriteString(line + "\n"); err != nil {
//				return nil, false, err
//			}
//
//			//switch state.Get() {
//			//case migrateUp, statementBeginUp, statementEndUp:
//			//
//			//}
//
//			if checkOnStatementEnds(line) {
//				statements = append(statements, statementBuffer.String())
//				statementBuffer.Reset()
//				continue
//			}
//		}
//
//		return
//	}
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
