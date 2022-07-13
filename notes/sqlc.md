# sqlc 

### 安装
1. `sudo snap install sqlc`
1. `sqlc help` 查看命令帮助
1. `sqlc init` 在项目根目录初始化 -> 生存 sqlc.yaml 配置文件

### sqlc 使用
1. 目录结构
``` txt
db
   ├── migrations - 迁移 SQL
   ├── query      - SQL查询语句
   └── sqlc       - 存放 sqlc 生成代码
```
2. sqlc.yaml配置
``` yaml
version: "1"
packages:
  - name: "db"                   # go代码 package name
    path: "./db/sqlc"            # 生存代码存放目录
    queries: "./db/query/"       # 查询语句存放目录
    schema: "./db/migrations/"   # 迁移文件存放目录
    engine: "postgresql"         # 数据库引擎，postgresql、mysql...
    ...
    # 其他配置参考 sqlc 官网
    # https://docs.sqlc.dev/en/latest/reference/config.html
```
3. 编写 CRUD 查询语句 
db/query/accounts.sql
``` sql
-- name: CreateAccount :one
insert into accounts (
    owner, balance, currency
) values (
    $1, $2, $3
)
returning *;
```

4. 编写Makefiel命令并执行 `make sqlc`
``` Makefile
## sqlc:
sqlc:
	sqlc generate
```
- **至此，sqlc使用方式如是这般**
    

### sqlc.arg 用法
- SQL 语句语法
``` sql 
-- name: AddAccountBalance :one
update accounts set balance = balance + sqlc.arg(amount) where id = sqlc.arg(id) returning *;
```
- 生成的 go 代码
``` go
type AddAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
```
