# TODO

## C 端架构（app / app-plugin）

- [ ] 在文档中固定术语：**package** = 后台代码生成包；**plugin** = 管理端扩展（Casbin/菜单）；**app** = C 端核心；**app-plugin** = 仅挂 `/app/v1`、仅用 `AppJWT` 的 C 端扩展。
- [ ] 约定目录：C 端业务模块使用 `app_` 前缀或独立 `server/app/`、`server/appplugin/`，与 `server/plugin/`（后台插件）区分。
- [ ] （可选）将 `plugin/appclient` 抽到 `server/app/` 或 `server/app/core`，当 C 端模块增多时再重构。
- [ ] （可选）代码生成/脚手架增加 `target: admin | app`，生成路径与默认中间件（Casbin vs AppJWT）分离。

## 其它

- [ ] 生产环境更换 `config.yaml` 中 `app-jwt.signing-key` 为强随机密钥。


```
在概念上把「管理端扩展」和「C 端扩展」分开，一般会清晰很多；是否立刻拆成独立目录，可以按团队规模再定。

现有两套在说什么
概念	典型含义（你们项目里）
package	代码生成用的业务包（api/v1/<pkg> / service / router），偏 B 端 CRUD
plugin	可选安装、可带 initialize（菜单/API/字典）、走 JWT + Casbin 的 管理端能力
C 端（app） 和它们不同：身份是 AppUser、鉴权是 app-jwt、不进 Casbin / 不进后台菜单。若仍统称「plugin」，新人容易以为也要配 sys_api、角色权限。

「app / app-plugin」单独设一层，有什么好处
心智模型：package/plugin = 后台；app / app-plugin = 终端用户。
代码生成 / 脚手架：以后可以加 target: admin | app，生成到不同路径、默认不加 Casbin。
依赖边界：C 端模块禁止依赖 SysUser、菜单、Casbin；检查脚本也好写。
推荐怎么落地（由轻到重）
轻量（先够用）

文档里固定用语：「C 端插件」 = 只挂 /app/v1、只用 AppJWT 的扩展；「后台插件」 = 现有 plugin。
目录仍可在 plugin/appclient 下用子包：plugin/appclient/xxx，或 plugin/app_order 这类 前缀约定（app_ 表示 C 端）。
加重（模块多了再做）

顶层目录区分，例如：
server/app/：C 端公共层（认证、中间件、路由前缀）
server/app/plugin/ 或 server/appplugin/：只注册 C 端路由的「app 插件」
后台继续：server/plugin/ + 自动生成 package。
结论
是的，除了 package 和 plugin 之外，单独建立「app（C 端核心）」和「app-plugin（C 端可插拔业务）」的概念会更好，主要是命名与边界，不必一开始就大搬迁；当前把 C 端放在 plugin/appclient 作为第一个 app 核心是合理的，等 C 端功能变多，再抽到 server/app/ 或统一 app_* 前缀即可。
```