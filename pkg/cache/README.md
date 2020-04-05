这里主要是关于缓存的一些库，主要基于内存型和NoSQL的

内存型的有：memory和big_cache
NoSQL的主要有: redis

各类库只要实现了cache定义的接口(driver)即可。
> 这里的接口driver命名参考了Go官方mysql接口的命名规范

## Reference
- bigcache: https://github.com/allegro/bigcache
- freecache: https://github.com/coocood/freecache
- concurrent_map: https://github.com/easierway/concurrent_map
- gocache: https://github.com/eko/gocache (Built-in stores, eg: bigcache,memcache,redis)