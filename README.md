## intranet dns
基于gin框架开发的内网dns管理系统后端demo, 支持动态dns修改(RFC 2136标准). 提供dns解析的基础服务为bind. 支持多节点部署.

启动方式: `go run cmd/main.go`, [本地swagger查看api详情](http://localhost:16789/swagger/index.html)

项目项目结构如下:
```bash
tree -d intranet-dns/
intranet-dns/
├── apis                # 控制器
├── cmd                 # main.go
├── config              # 配置文件
├── ctx                 # gin上下文
├── database            # 数据库初始化
├── docs                # swagger
├── lib                 # 依赖库
│   ├── ansible
│   └── redis           # redis: cache & lock
├── middleware          # 中间件: 日志, 审计, 鉴权, 跨域等
├── models              # 数据库表结构, dao函数
├── router              # 路由控制
├── service             # 服务代码
└── utils               # 工具函数
```

表结构:
```bash
mysql> show tables;
+-----------------------+
| Tables_in_dns_service |
+-----------------------+
| apis                  |
| dns_records           |
| sys_roles             |
| sys_users             |
+-----------------------+
```

功能点:
- api管理: 路由自动录入数据库, 支持禁用单个api, 支持api限速, 接入go-swagger注解
- 用户管理: 用户可绑定多个角色, 角色可绑定多个api接口, 实现rbac权限模型管控. 单个用户可禁用, 统计登录次数以及登录时间
- 日志审计: os.stdout支持输出api访问日志, 包含来源ip、请求耗时等. 数据库存储具体的body日志用于审计 
- 定时任务管理: 支持动态的增删改查定时任务, 定时任务可控制是否启动, 可查看最近的运行结果
- dns管理: bind9+go miekg/dns实现dns动态增删改查(注意二级域和三级域的区分, 域名保持{name}.{zone格式}, 支持A AAAA CNAME记录). dns探测等

## 相关组件
**golang**: gin, gorm, go-swagger, go-jwt, go-redis, miekg/dns

**database**: mysql, redis

**nameserver**: [dns和bind](https://www.junmajinlong.com/linux/dns_bind/index.html), [bind key配置](https://www.cnblogs.com/RichardLuo/p/DNS_P3.html)
```bash

```
