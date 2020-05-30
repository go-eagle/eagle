# Model

Model层，或者叫 `Entity`，实体层，用于存放我们的实体类，与数据库中的属性值基本保持一致。

## 数据库约定

这里默认使用 `MySQL` 数据库，尽量使用 `InnoDB` 作为存储引擎。

### 相关表采用统一前缀

比如和用户相关的，使用 `user_` 作为表前缀：

```bash
user            // 用户表
user_follow     // 用户关注表
user_fans       // 用户粉丝表
user_stat       // 用户统计表
```

### 统一字段名
 
 一个表中需要包含的三大字段：主键(id)，创建时间(created_at)，更新时间(updated_at)  
 如果需要用户id，一般用user_id即可。