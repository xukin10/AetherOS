# AetherOS Execution Contract Layer v1.0

## 1. System Overview

AetherOS uses a multi-layer AI code generation pipeline:

Architecture Layer → Orchestration Layer → Execution Layer → GitHub

- Architecture Layer: ChatGPT (design authority)
- Orchestration Layer: Codex (task decomposition & control)
- Execution Layer: DeepSeek API (code generation engine)

---

## 2. Authority Hierarchy (NON-OVERRIDABLE)

Highest → Lowest:

1. ChatGPT (Architecture Authority)
2. Existing Repository Code
3. Codex Orchestration Logic
4. DeepSeek Code Generation

DeepSeek cannot override architecture decisions.

---

## 3. Execution Flow

All code generation MUST follow this pipeline:

1. ChatGPT produces RFC / SPEC
2. Codex converts SPEC into structured prompt
3. Codex sends prompt to DeepSeek API
4. DeepSeek returns code only
5. Codex validates structure
6. Output is committed to GitHub

NO direct AI-to-repo writes allowed.

---

## 4. DeepSeek Role Definition

DeepSeek is ONLY a deterministic code generator.

Allowed:
- Generate Go code
- Generate unit tests
- Follow given interfaces strictly

Forbidden:
- Architecture design
- System modification
- Adding modules
- Changing interfaces
- Introducing new patterns

---

## 5. Codex Role Definition

Codex is responsible for:

- Translating RFC → structured prompt
- Ensuring architecture constraints are included
- Validating DeepSeek output format
- Rejecting invalid responses
- Maintaining repo consistency

---

## 6. Input Contract (Codex → DeepSeek)

Codex MUST provide:

- File structure
- Interfaces
- Function signatures
- Constraints
- Forbidden rules

No ambiguity allowed.

---

## 7. Output Contract (DeepSeek → Codex)

DeepSeek MUST return:

- Pure Go code only
- No explanations
- No alternative designs
- No architectural commentary

---

## 8. Determinism Rule

Same input prompt MUST produce structurally consistent output.

If ambiguous:
→ DeepSeek must choose minimal valid implementation
→ NOT invent architecture

---

## 9. Failure Handling

If DeepSeek output violates contract:

Codex must:
- reject output
- regenerate with stricter prompt

---

## 10. System Stability Rule

AetherOS architecture must remain stable across all AI agents.

No agent is allowed to:

- introduce new core subsystems
- modify runtime lifecycle model
- change dependency model

without RFC approval from ChatGPT.

---

## 11. Golden Rule

AetherOS is NOT a multi-agent experimental system.

It is a controlled AI software engineering pipeline.
