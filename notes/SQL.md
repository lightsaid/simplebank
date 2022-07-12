# SQL  

### 在转账交易业务中发生并发问题
1. 场景：账户A 并发（同时2个或以上）向账户B 转账，如果不作任何处理，转账完成后，账户A/B的balance是对应的不上。
    - 导致这个根源就在于并发时，事务1更新账户A的banlace但是没有commit或者rollback，而此时
    事务2读取了账户A的balance，得到并不是最新的数据了，这就导致了最终结果和期望不一致。

    - 当数据库开启2或2个以上两个事务时，一个事务在更新同一行数据，另一个事务在读同一行数据，就导致了最终错误。
1. 已知道问题所在，解决方案可以使用排他锁：`for update`, 在查询语句 `select * from accounts for update`,
    那么事务2在查询的时候，如果事务1没有commit/rollback, select 语句就会停在那里等待，直到事务1commit/rollback；

1. 即使使用了`for update` 问题并没有就此解决。


### for update
1. for update 是一种行级锁，又叫排它锁。


### for no key update 更新不会影响到id（外键）
