## 架构说明
使用nginx作为前端代理, intranet-dns-web为前端项目, intranet-dns-backend为后端项目, named为权威dns服务器. intranet-dns-backend遵循rfc2136等规范向named权威dns服务器发起动态dns请求, 二者之间定期进行dns记录的校验同步. <br>
![服务整体架构](./introduce-assets/framework.png)



## 服务介绍

### intranet-dns-web
基于Vue3、Vite、Vue-Router4、Pina、ant-design-vue3开发的内网dns管理系统前端界面.


系统管理演示: <br>
![系统管理演示](./introduce-assets/system-manager.gif)


<br>

### intranet-dns-backend
基于gin框架开发的内网dns管理系统后端demo, 支持动态dns修改(RFC 2136标准)、dns探测等功能. 提供dns解析的基础服务为bind. 支持多节点部署.

后端服务启动后, [本地swagger查看api详情](http://localhost:16789/swagger/index.html)
<br><img src="./intranet-dns-backend/docs/swagger.png" alt="swagger示例">

项目项目结构如下:
```bash
intranet-dns-backend/
├── apis            # 控制器
├── cmd             # main.go
├── config          # 配置文件
├── ctx             # gin上下文
├── database        # 数据库初始化
├── docs            # swagger
├── middleware      # gin中间件
├── models          # 数据库表结构, dao函数
├── router          # 路由控制
├── service         # 服务代码
│   ├── cronjob     # 定时任务
│   ├── dnslib      # dns服务
│   └── redis       # redis服务    
└── utils           # 工具函数
```

表结构:
```bash
mysql> show tables;
+-----------------------+
| Tables_in_dns_service |
+-----------------------+
| apis                  | # 路由
| audit_logs            | # 审计日志
| cronjobs              | # 定时任务
| dns_probes            | # dns探测
| dns_records           | # dns rr
| dns_zones             | # dns zone
| sys_roles             | # 系统角色
| sys_users             | # 系统用户
+-----------------------+
```

功能点:
- api管理: 路由自动录入数据库, 支持禁用单个api, 支持api限速, 接入go-swagger注解
- 用户管理: 用户可绑定多个角色, 角色可绑定多个api接口, 实现rbac权限模型管控. 单个用户可禁用, 统计登录次数以及登录时间. 利用go-jwt token来实现登录鉴权等.
- 日志审计: os.stdout支持输出api访问日志, 包含来源ip、请求耗时等. 数据库存储具体的body日志用于审计 
- 定时任务管理: 支持动态的增删改查定时任务, 定时任务可控制是否启动, 可查看最近的运行结果
- dns管理: bind9+go miekg/dns实现dns动态增删改查, dns探测等. [dynamic dns rfc2136](https://datatracker.ietf.org/doc/html/rfc2136)

#### 相关组件
**golang**: gin, gorm, go-swagger, go-jwt, go-redis, miekg/dns

**database**: mysql, redis

**nameserver**: [bind dynamic update](https://bind9.readthedocs.io/en/v9.16.18/advanced.html), [bind key配置](https://www.cnblogs.com/RichardLuo/p/DNS_P3.html)
```bash
tsig-keygen funny.com > /etc/funny.com.key

cat /etc/named.conf |grep -A 5 'funny.com'
include "/etc/funny.com.key";
zone "funny.com." {
        type master;
        file "funny.com.zone";
        allow-update { key funny.com; };
};
```

#### 使用举例
如增加dns记录:
```bash
# dig不存在的记录, NXDOMAIN
dig a.b.c.funny.com @your_dns_server
"->>HEADER<<- opcode: QUERY, status: NXDOMAIN, id: 24423
flags: qr aa rd ra; QUERY: 1, ANSWER: 0, AUTHORITY: 1, ADDITIONAL: 1"

# call api to add dns record
curl -X 'POST' \
  'http://localhost:8008/api/v1/dns/records' \
  -H 'accept: application/json' \
  -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzM2NTExMDc4LCJpc3MiOiJpbnRyYW5ldC1kbnMifQ.Gqubkw1gNICcH-X90FDVWMFHClJhhyn0ljUafRh7xZk' \
  -H 'Content-Type: application/json' \
  -d '{
  "record_content": "110.110.110.110",
  "record_name": "a.b.c.funny.com",
  "record_ttl": 120,
  "record_type": "A",
  "zone": "funny.com"
}'

# 内网dns记录已经动态添加成功
dig a.b.c.funny.com @your_dns_server +short
110.110.110.110

# 查看named zone文件
cat funny.com.zone|grep a.b.c
a.b.c			A	110.110.110.110
```

## 使用说明
1.准备工作: clone仓库, 安装redis和mysql, 安装named服务并配置TSIG key.
```bash
git clone git@github.com:tswcbyy1107/intranet-dns.git
```

2.启动后端: 修改config.ini配置, 填入自己的mysql和redis地址密码.
```bash
cd intranet-dns-backend
vim config/config.ini 
go run cmd/main.go
```

3.启动前端: 配置了多个环境变量
```bash
cd intranet-dns-web
npm i
npm run dev
```
