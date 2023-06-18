package storage

// keyvalStore is an interface for a key-value store. It is used to abstract
// the underlying storage implementation from the service. This allows us to
// easily swap out the storage implementation without having to change the
// service. It can be used for usage with storage liked Redis for scalability.
type keyvalStore interface {
	Set(key string, value uint) error
	Delete(key string) error
	Get(key string) (uint, error)
	Increment(key string)
}
