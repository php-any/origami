# 标准库伪代码参考

Origami 标准库的伪代码接口定义。

## 模块列表


### [functions](./functions.zy)

标准库函数


### [functions](./context/functions.zy)

标准库函数


### [functions](./Database\Sql/functions.zy)

标准库函数


### [log](./log.zy)

Log 类


### [exception](./exception.zy)

Exception 类


### [os](./os.zy)

OS 类


### [reflect](./reflect.zy)

Reflect 类


### [request](./Net\Http/request.zy)

Net\Http\Request 类


### [response](./Net\Http/response.zy)

Net\Http\Response 类


### [server](./Net\Http/server.zy)

Net\Http\Server 类


### [channel](./channel.zy)

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



[查看伪代码](./functions.zy)

### functions

标准库函数

**主要功能：**

- 函数：backgroundwithCancelwithTimeoutwithValue



[查看伪代码](./context/functions.zy)

### functions

标准库函数

**主要功能：**

- 函数：open



[查看伪代码](./Database\Sql/functions.zy)

### log

Log 类

**主要功能：**


- 类：Log


[查看伪代码](./log.zy)

### exception

Exception 类

**主要功能：**


- 类：Exception


[查看伪代码](./exception.zy)

### os

OS 类

**主要功能：**


- 类：OS


[查看伪代码](./os.zy)

### reflect

Reflect 类

**主要功能：**


- 类：Reflect


[查看伪代码](./reflect.zy)

### request

Net\Http\Request 类

**主要功能：**


- 类：Net\Http\Request


[查看伪代码](./Net\Http/request.zy)

### response

Net\Http\Response 类

**主要功能：**


- 类：Net\Http\Response


[查看伪代码](./Net\Http/response.zy)

### server

Net\Http\Server 类

**主要功能：**


- 类：Net\Http\Server


[查看伪代码](./Net\Http/server.zy)

### channel

Channel 类

**主要功能：**


- 类：Channel


[查看伪代码](./channel.zy)

