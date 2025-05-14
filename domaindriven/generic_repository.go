package domaindriven

type GenericRepository[T any] interface {
	Save(entity string, record any, t T) error
	SaveBulk(entity string, record []T, t T) error
	Find(entity string, fieldValues map[string]interface{}, t T) (T, error)
	Get(entity string, fieldValues map[string]interface{}, t T) ([]T, error)
	Update(entity string, conditions, fieldValues map[string]interface{}, t T) error
	Delete(entity string, fieldValues map[string]interface{}, t T) error
}
