# 标准库伪代码参考

Origami 标准库的伪代码接口定义。

## 模块列表


### [functions](./functions.php)

标准库函数


### [functions](./context/functions.php)

标准库函数


### [functions](./Annotation/functions.php)

标准库函数


### [functions](./Database\Sql/functions.php)

标准库函数


### [functions](./Database/functions.php)

标准库函数


### [functions](./Net\Http/functions.php)

标准库函数


### [inject](./Annotation/inject.php)

Annotation\Inject 类


### [channel](./channel.php)

Channel 类


### [column](./Database\Annotation/column.php)

Database\Annotation\Column 类


### [generatedvalue](./Database\Annotation/generatedvalue.php)

Database\Annotation\GeneratedValue 类


### [id](./Database\Annotation/id.php)

Database\Annotation\Id 类


### [table](./Database\Annotation/table.php)

Database\Annotation\Table 类


### [db](./Database/db.php)

Database\DB 类


### [conn](./Database\Sql/conn.php)

Database\Sql\Conn 类


### [db](./Database\Sql/db.php)

Database\Sql\DB 类


### [row](./Database\Sql/row.php)

Database\Sql\Row 类


### [rows](./Database\Sql/rows.php)

Database\Sql\Rows 类


### [stmt](./Database\Sql/stmt.php)

Database\Sql\Stmt 类


### [tx](./Database\Sql/tx.php)

Database\Sql\Tx 类


### [txoptions](./Database\Sql/txoptions.php)

Database\Sql\TxOptions 类


### [datetime](./datetime.php)

DateTime 类


### [datetimezone](./datetimezone.php)

DateTimeZone 类


### [exception](./exception.php)

Exception 类


### [hashmap](./hashmap.php)

HashMap 类


### [list](./list.php)

List 类


### [log](./log.php)

Log 类


### [application](./Net\Annotation/application.php)

Net\Annotation\Application 类


### [controller](./Net\Annotation/controller.php)

Net\Annotation\Controller 类


### [deletemapping](./Net\Annotation/deletemapping.php)

Net\Annotation\DeleteMapping 类


### [getmapping](./Net\Annotation/getmapping.php)

Net\Annotation\GetMapping 类


### [postmapping](./Net\Annotation/postmapping.php)

Net\Annotation\PostMapping 类


### [putmapping](./Net\Annotation/putmapping.php)

Net\Annotation\PutMapping 类


### [route](./Net\Annotation/route.php)

Net\Annotation\Route 类


### [cookie](./Net\Http/cookie.php)

Net\Http\Cookie 类


### [handler](./Net\Http/handler.php)

Net\Http\Handler 类


### [request](./Net\Http/request.php)

Net\Http\Request 类


### [response](./Net\Http/response.php)

Net\Http\Response 类


### [server](./Net\Http/server.php)

Net\Http\Server 类


### [os](./os.php)

OS 类


### [reflect](./reflect.php)

Reflect 类


### [reflectionexception](./reflectionexception.php)

ReflectionException 类



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

- 函数：booldumpfloatggintobjectstring



[查看伪代码](./functions.php)

### functions

标准库函数

**主要功能：**

- 函数：backgroundwithCancelwithCancelCausewithDeadlinewithDeadlineCausewithTimeoutwithTimeoutCausewithValuewithoutCancel



[查看伪代码](./context/functions.php)

### functions

标准库函数

**主要功能：**

- 函数：__internal_spring_inline



[查看伪代码](./Annotation/functions.php)

### functions

标准库函数

**主要功能：**

- 函数：open



[查看伪代码](./Database\Sql/functions.php)

### functions

标准库函数

**主要功能：**

- 函数：getConnectiongetDefaultConnectionlistConnectionsregisterConnectionregisterDefaultConnectionremoveConnection



[查看伪代码](./Database/functions.php)

### functions

标准库函数

**主要功能：**

- 函数：app



[查看伪代码](./Net\Http/functions.php)

### inject

Annotation\Inject 类

**主要功能：**


- 类：Annotation\Inject


[查看伪代码](./Annotation/inject.php)

### channel

Channel 类

**主要功能：**


- 类：Channel


[查看伪代码](./channel.php)

### column

Database\Annotation\Column 类

**主要功能：**


- 类：Database\Annotation\Column


