-- 从业务角度，转账记录不应该有 删除或者修改

-- name: CreateEntry :one
insert into entries (
    account_id, amount
) values (
    $1, $1
) returning *;

-- name: GetEntry :one
select * from entries where id = $1 limit 1;

-- name: GetEntries :many
select * from entries 
where account_id = $1
order by id 
limit $2 offset $3;

