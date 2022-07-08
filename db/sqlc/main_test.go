package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// 整个包下测试主入口

const (
	dbDriver = "postgres"
	dbDSN    = "postgresql://postgres:abc123@localhost:5432/simple_bank?sslmode=disable"
)

// 创建一个 db queries
var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbDSN)
	if err != nil {
		log.Fatal("cannot connect to postgre: ", err)
	}
	// NOTE: *sql.DB 实现了 DBTX interface 所有方法
	testQueries = New(conn)

	os.Exit(m.Run())
}
