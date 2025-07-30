# 反射与注解系统

本文档介绍如何使用 Origami 的反射系统来读取和分析注解信息。

## 概述

反射系统允许你在运行时检查类、方法和属性的注解信息，这对于构建框架、依赖注入系统和路由系统非常有用。

## 反射类 (Reflect)

### 基本用法

```php
// 创建反射实例
$reflect = new Reflect();
```

### 注解相关方法

#### 1. getAllAnnotations(className)

获取类的所有注解信息，包括类注解、属性注解和方法注解。

```php
$reflect = new Reflect();
$annotations = $reflect->getAllAnnotations("App\\Controller\\UserController");
echo $annotations;
```

**输出示例：**

```
=== App\Controller\UserController 类的完整注解信息 ===

类注解:
  1. Annotation\Controller
  2. Annotation\Route

属性注解:
  userService:
    1. Annotation\Inject

方法注解:
  getUserList:
    1. Annotation\GetMapping
```

#### 2. getClassAnnotations(className)

获取类的注解信息。

```php
$classAnnotations = $reflect->getClassAnnotations("App\\Controller\\UserController");
echo $classAnnotations;
```

**输出示例：**

```
类 App\Controller\UserController 的注解信息:
  1. Annotation\Controller
  2. Annotation\Route
```

#### 3. getPropertyAnnotations(className, propertyName)

获取指定属性的注解信息。

```php
$propertyAnnotations = $reflect->getPropertyAnnotations("App\\Controller\\UserController", "userService");
echo $propertyAnnotations;
```

**输出示例：**

```
属性 App\Controller\UserController::userService 的注解信息:
  1. Annotation\Inject
```

#### 4. getMethodAnnotations(className, methodName)

获取指定方法的注解信息。

```php
$methodAnnotations = $reflect->getMethodAnnotations("App\\Controller\\UserController", "getUserList");
echo $methodAnnotations;
```

**输出示例：**

```
方法 App\Controller\UserController::getUserList 的注解信息:
  1. Annotation\GetMapping
```

#### 5. getAnnotationDetails(className, memberType, memberName)

获取指定成员的详细注解信息。

```php
// 获取类注解详细信息
$classDetails = $reflect->getAnnotationDetails("App\\Controller\\UserController", "class", "UserController");

// 获取属性注解详细信息
$propertyDetails = $reflect->getAnnotationDetails("App\\Controller\\UserController", "property", "userService");

// 获取方法注解详细信息
$methodDetails = $reflect->getAnnotationDetails("App\\Controller\\UserController", "method", "getUserList");
```

**参数说明：**

- `className`: 完整的类名（包含命名空间）
- `memberType`: 成员类型（"class", "property", "method"）
- `memberName`: 成员名称

## 完整示例

### 1. 基本注解反射

```php
namespace App\Controller;

use Annotation\Route;
use Annotation\Controller;
use Annotation\GetMapping;
use Annotation\Inject;

@Controller
@Route(prefix: "/api/users")
class UserController {
    @Inject(service: "UserService")
    public $userService;

    @GetMapping(path: "/list")
    public function getUserList() {
        return "Hello from getUserList";
    }
}

echo "=== 反射读取注解信息 ===\n";

// 创建反射实例
$reflect = new Reflect();

// 获取类的所有注解信息
echo "--- UserController 类的完整注解信息 ---\n";
echo $reflect->getAllAnnotations("App\\Controller\\UserController");

// 获取类注解
echo "\n--- UserController 类注解 ---\n";
echo $reflect->getClassAnnotations("App\\Controller\\UserController");

// 获取属性注解
echo "\n--- userService 属性注解 ---\n";
echo $reflect->getPropertyAnnotations("App\\Controller\\UserController", "userService");

// 获取方法注解
echo "\n--- getUserList 方法注解 ---\n";
echo $reflect->getMethodAnnotations("App\\Controller\\UserController", "getUserList");
```

### 2. 路由系统示例

```php
// 模拟路由系统
class Router {
    private $routes = [];

    public function registerController($className) {
        $reflect = new Reflect();

        // 获取类信息
        $classInfo = $reflect->getClassInfo($className);
        echo "注册控制器: {$className}\n";

        // 获取类注解
        $classAnnotations = $reflect->getClassAnnotations($className);
        echo "类注解: {$classAnnotations}\n";

        // 获取所有方法
        $methods = $reflect->listMethods($className);
        echo "方法列表: {$methods}\n";

        // 检查特定方法是否有路由注解
        $getUserListAnnotations = $reflect->getMethodAnnotations($className, "getUserList");
        if ($getUserListAnnotations->indexOf("GetMapping") != -1) {
            echo "发现路由方法: getUserList\n";
            $this->routes[] = [
                'controller' => $className,
                'method' => 'getUserList',
                'type' => 'GET'
            ];
        }

        echo "注册完成，发现路由方法: getUserList\n";
    }

    public function getRoutes() {
        return $this->routes;
    }
}

// 使用路由系统
$router = new Router();
$router->registerController("App\\Controller\\UserController");
```

