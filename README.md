# PawMatch（宠物领养平台）

面向城市流浪动物救助机构、宠物收容所与潜在领养家庭的全栈领养撮合平台：待领养动物展示、救助机构入驻与认证、领养申请审核、领养回访、救助故事社区与公益捐赠。

## Docker Compose 一键启动（推荐）

```bash
cp .env.example .env
docker compose up -d --build
```

访问地址：

- 前端：http://localhost:8101
- 后端 API：http://localhost:3101
- MinIO 控制台：http://localhost:9001（minioadmin / minioadmin）
- 健康检查：http://localhost:3101/healthz

默认账号：`admin / admin123`（管理员）、`shelter / org123`（救助机构）、`adopter / user123`（领养人）。

关闭并清理：

```bash
docker compose down -v --remove-orphans
```

## 本地开发（备选）

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
go build ./...
go test ./...
```

前端：

```bash
cd frontend
npm install
npm run dev
npm run build
```

## 技术栈

| 分层 | 技术 |
| --- | --- |
| 前端 | React 18 + TypeScript + Ant Design + Vite + Zustand + React Router |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 16 + Redis 7（热点缓存，go-redis/v9）+ MinIO（对象存储，minio-go/v7） |
| 认证 | JWT（github.com/golang-jwt/jwt/v5）+ RBAC |
| 其他 | validator/v10、log/slog、Nginx |

## 项目目录结构

```
gb-30/
├── docker-compose.yml
├── .env.example
├── README.md
├── database/init.sql
├── backend/
│   ├── cmd/server/            # main.go + migrate/seed
│   └── internal/
│       ├── config/            # DB/Redis/MinIO/JWT/限流配置
│       ├── model/             # 10 个实体
│       ├── repository/        # 按实体分文件
│       ├── service/           # 按实体分文件 + 状态机 + stats
│       ├── handler/           # 按实体分文件 + upload/home
│       ├── router/            # router.go + 按实体分文件
│       ├── middleware/        # auth/rbac/rate_limiter/error_handler/logger/cors
│       ├── dto/
│       ├── constants/         # application/pet/organization/error_codes/log_templates/messages
│       └── util/              # jwt/logger/formatters/minio_client/redis_client
└── frontend/
    ├── nginx.conf
    └── src/
        ├── api/               # user/organization/pet/application/review/post/donation/favorite
        ├── stores/            # authStore/userStore/petStore/applicationStore/reviewStore
        ├── components/common/ # PetCard/ApplicationStatusBadge/OrgCard/PostCard/ReviewTaskList/...
        ├── hooks/             # useAuth/useAdoptionStats/useReviewTasks/useDonationStats
        ├── pages/             # Home/PetList/PetDetail/OrgList/OrgDetail/Apply/.../Login
        ├── router/            # index.tsx
        ├── utils/             # request/dateFormat/fileUpload
        └── constants/         # application/pet/organization/errorCodes
