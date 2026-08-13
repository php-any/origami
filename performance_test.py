#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import time
import sys

def performance_test():
    print("=== 性能测试：一百万次赋值 ===")
    
    # 记录开始时间
    start_time = time.time()
    
    # 执行一百万次赋值操作
    for i in range(1, 1000001):
        value = i
        result = value * 2
        sum_val = value + result
    
    # 记录结束时间
    end_time = time.time()
    
    # 计算执行时间
    execution_time = end_time - start_time
    
    print("执行完成！")
    print(f"总执行时间: {execution_time:.4f} 秒")
    print(f"平均每次操作时间: {execution_time / 1000000 * 1000000:.6f} 微秒")
    print(f"每秒操作次数: {1000000 / execution_time:.0f} 次/秒")
    
    print(f"时间信息：{start_time} -> {end_time}")
    print(f"循环次数：{i}")
    print(f"循环内 value：{value}")
    print(f"循环内 result：{result}")
    print(f"循环内 sum：{sum_val}")
    
    print("\n=== 性能测试完成 ===")

if __name__ == "__main__":
    performance_test() 