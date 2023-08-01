#vela-crontab 
# 磐石系统定时任务

## 时间格式
与Linux 中crontab命令相似，cron库支持用 5 个空格分隔的域来表示时间。这 5 个域含义依次为：

- Minutes：分钟，取值范围[0-59]，支持特殊字符* / , -；
- Hours：小时，取值范围[0-23]，支持特殊字符* / , -；
- Day of month：每月的第几天，取值范围[1-31]，支持特殊字符* / , - ?；
- Month：月，取值范围[1-12]或者使用月份名字缩写[JAN-DEC]，支持特殊字符* / , -；
- Day of week：周历，取值范围[0-6]或名字缩写[JUN-SAT]，支持特殊字符* / , - ?。
 
注意，月份和周历名称都是不区分大小写的，也就是说SUN/Sun/sun表示同样的含义（都是周日）。

特殊字符含义如下：

- \*：使用*的域可以匹配任何值，例如将月份域（第 4 个）设置为*，表示每个月；
- /：用来指定范围的步长，例如将小时域（第 2 个）设置为3-59/15表示第 3 分钟触发，以后每隔 15 分钟触发一次，因此第 2 次触发为第 18 分钟，第 3 次为 33 分钟。。。直到分钟大于 59；
- ,：用来列举一些离散的值和多个范围，例如将周历的域（第 5 个）设置为MON,WED,FRI表示周一、三和五；
- -：用来表示范围，例如将小时的域（第 1 个）设置为9-17表示上午 9 点到
- ?：只能用在月历和周历的域中，用来代替*，表示每月/周的任意一天。
 
预定义时间规则

为了方便使用，cron预定义了一些时间规则：

- @yearly：也可以写作@annually，表示每年第一天的 0 点。等价于0 0 1 1 *；
- @monthly：表示每月第一天的 0 点。等价于0 0 1 * *；
- @weekly：表示每周第一天的 0 点，注意第一天为周日，即周六结束，周日开始的那个 0 点。等价于0 0 * * 0；
- @daily：也可以写作@midnight，表示每天 0 点。等价于0 0 * * *；
- @hourly：表示每小时的开始。等价于0 * * * *。

固定时间间隔

cron支持固定时间间隔，格式为：
- @every <duration>   如: @every 1s , @every 1h
 
## vela.crontab(string)
参数为定时任务名称
```lua
    local cron = vela.crontab("metric")
    
    cron.task("@every 1s" , "获取系统信息1s" , function() end)
    cron.task("@every 5s" , "获取系统信息5s" , function() end)
    cron.task("@every 5s" , "获取系统信息5s" , function() end)
```

## cron.task(spec , title , function(t))
启动和添加任务函数
```lua
    cron.task("@every 1s" , "获取系统信息1s" , function() 
        --todo    
    end)
```