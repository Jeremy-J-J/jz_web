# 电子资源展示平台（无支付版）规范文档

## 1. 项目概述

- **项目名称**: JZ Web 电子资源展示平台
- **项目类型**: 前后端分离的Web应用
- **核心功能**: 仅展示电子资源，游客可浏览/点击跳转外部链接，管理员可增删改资源
- **目标用户**: 管理员（自己）+ 游客（所有人）

## 2. 技术架构

| 层级 | 技术栈 | 说明 |
|------|--------|------|
| 访问层 | Nginx (8888端口) | 静态页面托管 + API反向代理 |
| 前端层 | Vue3 + Vite | 纯静态构建产物 |
| 后端层 | Go + Gin | 编译为单一二进制，内存<150MB |
| 存储层 | SQLite3 + 本地文件 | 数据存储 + 封面图存储 |

## 3. 目录结构

```
/workspace/pro/jz_web
├── frontend/        # 前端源码（Vue3 + Vite）
│   └── dist/        # 构建产物（Nginx托管）
├── backend/         # 后端源码（Go + Gin）
│   └── app          # 编译后二进制文件
├── conf/
│   ├── nginx/       # Nginx配置文件
│   └── systemd/     # systemd服务配置
├── data/
│   ├── db/          # SQLite数据库文件
│   └── uploads/     # 资源封面图片
└── runtime/
    └── logs/        # 日志目录
```

## 4. 数据模型

### 4.1 管理员表 (admin)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| username | TEXT | 用户名，唯一 |
| password | TEXT | 密码（bcrypt加密） |
| created_at | DATETIME | 创建时间 |

### 4.2 类目表 (category)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| name | TEXT | 类目名称 |
| description | TEXT | 类目描述 |
| sort | INTEGER | 排序（数字越小越靠前） |
| created_at | DATETIME | 创建时间 |

### 4.3 资源表 (resource)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| category_id | INTEGER | 关联类目ID |
| title | TEXT | 资源标题 |
| cover | TEXT | 封面图片路径 |
| description | TEXT | 资源描述 |
| link | TEXT | 跳转外部链接 |
| status | INTEGER | 状态：0=下架，1=上架 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

## 5. API接口设计

### 5.1 公开接口（游客+管理员）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/categories | 获取所有类目 |
| GET | /api/resources | 获取资源列表（仅上架） |
| GET | /api/resources/:id | 获取资源详情 |
| GET | /api/search | 搜索资源（支持keyword和category_id过滤） |

### 5.2 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/admin/login | 管理员登录 |

### 5.3 管理接口（需JWT认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/admin/categories | 获取所有类目 |
| POST | /api/admin/categories | 新增类目 |
| PUT | /api/admin/categories/:id | 编辑类目 |
| DELETE | /api/admin/categories/:id | 删除类目 |
| GET | /api/admin/resources | 获取所有资源（含下架） |
| POST | /api/admin/resources | 新增资源 |
| PUT | /api/admin/resources/:id | 编辑资源 |
| DELETE | /api/admin/resources/:id | 删除资源 |
| PUT | /api/admin/resources/:id/toggle | 切换上下架状态 |

## 6. 前端页面设计

### 6.1 前台页面（游客）

- **首页** (`/`): 类目导航 + 资源卡片列表网格展示 + 搜索框
- **详情页** (`/resource/:id`): 资源详情，含跳转按钮

### 6.2 管理后台（管理员）

- **登录页** (`/admin/login`): 管理员登录
- **管理面板** (`/admin`): 标签页切换（资源管理 / 类目管理）
- **新增/编辑页** (`/admin/resource/new`, `/admin/resource/:id/edit`): 资源表单（含类目选择）

## 7. 权限设计

| 角色 | 前台浏览 | 点击跳转 | 关键词检索 | 登录管理 | 增删改资源 | 增删改类目 |
|------|----------|----------|------------|----------|------------|------------|
| 游客 | ✓ | ✓ | ✓ | ✗ | ✗ | ✗ |
| 管理员 | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

## 8. 安全规则

1. 仅存在一个管理员账户（初始账号由部署时配置）
2. 管理接口全部需要JWT认证
3. 上传仅允许图片格式（jpg/jpeg/png/gif/webp）
4. 游客无法访问管理接口
5. 跳转链接由管理员完全控制
6. 删除类目时，该类目下的资源自动移至"未分类"

## 9. 环境要求

- 服务器: 2核2G内存，50GB存储
- 端口: 8888（对外），8080（后端内部）
- 无Docker，纯原生部署

## 10. 验收标准

- [x] 前台可浏览资源列表
- [x] 前台可查看资源详情
- [x] 前台点击跳转外部链接
- [x] 管理员可登录
- [x] 管理员可新增/编辑/删除资源
- [x] 管理员可控制资源上下架
- [x] 封面图上传功能
- [x] JWT鉴权机制
- [x] Nginx 8888端口配置
- [x] systemd自启服务配置
- [x] 总资源占用 < 150MB
- [x] 类目管理功能（增删改）
- [x] 类目导航筛选
- [x] 关键词搜索功能