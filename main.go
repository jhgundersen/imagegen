package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	apiBase      = "https://api.defapi.org"
	pollInterval = 5 * time.Second
)

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// taskData handles both object and array result shapes across models.
type taskData struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
	// Result is kept raw so we can handle object vs array variants per model.
	Result       json.RawMessage `json:"result"`
	StatusReason struct {
		Message *string `json:"message"`
	} `json:"status_reason"`
}

// extractImageURL handles the three result shapes seen across models:
//   - object with "image" field        (Wan 2.7)
//   - object with "big_image_url"      (Midjourney)
//   - array of objects with "image"    (GPT-Image, Nano Banana)
func extractImageURL(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	// Try array first.
	var arr []struct {
		Image string `json:"image"`
	}
	if json.Unmarshal(raw, &arr) == nil && len(arr) > 0 {
		return arr[0].Image
	}
	// Try object.
	var obj struct {
		Image       string `json:"image"`
		BigImageURL string `json:"big_image_url"`
	}
	json.Unmarshal(raw, &obj)
	if obj.Image != "" {
		return obj.Image
	}
	return obj.BigImageURL
}

func apiKey() string {
	key := os.Getenv("DEFAPI_API_KEY")
	if key == "" {
		fmt.Fprintln(os.Stderr, "error: DEFAPI_API_KEY environment variable not set")
		os.Exit(1)
	}
	return key
}

func post(endpoint string, body map[string]any, key string) json.RawMessage {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", apiBase+endpoint, strings.NewReader(string(b)))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	return readResponse(resp)
}

func get(endpoint string, key string) json.RawMessage {
	req, _ := http.NewRequest("GET", apiBase+endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	return readResponse(resp)
}

func readResponse(resp *http.Response) json.RawMessage {
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "HTTP %d: %s\n", resp.StatusCode, string(raw))
		os.Exit(1)
	}
	var ar apiResponse
	if err := json.Unmarshal(raw, &ar); err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	if ar.Code != 0 {
		fmt.Fprintf(os.Stderr, "API error %d: %s\n", ar.Code, ar.Message)
		os.Exit(1)
	}
	return ar.Data
}

func extractTaskID(data json.RawMessage) string {
	var d struct {
		TaskID string `json:"task_id"`
	}
	json.Unmarshal(data, &d)
	return d.TaskID
}

type stringListFlag []string

func (f *stringListFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *stringListFlag) Set(value string) error {
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			*f = append(*f, item)
		}
	}
	return nil
}

func normalizeGPTImageModel(model string) string {
	return strings.TrimPrefix(model, "openai/")
}

func contains(items []string, item string) bool {
	for _, candidate := range items {
		if candidate == item {
			return true
		}
	}
	return false
}

func poll(taskID, key string) string {
	fmt.Printf("Task submitted: %s\nPolling", taskID)
	for {
		time.Sleep(pollInterval)
		data := get("/api/task/query?task_id="+taskID, key)
		var td taskData
		json.Unmarshal(data, &td)

		switch td.Status {
		case "success":
			fmt.Println(" done.")
			imageURL := extractImageURL(td.Result)
			if imageURL == "" {
				fmt.Fprintln(os.Stderr, "error: no image URL in response")
				os.Exit(1)
			}
			fmt.Printf("Image URL: %s\n", imageURL)
			return imageURL
		case "failed":
			msg := "unknown reason"
			if td.StatusReason.Message != nil && *td.StatusReason.Message != "" {
				msg = *td.StatusReason.Message
			}
			pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
			fmt.Fprintf(os.Stderr, "\ngeneration failed: %s\n%s\n", msg, pretty)
			os.Exit(1)
		default:
			fmt.Print(".")
		}
	}
}

func guessExt(url string) string {
	lower := strings.ToLower(url)
	for _, ext := range []string{".png", ".jpg", ".jpeg", ".webp"} {
		if strings.Contains(lower, ext) {
			return ext
		}
	}
	return ".png"
}

func download(imageURL, taskID string) string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "Downloads")
	os.MkdirAll(dir, 0755)
	dest := filepath.Join(dir, "imagegen_"+taskID+guessExt(imageURL))

	fmt.Println("Downloading...")
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "download error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	io.Copy(f, resp.Body)
	return dest
}

func printResult(dest string) {
	fmt.Printf("\nSaved to: \033]8;;file://%s\033\\%s\033]8;;\033\\\n", dest, dest)
}

