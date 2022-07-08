package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
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

// EqualStruct 两个结构体对比, s1, s2 必须是指针
func EqualStruct(t *testing.T, s1 any, s2 any, keys ...string) {
	v1 := reflect.ValueOf(s1)
	v2 := reflect.ValueOf(s2)

	require.Equal(t, v1.Kind(), reflect.Struct)
	require.Equal(t, v2.Kind(), reflect.Struct)

	if len(keys) > 0 {
		for _, key := range keys {
			require.Equal(t, fmt.Sprintf("%v", v1.FieldByName(key)), fmt.Sprintf("%v", v2.FieldByName(key)))
		}
	} else {
		for i := 0; i < v1.NumField(); i++ {
			require.Equal(t, fmt.Sprintf("%v", v1.Field(i)), fmt.Sprintf("%v", v2.Field(i)))
		}
	}

}

func TestEqualStruct(t *testing.T) {
	var e1 = Entry{
		ID:     1,
		Amount: 10,
	}
	var e2 = Entry{
		ID:     1,
		Amount: 10,
	}
	EqualStruct(t, e1, e2)
	EqualStruct(t, e1, e2, "Amount")
}
