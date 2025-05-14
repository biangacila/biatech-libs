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
func (g *GenericServiceImpl[T]) Save(entity string, record any, t T) error {
	return g.repo.Save(entity, record, t)
}
func (g *GenericServiceImpl[T]) SaveBulk(entity string, record []T, t T) error {
	return g.repo.SaveBulk(entity, record, t)
}
func (g *GenericServiceImpl[T]) Find(entity string, fieldValues map[string]interface{}, t T) (T, error) {
	return g.repo.Find(entity, fieldValues, t)
}
func (g *GenericServiceImpl[T]) Get(entity string, fieldValues map[string]interface{}, t T) ([]T, error) {
	return g.repo.Get(entity, fieldValues, t)
}
func (g *GenericServiceImpl[T]) Update(entity string, conditions, fieldValues map[string]interface{}, t T) error {
	return g.repo.Update(entity, conditions, fieldValues, t)
}
func (g *GenericServiceImpl[T]) Delete(entity string, fieldValues map[string]interface{}, t T) error {
	return g.repo.Delete(entity, fieldValues, t)
}
