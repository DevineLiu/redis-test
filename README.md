# Guide

测试集群模式（Test Redis Cluster）
```
redis-test --arch cluster 127.0.0.1:6379  127.0.0.1:6380 127.0.0.1:6381
```

测试哨兵模式(TEST Redis Sentinel)
```
redis-test --arch sentinel 127.0.0.1:26379  127.0.0.1:26380 127.0.0.1:26381
```

循环测试 (Loop Test)
```
redis-test --arch cluster  -l true 127.0.0.1:6379  127.0.0.1:6380 127.0.0.1:6381
```

