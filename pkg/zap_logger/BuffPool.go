package zapLogger

import "go.uber.org/zap/buffer"

var (
	_pool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	GetBuff = _pool.Get
)
