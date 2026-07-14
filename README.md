# BlctekIP IP-Chain

AI训练数据登记存证、合规辅助审查与受控交易平台。

## 当前施工状态

本仓库正在按照 PRD v7.1-r1 与 CODE PLAN v1.0 从零建设。
当前分支完成 Iteration 1 的第一批工程基线：

- Go API 与 Worker 可编译骨架；
- Request ID、统一问题响应、幂等键校验与 Outbox 领域模型；
- MySQL 可靠性基础表；
- Next.js 前端目录骨架；
- 本地 MySQL、Redis、MinIO 编排；
- CI、ADR 和当前状态文档。

## 本地验证

```bash
cp .env.example .env
make check
make run-api
```

健康检查：

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

## 施工规则

- 不直接在 `main` 持续开发；
- 每个迭代通过独立分支和 Draft PR 交付；
- 支付、结算、账务、权限和加密代码必须有第二次审查；
- 生产环境禁止启用 Mock Provider；
- 已执行的数据库迁移不得修改。
