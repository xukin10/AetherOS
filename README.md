# AetherOS

**Enterprise AI Control Plane**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-0.1.0--alpha-green.svg)](VERSION)

## What is AetherOS?

AetherOS is an open-source **Enterprise AI Control Plane** for orchestrating AI Agents, Workflows, Knowledge, Tools, and enterprise systems.

**Think of it as: Kubernetes for AI.**

## Vision

- Enterprise R&D Teams
- Enterprise Operations Teams
- AI-Powered Development
- Workflow Marketplace
- Plugin Marketplace
- MCP Support
- Multi-Model Support
- Enterprise Private Deployment

## Design Principles

| Principle | Description |
|-----------|-------------|
| **P1: Core First** | Core is always stable. Plugins extend infinitely. |
| **P2: Everything is Plugin** | Agents, Workflows, Tools, Policies, Knowledge - all plugins. |
| **P3: Workflow over Agent** | Workflow defines flow. Agent is just executor. |
| **P4: Event Driven** | All modules communicate via Events. No direct calls. |
| **P5: Controller Pattern** | Controller for business logic. Runtime for execution. |
| **P6: DDD** | Domain-Driven Design. No DB, HTTP, or Redis in domain. |
| **P7: Hexagonal Architecture** | Business logic independent of DB, LLM, HTTP. |

## Tech Stack

| Module | Technology |
|--------|------------|
| Language | Go 1.24+ |
| API | Chi |
| Database | PostgreSQL |
| Cache | Redis |
| Event Bus | NATS |
| Vector DB | Qdrant |
| Object Storage | MinIO |
| UI | Next.js + React |
| SDK | Go / Python / TypeScript |
| Observability | OpenTelemetry |
| Metrics | Prometheus |
| Dashboard | Grafana |
| Deploy | Docker + Helm |

## Project Structure

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed project structure.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.