```

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | gbadopt | Compose 项目名/容器前缀 |
| DB_NAME / DB_USER / DB_PASSWORD | gbadopt_db / gbadopt_user / gbadopt_pwd | PostgreSQL |
| REDIS_PORT | 6501 | Redis 端口 |
| MINIO_ROOT_USER / MINIO_ROOT_PASSWORD | minioadmin / minioadmin | MinIO 账号 |
| MINIO_PORT / MINIO_CONSOLE_PORT | 9000 / 9001 | MinIO API/控制台 |
| JWT_SECRET | change_me_to_a_long_random_string | JWT 密钥（生产必改） |
| FRONTEND_PORT / BACKEND_PORT / DB_PORT | 8101 / 3101 / 5601 | 端口映射 |

## Docker 部署说明

- 端口：前端 `8101:80`、后端 `${BACKEND_PORT:-3101}:8080`、DB `${DB_PORT:-5601}:5432`、Redis `${REDIS_PORT:-6501}:6379`、MinIO `${MINIO_PORT:-9000}:9000` + `${MINIO_CONSOLE_PORT:-9001}:9001`
- 数据卷：`db_data`、`redis_data`、`minio_data`（命名卷持久化）
- 依赖顺序：db/redis/minio healthcheck → backend `depends_on: service_healthy` → frontend
- 常见问题：端口冲突改 `.env` 后重新 up；重置数据 `docker compose down -v`；上传图片通过后端 `/api/v1/files/:key` 代理 MinIO 读取

## API 接口清单

> 后端统一前缀 `/api/v1`，响应统一为 `{ "code": 0, "message": "ok", "data": ... }`。标注「登录」的接口需携带 `Authorization: Bearer <JWT>`。

| 方法 | 路径 | 权限 | 说明 |
| --- | --- | --- | --- |
| GET | /healthz | 公开 | 健康检查 |
| GET | /api/v1/home/overview | 公开 | 首页聚合（Redis 缓存） |
| GET | /api/v1/files/:key | 公开 | 读取 MinIO 上传文件 |
| POST | /api/v1/uploads | 登录（限流） | 上传文件到 MinIO |
| POST | /api/v1/users/register | 公开（限流） | 注册并返回 JWT |
| POST | /api/v1/users/login | 公开（限流） | 登录并返回 JWT |
| GET | /api/v1/users/me | 登录 | 获取当前用户 |
| PUT | /api/v1/users/me | 登录 | 更新当前用户资料 |
| GET | /api/v1/orgs | 公开 | 机构分页列表 |
| GET | /api/v1/orgs/:id | 公开 | 机构详情 |
| POST | /api/v1/orgs | org（限流） | 注册机构 |
| PUT | /api/v1/orgs/:id/review | admin | 机构认证审核 |
| GET | /api/v1/pets | 公开 | 宠物分页列表/筛选 |
| GET | /api/v1/pets/:id | 公开 | 宠物详情 |
| POST | /api/v1/pets | org（限流） | 发布宠物 |
| PUT | /api/v1/pets/:id/status | org | 宠物状态变更 |
| POST | /api/v1/applications | 登录（限流） | 提交领养申请（事务：创建申请+宠物置为待领养） |
| GET | /api/v1/applications/me | 登录 | 我的申请列表 |
| GET | /api/v1/applications/org | org | 机构收到的申请 |
| PUT | /api/v1/applications/:id/status | 登录 | 申请状态流转（approved 时事务更新宠物为已领养） |
| GET | /api/v1/reviews/me | 登录 | 我的回访记录 |
| GET | /api/v1/reviews/org | org | 机构回访记录 |
| POST | /api/v1/reviews | org（限流） | 创建回访计划 |
| PUT | /api/v1/reviews/:id/submit | 登录 | 提交回访结果 |
| GET | /api/v1/posts | 公开 | 社区帖子列表 |
| GET | /api/v1/posts/:id | 公开 | 帖子详情 |
| GET | /api/v1/posts/:id/comments | 公开 | 帖子评论列表 |
| POST | /api/v1/posts | 登录（限流） | 发布帖子 |
| POST | /api/v1/posts/:id/comments | 登录（限流） | 发表评论（事务：写评论+自增评论数） |
| PUT | /api/v1/posts/:id/like | 登录 | 帖子点赞 |
| DELETE | /api/v1/comments/:id | 登录 | 删除自己的评论 |
| POST | /api/v1/donations | 登录（限流） | 捐款 |
| GET | /api/v1/donations/me | 登录 | 我的捐款记录 |
| POST | /api/v1/donations/usage | org | 发布善款使用记录 |
| GET | /api/v1/orgs/:id/usages | 公开 | 机构善款使用列表 |
| GET | /api/v1/orgs/:id/donation-stats | 公开 | 机构捐款统计 |
| GET | /api/v1/favorites | 登录 | 收藏列表 |
| POST | /api/v1/favorites | 登录（限流） | 添加收藏 |
| DELETE | /api/v1/favorites/:targetType/:targetId | 登录 | 取消收藏 |

## 枚举出现位置清单

### ApplicationStatus（submitted/org_review/communicating/confirmed/offline_interview/approved/rejected）

- 后端：`internal/constants/application.go`（定义+状态机）、`internal/model/adoption_application.go`（模型）、`internal/service/application_service.go`（流转校验）、`internal/util/formatters.go`（AppStatusText）、`internal/constants/log_templates.go`、`database/init.sql`
- 前端：`src/constants/application.ts`（定义）、`src/components/common/ApplicationStatusBadge.tsx`、`src/pages/Applications.tsx`（进度列表/筛选）、`src/hooks/useAdoptionStats.ts`

### PetSpecies（dog/cat/rabbit/other）

- 后端：`internal/constants/pet.go`（定义）、`internal/model/pet.go`（模型）、`internal/service/pet_service.go`（校验）、`internal/util/formatters.go`（PetSpeciesText）、`internal/constants/log_templates.go`、`database/init.sql`
- 前端：`src/constants/pet.ts`（定义）、`src/components/common/PetCard.tsx`、`src/pages/PetList.tsx`（筛选器）、`src/pages/PetDetail.tsx`

### PetStatus（available/pending/adopted）

- 后端：`internal/constants/pet.go`（定义）、`internal/model/pet.go`、`internal/service/pet_service.go`、`internal/service/application_service.go`（提交后置 pending/通过后置 adopted）、`internal/util/formatters.go`、`database/init.sql`
- 前端：`src/constants/pet.ts`、`src/components/common/PetCard.tsx`、`src/pages/PetList.tsx`、`src/pages/PetDetail.tsx`

### OrganizationStatus（pending/approved/rejected）

- 后端：`internal/constants/organization.go`（定义）、`internal/model/organization.go`、`internal/service/organization_service.go`、`internal/service/pet_service.go`（发布前校验认证）、`internal/util/formatters.go`、`database/init.sql`
- 前端：`src/constants/organization.ts`（定义）、`src/components/common/OrgCard.tsx`、`src/pages/OrgDetail.tsx`

## 横切关注点

- 认证授权（JWT + RBAC）：`middleware/auth.go`、`rbac.go`、`util/jwt.go`、路由中 org/admin 专属接口、前端 `RoleGuard.tsx` + 路由守卫
- 接口限流：`middleware/rate_limiter.go`（登录/申请/捐赠/上传启用）
- 全局错误处理：`middleware/error_handler.go`、`util/app_error.go`、`constants/error_codes.go`、前端 `utils/request.ts`
- 文件上传（MinIO）：`handler/upload_handler.go`、`util/minio_client.go`、`utils/fileUpload.ts`

## License

MIT
