package database

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadPlayersExcludesCharactersScheduledForDeletion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := regexp.QuoteMeta(`SELECT name, level, sex, vocation, looktype, lookhead, lookbody, looklegs,
			lookfeet, lookaddons, lastlogin from players where account_id = "42" AND deletion = 0`)
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{
		"name", "level", "sex", "vocation", "looktype", "lookhead", "lookbody", "looklegs",
		"lookfeet", "lookaddons", "lastlogin",
	}))

	players, err := LoadPlayers(db, &Account{ID: 42})

	assert.NoError(t, err)
	assert.Empty(t, players)
	assert.NoError(t, mock.ExpectationsWereMet())
}
