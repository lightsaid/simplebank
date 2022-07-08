# 数据库
**数据库相关笔记**

### 本地已经存在postgresql容器了，因此不在构建
原构建postgres 的 docker-compose.yaml
``` yaml
version: '3.8'

services:
  postgres:
    image: 'postgres:14.4-alpine'
    container_name: 'postgres-14-alpine'
    ports:
      - "5432:5432"
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: abc123
      POSTGRES_DB: postgres
    volumes:
      - /home/xqq/developEnv/data/postgres/:/var/lib/postgresql/data
```

### 创建数据库
1. `docker ps` 找到 postgres-14-alpine 此容器
1. `docker exec -it postgres-14-alpine  /bin/sh` - 通过交互式终端进入容器
1. `createdb --username=postgres --owner=postgres simple_bank`
    - 创建 simple_bank 数据库，创建人和属主都是postgres用户
1. `psql simple_bank` 尝试访问
    - 此时报错： psql: error: connection to server on socket "/var/run/postgresql/.s.PGSQL.5432" failed: FATAL:  role "root" does not exist
    - 原因psql默认使用系统用户名访问，当前系统用户是root，而psql没有创建root用户，只有postgres用户
    - 尝试解决：`su postgres && psql simple_bank`, 如果系统没有postgrs用户则需要茶创建 
    - 至此，问题解决
1. `dropdb simple_bank`
