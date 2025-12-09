# {{.ProjectName}}

{{.Description}}

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+ (可选)

### 安装

1. 克隆项目
```bash
git clone {{.ModulePath}}
cd {{.ProjectName}}
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置数据库
编辑 `config/conf.yml` 文件，设置数据库连接信息

4. 运行项目
```bash
go run .
```

## API 文档

### 认证接口

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/refresh` - 刷新令牌

### 用户接口

- `GET /api/user/profile` - 获取当前用户信息
- `GET /api/user/:id` - 根据ID获取用户信息

## 项目结构

```
{{.ProjectName}}/
├── cmd/              # 应用程序入口
├── config/           # 配置文件
├── infra/            # 基础设施层
├── logic/            # 业务逻辑层
├── web/              # Web层
└── README.md
```

## 作者

{{.AuthorName}} <{{.AuthorEmail}}>

## 创建时间

{{.GeneratedAt}}