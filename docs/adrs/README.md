# Architecture Decision Records (ADRs)

本目录包含 AetherOS 项目的所有架构决策记录。

## 什么是 ADR？

架构决策记录（ADR）是记录项目中做出的重要技术决策及其背景的文档。每个 ADR 都是原子的（单一决策），并包含：

- **背景** - 为什么需要这个决策
- **决策** - 我们决定做什么
- **后果** - 做这个决策的影响

## 📋 ADR 索引

| # | 标题 | 状态 | 日期 |
|---|------|------|------|
| 001 | [示例：异步事件系统架构](./ADR-001-async-event-system.md) | Proposed | TBD |

## 状态说明

- **Proposed** - 新提案，正在讨论
- **Accepted** - 已批准，可以实施
- **Deprecated** - 仍然有效，但不推荐用于新项目
- **Superseded** - 被更新的 ADR 替代

## 创建新 ADR

1. 复制 `/docs/adrs/ADR-TEMPLATE.md`
2. 命名为 `ADR-XXX-[title].md`（按顺序递增）
3. 填充模板内容
4. 创建 Pull Request
5. 审查后合并到 main

## 查看相关文件

- **RFC 提案** - `/docs/rfcs/`
- **治理指南** - `/docs/GOVERNANCE.md`
- **讨论** - GitHub Discussions tab

---

**如何引用 ADR：**

在代码注释、PR 描述或其他文档中引用 ADR：

```markdown
See ADR-001 for the rationale behind this architecture.
```

---

**最后更新：** 2026-07-01
