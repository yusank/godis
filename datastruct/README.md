# datastruct

This module implement all five data structrer of redis supported.

- string
- list
- set
- sorted set
- hash map

**尽可能减少该包的逻辑,不要考虑外部如何写 response,只考虑返回结果和遇到的 error, 外部处理 ErrNil 等错误, 除非像 mget 这种命令需要在内部把 ErrNil 消耗的操作.**