## 构建 Spring 示例

```bash
zy compile ./examples/spring --build
```

## 文件上传示例

`UploadController` 演示两种 multipart/form-data 自动绑定方式：

### 1. DTO（文本 + 文件）

```bash
curl -X POST http://localhost:8080/api/upload \
  -F "title=示例文档" \
  -F "file=@/path/to/hello.txt"
```

### 2. 单文件形参

```bash
curl -X POST http://localhost:8080/api/upload/avatar \
  -F "avatar=@/path/to/photo.png"
```

上传文件保存在 `storage/uploads/` 目录。
