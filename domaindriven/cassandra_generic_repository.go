package domaindriven

import (
	"errors"
	"fmt"
	"github.com/biangacila/biatech-libs/utils"
	"github.com/gocql/gocql"
)

type CassandraGenericRepository[T any] struct {
	session *gocql.Session
}

func (c *CassandraGenericRepository[T]) Find(dbName, entity string, fieldValues map[string]interface{}, t T) (T, error) {
	records, err := c.Get(dbName, entity, fieldValues, t)
	if err != nil {
		return t, err
	}
	if len(records) == 0 {
		return t, errors.New("no records found")
	}
	return records[(len(records) - 1)], nil
}

func (c *CassandraGenericRepository[T]) Get(DbName, entity string, fieldValues map[string]interface{}, t T) ([]T, error) {
	return utils.FetchRecordWithConditions(c.session, DbName, entity, fieldValues, t, " ALLOW FILTERING ")
}

func (c *CassandraGenericRepository[T]) Update(DbName, entity string, conditions, fieldValues map[string]interface{}, t T) error {
	valuesWhere, err := utils.WhereClauseBuilder(conditions)
	if err != nil {
		return err
	}
	valuesToUpdate, err := utils.UpdateClauseBuilder(fieldValues)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s.%s SET %s  %s", DbName, entity, valuesToUpdate, valuesWhere)

	return c.session.Query(query).Exec()
}

func (c *CassandraGenericRepository[T]) Delete(DbName, entity string, fieldValues map[string]interface{}, t T) error {
	valuesWhere, err := utils.WhereClauseBuilder(fieldValues)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s.%s  %s", DbName, entity, valuesWhere)
	return c.session.Query(query).Exec()
}

func NewCassandraGenericRepository[T any](session *gocql.Session, t T) *CassandraGenericRepository[T] {
	return &CassandraGenericRepository[T]{
		session: session,
	}
}
func (c *CassandraGenericRepository[T]) Save(DbName, entity string, record any, t T) error {
	return utils.InsertRecord(c.session, DbName, entity, record)
}
func (c *CassandraGenericRepository[T]) SaveBulk(DbName, entity string, records []T, t T) error {
	return utils.InsertBulkRecord(c.session, DbName, entity, records)
}
