# Origami Music — Fyne 音乐播放器示例

完整的 Android/iOS/桌面音乐 App，使用 Origami + Fyne 构建。

## 功能

- **用户登录认证** — 账号密码登录，模拟认证服务
- **歌单浏览** — 6 个预设歌单，21 首歌曲
- **歌曲搜索** — 按歌名/歌手搜索
- **音乐播放控制** — 播放、暂停、上一首、下一首
- **播放进度条** — 实时显示播放进度
- **多页面导航** — 登录 → 首页 → 歌单详情 → 正在播放

测试账号：`admin` / `123456` 或 `demo` / `demo`

## 项目结构

```
examples/music-app/
├── app.php                 ← 入口
├── pages/                  ← 页面组件
│   ├── login.php
│   ├── home.php
│   ├── playlist_detail.php
│   └── now_playing.php
├── services/               ← 业务逻辑
│   ├── auth.php
│   ├── music_data.php
│   └── player.php
├── main.go                 ← Go 入口（公共代码）
├── main_desktop.go         ← 桌面端：文件系统加载（!android）
├── main_android.go         ← 安卓端：go:embed 嵌入（android）
└── scripts/
    ├── build_desktop.bat    ← 构建桌面版
    ├── build_apk.bat        ← 构建 APK
    └── gen_icon.go          ← 图标生成器
```

## 构建策略

| 平台 | 加载方式 | 说明 |
|------|----------|------|
| **桌面** | 文件系统 | PHP 源码复制到 `build/php_src/`，可修改后直接重启生效 |
| **安卓** | `go:embed` | PHP 编译进 APK，无需外部文件 |

## 构建

### 桌面版

```bash
cd extensions/fyne
scripts\build_music_desktop.bat
```

构建脚本自动完成：
1. `go build` 编译 Go 二进制文件
2. 复制 `pages/`、`services/`、`app.php` 到 `build/php_src/`

桌面调试时可修改 `build/php_src/` 下的 `.php` 文件，重启 exe 即可生效，无需重新编译。

### APK (arm64)

```bash
cd extensions/fyne\examples\music-app
scripts\build_apk.bat
adb install -r build\music_app.apk
```
