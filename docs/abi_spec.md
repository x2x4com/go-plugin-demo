# 插件ABI规范 v1.0

## 文件格式
```json
{
  "name": "插件名称",
  "version": "语义化版本",
  "exports": {
    "methods": {
      "方法名": {
        "params": ["参数类型"],
        "returns": "返回类型"
      }
    }
  }
}
```

## 类型系统
| 类型        | Go对应类型          |
|-----------|-------------------|
| int       | int               |
| float     | float64           |
| string    | string            |
| bool      | bool              |
| any       | interface{}       |

## 示例：计算器插件
```json
{
  "name": "calculator",
  "version": "1.0.0",
  "exports": {
    "methods": {
      "Add": {
        "params": ["int", "int"],
        "returns": "int"
      },
      "Subtract": {
        "params": ["int", "int"],
        "returns": "int"
      }
    }
  }
}