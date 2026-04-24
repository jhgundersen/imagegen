#compdef imagegen

_imagegen() {
    local state

    _arguments \
        '1: :->model' \
        '*: :->args'

    case $state in
        model)
            _values 'model' \
                'wan[Alibaba Wan 2.7 Image]' \
                'mj[Midjourney]' \
                'gpt[OpenAI GPT-Image-1.5/2]' \
                'gpt2[OpenAI GPT-Image-2]' \
                'google[Google image models]'
            ;;
        args)
            case ${words[2]} in
                wan)
                    _arguments \
                        '--ratio[Aspect ratio]:ratio:(1:1 16:9 4:3 21:9 3:4 9:16 8:1)' \
                        '--open[Open image after download]' \
                        '*:prompt:'
                    ;;
                mj|midjourney)
                    _arguments \
                        '--speed[Processing speed]:speed:(fast relax)' \
                        '--bot[Bot type]:bot:(MID_JOURNEY NIJI_JOURNEY)' \
                        '--image[Image URL for editing]:url:' \
                        '--open[Open image after download]' \
                        '*:prompt:'
                    ;;
                gpt|gpt2)
                    _arguments \
                        '--model[Model]:model:(gpt-image-1.5 gpt-image-2)' \
                        '--size[Output size]:size:(auto 1024x1024 1536x1024 1024x1536 1:1 16:9 9:16)' \
                        '--quality[Quality]:quality:(auto high medium low)' \
                        '--background[Background]:background:(auto opaque transparent)' \
                        '--format[Output format]:format:(png jpeg webp)' \
                        '--image[Reference image URL for gpt-image-2]:url:' \
                        '--open[Open image after download]' \
                        '*:prompt:'
                    ;;
                google)
                    _arguments \
                        '--model[Model]:model:(nano-banana nano-banana-pro nano-banana-2 gemini-2.5-flash-image gemini-3.1-flash-image-preview)' \
                        '--ratio[Aspect ratio]:ratio:(auto 1:1 16:9 21:9 2:3 3:2 3:4 4:3 4:5 5:4 9:16)' \
                        '--size[Output resolution]:size:(1k 2k 4k)' \
                        '--open[Open image after download]' \
                        '*:prompt:'
                    ;;
            esac
            ;;
    esac
}

_imagegen
