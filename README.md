# S3Tool - S3 命令行管理工具

一个功能强大的 S3 命令行管理工具，支持存储桶管理、对象操作、预签名URL、多部分上传、版本控制、生命周期规则等功能。

## 功能特性

- ✅ **存储桶管理** - 创建、删除、列出、清空存储桶
- ✅ **对象操作** - 上传、下载、复制、移动、删除对象
- ✅ **预签名URL** - 生成临时访问URL
- ✅ **多部分上传** - 大文件分段上传
- ✅ **版本控制** - 管理存储桶版本控制配置
- ✅ **生命周期规则** - 配置对象生命周期管理
- ✅ **访问控制列表 (ACL)** - 管理桶和对象的访问权限
- ✅ **CORS配置** - 跨域资源共享配置
- ✅ **静态网站托管** - 将存储桶配置为静态网站
- ✅ **服务器端加密** - 配置存储桶加密
- ✅ **标签管理** - 管理存储桶和对象标签
- ✅ **日志记录** - 配置服务器访问日志
- ✅ **复制配置** - 跨区域复制配置
- ✅ **事件通知** - 配置存储桶事件通知

## 安装

```bash
go build -o s3tool main.go
sudo mv s3tool /usr/local/bin/
```

## 快速开始

### 1. 初始化配置

```bash
s3tool init
```

这将在 `~/.s3tool.yaml` 生成默认配置文件，你需要编辑该文件添加你的 AWS 凭证：

```yaml
region: us-east-1
access_key: YOUR_AWS_ACCESS_KEY
secret_key: YOUR_AWS_SECRET_KEY
# 可选：自定义端点（用于兼容其他S3兼容存储）
# endpoint: https://s3.example.com
# path_style: true
```

### 2. 列出所有存储桶

```bash
s3tool bucket list
```

## 命令详解

### 存储桶管理 (bucket)

| 命令 | 说明 | 用法 |
|------|------|------|
| `list` | 列出所有存储桶 | `s3tool bucket list` |
| `create` | 创建存储桶 | `s3tool bucket create <bucket-name>` |
| `delete` | 删除存储桶（必须为空） | `s3tool bucket delete <bucket-name>` |
| `info` | 获取存储桶信息 | `s3tool bucket info <bucket-name>` |
| `location` | 获取存储桶区域 | `s3tool bucket location <bucket-name>` |
| `empty` | 清空存储桶（删除所有对象） | `s3tool bucket empty <bucket-name>` |

**示例：**

```bash
# 创建存储桶
s3tool bucket create my-new-bucket

# 清空存储桶
s3tool bucket empty my-old-bucket
```

### 对象管理 (object)

| 命令 | 说明 | 用法 |
|------|------|------|
| `list` | 列出存储桶中的对象 | `s3tool object list <bucket> [flags]` |
| `put` | 上传本地文件到S3 | `s3tool object put <bucket> <key> <local-file> [flags]` |
| `get` | 从S3下载对象到本地 | `s3tool object get <bucket> <key> <local-file>` |
| `delete` | 删除指定对象 | `s3tool object delete <bucket> <key>` |
| `copy` | 在存储桶之间复制对象 | `s3tool object copy <source-bucket> <source-key> <dest-bucket> <dest-key>` |
| `info` | 获取对象详细信息 | `s3tool object info <bucket> <key>` |
| `url` | 生成对象访问URL | `s3tool object url <bucket> <key>` |
| `mv` | 移动/重命名对象 | `s3tool object mv <bucket> <source-key> <dest-key>` |

**Flags:**

- `list` 命令:
  - `-p, --prefix string` - 对象前缀过滤
  - `-r, --recursive` - 递归列出所有对象
  - `--max-keys int` - 最大返回对象数（默认1000）

- `put` 命令:
  - `-t, --content-type string` - 内容类型（默认 "application/octet-stream"）
  - `-s, --storage-class string` - 存储类型（如 STANDARD, REDUCED_REDUNDANCY）

**示例：**

```bash
# 列出存储桶中所有对象（递归）
s3tool object list my-bucket --recursive

# 上传文件，指定内容类型
s3tool object put my-bucket docs/report.pdf ./report.pdf --content-type application/pdf

# 下载对象
s3tool object get my-bucket docs/report.pdf ./downloaded_report.pdf
```

