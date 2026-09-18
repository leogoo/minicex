# NOTES — 教学偏好（务必在每课遵守）

1. **内容 spine = Web3 全栈进阶 0→8**（用户指定，不可偏离）：0 巩固 CEX → 1 区块链基础 → 2 Go 链上交互 → 3 钱包与密钥 → 4 充值 Worker → 5 提现 Pipeline → 6 链上对账 → 7 前端 Web3 → 8 端到端闭环。
2. **Go 从 0，但基础补强融入 Phase 0**，不另设独立「学 Go」阶段。
3. **先图后码**：每课先给流程图（组件 / 数据流），再讲思考，最后引导敲码。
4. **不堆大块代码**：代码分小步给出，每步给目标 + 提示，不让学习者陷入「抄代码」。重点是理解整体思考过程。
5. **分模块**：按 Phase + 模块组织，每模块单一职责，循序渐进。
6. **核心是跑通整体流程**：端到端最小闭环优先于局部深度。
7. **HTTP 框架用 Gin**（已确认）。
8. **生产扩展（C 阶段）在 B 全跑完后再进**（已确认）：幂等 / 风控 / Redis / Kafka / 链上安全专题。每课末尾可放「生产 reality」可选框，但深讲留到 C。
9. **工具最小逻辑**：Redis / Kafka 等只在最简用法出现，不转移课程重点。
10. **前端类比**：善用前端概念类比（JWT≈路由守卫、goroutine≈async、限流≈throttle、幂等≈请求去重、Gin handler≈controller、middleware≈Express/Koa 中间件）帮助理解。
11. **语言 / 格式**：中文为主；每课一份独立 HTML 存 `lessons/`（干净排版，Tufte 风格）。
12. **全新项目**：从零搭一个最小 Web3 交易所后端，不复用此前任何课程 / 项目代码。
13. **骨架不得悬空引用（重点）**：凡是课程骨架里出现的调用，必须在同一课给出被调用方的**方法定义 / 签名**（内部未导出方法同样要求）。方法体可以留白让学习者填，但签名与用途必须写明——否则学习者会卡在"这方法哪来的"。引出未定义的方法 = 课程 bug，需当场补齐（如第 4 课补 `Book.Snapshot()`）。

---

## 目录结构约定（模块式，用户确认）

- 每个业务域一个 `internal/<module>/` 包，**自带 model.go + service.go + handler.go（同 package）**，并暴露 `RegisterRoutes(r, svc)` 自挂载路由。
- `cmd/server/main.go` 只做装配：建引擎 → 各模块 `RegisterRoutes` → Run，不写逻辑。
- `internal/server` 留给「引擎创建 + 全局中间件（JWT/auth）」，**不装具体 handler**。
- 禁止 `internal/server/handler.go` 这种按技术角色集中堆放——多模块后会退化成上帝文件。按模块就近放才能 scale（account / match / deposit / withdraw / reconcile）。
