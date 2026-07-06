FROM tinygo/tinygo-dev

RUN apt-get update && \
    apt-get install -y --no-install-recommends fish && \
    rm -rf /var/lib/apt/lists/*

RUN go install golang.org/x/tools/gopls@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/sago35/tinygo-autocmpl@latest

RUN echo 'eval "$(tinygo-autocmpl --completion-script-bash)"' >> ~/.bashrc && \
    mkdir -p ~/.config/fish && \
    echo 'tinygo-autocmpl --completion-script-fish | source' >> ~/.config/fish/config.fish
