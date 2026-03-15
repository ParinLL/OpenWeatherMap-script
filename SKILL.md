---
name: openweathermap-cli
description: Use this skill when the user wants to run, troubleshoot, or extend the owget CLI for geocoding, current weather, and 5-day forecasts with OpenWeatherMap.
homepage: https://github.com/ParinLL/OpenWeatherMap-script
metadata: {"openclaw":{"homepage":"https://github.com/ParinLL/OpenWeatherMap-script","requires":{"env":["OPENWEATHER_API_KEY"],"binaries":["go"]},"primaryEnv":"OPENWEATHER_API_KEY"}}
---

# OpenWeatherMap CLI Skill

Instruction-only skill document for using and troubleshooting `owget` (OpenWeatherMap CLI).

## Skill Purpose and Trigger Scenarios

- The user wants current weather, forecast, or geocoding (`geo`) results.
- The user asks how to run `owget` commands or use parameters.
- The user reports API key issues, HTTP errors, or city lookup failures.

## Installation Command (or GitHub Link to Install Section)

- GitHub: https://github.com/ParinLL/OpenWeatherMap-script

Recommended before installation (to reduce supply-chain risk):

```bash
git clone git@github.com:ParinLL/OpenWeatherMap-script.git
cd OpenWeatherMap-script
git log --oneline -n 5
```

Review recent commits first, and pin to a trusted commit/tag before installing when possible.

Non-privileged install (recommended):

```bash
go install .
```

System-wide install (optional, requires `sudo`):

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o owget .
sudo install owget /usr/local/bin/
```

## Required Environment Variables / Permissions

Required environment variable:

```bash
export OPENWEATHER_API_KEY="your-api-key"
```

- Requires the `go` toolchain for build/install.
- If using `sudo install ... /usr/local/bin/`, system admin privileges are required, and source review should be completed first.
- Never expose full API keys in outputs; debug request logs should redact credential query params (for example, `appid`).

## Common Troubleshooting

- `error: OPENWEATHER_API_KEY env is required`
  - The env var is not set. Run `export OPENWEATHER_API_KEY="..."` first.
- API returns `401`
  - API key is invalid, expired, or mistyped. Re-check your OpenWeatherMap key.
- API returns `404` or city not found
  - Use `City,Country` format (for example, `Taipei,TW`) and verify with `owget geo "<query>"` first.
- Concern about credential leakage while using debug mode
  - Debug request URLs are redacted for sensitive params, but avoid long-running debug in shared/logged environments.
