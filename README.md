🎬 灵果AI (LingGuo AI)
全球领先的开源 AI 短剧与 AI 动漫自动化创作平台 | Open Source AI Short Drama & Anime Generator

灵果AI (LingGuo AI) 是一个专注于 AIGC 视频生成的开源项目，致力于为创作者提供从“文本到视频 (Text-to-Video)”的一站式解决方案。无论你是想制作跌宕起伏的AI短剧，还是想创作画风精美的AI动漫，灵果AI 都能通过自动化分镜、角色一致性保持、AI配音与视频合成技术，帮你将创意瞬间转化为视觉盛宴。

📖 官方文档 (Docs) | 🌐 在线体验 (Demo) | 💬 加入社区 (Discord/WeChat)

🔍 为什么选择灵果AI？(Why LingGuo AI?)
在 AIGC 爆发的时代，制作短剧和动漫的门槛正在被打破。灵果AI 作为一款完全开源免费的工具，解决了传统 AI 视频创作中流程繁琐、角色不连贯、分镜难控制等痛点。

✨ 核心特性 (Key Features)
📝 剧本到视频 (Script-to-Video)：输入大纲或剧本，LLM 自动拆解分镜提示词（Prompt）。

🎭 极致的角色一致性 (Character Consistency)：内置先进的 LoRA 和 ControlNet 调度策略，确保同角色在不同镜头下长相、服装一致。

🎨 专精 AI 短剧与动漫模型：集成针对竖屏短剧、日系动漫 (Anime)、国风3D等场景优化的视觉生成管线。

🗣️ 智能情感配音 (AI Voiceover & Lipsync)：一键生成带有情感的 TTS 语音，并支持唇形同步 (Wav2Lip/SadTalker)。

🎞️ 自动化后期合成 (Auto Video Editing)：自动匹配转场特效、背景音乐 (BGM) 和音效，一键导出 MP4。

🧩 高度可扩展的插件系统：轻松接入 ComfyUI 工作流、Stable Video Diffusion (SVD)、Runway、Sora 等底层大模型 API。

💡 应用场景 (Use Cases)
🎬 爆款 AI 短剧制作：快速批量生产网文改短剧，抢占短视频平台流量红利。

🌸 个人 AI 动漫连载：画渣也能做导演，独立完成原创连载动漫的画面与配音。

📚 小说推文视频化：小说作者/推文博主实现“图文转视频”的高效变现。

📺 商业广告与 MV：低成本制作产品宣发视频、音乐 MV 视觉切片。

🗺️ 发展路线图 (Roadmap)
[x] V 0.1: 核心 Text-to-Image 工作流搭建完成。

[x] V 0.2: 接入 AI 配音引擎与一键视频合成模块。

[ ] V 0.5: 推出专属“角色一致性”解决方案及动漫专用 Checkpoint。

[ ] V 1.0: 完善多语言支持，推出可视化分镜编辑器 (Storyboard Editor)。

[ ] V 2.0: 深度集成 AI 视频大模型 (如 Veo/Sora 类模型 API)，实现全时长动态生成。

🚀 快速开始 (Quick Start)灵果 AI 采用现代化的前后端分离架构。后端基于 Go 构建，保障高并发下的视频合成与异步任务调度；前端基于 Vue 3 + Vite，提供丝滑的专业级视频剪辑轨道交互体验。📋 环境要求 (Prerequisites)依赖/软件版本要求说明Go1.20+后端核心运行环境Node.js18+前端构建与开发环境 (推荐使用 pnpm)MySQL8.0+核心业务数据库Redis6.0+必需，用于 Asynq 异步任务队列与缓存FFmpeg4.0+必需，用于底层的视频合并、裁剪与音频处理🔧 安装 FFmpeg视频合成强依赖 FFmpeg，请确保安装并将其添加到系统环境变量中：macOS: brew install ffmpegUbuntu/Debian: sudo apt update && sudo apt install ffmpegWindows: 从 FFmpeg 官网 下载并配置环境变量。验证安装: 在终端运行 ffmpeg -version📦 1. 获取代码Bash# 克隆灵果AI仓库
git clone https://github.com/YourUsername/LingGuo-AI.git
cd LingGuo-AI
⚙️ 2. 后端部署 (Server Setup)后端包含 RESTful API 接口与基于 Asynq 的异步视频生成/合并任务引擎。Bash# 进入后端目录
cd server

# 1. 复制环境变量配置文件
cp .env.example .env
编辑 .env 文件，配置您的数据库、Redis 和大模型 API Key：注意：请先在您的 MySQL 中手动创建一个名为 spirit_fruit 的空数据库 (CREATE DATABASE spirit_fruit CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;)。Ini, TOML# 数据库配置
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


启动后端服务：Bash# 下载 Go 依赖
go mod tidy

# 启动服务
go run main.go serve
💡 自动化建表提示：项目内置了 GORM 自动迁移（AutoMigrate）。只要数据库连接成功，首次启动时会自动创建所有数据表，并播种默认的管理员账号与系统菜单。默认超级管理员账号：admin默认超级管理员密码：123456🎨 3. 前端部署 (Web Setup)请打开一个新的终端窗口，进入前端目录：Bash# 进入前端目录
cd web

# 1. 安装前端依赖 (推荐使用 pnpm，如果没有请先 npm install -g pnpm)
npm install

# 2. 启动开发服务器
npm run dev
启动成功后，浏览器访问终端提示的地址（通常是 http://localhost:5173）。输入刚才自动生成的默认账号 admin 和密码 123456，即可开启您的 AI 短剧创作之旅！

🤝 参与贡献 (Contributing)
灵果AI 是一个由社区驱动的开源项目。我们非常欢迎各位开发者、提示词工程师和视频创作者加入我们！
如果你有任何好点子，欢迎提交 Pull Requests 或发布 Issues。

请在贡献前阅读我们的 贡献指南 (CONTRIBUTING.md)。

📄 开源协议 (License)
本项目基于 MIT License 开源。欢迎免费用于商业闭源项目，但请保留原作者署名。