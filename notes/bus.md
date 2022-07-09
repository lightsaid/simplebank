# 业务逻辑

/**
转账交易，假设 a 账户 往 b 账户转 10 元
1. 创建一条转账记录(transfer)，amount = 10
2. a 账户创建一条记录(entry) = -10
3. b 账户创建一条记录（entry） = +10
4. a 账户 balance - 10
5. b 账户 balance + 10
*/