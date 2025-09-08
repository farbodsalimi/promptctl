## Project Plan: Prompt Vault CLI

### 1. **Core Goals**

- Store prompts in a versioned vault.
- Support prompt templating using Go’s `text/template`.
- Keep metadata (name, versions, timestamps, provider config, etc.) in **SQLite**.
- Provide a simple, intuitive and interactive CLI interface.

---

### 2. **Data Model (SQLite schema)**

**Tables:**

1. **vaults**

   - `id` (PK, UUID or int autoincrement)
   - `name` (unique)
   - `created_at`

2. **prompts**

   - `id` (PK)
   - `vault_id` (FK → vaults.id)
   - `name` (unique within vault)
   - `created_at`

3. **prompt_versions**

   - `id` (PK)
   - `prompt_id` (FK → prompts.id)
   - `version` (int, auto-increment per prompt)
   - `content` (text, Go template string)
   - `created_at`

4. **runs**

   - `id` (PK)
   - `prompt_version_id` (FK → prompt_versions.id)
   - `provider` (enum: openai, google, anthropic)
   - `params` (json: input variables + config like model, temperature, etc.)
   - `response` (text)
   - `created_at`

---

### 3. **CLI Commands**

#### Vault management

- `vault create <name>` — create a new vault.
- `vault list` — list all vaults.
- `vault delete <name>` — delete a vault.

#### Prompt management

- `prompt add --vault=<vault> --name=<name> --prompt=<prompt>` — add a new prompt (store as a template).
- `prompt update --vault=<vault> --name=<name> --prompt=<prompt>` — update an existing prompt → creates new version.
- `prompt list --vault=<vault>` — list prompts in a vault.
- `prompt history --vault=<vault> --name=<name>` — list all versions with timestamps.
- `prompt show --vault=<vault> --name=<name> [--version=N]` — view specific version or latest.

#### Running prompts

- `prompt run --vault=<vault> --name=<name> [--version=N] --provider=openai --model=gpt-4 --vars "key1=val1,key2=val2"`

  - Renders Go template with provided vars.
  - Sends request to chosen provider.
  - Stores the run in DB for reproducibility.

#### Run history

- `run list [--prompt=<name>] [--vault=<vault>]`
- `run show <id>`

---

### 4. **LLM Providers**

- `provider add <name> <api key>` — add a new LLM provider.
- `provider update <name> <api key>` — update a LLM provider.
- `provider list <name> <api key>` — list a LLM provider.
- `provider delete <name> <api key>` — delete a LLM provider.

Providers can be configured via env vars or `~/.promptctl/config.json`.
Example config:

```json
{
  "openai": { "api_key": "sk-..." },
  "anthropic": { "api_key": "sk-..." },
  "google": { "api_key": "..." }
}
```

---

### 5. **Prompt Storage & Rendering**

- Store raw template text in `prompt_versions.content`.
- At runtime:

  ```go
  // read prompt template and it variables from database
  tmpl, _ := template.New("prompt").Parse(content)
  var buf bytes.Buffer
  tmpl.Execute(&buf, vars)
  finalPrompt := buf.String()
  ```

- This ensures flexible prompt parameterization.

---

### 6. **Versioning Strategy**

- Every update → new row in `prompt_versions`.
- Version is incremental, scoped per prompt.
- CLI defaults to **latest version**, but user can request `--version=N`.

---

### 7. **CLI Framework**

- Use **Cobra** (Go’s standard CLI lib).
- Supports subcommands, flags, autocomplete.
- Structure:

  ```
  cmd/
    vault.go
    prompt.go
    run.go
    root.go
  ```

---

### 8. **Future Extensions**

- Git-style diff: show changes between prompt versions.
- Export/import vaults as JSON or YAML.
- Sharing vaults across team (sync to Git repo or cloud DB).
- Optional caching of responses.
- Integration with `fzf` for interactive prompt selection.
