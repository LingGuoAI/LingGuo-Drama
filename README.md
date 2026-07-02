# 灵果短剧 / spiritFruit

灵果短剧是一个面向 AI 短剧、动漫和视频创作的开源项目。项目提供从剧本、角色、场景、分镜到图片/视频生成和片段合并的后台工作台，适合用来搭建可二次开发的 AIGC 视频生产流程。

## 功能概览

- 剧本管理：支持剧本录入、AI 生成和章节化创作。
- 角色/场景/道具：从剧本中提取实体并生成对应素材。
- 分镜工作流：拆分分镜、生成帧提示词、生成分镜图片和视频。
- 视频合成：基于 FFmpeg 合并分镜视频片段。
- 异步任务：使用 Redis + Asynq 处理 AI 生成和视频合成任务。
- 管理后台：Vue 3 + Vite + TDesign 前端工作台。

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 后端 | Go、Gin、GORM、Cobra |
| 前端 | Vue 3、Vite、TDesign、Pinia |
| 数据库 | MySQL 8+ |
| 队列/缓存 | Redis、Asynq |
| 视频处理 | FFmpeg |
| 部署 | Docker Compose、Nginx |

## Docker 快速启动

推荐先用 Docker Compose 启动完整环境：

```bash
cp .env.example .env
docker compose up -d --build
```

启动后访问：

- 前端：http://localhost
- 后端健康检查：http://localhost:8080/healthz
- 后端就绪检查：http://localhost:8080/readyz
- API 文档入口：http://localhost/docs/index.html

默认管理员：

- 账号：`admin`
- 密码：`123456`

后端首次启动会自动执行 GORM AutoMigrate，并在空库中初始化默认管理员和系统菜单。

更多部署、备份、更新和排障命令见 [Docker 部署指南](docs/deployment.md)。

## 本地开发

### 环境要求

如果使用 Docker Compose 部署，MySQL、Redis、FFmpeg 会随容器一起准备好；如果本地开发或直接运行后端，请先安装以下软件：

| 依赖/软件 | 版本要求 | 说明 |
| --- | --- | --- |
| Go | 1.25+ | 后端运行与构建环境 |
| Node.js | 18+ | 前端依赖安装与 Vite 开发服务 |
| MySQL | 8.0+ | 核心业务数据库 |
| Redis | 6.0+ | 缓存与 Asynq 异步任务队列 |
| FFmpeg | 4.0+ | 视频合并、转码、时长分析等视频处理能力 |

### 安装 FFmpeg

视频生成和合成流程依赖 FFmpeg，本地运行后端或 Worker 前请确保 `ffmpeg` 已加入系统 PATH。

macOS:

```bash
brew install ffmpeg
```

Ubuntu/Debian:

```bash
sudo apt update
sudo apt install -y ffmpeg
```

CentOS/RHEL:

```bash
sudo dnf install -y ffmpeg
```

Windows:

从 [FFmpeg 官网](https://ffmpeg.org/download.html) 下载构建包，解压后将 `bin` 目录加入系统环境变量 `PATH`。

验证安装：

```bash
ffmpeg -version
ffprobe -version
```

### 版本检查

```bash
go version
node -v
npm -v
mysql --version
redis-server --version
ffmpeg -version
```

### 后端

```bash
cd server
cp .env.example .env
go mod download
go run main.go serve
```

后端默认监听 `8080`，配置来自 `server/.env`。

异步任务 Worker 需单独启动：

```bash
cd server
go run main.go worker
```

请确保 MySQL 中存在数据库：

```sql
CREATE DATABASE spirit_fruit CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 前端

```bash
cd web
npm install
npm run dev
```

前端开发服务默认监听 `3002`，开发环境接口配置在 `web/.env.development`：

```ini
VITE_API_URL=http://localhost:8080
VITE_API_URL_PREFIX=/admin/v1
```

## 关键配置

后端常用配置在 `server/.env.example` 和根目录 `.env.example` 中维护：

| 变量 | 说明 |
| --- | --- |
| `APP_PORT` | 后端监听端口 |
| `APP_KEY` | 应用密钥，生产环境必须修改 |
| `DB_HOST` / `DB_DATABASE` | MySQL 地址和库名 |
| `REDIS_HOST` / `REDIS_PASSWORD` | Redis 地址和密码 |
| `REDIS_MAIN_DB` | 业务 Redis DB |
| `REDIS_CACHE_DB` | 缓存 Redis DB |
| `REDIS_ASYNC_DB` | Asynq 队列 Redis DB |
| `AI_PROVIDER` | 默认文本/图片模型提供商 |
| `VIDEO_PROVIDER` | 默认视频生成提供商 |

AI 服务相关 Key 可按实际 provider 填写，例如 `OPENAI_API_KEY`、`GETGOAPI_API_KEY`、`VOLCES_API_KEY`、`MINIMAX_API_KEY`。

## 创作链路

1. 创建短剧项目。
2. 创建或 AI 生成剧本。
3. 提取角色、场景、道具。
4. 生成角色图、场景图、道具图。
5. 拆分分镜并生成帧提示词。
6. 生成分镜图片和视频。
7. 合并片段并导出成片。

项目流程梳理和可落地的改进建议见 [项目流程梳理与改进方案](docs/process-improvements.md)。

## 常用命令

```bash
# 查看容器状态
docker compose ps

# 查看后端日志
docker compose logs -f server

# 查看异步任务日志
docker compose logs -f worker

# 停止服务
docker compose down

# 本地后端启动
cd server && go run main.go serve

# 本地前端启动
cd web && npm run dev
```

## 安全提醒

- 不要提交 `.env`、API Key、数据库备份和生成文件。
- 生产环境必须修改默认的 `APP_KEY`、数据库密码、Redis 密码和默认管理员密码。
- 如果密钥曾经提交或泄漏，请立即吊销并重新生成。

## API 配置站点

2 分钟完成配置：[API 聚合站点](https://www.zdom.cn/pricing)

## 项目社群

灵果短剧 (Lingo Short Video) 项目社群：

![灵果短剧交流群二维码](https://img.aiaiapi.com/zhoudao.png)

## 参与贡献

欢迎提交 Issue 和 Pull Request。建议优先补充部署体验、任务状态可观测性、API 契约、测试和工作台流程体验。

## License

MIT
