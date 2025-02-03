package mapper

type Mapper[T, K any] interface {
	MapFromEntity(K) (T, error)
	MapToEntity(T) (K, error)
}
