# Vibe Music Server

一个为 Vibe Music 音乐流媒体应用提供支持的后端服务器。

## 🚀 简介

Vibe Music Server 是一个使用 Go 语言构建的高性能后端服务。它为音乐播放、用户管理、播放列表等功能提供 RESTful API。

## ✨ 功能特性

-   🔐 **用户认证**: 安全的用户注册和登录 (JWT)。
-   🎵 **音乐库管理**: 管理歌曲、歌手和专辑信息。
-   🎶 **播放列表**: 创建、查看、更新和删除个人播放列表。
-   🎧 **音乐流**: 高效的音乐文件流式传输。
-   🔍 **搜索**: 按歌曲、歌手或专辑进行搜索。

## 🛠️ 技术栈

-   **语言**: [Go](https://golang.org/)
-   **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
-   **数据库**: [PostgreSQL](https://www.postgresql.org/)
-   **ORM**: [GORM](https://gorm.io/)
-   **缓存**: [Redis](https://redis.io/)
-   **配置**: [Viper](https://github.com/spf13/viper)
-   **认证**: [JWT (JSON Web Tokens)](https://jwt.io/)

## 📋 环境准备

在开始之前，请确保您已安装以下软件：

-   [Go](https://golang.org/dl/) (版本 1.18 或更高)
-   [PostgreSQL](https://www.postgresql.org/download/)
-   [Redis](https://redis.io/download)
-   [Docker](https://www.docker.com/get-started) (可选, 用于快速启动数据库)

## ⚙️ 安装与运行

1.  **克隆仓库**
    ```bash
    git clone https://github.com/your-username/vibe-music-server.git
    cd vibe-music-server
    ```

2.  **安装依赖**
    ```bash
    go mod tidy
    ```

3.  **修改配置文件**
    
    项目使用 `config/templateConfig.yml` 作为配置模板。建议您在 `config` 目录下创建一个 `config.yml` 文件来覆盖默认配置并配置您的敏感信息（如数据库密码、密钥等）。
    
    `config.yml` 的配置项会覆盖 `templateConfig.yml` 中的同名配置项。
    
    下面是一个 `config.yml` 的示例，您可以根据实际情况进行修改：
    ```yaml
    # config/config.yml

    # 配置数据源
    database:
      host: 127.0.0.1
      port: 3306
      name: vibe_music
      username: root
      password: YOUR_DB_PASSWORD

    # Redis服务连接配置
    redis:
      host: 127.0.0.1
      port: 6379
      password: YOUR_REDIS_PASSWORD

    minio:
      endpoint: http://127.0.0.1:9090
      accessKey: minioadmin
      secretKey: minioadmin
      bucket: vibe-music

    # 配置邮件服务
    mail:
      host: YOUR_SMTP_HOST
      port: 465
      user: YOUR_EMAIL_ADDRESS
      password: YOUR_EMAIL_PASSWORD

    # JWT 密钥
    jwt:
      secret: YOUR_VERY_SECRET_JWT_KEY
    ```
    **注意**: `config.yml` 已被添加到 `.gitignore` 中，以避免将敏感信息提交到版本控制系统。

4.  **运行数据库迁移** (如果使用 GORM)
    ```bash
    # 您可能需要一个迁移命令或在应用启动时自动迁移
    go run ./cmd/migrate
    ```

5.  **启动服务**
    ```bash
    go run ./cmd/server/main.go
    ```
    服务器将在 `http://localhost:8080` 上运行。

## 📡 API 端点

以下是 API 的一些主要端点示例，更多详情请参阅 [API 文档](./docs/vibe-music.openapi.json)。

### 用户 (`/user`)
-   `POST /user/register`: 注册新用户
-   `POST /user/login`: 用户登录
-   `GET /user/getUserInfo`: 获取当前用户信息 (需要认证)
-   `PUT /user/updateUserInfo`: 更新用户信息 (需要认证)

### 歌曲 (`/song`)
-   `POST /song/getAllSongs`: 获取歌曲列表（支持分页和搜索）
-   `GET /song/getRecommendedSongs`: 获取推荐歌曲
-   `GET /song/getSongDetail/{id}`: 获取单首歌曲详情

### 歌手 (`/artist`)
-   `POST /artist/getAllArtists`: 获取歌手列表（支持分页和搜索）
-   `GET /artist/getRandomArtists`: 获取随机歌手
-   `GET /artist/getArtistDetail/{id}`: 获取歌手详情及其歌曲

### 歌单 (`/playlist`)
-   `POST /playlist/getAllPlaylists`: 获取歌单列表（支持分页和搜索）
-   `GET /playlist/getRecommendedPlaylists`: 获取推荐歌单
-   `GET /playlist/getPlaylistDetail/{id}`: 获取歌单详情

### 收藏 (`/favorite`)
-   `POST /favorite/collectSong`: 收藏歌曲 (需要认证)
-   `DELETE /favorite/cancelCollectSong`: 取消收藏歌曲 (需要认证)
-   `POST /favorite/getFavoriteSongs`: 获取收藏的歌曲列表 (需要认证)
-   `POST /favorite/collectPlaylist`: 收藏歌单 (需要认证)
-   `DELETE /favorite/cancelCollectPlaylist`: 取消收藏歌单 (需要认证)
-   `POST /favorite/getFavoritePlaylists`: 获取收藏的歌单列表 (需要认证)

### 评论 (`/comment`)
-   `POST /comment/addSongComment`: 新增歌曲评论 (需要认证)
-   `POST /comment/addPlaylistComment`: 新增歌单评论 (需要认证)
-   `PATCH /comment/likeComment/{id}`: 点赞评论 (需要认证)

## 🤝 贡献

欢迎各种形式的贡献！如果您想为这个项目做出贡献，请遵循以下步骤：

1.  Fork 本仓库
2.  创建您的分支 (`git checkout -b feature/AmazingFeature`)
3.  提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4.  推送到分支 (`git push origin feature/AmazingFeature`)
5.  提交一个 Pull Request

## 📄 许可证

该项目使用 MIT 许可证。有关详细信息，请参阅 [`LICENSE`](./LICENSE) 文件。
