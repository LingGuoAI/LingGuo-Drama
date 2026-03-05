markdown
# 🎬 灵果AI (LingGuo AI) 

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Stars](https://img.shields.io/github/stars/YourUsername/LingGuo-AI?style=social)](https://github.com/YourUsername/LingGuo-AI)
[![Forks](https://img.shields.io/github/forks/YourUsername/LingGuo-AI?style=social)](https://github.com/YourUsername/LingGuo-AI)

> **全球领先的开源 AI 短剧与 AI 动漫自动化创作平台** | **Open Source AI Short Drama & Anime Generator**

**灵果AI (LingGuo AI)** 是一个专注于 **AIGC 视频生成**的开源项目，致力于为创作者提供从“文本到视频 (Text-to-Video)”的一站式解决方案。无论你是想制作跌宕起伏的**AI短剧**，还是想创作画风精美的**AI动漫**，灵果AI 都能通过自动化分镜、角色一致性保持、AI配音与视频合成技术，帮你将创意瞬间转化为视觉盛宴。

[📖 官方文档 (Docs)]() | [🌐 在线体验 (Demo)]() | [💬 加入社区 (Discord/WeChat)]()

---

## 🔍 为什么选择灵果AI？(Why LingGuo AI?)

在 AIGC 爆发的时代，制作短剧和动漫的门槛正在被打破。灵果AI 作为一款完全**开源免费**的工具，解决了传统 AI 视频创作中流程繁琐、角色不连贯、分镜难控制等痛点。

### ✨ 核心特性 (Key Features)

* 📝 **剧本到视频 (Script-to-Video)**：输入大纲或剧本，LLM 自动拆解分镜提示词（Prompt）。
* 🎭 **极致的角色一致性 (Character Consistency)**：内置先进的 LoRA 和 ControlNet 调度策略，确保同角色在不同镜头下长相、服装一致。
* 🎨 **专精 AI 短剧与动漫模型**：集成针对**竖屏短剧**、**日系动漫 (Anime)**、**国风3D**等场景优化的视觉生成管线。
* 🗣️ **智能情感配音 (AI Voiceover & Lipsync)**：一键生成带有情感的 TTS 语音，并支持唇形同步 (Wav2Lip/SadTalker)。
* 🎞️ **自动化后期合成 (Auto Video Editing)**：自动匹配转场特效、背景音乐 (BGM) 和音效，一键导出 MP4。
* 🧩 **高度可扩展的插件系统**：轻松接入 ComfyUI 工作流、Stable Video Diffusion (SVD)、Runway、Sora 等底层大模型 API。

### 💡 应用场景 (Use Cases)

* 🎬 **爆款 AI 短剧制作**：快速批量生产网文改短剧，抢占短视频平台流量红利。
* 🌸 **个人 AI 动漫连载**：画渣也能做导演，独立完成原创连载动漫的画面与配音。
* 📚 **小说推文视频化**：小说作者/推文博主实现“图文转视频”的高效变现。
* 📺 **商业广告与 MV**：低成本制作产品宣发视频、音乐 MV 视觉切片。

### 🗺️ 发展路线图 (Roadmap)

* [x] V 0.1: 核心 Text-to-Image 工作流搭建完成。
* [x] V 0.2: 接入 AI 配音引擎与一键视频合成模块。
* [ ] V 0.5: 推出专属“角色一致性”解决方案及动漫专用 Checkpoint。
* [ ] V 1.0: 完善多语言支持，推出可视化分镜编辑器 (Storyboard Editor)。
* [ ] V 2.0: 深度集成 AI 视频大模型 (如 Veo/Sora 类模型 API)，实现全时长动态生成。

---

## 🚀 快速开始 (Quick Start)

灵果 AI 采用现代化的前后端分离架构。后端基于 **Go** 构建，保障高并发下的视频合成与异步任务调度；前端基于 **Vue 3 + Vite**，提供丝滑的专业级视频剪辑轨道交互体验。

### 📋 环境要求 (Prerequisites)

| 依赖/软件 | 版本要求 | 说明 |
| :--- | :--- | :--- |
| **Go** | 1.20+ | 后端核心运行环境 |
| **Node.js** | 18+ | 前端构建与开发环境 (推荐使用 `pnpm`) |
| **MySQL** | 8.0+ | 核心业务数据库 |
| **Redis** | 6.0+ | **必需**，用于 Asynq 异步任务队列与缓存 |
| **FFmpeg** | 4.0+ | **必需**，用于底层的视频合并、转场与音频处理 |

#### 🔧 安装 FFmpeg
视频合成强依赖 FFmpeg，请确保安装并将其添加到系统环境变量中：
* **macOS**: `brew install ffmpeg`
* **Ubuntu/Debian**: `sudo apt update && sudo apt install ffmpeg`
* **Windows**: 从 [FFmpeg 官网](https://ffmpeg.org/download.html) 下载并配置环境变量。
* **验证安装**: 在终端运行 `ffmpeg -version`

---

### 📦 1. 获取代码

```bash
# 克隆灵果AI仓库
git clone [https://github.com/YourUsername/LingGuo-AI.git](https://github.com/YourUsername/LingGuo-AI.git)
cd LingGuo-AI

```

### ⚙️ 2. 后端部署 (Server Setup)

后端包含 RESTful API 接口与基于 Asynq 的异步视频生成/合并任务引擎。

```bash
# 进入后端目录
cd server

# 1. 复制环境变量配置文件
cp .env.example .env

```

**编辑 `.env` 文件**，配置您的数据库、Redis 和大模型 API Key：

> **注意**：请先在您的 MySQL 中手动创建一个名为 `spirit_fruit` 的空数据库 (`CREATE DATABASE spirit_fruit CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`)。

```ini
# 数据库配置
DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=spirit_fruit
DB_USERNAME=root
DB_PASSWORD=您的数据库密码

# Redis 配置 (用于任务队列)
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=您的Redis密码

# AI 模型配置 (根据需要填入)
OPENAI_API_KEY=sk-xxxxxx
VOLCES_API_KEY=xxxxxx

```

**启动后端服务：**

```bash
# 下载 Go 依赖
go mod tidy

# 启动服务
go run main.go serve

```

> **💡 自动化建表提示**：项目内置了 GORM 自动迁移（AutoMigrate）。只要数据库连接成功，**首次启动时会自动创建所有数据表，并播种默认的管理员账号与系统菜单**。
> * 默认超级管理员账号：`admin`
> * 默认超级管理员密码：`123456`
>
>

---

### 🎨 3. 前端部署 (Web Setup)

请打开一个**新的终端窗口**，进入前端目录：

```bash
# 进入前端目录
cd web

# 1. 安装前端依赖 (推荐使用 pnpm，如果没有请先 npm install -g pnpm)
pnpm install

# 2. 启动开发服务器
pnpm run dev

```

启动成功后，浏览器访问终端提示的地址（通常是 `http://localhost:3002`）。
输入刚才自动生成的默认账号 `admin` 和密码 `123456`，即可开启您的 AI 短剧创作之旅！

---

## 🏭 生产环境部署 (Production Deployment)

在生产环境中，建议将前后端分别构建并使用 Nginx 进行反向代理。

### 1. 编译构建

```bash
# 构建前端
cd web
pnpm run build
# 构建产物将在 web/dist 目录下

# 编译后端
cd ../server
go build -o spirit-fruit-server main.go

```

### 2. 使用 Systemd 守护后端进程

创建 `/etc/systemd/system/spirit-fruit.service`：

```ini
[Unit]
Description=LingGuo AI Server
After=network.target mysql.service redis.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/LingGuo-AI/server
ExecStart=/opt/LingGuo-AI/server/spirit-fruit-server serve
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

```

运行：`sudo systemctl enable --now spirit-fruit`

### 3. Nginx 反向代理配置参考

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /opt/LingGuo-AI/web/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /api/ {
        proxy_pass [http://127.0.0.1:8080/](http://127.0.0.1:8080/);
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 本地生成的视频与图片资源访问
    location /uploads/ {
        alias /opt/LingGuo-AI/server/uploads/;
        # 允许跨域以便视频播放
        add_header Access-Control-Allow-Origin *; 
    }
}

```

---

## 🤝 参与贡献 (Contributing)

灵果AI 是一个由社区驱动的开源项目。我们非常欢迎各位开发者、提示词工程师和视频创作者加入我们！
如果您有任何好点子或发现了 Bug，欢迎提交 Pull Requests 或发布 Issues。

## 📄 开源协议 (License)

本项目基于 MIT License 开源。欢迎免费用于商业或个人项目，但请保留原作者署名。

