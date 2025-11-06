```go
import (
	"github.com/shardul-d/kv-store"
)

kvstore, _ := kvstore.Init(kvstore.WithDir("data/"))

// Set a key.
kvstore.Put("hello", []byte("world"))

// Fetch the key.
v, _ := kvstore.Get("hello")

// Delete a key.
kvstore.Delete("hello")

// Set with expiry.
kvstore.PutEx("hello", []byte("world"), time.Second * 5)
```

For a complete example, visit [examples](./examples/main.go).
