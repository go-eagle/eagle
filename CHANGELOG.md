## Changelog

### v1.0.4
- 新增 cache 可以通过配置文件进行设定，支持内存、redis缓存，默认使用memory
- 更新 gin框架升级到1.6.3
- 修复 日志行数显示问题

### v1.0.3
- 新增 user_repo 增加单元测试
- 新增 支持单独、docker(docker-compose方式)、Supervisord 三种部署方式
- 优化 使用zap日志库替换原有日志库
- 移除 BaseModel，ID由各model单独管理

### v1.0.2
- 新增 健康检查接口，便于对应用进行探活检测
- 新增 metrics接口，便于prometheus进行监控
- 新增 scripts目录，将各种脚本文件统一放到此目录
- 新增 schedule目录，可以配置需要跑的计划任务
- 新增 支持发送邮件模块
- 优化 main.go
- 更新 在repo README 中增加使用建议

### v1.0.1
- 优化 service的接口命名及导出命名
- 更新 README 新增开发约定

### v1.0.0
- 新增 第一个可用版本发布