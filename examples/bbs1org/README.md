# bbs1org

一个极简 PHP 论坛。由一个仅100多KB大小的PHP文件构建。纯原生、无框架、无依赖，支持 SQLite、MySQL 和 PostgreSQL。适合社区站点、低成本部署和 AI 二次开发。

## 特点

- 纯原生 PHP，单核心代码文件，无框架和 Composer 依赖，部署与维护简单
- 支持 SQLite、MySQL 和 PostgreSQL，数据库结构和搜索能力保持跨引擎兼容
- 包含首页、版块、主题、回帖、收藏、个人主页和后台管理等完整论坛功能
- 支持用户组、版块权限、站点设置、注册控制、发帖限制和附件管理
- 插件机制支持 Hook、路由、后台页面、资源合并、在线安装更新和统一计划任务
- 站点设置、版块和用户组按需懒加载，数据库结构简单，负载能力强
- 支持 AJAX 交互和响应式布局，兼顾 PC 与移动端使用体验

## 环境

- PHP 8.1+
- PDO SQLite、PDO MySQL 或 PDO PostgreSQL 扩展

## 演示

https://bbs1.org

## Docker 源码部署

服务器需先安装 Docker `curl -fsSL https://get.docker.com -o install-docker.sh && sudo sh install-docker.sh`

源码仓库和 Docker 配置仓库应放在同一目录下。

```bash
git clone https://github.com/bbs1org/bbs1org.git
git clone https://github.com/bbs1org/bbs1org_docker.git
cd bbs1org_docker
cp .env.example .env
# 如需修改端口或数据库模式，请先编辑 .env
docker compose up -d
```

容器启动后访问默认8080端口：

```text
http://服务器地址:8080
```

首次访问会进入网页安装程序，请在页面中设置站点、默认版块和管理员账号。使用 MySQL 或 PostgreSQL 时，先通过 `docker compose ps` 确认数据库状态为 `healthy`。

默认使用 `SQLite`。如需修改 `8080` 端口，或者启用 `MySQL`、`PostgreSQL`，请在启动前修改 `.env`。数据库容器会按 `.env` 中的 `DB_NAME`、`DB_USER`、`DB_PASSWORD` 初始化；网页安装时填写同样的数据库名、用户名和密码。

| 数据库 | `COMPOSE_PROFILES` | 数据库地址 | 端口 |
| --- | --- | --- | --- |
| SQLite | `sqlite` | 无需填写 | 无需填写 |
| MySQL | `mysql` | `mysql` | `3306` |
| PostgreSQL | `pgsql` | `postgres` | `5432` |

常用操作（均在 `bbs1org_docker` 目录执行）：

```bash
docker compose ps                 # 查看状态
docker compose logs -f            # 查看日志
docker compose restart            # 重启
docker compose down               # 停止并保留数据卷
```

## 虚拟机部署（已有 Nginx/Apache + PHP8 环境）

- 下载项目压缩包 `https://github.com/bbs1org/bbs1org/archive/refs/heads/main.zip`
- 解压缩后通过 FTP 上传到网站目录，确保 `index.php` 位于根目录。
- 访问站点域名，根据指示进行安装即可。
- 使用第三方服务 (比如: cron-job.org) 定时请求服务，启用定时任务。
每分钟以 `GET` 方式访问：`https://你的域名/index.php?a=cron`

## 面板部署

宝塔和 1Panel 在面板 终端 执行Docker 源码部署代码即可。

## 在线升级

在后台设置底部点击“升级”，检测更新后选择文件并执行“在线升级”。升级前请先备份数据库、附件、头像和插件目录。

## 数据库转换和迁移

先在新数据库完成安装并登录管理员账号，再从升级页进入“数据迁入”，或访问 `index.php?a=migrate`。选择旧数据库类型并填写连接信息，程序会迁入旧库的全部普通数据表；当前库没有的表会自动复制字段、主键和索引后再导入数据，同名表则清空后替换，并保留原 ID。

插件数据表会一并迁入。附件、头像和插件程序文件不在数据库中，需要另外复制 `app/upload/`、`app/avatars/` 和 `app/plugins/`。迁移前请备份新旧数据库。

## 数据备份

SQLite数据目录 `app/data/ `、附件目录 `app/upload/`、插件目录 `app/plugins/` 需要定期备份。

## AI插件开发指南

推荐使用AI为本项目开发插件。插件开发规范、最小示例、资源、计划任务、数据库与 Hook 说明请阅读 [`.ai-rules.md`](.ai-rules.md)。
本文件用于约束 AI 为 bbs1org 开发或修改插件，规则按优先级排列且必须逐条遵守。
