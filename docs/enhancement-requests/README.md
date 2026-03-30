# Enhancement Requests

This directory contains enhancement requests for W3Pilot based on real-world usage feedback.

## Status Legend

- **Proposed**: Initial draft, under consideration
- **Accepted**: Approved for implementation
- **In Progress**: Currently being implemented
- **Complete**: Implemented and released
- **Declined**: Not accepted (with rationale)

## Current Requests

| ID | Title | Priority | Status | Version |
|----|-------|----------|--------|---------|
| [mcp-enhancements-2026-03-29](mcp-enhancements-2026-03-29.md) | MCP Server Enhancements | Mixed | Complete | v0.8.0 |

## Naming Convention

Files follow the pattern: `<topic>-YYYY-MM-DD.md`

- Date indicates when the request was received/created
- Topic should be descriptive but concise

## Historical Note

This directory previously contained enhancement requests for upstream dependencies (VibiumDev/clicker).
Those were deleted after W3Pilot adopted a dual-protocol architecture (BiDi + CDP), making upstream
CDP passthrough unnecessary. W3Pilot connects directly to Chrome's CDP endpoint for profiling features.
