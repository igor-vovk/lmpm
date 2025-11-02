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
      - source: local-dir # Reference to source name
        files:            # List of file paths to include
          - file1.txt
          - folder/file2.txt
      - source: git-repo
        files:
          - README.md
```

#### Minimal Configuration Example

Since the `working_dir` source is automatically available and `source` defaults to `working_dir`, you can write a minimal configuration:

```yaml
version: 1

targets:
  - name: my-target
    output: ./output
    include:
      - files:              # source defaults to working_dir
          - prompts/system.txt
          - prompts/user.txt
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
      - source: working_dir
        files:
          - prompts/system.txt
          - prompts/user.txt
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
      - files:
          - prompts/system.txt
          - prompts/user.txt
          - deep/nested/config.yaml

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
      - files:
          - prompts/system.txt
          - prompts/user.txt
          - deep/nested/config.yaml

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
      - files:
          - docs/intro.md
          - docs/guide.md
          - docs/api.md

# Results in a single file: ./all-docs.md
# Content format:
#
# # File: docs/intro.md
#
# <content of intro.md>
#
# # File: docs/guide.md
#
# <content of guide.md>
#
# # File: docs/api.md
#
# <content of api.md>
```

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
- `include`: List of includes from sources
  - `source`: Reference to a source name (optional, defaults to `working_dir`)
  - `files`: List of file paths to include from that source

### Configuration Location
- Default: `pim.yaml` or `.pim.yaml` in the current directory
- Can be overridden with `--config` flag

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
