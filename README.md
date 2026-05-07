# Port Monitor

开箱即用的本机端口监控小工具。Go 后端 + Vue 3 前端，前端构建产物用 `embed` 内嵌进单一可执行文件，双击即用。

- 默认监听 `0.0.0.0:7000`
- 默认账号 `admin / admin`
- 默认每 1s 扫描一次本机端口（前端可改，最小 200ms）
- 实时通过 WebSocket 推送

## 直接使用

构建产物位于仓库根目录：`port-monitor.exe`（Windows）。双击运行后浏览器打开 <http://localhost:7000>，登录 `admin / admin`。

首次启动会在 exe 同目录写出 `config.json`，可手动修改：

```json
{
  "listenAddr": "0.0.0.0:7000",
  "intervalMs": 1000,
  "highRiskPorts": [21, 22, 23, 135, 137, 138, 139, 445, 3389,
                    3306, 5432, 6379, 27017, 11211, 9200, 2375, 2376]
}
```

`listenAddr` 改完需重启；`intervalMs` 与 `highRiskPorts` 也可在前端"设置"对话框修改并立即生效。

## 项目结构

```
port-monitor/
├── backend/        Go 服务，gopsutil 扫端口、gorilla/websocket 推送、go:embed 内嵌前端
├── frontend/       Vue 3 + Vite + Element Plus
├── assets/         设计参考图（home.png / login.png）
├── build.ps1       Windows 一键构建脚本
├── build.sh        Linux 一键构建脚本
└── port-monitor.exe  构建产物
```

## 开发

后端：

```powershell
cd backend
go run .
```

前端（另开一个窗口，Vite 代理 `/api` 与 `/ws` 到 `:7000`）：

```powershell
cd frontend
npm install
npm run dev
```

浏览器打开 <http://localhost:5173>。

## 生产构建

Windows：

```powershell
.\build.ps1
```

Linux / macOS：

```bash
./build.sh
```

脚本流程：`npm run build` → 把 `frontend/dist` 复制到 `backend/web/dist` → `go build` 出单一二进制。

## Linux 生产启动

Linux 构建完成后，根目录会生成 `port-monitor` 二进制文件。首次启动会在二进制同目录生成 `config.json`，默认监听 `0.0.0.0:7000`。

前台验证：

```bash
chmod +x ./port-monitor
./port-monitor
```

浏览器打开 `http://服务器IP:7000`，使用 `admin / admin` 登录。

临时后台运行：

```bash
nohup ./port-monitor > port-monitor.log 2>&1 &
```

推荐生产环境使用 systemd 托管。示例安装目录为 `/opt/port-monitor`：

```bash
sudo mkdir -p /opt/port-monitor
sudo cp ./port-monitor /opt/port-monitor/
sudo chmod +x /opt/port-monitor/port-monitor
```

创建 `/etc/systemd/system/port-monitor.service`：

```ini
[Unit]
Description=Port Monitor
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/port-monitor
ExecStart=/opt/port-monitor/port-monitor
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

启动并设置开机自启：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now port-monitor
sudo systemctl status port-monitor
```

查看日志：

```bash
journalctl -u port-monitor -f
```

如需修改端口，编辑 `/opt/port-monitor/config.json` 里的 `listenAddr`，然后重启服务：

```bash
sudo systemctl restart port-monitor
```

## 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/login` | body `{username, password}` → `{token, expiresAt}` |
| GET  | `/api/ports` | 当前快照（带 Bearer token） |
| GET  | `/api/settings` | 当前配置 |
| PUT  | `/api/settings` | body `{intervalMs, highRiskPorts}`，立即生效并写盘 |
| GET  | `/api/ws` | WebSocket，query `?token=...`；连接时下发当前快照，之后随每次 tick 推送 `{type:"ports", ports, timestamp}` |

JWT 用 HS256，密钥进程启动时随机生成（重启后旧 token 失效）。

## 关于 Windows 下 "进程" 列显示 "-"

Windows 下读取系统/其他用户的进程信息需要管理员权限。以管理员身份运行 `port-monitor.exe` 即可显示完整进程名与启动时间，否则部分 PID 会显示 "-"。
