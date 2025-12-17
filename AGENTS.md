# Qwen Code Project Context

## Project Overview

go2rtc 是一个用 Go 语言编写的多功能、低延迟的摄像头流媒体应用程序。
它支持多种流媒体协议（如 RTSP、WebRTC、HomeKit、FFmpeg、RTMP 等），
旨在作为独立应用或与其他智能家庭平台（如 Home Assistant）集成。

主要功能包括：

- 支持多种输入源（RTSP、RTMP、HTTP、USB 摄像头、FFmpeg 等）。
- 支持多种输出流（WebRTC、RTSP、MSE/MP4、HomeKit、HLS、MJPEG 等）。
- 实现动态编解码器协商和转码（通过 FFmpeg）。
- 支持双向音频流（针对特定摄像头）。
- 零依赖、零配置的二进制文件（支持多平台）。
- 集成 Web API，支持自定义 Web 界面。
- 模块化架构，支持高度可扩展性。

项目使用 Go 1.24.0，依赖 Pion WebRTC 库及其他第三方库。

## Project Structure

- `main.go`: 主入口点，负责初始化所有模块。
- `go.mod` / `go.sum`: Go 依赖管理文件。
- `internal/`: 核心内部模块实现，每个目录对应一个功能模块（如 `rtsp`, `webrtc`, `streams` 等）。
- `pkg/`: 可重用的包，包含处理不同流格式和协议的代码，以及核心功能（如 `core`, `shell` 等）。
- `api/`: HTTP API 实现，用于与外部系统交互。
- `assets/`: 项目资源，如徽标图片。
- `docker/`: Docker 配置文件。
- `examples/`: 配置示例。
- `www/`: Web 用户界面前端代码。

## Building and Running

### Build from Source

1. 确保已安装 Go 1.24.0 或更高版本。
2. 克隆仓库。
3. 在项目根目录运行以下命令进行构建：

    ```sh
    go build -o go2rtc main.go
    ```

### Run

1. 构建完成后，使用以下命令运行应用程序：
    ```sh
    ./go2rtc -c go2rtc.yaml
    ```
    默认情况下，它会查找当前工作目录下的 `go2rtc.yaml` 配置文件。
    Web UI 默认运行在 `http://localhost:1984/`。

2. 可通过 `-c` 或 `--config` 标志指定配置文件路径。

### Docker

也可以使用 Docker 运行：

```sh
docker run -it --rm -p 1984:1984 -p 8554:8554 -p 8555:8555 alexxit/go2rtc
```

## Development Conventions

- Go 代码风格遵循 Go 官方规范。
- 模块化设计：功能被拆分到 `internal/` 和 `pkg/` 下的独立包中。
- 配置使用 YAML 格式。
- 代码结构清晰，每个协议或功能有对应的目录和文件。
- 通过 `pkg/README.md` 可以了解各种输入输出格式、协议和编解码器的命名约定。