# GitHub Copilot Custom Instructions

## Primary Context & Single Source of Truth
All AI Copilot and Agent behaviors MUST be strictly governed by [`AGENTS.md`](../AGENTS.md).

## Rules Summary
1. **TDD First**: Generate failing unit tests prior to feature implementation code.
2. **Local Sandbox Check**: Validate local build and tests (`make test`) before committing.
3. **Atomic Micro-Commits**: Keep commits single-purpose and adhere to Conventional Commits.
