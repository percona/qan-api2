// qan-api2
// Copyright (C) 2019 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mysqlParser(t *testing.T) {
	expectedQuery := "insert into pmm.Person(`Name`, Addr) values (:1, :2)"
	query, placeholders, err := mysqlParser("/* comment */ INSERT INTO pmm.Person(Name, Addr) VALUES('Test', 'Test2')")
	assert.NoError(t, err)
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, uint32(2), placeholders)

	expectedQuery = "select * from cities"
	query, placeholders, err = mysqlParser("/* comment */ SELECT * FROM cities")
	assert.NoError(t, err)
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, uint32(0), placeholders)
}
