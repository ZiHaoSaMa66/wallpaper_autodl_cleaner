# wp-cleaner

> 清理 Steam Wallpaper Engine 跨账号登录时自动下载的壁纸

[English](README.md) 中文

## 使用方法

### 1. 获取 Steam Web API Key
前往 https://steamcommunity.com/dev/apikey 注册一个密钥
（免费，任意域名都行）

### 2. 运行工具预览
```cmd
wp-cleaner.exe -api-key=你的KEY -dry-run
```

此命令将：
- 自动检测 Steam 安装路径
- 识别当前登录的 Steam 用户
- 扫描 workshop/content/431960/ 下所有已下载壁纸
- 从 Steam 获取壁纸元数据（标题、作者等）
- 对比你的订阅列表
- 显示哪些壁纸可以安全清理

### 3. 执行清理
```cmd
wp-cleaner.exe -api-key=你的KEY -dry-run=false
```

将所有非当前用户订阅的壁纸移至隐藏备份文件夹（`.trash-*`）。

### 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-api-key` | `""` | Steam Web API Key（订阅检测必需） |
| `-steam-id` | `""` | SteamID64（留空自动从 loginusers.vdf 检测） |
| `-dry-run` | `true` | 预览模式，不实际执行清理 |
| `-force` | `false` | 跳过确认提示 |

## 工作原理

1. 从 Windows 注册表读取 Steam 安装路径
2. 解析 `loginusers.vdf` 获取本机所有 Steam 账户
3. 扫描 `steamapps/workshop/content/431960/` 已下载壁纸文件夹
4. 调用 `IPublishedFileService/GetUserFiles?type=mysubscriptions` 获取**你的**订阅列表
5. 调用 `GetPublishedFileDetails`（公开 API）获取壁纸标题
6. 对比列表，找出**未订阅**的壁纸
7. 将未订阅文件夹重命名为 `.trash-*` 前缀，以便安全审查/删除

## 从源码构建

```bash
go build -o wp-cleaner.exe .
```

环境要求：Go 1.22+, Windows 10/11
