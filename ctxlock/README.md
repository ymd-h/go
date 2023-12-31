# ctxlock: Context-aware Lock

`ctxlock` provides context-aware locks (`Lock` / `SharableLock`).
User can safely cancel a wating to acquire the lock through `context.Context`.

Unlike standard `sync.Mutex`, there is no unlock method.
An unlock function (`ctxlock.UnlockFunc`) is returned from its lock methods instead.

The unlock function works only once and it is safe to call it multiple times.

The unlock function is not connected to the `context.Context`.
If you want to call it when a context is canceled,
`UnlockFunc.UnlockOnCancel(ctx)` method can be used.

`ctxlock` is based on goroutine and channel,
so that it might have performance overhead than that of `sync.Mutex`.


## (Pseudo) Example

### Ordinary Exclusive Lock
```Go
L := ctxlock.NewLock()

// Lock() method tries to lock and returns unlock function.
unlock, err := L.Lock(context.Background())
```

### Reader-Writer Lock
```Go
L := ctxlock.NewSharableLock()

// Reader Lock
unlock, err := L.SharedLock(context.Background())

// Writer Lock
unlock, err := L.ExclusiveLock(context.Background())
```
