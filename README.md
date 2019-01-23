# session-store

[![Build Status](https://img.shields.io/travis/vicanso/sessionstore.svg?label=linux+build)](https://travis-ci.org/vicanso/sessionstore)


session store for cod, it supports redis and memory.

## redis store

```go
client := NewRedisClient(&redis.Options{
  Addr: "localhost:6379",
})
rs = NewRedisStore(client)
rs.SetOptions(&Options{
  TTL: 3600 * time.Second,
})
```

## memory store

```go
ms := NewMemoryStore(1024)
ms.SetOptions(&Options{
  TTL: 3600 * time.Second,
})
```

