## dns-service
基于gin框架开发的内网dns管理系统后端demo, 支持动态dns修改. 提供dns解析的基础服务为bind.

启动方式: `go run cmd/main.go`, [api详情](http://localhost:16789/swagger/index.html)


## 相关组件
**golang**: gin, gorm, go-swagger, go-jwt

**other**: mysql, redis, [bind9 doc](https://downloads.isc.org/isc/bind9/9.18.25/doc/arm/html/index.html)

bind9安装和TSIG key配置:
```bash
yum install bind
mkdir -p /var/named
chown -R named.named /var/named
tsig-keygen dns-service > /etc/named/dns-service.key

# 编辑key以及其他配置
vim /etc/named.conf
systemctl start named
```

<details><summary>bind9配置</summary>

```bash
# 主配置
cat /etc/named.conf
acl trust {
     127.0.0.1/32;
     172.24.8.122/32;
     10.89.254.10/32;
};

acl server_master {
    172.24.8.122;
};

options {
    directory "/var/named";
    dump-file "/var/named/data/cache_dump.db";
    recursion yes;
    version "MediaV DNS 1.0";
    auth-nxdomain no;
    zone-statistics yes;
    statistics-file "/var/named/data/named_stats.txt";
    listen-on-v6 { none; };
    allow-recursion { trust; };
    max-recursion-queries 100;
};

include "/etc/rndc.key";
include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";
include "/etc/named/dns-service.key";

controls {
	inet 127.0.0.1 port 953
 	allow { 127.0.0.1; } keys { "rndc-key"; };
};

logging {
    channel query_log {
        file "/data0/log/named/query.log"  versions 10 size 100m;
        severity        info;
        print-time        yes;
        print-category yes;
    };
    category queries {
        query_log;
    };

    channel update_log {
        file "/data0/log/named/update.log"  versions 10 size 100m;
        severity        info;
        print-time        yes;
        print-category yes;
    };
    category update {
        update_log;
    };

    channel general_log {
        file "/data0/log/named/general.log"  versions 10 size 100m;
        severity        info;
        print-time        yes;
        print-category yes;
    };
    category general { general_log; };

    channel xfer_log {
         file "/data0/log/named/xfer.log" versions 10 size 100m;
         severity info;
         print-category yes;
         print-severity yes;
         print-time yes;
    };
    category xfer-in { xfer_log; };
    category xfer-out { xfer_log; };
};

zone "." IN {
	type hint;
	file "named.ca";
};

zone "test.com" {
	type master;
	file "test.com.zone";
	allow-query { trust; };
	allow-update { key dns-service; };
};

# 测试zone配置
cat /var/named/test.com.zone
$ORIGIN .
$TTL 86400	; 1 day
test.com		IN SOA	ns1.test.com. admin.test.com. (
				2020021392 ; serial
				1800       ; refresh (30 minutes)
				600        ; retry (10 minutes)
				2592000    ; expire (4 weeks 2 days)
				3600       ; minimum (1 hour)
				)
			NS	ns1.test.com.

$ORIGIN test.com.
$TTL 60	; 1 minute
ns1 A 172.24.8.122
```
</details>