[查看伪代码](./Database\Annotation/column.php)

### generatedvalue

Database\Annotation\GeneratedValue 类

**主要功能：**


- 类：Database\Annotation\GeneratedValue


[查看伪代码](./Database\Annotation/generatedvalue.php)

### id

Database\Annotation\Id 类

**主要功能：**


- 类：Database\Annotation\Id


[查看伪代码](./Database\Annotation/id.php)

### table

Database\Annotation\Table 类

**主要功能：**


- 类：Database\Annotation\Table


[查看伪代码](./Database\Annotation/table.php)

### db

Database\DB 类

**主要功能：**


- 类：Database\DB


[查看伪代码](./Database/db.php)

### conn

Database\Sql\Conn 类

**主要功能：**


- 类：Database\Sql\Conn


[查看伪代码](./Database\Sql/conn.php)

### db

Database\Sql\DB 类

**主要功能：**


- 类：Database\Sql\DB


[查看伪代码](./Database\Sql/db.php)

### row

Database\Sql\Row 类

**主要功能：**


- 类：Database\Sql\Row


[查看伪代码](./Database\Sql/row.php)

### rows

Database\Sql\Rows 类

**主要功能：**


- 类：Database\Sql\Rows


[查看伪代码](./Database\Sql/rows.php)

### stmt

Database\Sql\Stmt 类

**主要功能：**


- 类：Database\Sql\Stmt


[查看伪代码](./Database\Sql/stmt.php)

### tx

Database\Sql\Tx 类

**主要功能：**


- 类：Database\Sql\Tx


[查看伪代码](./Database\Sql/tx.php)

### txoptions

Database\Sql\TxOptions 类

**主要功能：**


- 类：Database\Sql\TxOptions


[查看伪代码](./Database\Sql/txoptions.php)

### datetime

DateTime 类

**主要功能：**


- 类：DateTime


[查看伪代码](./datetime.php)

### datetimezone

DateTimeZone 类

**主要功能：**


- 类：DateTimeZone


[查看伪代码](./datetimezone.php)

### exception

Exception 类

**主要功能：**


- 类：Exception


[查看伪代码](./exception.php)

### hashmap

HashMap 类

**主要功能：**


- 类：HashMap


[查看伪代码](./hashmap.php)

### list

List 类

**主要功能：**


- 类：List


[查看伪代码](./list.php)

### log

Log 类

**主要功能：**


- 类：Log


[查看伪代码](./log.php)

### application

Net\Annotation\Application 类

**主要功能：**


- 类：Net\Annotation\Application


[查看伪代码](./Net\Annotation/application.php)

### controller

Net\Annotation\Controller 类

**主要功能：**


- 类：Net\Annotation\Controller


[查看伪代码](./Net\Annotation/controller.php)

### deletemapping

Net\Annotation\DeleteMapping 类

**主要功能：**


- 类：Net\Annotation\DeleteMapping


[查看伪代码](./Net\Annotation/deletemapping.php)

### getmapping

Net\Annotation\GetMapping 类

**主要功能：**


- 类：Net\Annotation\GetMapping


[查看伪代码](./Net\Annotation/getmapping.php)

### postmapping

Net\Annotation\PostMapping 类

**主要功能：**


- 类：Net\Annotation\PostMapping


[查看伪代码](./Net\Annotation/postmapping.php)

### putmapping

Net\Annotation\PutMapping 类

**主要功能：**


- 类：Net\Annotation\PutMapping


[查看伪代码](./Net\Annotation/putmapping.php)

### route

Net\Annotation\Route 类

**主要功能：**


- 类：Net\Annotation\Route


[查看伪代码](./Net\Annotation/route.php)

### cookie

Net\Http\Cookie 类

**主要功能：**


- 类：Net\Http\Cookie


[查看伪代码](./Net\Http/cookie.php)

### handler

Net\Http\Handler 类

**主要功能：**


- 类：Net\Http\Handler


[查看伪代码](./Net\Http/handler.php)

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

### server

Net\Http\Server 类

**主要功能：**


- 类：Net\Http\Server


[查看伪代码](./Net\Http/server.php)

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

### reflectionexception

ReflectionException 类

**主要功能：**


- 类：ReflectionException


[查看伪代码](./reflectionexception.php)

