FROM tinygo/tinygo-dev

RUN go install golang.org/x/tools/gopls@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/sago35/tinygo-autocmpl@latest

RUN echo 'eval "$(tinygo-autocmpl --completion-script-bash)"' >> ~/.bashrc
