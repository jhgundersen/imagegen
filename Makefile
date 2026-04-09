BINARY      := imagegen
INSTALL_DIR := $(HOME)/.local/bin
BASH_COMP   := $(HOME)/.local/share/bash-completion/completions
ZSH_COMP    := $(HOME)/.local/share/zsh/site-functions

.PHONY: build install install-completions clean

build:
	go build -o $(BINARY) .

install: build install-completions
	install -Dm755 $(BINARY) $(INSTALL_DIR)/$(BINARY)

install-completions:
	install -Dm644 completions/$(BINARY).bash $(BASH_COMP)/$(BINARY)
	install -Dm644 completions/$(BINARY).zsh  $(ZSH_COMP)/_$(BINARY)

clean:
	rm -f $(BINARY)
