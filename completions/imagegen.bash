_imagegen() {
    local cur prev words cword
    _init_completion || return

    if [[ $cword -eq 1 ]]; then
        COMPREPLY=($(compgen -W "wan mj gpt gpt2 google" -- "$cur"))
        return
    fi

    case "${words[1]}" in
        wan)
            case "$prev" in
                --ratio) COMPREPLY=($(compgen -W "1:1 16:9 4:3 21:9 3:4 9:16 8:1" -- "$cur")) ;;
                --output|-o) COMPREPLY=($(compgen -f -- "$cur")) ;;
                *)       COMPREPLY=($(compgen -W "--ratio --output -o --open" -- "$cur")) ;;
            esac ;;
        mj|midjourney)
            case "$prev" in
                --speed) COMPREPLY=($(compgen -W "fast relax" -- "$cur")) ;;
                --bot)   COMPREPLY=($(compgen -W "MID_JOURNEY NIJI_JOURNEY" -- "$cur")) ;;
                --output|-o) COMPREPLY=($(compgen -f -- "$cur")) ;;
                *)       COMPREPLY=($(compgen -W "--speed --bot --image --output -o --open" -- "$cur")) ;;
            esac ;;
        gpt|gpt2)
            case "$prev" in
                --model)      COMPREPLY=($(compgen -W "gpt-image-1.5 gpt-image-2" -- "$cur")) ;;
                --size)       COMPREPLY=($(compgen -W "auto 1024x1024 1536x1024 1024x1536 1:1 16:9 9:16" -- "$cur")) ;;
                --quality)    COMPREPLY=($(compgen -W "auto high medium low" -- "$cur")) ;;
                --background) COMPREPLY=($(compgen -W "auto opaque transparent" -- "$cur")) ;;
                --format)     COMPREPLY=($(compgen -W "png jpeg webp" -- "$cur")) ;;
                --output|-o)  COMPREPLY=($(compgen -f -- "$cur")) ;;
                *)            COMPREPLY=($(compgen -W "--model --size --quality --background --format --image --output -o --open" -- "$cur")) ;;
            esac ;;
        google)
            case "$prev" in
                --model) COMPREPLY=($(compgen -W "nano-banana nano-banana-pro nano-banana-2 gemini-2.5-flash-image gemini-3.1-flash-image-preview" -- "$cur")) ;;
                --ratio) COMPREPLY=($(compgen -W "auto 1:1 16:9 21:9 2:3 3:2 3:4 4:3 4:5 5:4 9:16" -- "$cur")) ;;
                --size)  COMPREPLY=($(compgen -W "1k 2k 4k" -- "$cur")) ;;
                --output|-o) COMPREPLY=($(compgen -f -- "$cur")) ;;
                *)       COMPREPLY=($(compgen -W "--model --ratio --size --output -o --open" -- "$cur")) ;;
            esac ;;
    esac
}

complete -F _imagegen imagegen