### 预签名URL (presign)

生成临时访问URL，可用于分享私有对象。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 生成GET请求预签名URL | `s3tool presign get <bucket> <key> [flags]` |
| `put` | 生成PUT请求预签名URL | `s3tool presign put <bucket> <key> [flags]` |
| `delete` | 生成DELETE请求预签名URL | `s3tool presign delete <bucket> <key> [flags]` |

**Flags:**
- `--expires int` - URL过期时间（秒，默认3600）

**示例：**

```bash
# 生成一个有效期为1小时的下载链接
s3tool presign get my-bucket private-file.txt --expires 60

# 生成上传链接，允许其他人上传文件
s3tool presign put my-bucket uploads/new-file.txt --expires 30
```

### 多部分上传 (multipart)

用于大文件上传，支持断点续传。

| 命令 | 说明 | 用法 |
|------|------|------|
| `init` | 初始化多部分上传 | `s3tool multipart init <bucket> <key>` |
| `upload` | 上传分片 | `s3tool multipart upload <bucket> <key> <upload-id> <part-number> <file>` |
| `list` | 列出已上传的分片 | `s3tool multipart list <bucket> <key> <upload-id>` |
| `complete` | 完成多部分上传 | `s3tool multipart complete <bucket> <key> <upload-id> <parts...>` |
| `abort` | 终止多部分上传 | `s3tool multipart abort <bucket> <key> <upload-id>` |

**示例：**

```bash
# 初始化多部分上传
s3tool multipart init my-bucket large-file.iso
# 返回: Upload ID: example-upload-id

# 上传分块
s3tool multipart upload my-bucket large-file.iso example-upload-id 1 part1.bin
s3tool multipart upload my-bucket large-file.iso example-upload-id 2 part2.bin

# 完成上传
s3tool multipart complete my-bucket large-file.iso example-upload-id 1:etag1 2:etag2
```

### 版本控制 (versioning)

管理存储桶的版本控制配置。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取版本控制状态 | `s3tool versioning get <bucket>` |
| `enable` | 启用版本控制 | `s3tool versioning enable <bucket>` |
| `suspend` | 暂停版本控制 | `s3tool versioning suspend <bucket>` |

**示例：**

```bash
# 启用版本控制
s3tool versioning enable my-bucket

# 查看版本控制状态
s3tool versioning get my-bucket
```

### 访问控制列表 (acl)

管理存储桶和对象的访问权限。

| 命令 | 说明 | 用法 |
|------|------|------|
| `bucket-get` | 获取存储桶ACL | `s3tool acl bucket-get <bucket>` |
| `bucket-set` | 设置存储桶ACL | `s3tool acl bucket-set <bucket> <acl>` |
| `object-get` | 获取对象ACL | `s3tool acl object-get <bucket> <key>` |
| `object-set` | 设置对象ACL | `s3tool acl object-set <bucket> <key> <acl>` |

**可用的ACL值：**
- `private` - 私有（默认）
- `public-read` - 公开读
- `public-read-write` - 公开读写
- `authenticated-read` - 认证用户读
- `bucket-owner-read` - 桶拥有者读
- `bucket-owner-full-control` - 桶拥有者完全控制

**示例：**

```bash
# 设置存储桶为公开可读
s3tool acl bucket-set my-public-bucket public-read

# 设置对象为私有
s3tool acl object-set my-bucket private-file.txt private
```

### 生命周期规则 (lifecycle)

管理对象生命周期规则，可用于自动归档或删除对象。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取生命周期规则 | `s3tool lifecycle get <bucket>` |
| `delete` | 删除所有生命周期规则 | `s3tool lifecycle delete <bucket>` |

**示例：**

```bash
# 查看生命周期规则
s3tool lifecycle get my-bucket
```

### CORS配置 (cors)

管理存储桶的跨域资源共享配置。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取CORS配置 | `s3tool cors get <bucket>` |
| `delete` | 删除CORS配置 | `s3tool cors delete <bucket>` |

**示例：**

```bash
# 查看CORS配置
s3tool cors get my-bucket
```

### 静态网站托管 (website)

将存储桶配置为静态网站。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取静态网站配置 | `s3tool website get <bucket>` |
| `delete` | 删除静态网站配置 | `s3tool website delete <bucket>` |

