package store


type Stabel[K comparable, V any] interface {
   Get(key K) (V, bool)
   Set(key K, value V)
   Delete(key K)
   Count() int64
   All(f func(key K, value V) bool)
}