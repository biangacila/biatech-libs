package domaindriven

// GenericServiceImpl is the concrete implementation of the GenericService interface
type GenericServiceImpl[T any] struct {
	repo GenericRepository[T]
	dto  any
}

func NewGenericServiceImpl[T any](repo GenericRepository[T]) *GenericServiceImpl[T] {
	return &GenericServiceImpl[T]{repo: repo}
}
func (g *GenericServiceImpl[T]) SetDto(dto any) {
	g.dto = dto
}
func (g *GenericServiceImpl[T]) Save(dbName, entity string, record any, t T) error {
	return g.repo.Save(dbName, entity, record, t)
}
func (g *GenericServiceImpl[T]) SaveBulk(dbName, entity string, record []T, t T) error {
	return g.repo.SaveBulk(dbName, entity, record, t)
}
func (g *GenericServiceImpl[T]) Find(dbName, entity string, fieldValues map[string]interface{}, t T) (T, error) {
	return g.repo.Find(dbName, entity, fieldValues, t)
}
func (g *GenericServiceImpl[T]) Get(dbName, entity string, fieldValues map[string]interface{}, t T) ([]T, error) {
	return g.repo.Get(dbName, entity, fieldValues, t)
}
func (g *GenericServiceImpl[T]) Update(dbName, entity string, conditions, fieldValues map[string]interface{}, t T) error {
	return g.repo.Update(dbName, entity, conditions, fieldValues, t)
}
func (g *GenericServiceImpl[T]) Delete(dbName, entity string, fieldValues map[string]interface{}, t T) error {
	return g.repo.Delete(dbName, entity, fieldValues, t)
}
