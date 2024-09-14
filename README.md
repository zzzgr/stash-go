## 通用存包服务

### 1. 简介

需要配合抓包工具转发功能使用

### 2. 使用方法

```shell
./stash-go --http_port 8080 --redis_password xxx
```

**支持参数**

| 参数             | 说明       | 默认值            |
|----------------|----------|----------------|
| http_port      | http端口   | 8080           |
| redis_addr     | redis地址  | 127.0.0.1:6379 |
| redis_username | redis用户名 | 空              |
| redis_password | redis密码  | 空              |
| redis_db       | redis数据库 | 0              |

接口

| 接口     | 请求方法 | 说明  |
|--------|------|-----|
| /get   | GET  | 获取包 |
| /clear | GET  | 清空  |
| /*     | ALL  | 存包  |