package user

import (
	"database/sql"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/1024casts/snake/model"
)

type Suite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	repository Repo
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	gdb, err := gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.db = gdb
	s.db.LogMode(true)

	s.repository = NewUserRepo()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

//nolint: golint
func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

//func (s *Suite) Test_repository_Create() {
//	var (
//		id       = 20
//		username = "test-name"
//	)
//
//	s.mock.ExpectQuery(regexp.QuoteMeta(
//		`INSERT INTO "users" ("id","username")
//       VALUES ($1,$2) RETURNING "users"."id"`)).
//		WithArgs(id, username).
//		WillReturnRows(
//			sqlmock.NewRows([]string{"id"}).AddRow(id))
//
//	user := model.UserModel{
//		ID:       uint64(id),
//		Username: username,
//	}
//	_, err := s.repository.Create(s.db, user)
//
//	require.NoError(s.T(), err)
//}

func (s *Suite) Test_repository_GetUserByID() {
	var (
		id       uint64 = 2
		username        = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE (id = $1)`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(id, username))

	res, err := s.repository.GetUserByID(s.db, id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&model.UserModel{ID: id, Username: username}, res))
}

func (s *Suite) Test_repository_GetUserByPhone() {
	var (
		id       uint64 = 2
		username        = "test-phone"
		phone           = 13011112222
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE (phone = $1)`)).
		WithArgs(phone).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "phone"}).AddRow(id, username, phone))

	res, err := s.repository.GetUserByPhone(s.db, phone)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&model.UserModel{ID: id, Username: username, Phone: phone}, res))
}

func (s *Suite) Test_repository_GetUserByEmail() {
	var (
		id       uint64 = 2
		username        = "test-email"
		email           = "test-email@test.com"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE (email = $1)`)).
		WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(id, username, email))

	res, err := s.repository.GetUserByEmail(s.db, email)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&model.UserModel{ID: id, Username: username, Email: email}, res))
}
