# imagegen

CLI tool for AI image generation via [defapi.org](https://defapi.org).

## Installation

**Pre-built binary** (Linux/macOS):

```sh
# Detect OS and arch, download latest release
OS=$(uname -s | tr '[:upper:]' '[:lower:]') ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') && \
curl -fsSL "https://github.com/jhgundersen/imagegen/releases/latest/download/imagegen-${OS}-${ARCH}" \
  -o ~/.local/bin/imagegen && chmod +x ~/.local/bin/imagegen
```

**Using Go:**

```sh
go install github.com/jhgundersen/imagegen@latest
```

**From source:**

```sh
make install
```

Installs to `~/.local/bin/imagegen`.

## Requirements

Set your API key:

```sh
export DEFAPI_API_KEY=your_key_here
```

## Usage

```sh
imagegen <model> [flags] <prompt>
```

### Models

#### `wan` — Alibaba Wan 2.7 Image

```sh
imagegen wan "a misty mountain lake at dawn"
imagegen wan "a portrait in watercolor style" --ratio 3:4
```

| Flag | Default | Options |
|------|---------|---------|
| `--ratio` | `1:1` | `1:1`, `16:9`, `4:3`, `21:9`, `3:4`, `9:16`, `8:1` |

#### `mj` — Midjourney

```sh
imagegen mj "a samurai cat --ar 3:4 --stylize 750"
imagegen mj "a cyberpunk city" --bot NIJI_JOURNEY --speed relax
imagegen mj "add snow to this scene" --image https://example.com/photo.jpg
```

Midjourney parameters (`--ar`, `--stylize`, `--chaos`, etc.) can be appended directly to the prompt string.

| Flag | Default | Options |
|------|---------|---------|
| `--speed` | `fast` | `fast`, `relax` |
| `--bot` | `MID_JOURNEY` | `MID_JOURNEY`, `NIJI_JOURNEY` |
| `--image` | | Image URL or base64 — switches to the edits endpoint |

#### `gpt` — OpenAI GPT-Image-1.5

```sh
imagegen gpt "a cozy coffee shop interior"
imagegen gpt "a product photo on white background" --quality high --size 1024x1024
imagegen gpt "a logo design" --background transparent --format png
```

| Flag | Default | Options |
|------|---------|---------|
| `--size` | `auto` | `auto`, `1024x1024`, `1536x1024`, `1024x1536` |
| `--quality` | `auto` | `auto`, `high`, `medium`, `low` |
| `--background` | `auto` | `auto`, `opaque`, `transparent` |
| `--format` | `png` | `png`, `jpeg`, `webp` |

#### `google` — Google image models

```sh
imagegen google "a sunset over the ocean"
imagegen google "a portrait" --model nano-banana-pro --size 4k --ratio 3:4
imagegen google "a wide landscape" --model gemini-2.5-flash-image --ratio 16:9
```

| Flag | Default | Options |
|------|---------|---------|
| `--model` | `nano-banana-2` | `nano-banana`, `nano-banana-pro`, `nano-banana-2`, `gemini-2.5-flash-image`, `gemini-3.1-flash-image-preview` |
| `--ratio` | `1:1` | `auto`, `1:1`, `16:9`, `21:9`, `2:3`, `3:2`, `3:4`, `4:3`, `4:5`, `5:4`, `9:16` |
| `--size` | | `1k`, `2k`, `4k` — only for `nano-banana-pro` and `nano-banana-2` |

## Output

Generated images are downloaded to `~/Downloads/imagegen_<task_id>.<ext>`. The path is printed as a clickable link in terminals that support OSC 8 hyperlinks (Kitty, Alacritty, WezTerm, GNOME Terminal 3.26+).
