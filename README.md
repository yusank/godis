# Godis

> Go  语言实现 `redis` 协议的功能

## TODO

- [ ] tcp 协议
  - [x] `decode/encode` 协议
  - [ ] 网络优化
  - [x] 优雅关闭,退出时等待未处理完成的 connection
- [ ] 五种数据结构
  - [x] string
  - [x] list
  - [x] set
  - [x] zset
  - [x] hash map
- [x] 大部分常用的命令
- [x] 可通过 redis-cli 连接且可用
- [ ] 压测
- [ ] 项目结构优化
- [ ] 完善单元测试 e2e 测试

## Not In Feature

- persistence data to local
- distribution (may be will support, not sure right now)

## Design

![data transfer](./static/godis_data_transfer.png)