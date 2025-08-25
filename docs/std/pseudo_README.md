# 标准库伪代码参考

Origami 标准库的伪代码接口定义。

## 模块列表


### [functions](./functions.php)

标准库函数


### [functions](./context/functions.php)

标准库函数


### [functions](./database\sql/functions.php)

标准库函数


### [log](./log.php)

Log 类


### [exception](./exception.php)

Exception 类


### [os](./os.php)

OS 类


### [reflect](./reflect.php)

Reflect 类


### [server](./Net\Http/server.php)

Net\Http\Server 类


### [request](./Net\Http/request.php)

Net\Http\Request 类


### [response](./Net\Http/response.php)

Net\Http\Response 类


### [channel](./channel.php)

Channel 类


### [conn](./database\sql/conn.php)

database\sql\Conn 类


### [db](./database\sql/db.php)

database\sql\DB 类


### [row](./database\sql/row.php)

database\sql\Row 类


### [rows](./database\sql/rows.php)

database\sql\Rows 类


### [stmt](./database\sql/stmt.php)

database\sql\Stmt 类


### [tx](./database\sql/tx.php)

database\sql\Tx 类


### [txoptions](./database\sql/txoptions.php)

database\sql\TxOptions 类



## 快速开始

`php
<?php
// 使用标准库函数
dump("Hello World");

// 使用标准库类
$log = new Log();
$log->info("Application started");

// 使用反射
$reflect = new Reflect();
$classInfo = $reflect->getClassInfo("MyClass");
`

## 模块说明


### functions

标准库函数

**主要功能：**

- 函数：dumpinclude



[查看伪代码](./functions.php)

### functions

标准库函数

**主要功能：**

- 函数：backgroundwithCancelwithTimeoutwithValue



[查看伪代码](./context/functions.php)

### functions

标准库函数

**主要功能：**

- 函数：open



[查看伪代码](./database\sql/functions.php)

### log

Log 类

**主要功能：**


- 类：Log


[查看伪代码](./log.php)

### exception

Exception 类

**主要功能：**


- 类：Exception


[查看伪代码](./exception.php)

### os

OS 类

**主要功能：**


- 类：OS


[查看伪代码](./os.php)

### reflect

Reflect 类

**主要功能：**


- 类：Reflect


[查看伪代码](./reflect.php)

### server

Net\Http\Server 类

**主要功能：**


- 类：Net\Http\Server


[查看伪代码](./Net\Http/server.php)

### request

Net\Http\Request 类

**主要功能：**


- 类：Net\Http\Request


[查看伪代码](./Net\Http/request.php)

### response

Net\Http\Response 类

**主要功能：**


- 类：Net\Http\Response


[查看伪代码](./Net\Http/response.php)

### channel

Channel 类

**主要功能：**


- 类：Channel


[查看伪代码](./channel.php)

### conn

database\sql\Conn 类

**主要功能：**


- 类：database\sql\Conn


[查看伪代码](./database\sql/conn.php)

### db

database\sql\DB 类

**主要功能：**


- 类：database\sql\DB


[查看伪代码](./database\sql/db.php)

### row

database\sql\Row 类

**主要功能：**


- 类：database\sql\Row


[查看伪代码](./database\sql/row.php)

### rows

database\sql\Rows 类

**主要功能：**


- 类：database\sql\Rows


[查看伪代码](./database\sql/rows.php)

### stmt

database\sql\Stmt 类

**主要功能：**


- 类：database\sql\Stmt


[查看伪代码](./database\sql/stmt.php)

### tx

database\sql\Tx 类

**主要功能：**


- 类：database\sql\Tx


[查看伪代码](./database\sql/tx.php)

### txoptions

database\sql\TxOptions 类

**主要功能：**


- 类：database\sql\TxOptions


[查看伪代码](./database\sql/txoptions.php)