**示例：**

```bash
# 查看静态网站配置
s3tool website get my-website-bucket
```

### 服务器端加密 (encryption)

配置存储桶的服务器端加密。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取加密配置 | `s3tool encryption get <bucket>` |
| `delete` | 删除加密配置 | `s3tool encryption delete <bucket>` |

**示例：**

```bash
# 查看加密配置
s3tool encryption get my-bucket
```

### 标签管理 (tags)

管理存储桶和对象的标签。

| 命令 | 说明 | 用法 |
|------|------|------|
| `bucket-get` | 获取存储桶标签 | `s3tool tags bucket-get <bucket>` |
| `bucket-put` | 设置存储桶标签 | `s3tool tags bucket-put <bucket> <key1=value1> [key2=value2...]` |
| `bucket-delete` | 删除存储桶标签 | `s3tool tags bucket-delete <bucket>` |
| `object-get` | 获取对象标签 | `s3tool tags object-get <bucket> <key>` |

**示例：**

```bash
# 设置存储桶标签
s3tool tags bucket-put my-bucket Environment=Production Project=MyApp

# 获取存储桶标签
s3tool tags bucket-get my-bucket
```

### 日志记录 (logging)

配置服务器访问日志记录。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取日志记录配置 | `s3tool logging get <bucket>` |
| `disable` | 禁用日志记录 | `s3tool logging disable <bucket>` |

**示例：**

```bash
# 查看日志配置
s3tool logging get my-bucket
```

### 复制配置 (replication)

配置跨区域复制。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取复制配置 | `s3tool replication get <bucket>` |
| `delete` | 删除复制配置 | `s3tool replication delete <bucket>` |

**示例：**

```bash
# 查看复制配置
s3tool replication get my-bucket
```

### 事件通知 (notification)

配置存储桶事件通知。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取事件通知配置 | `s3tool notification get <bucket>` |
| `delete` | 删除事件通知配置 | `s3tool notification delete <bucket>` |

**示例：**

```bash
# 查看事件通知配置
s3tool notification get my-bucket
```

### 存储桶策略 (policy)

管理存储桶策略。

| 命令 | 说明 | 用法 |
|------|------|------|
| `get` | 获取存储桶策略 | `s3tool policy get <bucket>` |
| `set` | 设置存储桶策略（从JSON文件） | `s3tool policy set <bucket> <policy-file>` |
| `delete` | 删除存储桶策略 | `s3tool policy delete <bucket>` |

**示例：**

```bash
# 查看当前策略
s3tool policy get my-bucket

# 设置新策略
s3tool policy set my-bucket policy.json
```

## 全局Flags

| Flag | 说明 |
|------|------|
| `-c, --config string` | 指定配置文件路径（默认: ~/.s3tool.yaml） |
| `-v, --verbose` | 详细输出模式 |
| `-h, --help` | 显示帮助信息 |

## 配置文件示例

```yaml
# AWS 区域
region: us-east-1

# AWS 凭证
access_key: AKIAXXXXXXXXXXXXXX
secret_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# 可选：自定义端点（用于S3兼容存储如MinIO）
# endpoint: http://localhost:9000

# 可选：使用路径风格寻址（MinIO等需要）
# path_style: true

# 可选：使用特定的AWS配置文件
# profile: default
```

## 命令补全

生成Shell自动补全脚本：

```bash
# Bash
s3tool completion bash > /etc/bash_completion.d/s3tool

# Zsh
s3tool completion zsh > "${fpath[1]}/_s3tool"

# Fish
s3tool completion fish > ~/.config/fish/completions/s3tool.fish
```

## 故障排除

### 常见错误

1. **"NoCredentialProviders"**: 请检查配置文件中的 `access_key` 和 `secret_key` 是否正确。

2. **"BucketAlreadyExists"**: 存储桶名称必须全局唯一，请尝试其他名称。

3. **"BucketNotEmpty"**: 删除存储桶前必须先清空所有对象，使用 `s3tool bucket empty <bucket-name>`。

4. **"AccessDenied"**: 请确保你的凭证具有执行该操作的权限。

5. **"SignatureDoesNotMatch"**: 请检查 `secret_key` 是否正确，以及系统时间是否同步。

## 许可证

MIT