func openFile(path string) {
	cmd := "xdg-open"
	if runtime.GOOS == "darwin" {
		cmd = "open"
	}
	if err := exec.Command(cmd, path).Start(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not open file: %v\n", err)
	}
}

// --- model subcommands ---

func cmdWan(args []string) {
	fs := flag.NewFlagSet("wan", flag.ExitOnError)
	ratio := fs.String("ratio", "1:1", "Aspect ratio: 1:1 16:9 4:3 21:9 3:4 9:16 8:1")
	open := fs.Bool("open", false, "Open the image after download")
	fs.Usage = func() {
		fmt.Println("Usage: imagegen wan [flags] <prompt>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if prompt == "" {
		fs.Usage()
		os.Exit(1)
	}

	key := apiKey()
	fmt.Printf("Model: wan-2.7-image | Ratio: %s\nPrompt: %s\n\n", *ratio, prompt)

	data := post("/api/wan-image/gen", map[string]any{
		"model":        "wan-2.7-image",
		"prompt":       prompt,
		"aspect_ratio": *ratio,
	}, key)

	taskID := extractTaskID(data)
	imageURL := poll(taskID, key)
	dest := download(imageURL, taskID)
	printResult(dest)
	if *open {
		openFile(dest)
	}
}

func cmdMidjourney(args []string) {
	fs := flag.NewFlagSet("mj", flag.ExitOnError)
	speed := fs.String("speed", "fast", "Processing speed: fast, relax")
	bot := fs.String("bot", "MID_JOURNEY", "Bot type: MID_JOURNEY, NIJI_JOURNEY")
	image := fs.String("image", "", "Image URL or base64 for editing (uses edits endpoint)")
	open := fs.Bool("open", false, "Open the image after download")
	fs.Usage = func() {
		fmt.Println("Usage: imagegen mj [flags] <prompt>")
		fmt.Println()
		fmt.Println("Midjourney parameters (--ar, --stylize, etc.) can be appended directly to the prompt.")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if prompt == "" {
		fs.Usage()
		os.Exit(1)
	}

	key := apiKey()

	if *image != "" {
		fmt.Printf("Model: midjourney/edits | Speed: %s\nPrompt: %s\nImage: %s\n\n", *speed, prompt, *image)
		data := post("/api/midjourney/edits", map[string]any{
			"prompt": prompt,
			"image":  *image,
			"speed":  *speed,
		}, key)
		taskID := extractTaskID(data)
		imageURL := poll(taskID, key)
		dest := download(imageURL, taskID)
		printResult(dest)
		if *open {
			openFile(dest)
		}
	} else {
		fmt.Printf("Model: midjourney/imagine | Bot: %s | Speed: %s\nPrompt: %s\n\n", *bot, *speed, prompt)
		data := post("/api/midjourney/imagine", map[string]any{
			"prompt":   prompt,
			"bot_type": *bot,
			"speed":    *speed,
		}, key)
		taskID := extractTaskID(data)
		imageURL := poll(taskID, key)
		dest := download(imageURL, taskID)
		printResult(dest)
		if *open {
			openFile(dest)
		}
	}
}

func cmdGPTImage(command string, args []string, defaultModel string) {
	fs := flag.NewFlagSet(command, flag.ExitOnError)
	model := fs.String("model", defaultModel, "Model: gpt-image-1.5, gpt-image-2")
	size := fs.String("size", "auto", "Output size: auto, 1024x1024, 1536x1024, 1024x1536 (gpt-image-1.5) or auto, 1:1, 16:9, 9:16 (gpt-image-2)")
	quality := fs.String("quality", "auto", "Quality: auto, high, medium, low")
	background := fs.String("background", "auto", "Background: auto, opaque, transparent")
	format := fs.String("format", "png", "Output format: png, jpeg, webp")
	var images stringListFlag
	fs.Var(&images, "image", "Reference image URL for gpt-image-2 editing (repeatable or comma-separated)")
	open := fs.Bool("open", false, "Open the image after download")
	fs.Usage = func() {
		fmt.Printf("Usage: imagegen %s [flags] <prompt>\n", command)
		fs.PrintDefaults()
	}
	fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if prompt == "" {
		fs.Usage()
		os.Exit(1)
	}

	*model = normalizeGPTImageModel(*model)
	body := map[string]any{
		"model":  "openai/" + *model,
		"prompt": prompt,
		"size":   *size,
	}

	switch *model {
	case "gpt-image-1.5":
		if len(images) > 0 {
			fmt.Fprintln(os.Stderr, "error: --image is only supported with gpt-image-2")
			os.Exit(1)
		}
		body["quality"] = *quality
		body["background"] = *background
		body["output_format"] = *format
	case "gpt-image-2":
		if !contains([]string{"auto", "1:1", "16:9", "9:16"}, *size) {
			fmt.Fprintln(os.Stderr, "error: gpt-image-2 --size must be one of: auto, 1:1, 16:9, 9:16")
			os.Exit(1)
		}
		if len(images) > 16 {
			fmt.Fprintln(os.Stderr, "error: gpt-image-2 supports at most 16 --image values")
			os.Exit(1)
		}
		if len(images) > 0 {
			body["images"] = []string(images)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown GPT image model: %s\n", *model)
		os.Exit(1)
	}

	fmt.Printf("Model: %s | Size: %s", *model, *size)
	if *model == "gpt-image-1.5" {
		fmt.Printf(" | Quality: %s | Format: %s", *quality, *format)
	}
	if len(images) > 0 {
		fmt.Printf(" | Images: %d", len(images))
	}
	fmt.Printf("\nPrompt: %s\n\n", prompt)

	key := apiKey()
	data := post("/api/gpt-image/gen", body, key)

	taskID := extractTaskID(data)
	imageURL := poll(taskID, key)
	dest := download(imageURL, taskID)
	printResult(dest)
	if *open {
		openFile(dest)
	}
}

var googleModels = []string{
	"nano-banana",
	"nano-banana-pro",
	"nano-banana-2",
	"gemini-2.5-flash-image",
	"gemini-3.1-flash-image-preview",
}

// sizeSupportedModels are the only ones that accept the image_size parameter.
var sizeSupportedModels = map[string]bool{
	"nano-banana-pro": true,
	"nano-banana-2":   true,
}

func cmdGoogle(args []string) {
	fs := flag.NewFlagSet("google", flag.ExitOnError)
	model := fs.String("model", "nano-banana-2", "Model: "+strings.Join(googleModels, ", "))
	ratio := fs.String("ratio", "1:1", "Aspect ratio: auto 1:1 16:9 21:9 2:3 3:2 3:4 4:3 4:5 5:4 9:16")
	size := fs.String("size", "", "Output resolution: 1k, 2k, 4k (only for nano-banana-pro and nano-banana-2)")
	open := fs.Bool("open", false, "Open the image after download")
	fs.Usage = func() {
		fmt.Println("Usage: imagegen google [flags] <prompt>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if prompt == "" {
		fs.Usage()
		os.Exit(1)
	}
	if *size != "" && !sizeSupportedModels[*model] {
		fmt.Fprintf(os.Stderr, "warning: --size is only supported by nano-banana-pro and nano-banana-2, ignoring\n")
		*size = ""
	}

	key := apiKey()
	fmt.Printf("Model: google/%s | Ratio: %s", *model, *ratio)
	if *size != "" {
		fmt.Printf(" | Size: %s", *size)
	}
	fmt.Printf("\nPrompt: %s\n\n", prompt)

	body := map[string]any{
		"model":        "google/" + *model,
		"prompt":       prompt,
		"aspect_ratio": *ratio,
	}
	if *size != "" {
		body["image_size"] = *size
	}

	data := post("/api/image/gen", body, key)
	taskID := extractTaskID(data)
	imageURL := poll(taskID, key)
	dest := download(imageURL, taskID)
	printResult(dest)
	if *open {
		openFile(dest)
	}
}

func usage() {
	fmt.Println(`Usage: imagegen <model> [flags] <prompt>

Models:
  wan     Alibaba Wan 2.7 Image (text-to-image)
  mj      Midjourney Imagine (text-to-image, or edit with --image)
  gpt     OpenAI GPT-Image-1.5/2 (text-to-image, or gpt-image-2 edit with --image)
  gpt2    OpenAI GPT-Image-2 shortcut
  google  Google image models via --model flag (default: nano-banana-2)
            nano-banana, nano-banana-pro, nano-banana-2,
            gemini-2.5-flash-image, gemini-3.1-flash-image-preview

Run 'imagegen <model> --help' for model-specific flags.`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "wan":
		cmdWan(os.Args[2:])
	case "mj", "midjourney":
		cmdMidjourney(os.Args[2:])
	case "gpt":
		cmdGPTImage("gpt", os.Args[2:], "gpt-image-1.5")
	case "gpt2":
		cmdGPTImage("gpt2", os.Args[2:], "gpt-image-2")
	case "google":
		cmdGoogle(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown model: %s\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}
