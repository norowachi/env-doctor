# 🩺 env-doctor

Everyone has a `.env.example`, nobody enforces it.  
`env-doctor` is a zero-dependency CLI that makes your `.env` a happy first-class citizen

> Validate your `.env` files against `.env.example`
> 
> Catch missing keys, wrong types, and placeholder values before they cause bugs.

---

## Features

- [x] Detect missing keys
- [x] Detect empty values
- [x] Flagging placeholder values (`changeme`, `xxx`, `secret`, etc.)
- [x] Type validation via comment annotations (`url`, `number`, `boolean`, `email`)
- [x] Custom regex type (`@type: /^[A-Z]{3}$/`)
- [x] Required field enforcement
- [x] Non-zero exit code on errors (CI-friendly)
- [x] JSON output for CI pipelines
- [x] Zero external dependencies

---

## Installation

```bash
# Build from source
git clone https://github.com/norowachi/env-doctor
cd env-doctor
go build ./cmd/env-doctor

# Move to PATH
mv env-doctor /usr/local/bin/
```

---

## Usage

```bash
# Basic check, compares .env.example and .env in the current dir
env-doctor

# Custom paths
env-doctor --example .env.schema --env .env.production

# JSON output for Scripts / CI
env-doctor --json
```

---

## Schema Annotations

Add annotations as comments directly above a key in `.env.example`:

```bash
# @required
# @type: url
# @desc: PostgreSQL connection string
DATABASE_URL=

# @required
# @type: number
PORT=3000

# @type: boolean
DEBUG=false

# @type: email
ADMIN_EMAIL=
```
Treated annotations: `@required` and `@type`


### Supported types

| Type | Validates |
|------|-----------|
| `url` | Well-formated URL with scheme and host |
| `number` / `int` | Integer value |
| `boolean` / `bool` | `true`, `false` / `1`, `0` / `yes`, `no` |
| `email` | Valid email address format |
| `/pattern/` | Custom regex pattern |

#### Regex examples

```bash
# AWS region: us-east-1, eu-west-2, etc.
# @type: /^[a-z]+-[a-z]+-[0-9]+$/
AWS_REGION=

# Semver: v1.2.3
# @type: /^v[0-9]+\.[0-9]+\.[0-9]+$/
APP_VERSION=

# ISO currency code: USD, EUR, GBP
# @type: /^[A-Z]{2,3}$/
CURRENCY=

# UUID v4
# @type: /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/
TENANT_ID=
```

---

## CI Integration

### GitHub Actions

```yaml
- name: Validate environment
  run: |
    env-doctor --env .env.ci --json
```

### Pre-commit hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/sh
env-doctor || exit 1
```

```bash
chmod +x .git/hooks/pre-commit
```

---

## Example output

```
🩺 env-doctor  →  .env.example vs .env

───────────────────────────────────────────────────────
✗ ERROR    DATABASE_URL                        Expected a valid URL, got: "not-a-valid-url"
⚠ WARN     APP_SECRET                          Value looks like a placeholder: "changeme"
✗ ERROR    PORT                                Expected a number, got: "abc"
✗ ERROR    DEBUG                               Expected a boolean (true/false/1/0), got: "maybe"
✗ ERROR    ADMIN_EMAIL                         Expected a valid email address, got: "not-an-email"
⚠ WARN     APP_VERSION                         Key is missing from your .env file
✗ ERROR    CURRENCY                            Key is missing from your .env file
⚠ WARN     STRIPE_KEY                          Value is empty
⚬ INFO     EXTRA_KEY                           Key exists in .env but not in .env.example
───────────────────────────────────────────────────────
  4 error(s)  2 warning(s)  1 info
```

---

## Roadmap

- [ ] `.env-doctor.yml` config file for project-level settings
- [ ] Homebrew tap & apt/aur package
- [ ] VS Code extension (inline squiggles in `.env` files)
- [ ] Multi-env checks support (`.env.staging`, `.env.production`)

---

## Contributing

PRs are welcome!
> Please open an issue first for large changes.

```bash
git clone https://github.com/norowachi/env-doctor
cd env-doctor
go test ./...
go build ./cmd/env-doctor
```

## License

MIT