## 注解类型

### 特性注解 (Feature Annotations)

特性注解只接收注解声明的参数，主要用于标记和元数据存储。

**示例：**

```php
@Controller(name: "UserController")
@Route(prefix: "/api/users")
@GetMapping(path: "/list")
```

### 宏注解 (Macro Annotations)

宏注解可以接收注解参数和被注解的节点，能够修改语法节点。

**示例：**

```php
@Inject(service: "UserService")
public $userService;
```

## 反射系统的应用场景

### 1. 框架开发

```php
// 自动注册控制器
class Framework {
    public function registerControllers() {
        $reflect = new Reflect();

        // 扫描所有类
        $classes = $reflect->listClasses();

        foreach ($classes as $className) {
            $classAnnotations = $reflect->getClassAnnotations($className);

            if ($classAnnotations->indexOf("Controller") != -1) {
                echo "发现控制器: $className\n";
                $this->registerController($className);
            }
        }
    }
}
```

### 2. 依赖注入

```php
// 自动注入依赖
class DependencyInjector {
    public function injectDependencies($object) {
        $reflect = new Reflect();
        $className = get_class($object);

        $properties = $reflect->listProperties($className);

        foreach ($properties as $propertyName) {
            $propertyAnnotations = $reflect->getPropertyAnnotations($className, $propertyName);

            if ($propertyAnnotations->indexOf("Inject") != -1) {
                echo "注入依赖: $propertyName\n";
                $this->injectProperty($object, $propertyName);
            }
        }
    }
}
```

### 3. 路由注册

```php
// 自动注册路由
class RouteRegistry {
    public function registerRoutes($controllerClass) {
        $reflect = new Reflect();

        $methods = $reflect->listMethods($controllerClass);

        foreach ($methods as $methodName) {
            $methodAnnotations = $reflect->getMethodAnnotations($controllerClass, $methodName);

            if ($methodAnnotations->indexOf("GetMapping") != -1) {
                echo "注册 GET 路由: $methodName\n";
                $this->registerGetRoute($controllerClass, $methodName);
            }

            if ($methodAnnotations->indexOf("PostMapping") != -1) {
                echo "注册 POST 路由: $methodName\n";
                $this->registerPostRoute($controllerClass, $methodName);
            }
        }
    }
}
```

## 最佳实践

### 1. 错误处理

```php
$reflect = new Reflect();

try {
    $annotations = $reflect->getAllAnnotations("NonExistentClass");
    echo $annotations;
} catch (Exception $e) {
    echo "类不存在或没有注解信息\n";
}
```

### 2. 性能优化

```php
// 缓存反射结果
class AnnotationCache {
    private $cache = [];

    public function getAnnotations($className) {
        if (!isset($this->cache[$className])) {
            $reflect = new Reflect();
            $this->cache[$className] = $reflect->getAllAnnotations($className);
        }

        return $this->cache[$className];
    }
}
```

### 3. 类型检查

```php
// 检查注解类型
function hasAnnotation($annotations, $annotationType) {
    return $annotations->indexOf($annotationType) != -1;
}

$reflect = new Reflect();
$classAnnotations = $reflect->getClassAnnotations("App\\Controller\\UserController");

if (hasAnnotation($classAnnotations, "Controller")) {
    echo "这是一个控制器类\n";
}
```

## 注意事项

1. **性能考虑**: 反射操作相对较慢，建议在应用启动时进行，而不是在请求处理过程中
2. **错误处理**: 总是检查类和方法是否存在，避免运行时错误
3. **命名空间**: 使用完整的类名（包含命名空间）进行反射操作
4. **注解顺序**: 注解的解析顺序可能与声明顺序不同

## 总结

反射系统为 Origami 提供了强大的元编程能力，可以：

- 在运行时检查类的结构和注解
- 实现自动化的框架功能
- 构建依赖注入系统
- 创建路由注册机制
- 实现各种元数据驱动的功能

通过合理使用反射和注解，可以大大简化框架开发和应用构建过程。
