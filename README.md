# utilsext

## 介绍

utilsext 是一个 Go 语言的工具库，提供了一些常用的工具函数和方法。
该库依赖于

-   [commonutils](https://github.com/qiuliaogit/commonutils) 库。
-   [redis](https://github.com/go-redis/redis/v8) 库。
-   [shopspring/decimal](https://github.com/shopspring/decimal) 库。

## 实现的功能有

-   1.0.0
    -   实现 redis 的工具类
        -   redis hset 工具类
        -   redis list 工具类
        -   redis set 工具类
        -   redis zset 工具类
        -   redis queue 工具类
    -   实现 decimal 的工具类
        -   decimal 的工具类, 类型转化和计算（加减乘除）
