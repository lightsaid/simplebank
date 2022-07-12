-- name: CreateAccount :one
insert into accounts (
    owner, balance, currency
) values (
    $1, $2, $3
)
returning *;

-- name: GetAccount :one
select * from accounts where id = $1 limit 1;

-- name: GetAccountForUpdate :one
select * from accounts where id = $1 limit 1 for no key update;

-- name: ListAccounts :many
select * from accounts 
order by id
limit $1 offset $2;

/*
没有返回
name: UpdateAccount :exec
update accounts set balance = $2 where id = $1;
*/

-- name: UpdateAccount :one
update accounts set balance = $2 where id = $1 returning *;

-- name: DeleteAccount :exec
delete from accounts where id = $1;