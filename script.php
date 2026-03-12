<?php

// 精确复现 Laravel 容器的多层嵌套问题

class Container {
    public $bindings = [];
    
    public function bind($abstract, $concrete = null) {
        if (! $concrete instanceof Closure) {
            $concrete = $this->getClosure($abstract, $concrete);
        }
        $this->bindings[$abstract] = ['concrete' => $concrete];
    }
    
    // 这就是 Laravel 的 getClosure 实现
    protected function getClosure($abstract, $concrete) {
        return function ($container, $parameters = []) use ($abstract, $concrete) {
            echo "[CLOSURE L1] abstract=" . gettype($abstract) . ", concrete=" . gettype($concrete) . "\n";
            
            if ($abstract == $concrete) {
                echo "[CLOSURE L1] abstract==concrete, calling build\n";
                return $container->build($concrete);
            }

            echo "[CLOSURE L1] calling resolve with concrete=" . gettype($concrete) . "\n";
            return $container->resolve($concrete, $parameters);
        };
    }
    
    public function build($concrete) {
        echo "[BUILD] concrete type: " . gettype($concrete) . "\n";
        var_dump($concrete);
        
        if (!is_string($concrete)) {
            throw new Exception("Target class [" . (is_bool($concrete) ? 'false' : gettype($concrete)) . "] does not exist.");
        }
        
        return "built_" . $concrete;
    }
    
    public function resolve($abstract, $parameters = []) {
        echo "[RESOLVE] abstract type: " . gettype($abstract) . "\n";
        var_dump($abstract);
        
        $concrete = $this->getConcrete($abstract);
        echo "[RESOLVE] after getConcrete, concrete type: " . gettype($concrete) . "\n";
        var_dump($concrete);
        
        if ($this->isBuildable($concrete, $abstract)) {
            echo "[RESOLVE] is buildable, executing closure\n";
            return $concrete($this, $parameters);
        }
        
        return $this->build($concrete);
    }
    
    protected function getConcrete($abstract) {
        if (isset($this->bindings[$abstract])) {
            return $this->bindings[$abstract]['concrete'];
        }
        return $abstract;
    }
    
    protected function isBuildable($concrete, $abstract) {
        return $concrete === $abstract || $concrete instanceof Closure;
    }
    
    public function make($abstract) {
        return $this->resolve($abstract);
    }
}

echo "=== Test 1: Direct binding (InterfaceA -> ClassA) ===\n";
$container = new Container();
$container->bind('InterfaceA', 'ClassA');

try {
    $result = $container->make('InterfaceA');
    echo "SUCCESS: " . $result . "\n\n";
} catch (Exception $e) {
    echo "ERROR: " . $e->getMessage() . "\n\n";
}

echo "=== Test 2: Nested closure scenario ===\n";
// 模拟更复杂的场景：一个闭包返回另一个闭包
$container2 = new Container();

// 第一次绑定：InterfaceB -> 闭包
$container2->bind('InterfaceB', function($container) {
    echo "[OUTER CLOSURE] called\n";
    
    // 这个闭包又调用了 make，触发另一层解析
    return $container->make('DependencyC');
});

// 第二次绑定：DependencyC -> ClassC
$container2->bind('DependencyC', 'ClassC');

try {
    $result = $container2->make('InterfaceB');
    echo "SUCCESS: " . $result . "\n\n";
} catch (Exception $e) {
    echo "ERROR: " . $e->getMessage() . "\n\n";
}

echo "=== Test 3: Three-level nesting ===\n";
// 三层嵌套
$container3 = new Container();

$container3->bind('Level1', function($container) {
    echo "[LEVEL 1] called\n";
    return $container->make('Level2');
});

$container3->bind('Level2', function($container) {
    echo "[LEVEL 2] called\n";
    return $container->make('Level3');
});

$container3->bind('Level3', 'FinalClass');

try {
    $result = $container3->make('Level1');
    echo "SUCCESS: " . $result . "\n\n";
} catch (Exception $e) {
    echo "ERROR: " . $e->getMessage() . "\n\n";
}
