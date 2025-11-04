# PIM

**PIM (Prompt Instruction Manager)** is a command-line utility for managing prompt instructions and related files from
multiple sources. Think of it as a package manager for AI prompts and instructions.

## Why PIM?

PIM solves common challenges when working with AI prompts and instructions across projects:

### 📝 **Instructions as Code**

Define your instructions once in a dedicated repository, then reuse them across multiple projects. Keep your AI prompts
version-controlled, reviewable, and maintainable just like your code.

### 🧩 **Modular Instruction Management**

Split large instruction files into smaller, focused components. PIM automatically merges and concatenates them into a
single file, making complex prompt systems easier to maintain and understand.

### 🔄 **Subscribe to External Prompt Libraries**

Pull instructions and prompts from external repositories
like [github/awesome-copilot](https://github.com/github/awesome-copilot). Stay up-to-date with community best practices
and organizational standards.

### 🏢 **Enterprise Governance**

Perfect for organizations managing multiple repositories. Centralize prompt governance, ensure consistency across teams,
and maintain compliance with organizational standards in multi-repo setups.

## Installation

### From Source

```bash
git clone https://github.com/hubblew/pim.git
cd pim
make install
```

This will install the `pim` binary to `$GOPATH/bin` (usually `~/go/bin`).

## Quick Start

1. Create a `pim.yaml` configuration file:

```yaml
version: 1

targets:
  - name: prompts
    output: .github/copilot-instructions.md
    include:
      - "prompts/system.txt"
      - "prompts/user.txt"
```

2. Run PIM:

```bash
pim install
```

## Use Cases

### Use Case 1: Reusable Instructions Across Projects

Create a central repository of prompts and reuse them across multiple projects:

```yaml
# In your shared prompts repository: github.com/myorg/ai-prompts
# prompts/
#   ├── code-review.md
#   ├── documentation.md
#   └── testing.md
```

```yaml
# In each project's pim.yaml
version: 1

sources:
  - name: org-prompts
    url: github.com/myorg/ai-prompts

targets:
  - name: copilot-instructions
    output: .github/copilot-instructions.md
    include:
      - "@org-prompts/prompts/code-review.md"
      - "@org-prompts/prompts/documentation.md"
```

### Use Case 2: Modular Instruction Files

Split complex instructions into maintainable components that PIM merges automatically:

```yaml
version: 1

targets:
  - name: combined-instructions
    output: .github/copilot-instructions.md  # .md extension triggers concat strategy
    include:
      - "instructions/base-rules.md"
      - "instructions/coding-style.md"
      - "instructions/security-guidelines.md"
      - "instructions/project-specific.md"
# Results in a single file with all instructions concatenated
```

### Use Case 3: Subscribe to Community Prompts

Stay updated with best practices from community repositories:

```yaml
version: 1

sources:
  - name: awesome-copilot
    url: github.com/github/awesome-copilot
  - name: org-standards
    url: github.com/myorg/engineering-standards

targets:
  - name: ai-instructions
    output: .github/copilot-instructions.md
    include:
      - "@awesome-copilot/prompts/best-practices.md"
      - "@org-standards/ai/code-quality.md"
      - "@org-standards/ai/security.md"
      - "docs/project-context.md"
```

### Use Case 4: Multi-Repo Governance

Ensure consistent AI behavior across organizational repositories:

```yaml
# Central governance repo: github.com/myorg/ai-governance
# governance/
#   ├── security-requirements.md
#   ├── code-standards.md
#   └── compliance.md

# Each team repository uses:
version: 1

sources:
  - name: governance
    url: github.com/myorg/ai-governance

targets:
  - name: copilot-setup
    output: .github/copilot-instructions.md
    include:
      - "@governance/governance/security-requirements.md"
      - "@governance/governance/code-standards.md"
      - "@governance/governance/compliance.md"
      - ".github/team-guidelines.md"
```

## Usage

### Commands

- `pim install [directory]` - Fetch files from sources to targets (defaults to current directory)
- `pim version` - Print version information
- `pim help` - Show help

### Configuration

PIM looks for `pim.yaml` or `.pim.yaml` in the current directory (or the directory specified as an argument).

#### Basic Configuration

```yaml
version: 1

sources:
  - name: local-prompts
    url: /path/to/prompts
  - name: shared-repo
    url: github.com/user/prompts-repo

targets:
  - name: my-project
    output: prompts/
    include:
      - "@local-prompts/system.txt"
      - "@local-prompts/user.txt"
      - "@shared-repo/templates/common.txt"
```

#### Strategy Examples

**Flatten strategy (default for directories)** - All files copied to output root:

```yaml
targets:
  - name: prompts
    output: output/
    strategy: flatten  # This is the default for directories
    include:
      - "prompts/system.txt"
      - "prompts/user.txt"
      - "deep/nested/file.txt"
# Result:
# output/
#   ├── system.txt
#   ├── user.txt
#   └── file.txt
```

**Preserve strategy** - Maintains directory structure:

```yaml
targets:
  - name: prompts
    output: output/
    strategy: preserve
    include:
      - "prompts/system.txt"
      - "prompts/user.txt"
      - "deep/nested/file.txt"
# Result:
# output/
#   ├── prompts/
#   │   ├── system.txt
#   │   └── user.txt
#   └── deep/
#       └── nested/
#           └── file.txt
```

**Concat strategy (default for .md/.txt files)** - Concatenates all files:

```yaml
targets:
  - name: combined-prompts
    output: all-prompts.md  # .md or .txt triggers concat by default
    include:
      - "prompts/system.txt"
      - "prompts/user.txt"
# Result: Single file
# all-prompts.md
```

#### Minimal Configuration

The `working_dir` source is automatically available and points to the current directory:

```yaml
version: 1

targets:
  - name: local-files
    output: output/
    include:
      - "file1.txt"
      - "file2.txt"
```

### Configuration Options

**Sources:**

- `name` - Unique identifier for the source
- `url` - Local directory path or Git repository URL
    - Local: `/absolute/path` or `./relative/path`
    - Git: `github.com/user/repo`

**Special Sources:**

- `working_dir` - Automatically added, points to current working directory

**Targets:**

- `name` - Target name
- `output` - Directory where files will be copied, or file path for concatenation
- `strategy` - How to organize copied files (optional, auto-detected)
    - `flatten` - Remove subdirectories, copy all files to output root (default for directories)
    - `preserve` - Maintain original directory structure
    - `concat` - Concatenate all files into a single output file (default for .md/.txt outputs)
- `include` - List of file paths to include
    - Format: `"path/to/file.txt"` for local files (from working_dir source)
    - Format: `"@source-name/path/to/file.txt"` for files from other sources
    - Multiple files can be included in one string separated by commas
    - Wildcards: Supports `*`, `?`, and `[...]` patterns (e.g., `"prompts/*.md"`, `"@source/docs/[a-z]*.txt"`)

## Development

### Running Tests

```bash
make test          # Run all tests
make test-verbose  # Run tests with verbose output
```

### Building

```bash
make build  # Build the binary
make clean  # Remove build artifacts
```

## License

See [LICENSE](LICENSE) file for details.

## Documentation

For detailed specification, see [SPEC.md](SPEC.md).

## Awesome Lists
- [Awesome Copilot](https://github.com/github/awesome-copilot/)
- [Contains Studio AI Agents](https://github.com/contains-studio/agents)