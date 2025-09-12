ARG NODE_IMAGE=node:lts-bookworm
FROM ${NODE_IMAGE} AS base

# Install build tools and runtime deps for Vim, plus tmux/screen
USER root
RUN set -eux; \
    apt-get update; \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
      ca-certificates curl git build-essential \
      libncurses5-dev libncursesw5-dev libacl1-dev \
      python3 python3-dev \
      golang-go \
      tmux screen \
      locales \
    && rm -rf /var/lib/apt/lists/*; \
    echo "en_US.UTF-8 UTF-8" > /etc/locale.gen; \
    locale-gen

# Build latest Vim from source
FROM base AS vim-build
RUN set -eux; \
    git clone --depth=1 https://github.com/vim/vim.git /tmp/vim; \
    cd /tmp/vim && \
    ./configure \
      --with-features=huge \
      --enable-multibyte \
      --enable-python3interp=yes \
      --enable-terminal \
      --enable-cscope \
      --prefix=/usr \
      --enable-fail-if-missing; \
    make -C /tmp/vim -j"$(nproc)"; \
    make -C /tmp/vim install; \
    rm -rf /tmp/vim

# Final image: Node + latest Vim + your tilde configs
FROM base AS final

# Copy the compiled Vim from builder stage
COPY --from=vim-build /usr/bin/vim /usr/bin/vim
COPY --from=vim-build /usr/share/vim /usr/share/vim

# Prepare workspace and user home
ENV LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    TERM=xterm-256color \
    WORKDIR=/work

WORKDIR $WORKDIR

# Copy tilde repo and install configs for the default user (node)
COPY ./ /opt/tilde/

# Install vim-plug into the node user's home and place dotfiles
RUN set -eux; \
    userhome=/home/node; \
    # Ensure plugin dirs exist with correct ownership
    install -d -o node -g node "$userhome/.vim/autoload" "$userhome/.vim/plugged"; \
    # Install vim-plug
    curl -fsSL https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim \
      -o "$userhome/.vim/autoload/plug.vim"; \
    # Copy dotfiles if present in tilde repo
    for f in .vimrc .tmux.conf .screenrc; do \
      if [ -f "/opt/tilde/$f" ]; then \
        cp "/opt/tilde/$f" "$userhome/$f"; \
      fi; \
    done; \
    chown -R node:node "$userhome/.vim" "$userhome/.vimrc" "$userhome/.tmux.conf" "$userhome/.screenrc" 2>/dev/null || true

# Switch to non-root by default
USER node

# Pre-install Vim plugins at build time for the node user
# Disable with: --build-arg INSTALL_PLUGINS=0
ARG INSTALL_PLUGINS=1
RUN if [ "$INSTALL_PLUGINS" = "1" ]; then \
      set -eux; \
      export HOME=/home/node; \
      export TERM=xterm-256color; \
      # Install plugins with better error handling and retry logic
      echo "Installing vim plugins at build time..."; \
      vim -E -s -u "$HOME/.vimrc" +"PlugInstall --sync" +qall || { \
        echo "First plugin install attempt failed, retrying..."; \
        sleep 2; \
        vim -E -s -u "$HOME/.vimrc" +"PlugInstall --sync" +qall || true; \
      }; \
      # Verify plugins were installed
      if [ -d "$HOME/.vim/plugged" ] && [ "$(ls -A $HOME/.vim/plugged 2>/dev/null)" ]; then \
        echo "Vim plugins successfully installed at build time"; \
        ls -la "$HOME/.vim/plugged"; \
      else \
        echo "Warning: Plugin installation may have failed, but continuing build"; \
      fi; \
    else \
      echo "Skipping plugin install at build time (INSTALL_PLUGINS=0)"; \
    fi

# Default to an interactive shell in the project workspace
CMD ["bash"]
