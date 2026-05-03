## Secret Book（密码本）

### 特性

- 根据密码加密保存
- 修改密码信息后不会完全删除记录，会保留历史版本方便查看
- 密码错误3次销毁数据（永久销毁，无法恢复）

### 快速开始

> 详细帮助运行 `sbook help`

```bash
# 添加一个平台的账号
# 运行这行命令后会出现 'Enter Book Password:' 字样，表示密码本的密码，如果是第一次运行，就是新建密码
# 然后会出现一个 'Enter Password:' 字样，表示这个记录的密码
sbook add <platform> <account> [remark]

# 查看所有平台
sbook list

# 查看所有qq的账号信息
sbook list qq
```

---

### 安装

#### 手动构建

##### 依赖

- [go 1.24.4](https://go.dev/)

```bash
make

# 或（没有 make 命令）

# Windows
go build -o sbook.exe

# Linux
go build -o sbook
```

---

### 环境变量

- `SECRET_BOOK_PASSWORD` - 存放密码

```bash
# Windows cmd
set SECRET_BOOK_PASSWORD=password

# Windows powershell
$env:SECRET_BOOK_PASSWORD="password"

# Linux
export SECRET_BOOK_PASSWORD=password

# 下面执行命令时就不需要输入密码本的密码了
```

### 相关文件

- 创建时默认会在程序同级目录下生成 `secrets` 文件，除非指定 `--secrets-path` 选项
- 读取时默认会读取程序同级下的 `secrets` 文件，除非指定 `--secrets-path` 选项
