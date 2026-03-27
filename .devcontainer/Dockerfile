FROM  mcr.microsoft.com/devcontainers/base:ubuntu

COPY --from=jdxcode/mise /usr/local/bin/mise /usr/local/bin/

RUN echo 'eval "$(mise activate bash)"' >> /home/vscode/.bashrc \
 && echo 'eval "$(mise activate zsh)"'  >> /home/vscode/.zshrc
RUN echo 'export TERM=xterm-256color' >> /home/vscode/.bashrc \
 && echo 'export TERM=xterm-256color' >> /home/vscode/.zshrc

RUN apt-get update &&\
    apt-get install -y xfsprogs &&\
    rm -rf /var/lib/apt/lists/*
