# imagegen

CLI tool for AI image generation via [defapi.org](https://defapi.org).

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/jhgundersen/imagegen/master/install.sh | sh
```

Or with Go: `go install github.com/jhgundersen/imagegen@latest`

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
imagegen wan "a misty mountain lake at dawn" -o lake.png
```

| Flag | Default | Options |
|------|---------|---------|
| `--ratio` | `1:1` | `1:1`, `16:9`, `4:3`, `21:9`, `3:4`, `9:16`, `8:1` |
| `--output`, `-o` | | Save image to this file path |

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
| `--output`, `-o` | | Save image to this file path |

#### `gpt` / `gpt2` — OpenAI GPT-Image

```sh
imagegen gpt "a cozy coffee shop interior"
imagegen gpt "a product photo on white background" --quality high --size 1024x1024
imagegen gpt "a logo design" --background transparent --format png
imagegen gpt --model gpt-image-2 "a wide editorial photo of a mountain road" --size 16:9
imagegen gpt2 "transform the scene to nighttime" --image https://example.com/photo.jpg
imagegen gpt2 "a simple black and white icon" -o icon.png
```

`gpt` defaults to `gpt-image-1.5` for compatibility. Use `--model gpt-image-2` or the `gpt2` shortcut to use GPT-Image-2.

| Flag | Default | Options |
|------|---------|---------|
| `--model` | `gpt-image-1.5` | `gpt-image-1.5`, `gpt-image-2` |
| `--size` | `auto` | `gpt-image-1.5`: `auto`, `1024x1024`, `1536x1024`, `1024x1536`; `gpt-image-2`: `auto`, `1:1`, `16:9`, `9:16` |
| `--quality` | `auto` | `gpt-image-1.5` only: `auto`, `high`, `medium`, `low` |
| `--background` | `auto` | `gpt-image-1.5` only: `auto`, `opaque`, `transparent` |
| `--format` | `png` | `gpt-image-1.5` only: `png`, `jpeg`, `webp` |
| `--image` | | `gpt-image-2` only: reference image URL, repeatable or comma-separated, up to 16 |
| `--output`, `-o` | | Save image to this file path |

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
| `--output`, `-o` | | Save image to this file path |

## Output

Generated images are downloaded to `~/Downloads/imagegen_<task_id>.<ext>` by default. Use `--output` or `-o` to choose the exact file path. Parent directories are created automatically. The path is printed as a clickable link in terminals that support OSC 8 hyperlinks (Kitty, Alacritty, WezTerm, GNOME Terminal 3.26+).
