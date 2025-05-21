package domaindriven

import (
	"errors"
	"fmt"
	"github.com/biangacila/biatech-libs/utils"
	"github.com/gocql/gocql"
)

type CassandraGenericRepository[T any] struct {
	session *gocql.Session
	dbName  string
}

func (c *CassandraGenericRepository[T]) Find(entity string, fieldValues map[string]interface{}, t T) (T, error) {
	records, err := c.Get(entity, fieldValues, t)
	if err != nil {
		return t, err
	}
	if len(records) == 0 {
		return t, errors.New("no records found")
	}
	return records[(len(records) - 1)], nil
}

func (c *CassandraGenericRepository[T]) Get(entity string, fieldValues map[string]interface{}, t T) ([]T, error) {
	return utils.FetchRecordWithConditions(c.session, c.dbName, entity, fieldValues, t, " ALLOW FILTERING ")
}

func (c *CassandraGenericRepository[T]) Update(entity string, conditions, fieldValues map[string]interface{}, t T) error {
	valuesWhere, err := utils.WhereClauseBuilder(conditions)
	if err != nil {
		return err
	}
	valuesToUpdate, err := utils.UpdateClauseBuilder(fieldValues)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s.%s SET %s  %s", c.dbName, entity, valuesToUpdate, valuesWhere)

	err = c.session.Query(query).Exec()
	if err != nil {
		return errors.New(fmt.Sprintf("error updating record: %v  -> %v", err.Error(), query))
	}
	return c.session.Query(query).Exec()
}

func (c *CassandraGenericRepository[T]) Delete(entity string, fieldValues map[string]interface{}, t T) error {
	valuesWhere, err := utils.WhereClauseBuilder(fieldValues)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s.%s  %s", c.dbName, entity, valuesWhere)
	return c.session.Query(query).Exec()
}

func NewCassandraGenericRepository[T any](session *gocql.Session, dbName string, t T) *CassandraGenericRepository[T] {
	return &CassandraGenericRepository[T]{
		session: session,
		dbName:  dbName,
	}
}
func (c *CassandraGenericRepository[T]) Save(entity string, record any, t T) error {
	return utils.InsertRecord(c.session, c.dbName, entity, record)
}
func (c *CassandraGenericRepository[T]) SaveBulk(entity string, records []T, t T) error {
	return utils.InsertBulkRecord(c.session, c.dbName, entity, records)
}
