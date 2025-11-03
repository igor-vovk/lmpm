# PIM Specification

## Overview
PIM (Prompt Instruction Manager) is a command-line utility for managing prompt instructions and related files.

## Configuration File

The tool uses a YAML configuration file to manage settings and package definitions.

### Configuration Format

```yaml
version: 1  # Configuration schema version (default: 1)

sources:
  - name: local-dir       # Unique identifier for this source
    url: /path/to/dir     # Local directory path or git repository URL
  - name: git-repo
    url: https://github.com/username/repo.git

targets:
  - name: my-target       # Target name
    output: ./output/dir  # Output directory for downloaded files
    strategy: flatten     # Optional: 'flatten' (default), 'preserve' or 'concat'
    include:
      - "@local-dir/file1.txt"            # File from local-dir source
      - "@git-repo/README.md"             # File from git-repo source
      - "local-file.txt"                  # File from working_dir (implicit)
```

#### Include Format

The `include` field uses a simple string format to specify files:

- **Local files** (from `working_dir`): `"path/to/file.txt"`
- **Files from other sources**: `"@source-name/path/to/file.txt"`
- **Multiple files** in one line: `"file1.txt, file2.txt, file3.txt"`
- **Wildcard patterns**: `"*.md"`, `"**/*.txt"`, `"??.yaml"`, `"[a-z]*.md"`

#### Minimal Configuration Example

Since the `working_dir` source is automatically available, you can write a minimal configuration:

```yaml
version: 1

targets:
  - name: my-target
    output: ./output
    include:
      - "prompts/system.txt"
      - "prompts/user.txt"
```

This is equivalent to:

```yaml
version: 1

sources:
  - name: working_dir     # Automatically added
    url: /current/working/directory

targets:
  - name: my-target
    output: ./output
    include:
      - "@working_dir/prompts/system.txt"
      - "@working_dir/prompts/user.txt"
```

### Strategy Examples

#### Flatten Strategy (Default for Directories)

With `strategy: flatten`, all files are copied directly to the output directory root, regardless of their original subdirectory structure:

```yaml
version: 1

targets:
  - name: my-target
    output: ./output
    strategy: flatten  # or omit for directories, as this is the default
    include:
      - "prompts/system.txt"
      - "prompts/user.txt"
      - "deep/nested/config.yaml"

# Results in:
# ./output/system.txt
# ./output/user.txt
# ./output/config.yaml
```

#### Preserve Strategy

With `strategy: preserve`, the original directory structure is maintained:

```yaml
version: 1

targets:
  - name: my-target
    output: ./output
    strategy: preserve
    include:
      - "prompts/system.txt"
      - "prompts/user.txt"
      - "deep/nested/config.yaml"

# Results in:
# ./output/prompts/system.txt
# ./output/prompts/user.txt
# ./output/deep/nested/config.yaml
```

#### Concat Strategy (Default for .md/.txt Files)

With `strategy: concat`, all files are concatenated into a single output file. This is automatically selected when the output path ends with `.md` or `.txt`:

```yaml
version: 1

targets:
  - name: combined-docs
    output: ./all-docs.md  # .md extension triggers concat automatically
    include:
      - "docs/intro.md"
      - "docs/guide.md"
      - "docs/api.md"

# Results in a single file: ./all-docs.md
```

### Configuration Elements

#### Sources
- `name`: Unique identifier for the source
- `url`: Either a local directory path or a git repository URL
  - Local directories: `/path/to/directory` or `./relative/path`
  - Git repositories: `https://github.com/user/repo.git` or `git@github.com:user/repo.git`

**Special Sources:**
- `working_dir`: Automatically added to all configurations, pointing to the current working directory. This source is always available even if not explicitly defined in the YAML.

#### Targets
- `name`: Name of the target
- `output`: Directory where files will be downloaded/copied, or file path for concatenation
- `strategy`: How to organize copied files (optional, auto-detected based on output)
  - `flatten`: Remove subdirectories, copy all files to output root directory (default for directories)
  - `preserve`: Maintain the original directory structure from the source
  - `concat`: Concatenate all files into a single output file (default when output ends with .md or .txt)
- `include`: List of file paths to include
  - Format: `"path/to/file.txt"` for local files (from working_dir source)
  - Format: `"@source-name/path/to/file.txt"` for files from named sources
  - Supports wildcards: `*` (any characters), `?` (single character), `[...]` (character class)
  - Example: `"prompts/*.md"`, `"@source/docs/**/*.txt"`, `"config/[a-z]*.yaml"`

### Configuration Location
- Default: `pim.yaml` or `.pim.yaml` in the current directory
- Can be overridden with `--config` flag

## Wildcard Support

PIM supports glob patterns for flexible file selection using the standard wildcard syntax:

- `*` - Matches any number of characters (but not directory separators)
- `?` - Matches exactly one character
- `[abc]` - Matches any character in the set
- `[a-z]` - Matches any character in the range
- `**` - Matches any number of directories (when using recursive patterns)

### Wildcard Examples

```yaml
version: 1

targets:
  - name: all-markdown
    output: ./docs/combined.md
    include:
      - "docs/*.md"                       # All .md files in docs/
      - "@org-lib/prompts/**/*.md"       # All .md files recursively
      - "config/?.yaml"                   # Single-char filenames with .yaml
      - "@source/templates/[a-z]*.txt"   # Lowercase-starting files
```

### Error Handling

- If a wildcard pattern matches no files, PIM will return an error
- Non-wildcard patterns (literal paths) must also exist, or an error is returned

## Features (To Be Defined)
- Package management
- Version control
- Dependency resolution
- Configuration management

## Future Work
- Define package structure
- Define repository format
- Add authentication mechanisms
- Add caching strategies
