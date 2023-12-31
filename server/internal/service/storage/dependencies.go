//go:generate mockgen -package storage -source dependencies.go -destination mocks.go

package storage

// keyValStore is an interface for a key-value store. It is used to abstract
// the underlying storage implementation from the service. This allows us to
// easily swap out the storage implementation without having to change the
// service. It can be used for usage with storage liked Redis for scalability.
type keyValStore interface {
	Set(key string, value uint)
	Delete(key string)
	Get(key string) (uint, error)
	Increment(key string)
}
