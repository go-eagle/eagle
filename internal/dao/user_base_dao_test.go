package dao

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/go-eagle/eagle/internal/model"
)

type Suite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	repository Dao
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	gdb, err := gorm.Open("mysql", db)
	require.NoError(s.T(), err)

	s.db = gdb
	s.db.LogMode(true)

	//s.repository = New(s.db)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

//nolint: golint
func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_repository_Create() {
	user := model.UserBaseModel{
		Username:  "test-name",
		Password:  "123456",
		Phone:     123455678,
		Email:     "test@test.com",
		Avatar:    "/statics/avatar/1.jpg",
		Sex:       1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	const sqlInsert = `INSERT INTO "user_base" ("username","password","phone","email","avatar","sex","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "user_base"."id"`
	const newID = 1

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
		WithArgs(user.Username, user.Password, user.Phone, user.Email, user.Avatar, user.Sex, user.CreatedAt, user.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newID))
	s.mock.ExpectCommit()

	_, err := s.repository.CreateUser(context.TODO(), &user)

	require.NoError(s.T(), err)
}

func (s *Suite) Test_repository_GetUserByID() {
	var (
		id       uint64 = 2
		username        = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_base" WHERE (id = $1)`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(id, username))

	res, err := s.repository.GetOneUser(context.TODO(), id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&model.UserBaseModel{ID: id, Username: username}, res))
}

func (s *Suite) Test_repository_GetUserByPhone() {
	var (
		id       uint64 = 2
		username        = "test-phone"
		phone    int64  = 13011112222
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_base" WHERE (phone = $1)`)).
		WithArgs(phone).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "phone"}).AddRow(id, username, phone))

	res, err := s.repository.GetUserByPhone(context.TODO(), phone)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&model.UserBaseModel{ID: id, Username: username, Phone: phone}, res))
}

func (s *Suite) Test_repository_GetUserByEmail() {
	var (
		id       uint64 = 2
		username        = "test-email"
		email           = "test-email@test.com"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_base" WHERE (email = $1)`)).
		WithArgs(email).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(id, username, email))

	res, err := s.repository.GetUserByEmail(context.TODO(), email)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&model.UserBaseModel{ID: id, Username: username, Email: email}, res))
}
