# 标准库伪代码参考

Origami 标准库的伪代码接口定义。

## 模块列表


### [functions](./functions.php)

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

