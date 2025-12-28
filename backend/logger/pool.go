package logger

import (
	"strings"
	"sync"
)

var (
	// bufferPool for log message buffers
	bufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, 8192) // 8KB pre-allocation
			return &buf
		},
	}

	// entryPool for LogEntry structs
	entryPool = sync.Pool{
		New: func() interface{} {
			return &LogEntry{}
		},
	}

	// stackBufPool for stack trace collection
	stackBufPool = sync.Pool{
		New: func() interface{} {
			buf := make([]uintptr, 128)
			return &buf
		},
	}

	// stringBuilderPool for string building operations
	stringBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	// timeBufferPool for timestamp formatting
	timeBufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, 64)
			return &buf
		},
	}
)
