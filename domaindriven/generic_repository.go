package domaindriven

type GenericRepository[T any] interface {
	Save(dbName, entity string, record any, t T) error
	SaveBulk(dbName, entity string, record []T, t T) error
	Find(dbName, entity string, fieldValues map[string]interface{}, t T) (T, error)
	Get(dbName, entity string, fieldValues map[string]interface{}, t T) ([]T, error)
	Update(dbName, entity string, conditions, fieldValues map[string]interface{}, t T) error
	Delete(dbName, entity string, fieldValues map[string]interface{}, t T) error
}
