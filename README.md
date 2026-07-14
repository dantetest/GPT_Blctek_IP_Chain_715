# BlctekIP IP-Chain

AI训练数据登记存证、合规辅助审查与受控交易平台。

## 当前施工状态

本仓库正在按照 PRD v7.1-r1 与 CODE PLAN v1.0 从零建设。
当前Draft PR已完成：

- Go API与Worker工程基线；
- Request ID、问题响应、幂等与Outbox基础；
- MySQL可靠性基础表；
- Dataset与不可变Dataset Version领域模型；
- Manifest v1流式哈希、NFC路径规范和Merkle Root；
- 静态Manifest测试向量与三操作系统CI；
- Next.js前端目录骨架；
- MySQL、Redis、MinIO本地编排；
- ADR与数据模型文档。

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

Manifest规范：`docs/MANIFEST_SPEC_V1.md`。

## 施工规则

- 不直接在 `main` 持续开发；
- 每个迭代通过独立分支和Draft PR交付；
- 支付、结算、账务、权限和加密代码必须有第二次审查；
- 生产环境禁止启用Mock Provider；
- 已执行的数据库迁移不得修改；
- 已发布Dataset Version不得覆盖修改。
