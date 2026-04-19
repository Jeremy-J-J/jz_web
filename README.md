# AI数据集分享平台

一个基于 Vue3 + Go + SQLite 的电子资源展示平台，支持资源分类、搜索、访问统计和管理后台。

## 功能特性

- **前台展示**
  - 资源卡片网格展示
  - 类目导航筛选
  - 关键词搜索
  - 访问统计（访问总数、今日访问）
  - 点击跳转外部链接

- **管理后台**
  - JWT 认证登录
  - 资源增删改查
  - 类目管理
  - 访问统计查看
  - 封面上传

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue3 + Vite |
| 后端 | Go + Gin |
| 数据库 | SQLite3 |
| Web服务器 | Nginx |

## 目录结构

```
jz_web/
├── frontend/           # 前端源码
│   ├── src/           # Vue组件源码
│   └── dist/          # 构建产物
├── backend/           # 后端源码
│   ├── handlers/      # 路由处理器
│   ├── middleware/     # 中间件
│   ├── models/        # 数据模型
│   ├── utils/         # 工具函数
│   └── main.go        # 入口文件
├── conf/
│   ├── nginx/         # Nginx配置
│   └── systemd/       # systemd服务配置
├── data/
│   ├── db/            # SQLite数据库
│   └── uploads/       # 封面上传目录
├── runtime/
│   └── logs/          # 日志目录
└── SPEC.md            # 项目规范文档
```

## 快速部署

### 前置要求

- Go 1.21+
- Node.js 18+
- Nginx
- SQLite3

### 构建后端

```bash
cd backend
go build -o app .
```

### 构建前端

```bash
cd frontend
npm install
npm run build
```

### 配置systemd服务

```bash
sudo cp conf/systemd/jz-web.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable jz-web
sudo systemctl start jz-web
```

### Nginx配置

```bash
sudo cp conf/nginx/jz_web.conf /etc/nginx/sites-available/
sudo ln -s /etc/nginx/sites-available/jz_web.conf /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 默认管理员

- 用户名: `jeremyj`
- 密码: `jiangcc8484`

## API接口

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/categories | 获取类目列表 |
| GET | /api/resources | 获取资源列表 |
| GET | /api/resources/:id | 获取资源详情 |
| GET | /api/search | 搜索资源 |
| GET | /api/stats | 获取访问统计 |
| POST | /api/visit | 记录访问 |

### 管理接口（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/admin/login | 管理员登录 |
| GET | /api/admin/categories | 获取类目 |
| POST | /api/admin/categories | 新增类目 |
| PUT | /api/admin/categories/:id | 编辑类目 |
| DELETE | /api/admin/categories/:id | 删除类目 |
| GET | /api/admin/resources | 获取资源列表 |
| POST | /api/admin/resources | 新增资源 |
| PUT | /api/admin/resources/:id | 编辑资源 |
| DELETE | /api/admin/resources/:id | 删除资源 |
| PUT | /api/admin/resources/:id/toggle | 切换状态 |
| POST | /api/admin/upload | 上传封面 |

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| DB_PATH | /workspace/pro/jz_web/data/db/app.db | 数据库路径 |

## 许可证

MIT
