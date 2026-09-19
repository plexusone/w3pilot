# CLAUDE.md — w3pilot

Go browser-automation library using WebDriver BiDi (+ CDP) for AI-assisted
automation: SDK, MCP server, CLI, script runner, session recording, and
`pagecapture`.

## Build & test

```bash
go build ./...
go test ./...
```

Browser-driving tests need Chrome/Chromium; `go build ./...` compiles without it.

## pagecapture subpackage

`pagecapture` captures normalized page evidence (screenshot, rendered HTML,
visible text, title/lang/links/headings) for reuse across accessibility,
functional/journey, and i18n consumers. See `docs/guide/page-capture.md`.

Key design point — **SPA-aware settling**: the default `SettleDOMStable` polls
the DOM until it stops changing *and* holds real content, gated by a
minimum-settle floor. Do not revert to waiting on the `load` event or network
idle; client-rendered apps never reliably reach those, and an early settle
captures a half-rendered shell (this was a real bug — the floor fixes it).

`CaptureCurrent` captures the already-loaded page without navigating (for
callers that already settled it, e.g. an audit engine) — avoids a redundant
navigation. `Capture` navigates+settles+captures; `CaptureSite` crawls same-host
`<a href>` links (note: SPA button/router navigation isn't followable this way —
that needs interaction, i.e. a journey runner).

## Conventions

- `pagecapture` depends only on the w3pilot root package + stdlib; keep it that
  way so consumers of the root aren't burdened.
- BiDi is the primary protocol (the standards bet); CDP is for
  profiling/debugging.
