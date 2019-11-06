package db

type DatabaseConfig interface {
	GetHost() string
	GetPort() uint
	GetUsername() string
	GetPassword() string
	GetDatabaseName() string
	GetVerbose() bool
}
