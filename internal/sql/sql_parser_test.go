package sql

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParseSQLMigration(t *testing.T) {
	testCases := []struct {
		sql                 string
		upStatementsCount   int
		downStatementsCount int
		withTransaction     bool
	}{
		{
			sql:                 multilineSQL,
			upStatementsCount:   4,
			downStatementsCount: 1,
			withTransaction:     true,
		},
		{
			sql:                 SQLWithStatements,
			upStatementsCount:   2,
			downStatementsCount: 3,
			withTransaction:     true,
		},
		{
			sql:                 multiUpDown,
			upStatementsCount:   1,
			downStatementsCount: 1,
			withTransaction:     true,
		},
		{
			sql:                 noTransactionSQL,
			upStatementsCount:   1,
			downStatementsCount: 1,
			withTransaction:     false,
		},
	}

	for _, testCase := range testCases {
		statements, useTx, err := ParseSQLMigration(strings.NewReader(testCase.sql), true)
		require.NoError(t, err)
		require.Len(t, statements, testCase.upStatementsCount)
		require.Equal(t, testCase.withTransaction, useTx)

		statements, _, err = ParseSQLMigration(strings.NewReader(testCase.sql), false)
		require.NoError(t, err)
		require.Len(t, statements, testCase.downStatementsCount)
		require.Equal(t, testCase.withTransaction, useTx)
	}
}

func TestParsingErrors(t *testing.T) {
	testCases := []string{
		unfinishedSQL,
		unfinishedUpStatementSQL,
		noStatementEndSQL,
	}

	for _, sql := range testCases {
		_, _, err := ParseSQLMigration(strings.NewReader(sql), true)
		require.Error(t, err)
	}
}

func TestSemicolons(t *testing.T) {
	testCases := []struct {
		line   string
		result bool
	}{
		{line: "END;", result: true},
		{line: "END; -- comment", result: true},
		{line: "END   ; -- comment", result: true},
		{line: "END -- comment", result: false},
		{line: "END -- comment ;", result: false},
		{line: "END \" ; \" -- comment", result: false},
	}

	for _, test := range testCases {
		r := checkOnStatementEnds(test.line)
		require.Equal(t, test.result, r)
	}
}

const multilineSQL = `


-- +dbshaker UpStart
CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	age INTEGER NOT NULL,
	name VARCHAR NOT NULL
);     -- 1st stmt

-- comment
SELECT 1;           -- 2nd stmt
SELECT 2; SELECT 2; -- 3rd stmt
SELECT 3;           -- 4th stmt

-- +dbshaker UpEnd

-- +dbshaker DownStart
-- comment
DROP TABLE users;

-- +dbshaker DownEnd
`

const SQLWithStatements = `-- +dbshaker UpStart
-- +dbshaker StatementBegin
	
	SELECT name
		   , salary
	FROM People
	WHERE name in (SELECT DISTINCT name 
				   FROM population 
				   WHERE country = "Canada"
						 AND city = "Toronto")
		  AND salary >= (SELECT AVG(salary)
						 FROM salaries
						 WHERE gender = "Female")

-- +dbshaker StatementEnd

SELECT EMP_ID, NAME FROM EMPLOYEE_TBL WHERE EMP_ID = '0000';

-- +dbshaker UpEnd


-- +dbshaker DownStart

-- +dbshaker StatementBegin
SELECT Name, Age FROM Patients WHERE Age > 40
GROUP BY Name, Age ORDER BY Name;
-- +dbshaker StatementEnd

-- +dbshaker StatementBegin
SELECT Student_ID FROM STUDENT;
-- +dbshaker StatementEnd

-- +dbshaker StatementBegin
SELECT SUM(Salary)FROM Employee WHERE Emp_Age < 30;
-- +dbshaker StatementEnd

-- +dbshaker DownEnd
`

const unfinishedSQL = `
-- +dbshaker UpStart
ALTER TABLE post

-- +dbshaker UpEnd
`

const unfinishedUpStatementSQL = `-- +dbshaker UpStart
-- This is just a comment`

const noStatementEndSQL = `-- +dbshaker UpStart
-- +dbshaker StatementBegin
INSERT INTO users(id, name, age)
VALUES (1, 'Stas', 69);

-- +dbshaker UpEnd`

const multiUpDown = `-- +dbshaker UpStart
CREATE TABLE locations(
	id INT PRIMARY KEY,
	x GEOPOINT NOT NULL,
	y GEOPOINT NOT NULL
);
-- +dbshaker UpEnd


-- +dbshaker DownStart
DROP TABLE locations;
-- +dbshaker DownEnd

-- +dbshaker UpStart
CREATE TABLE cars(
	id INT PRIMARY KEY,
	name VARCHAR NOT NULL,
	PRICE DECIMAL NOT NULL
);
-- +dbshaker UpEnd

`

const noTransactionSQL = `-- +dbshaker NO_TRANSACTION
-- +dbshaker UpStart
CREATE TABLE locations(
	id INT PRIMARY KEY,
	x GEOPOINT NOT NULL,
	y GEOPOINT NOT NULL
);
-- +dbshaker UpEnd


-- +dbshaker DownStart
DROP TABLE locations;
-- +dbshaker DownEnd
`
