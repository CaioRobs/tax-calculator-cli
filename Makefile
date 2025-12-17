# Makefile para testes e cobertura (Go)

GO := go
PKGS := ./...
COVER_OUT := cover.out
COVER_HTML := coverage.html
RM := rm -f

# Deleção cross-plataforma: usa cmd.exe no Windows e rm no Unix
ifeq ($(OS),Windows_NT)
	SHELL := cmd
	.SHELLFLAGS := /C
endif

.PHONY: test coverage cover-html coverage-all clean

## test: roda todos os testes
test:
	$(GO) test -v $(PKGS)

## coverage: gera perfil de cobertura (cover.out) e mostra resumo por função
coverage:
	$(GO) test $(PKGS) -coverprofile=$(COVER_OUT)
	$(GO) tool cover -func=$(COVER_OUT)

## cover-html: gera relatório HTML (coverage.html)
cover-html: coverage
	$(GO) tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

## coverage-all: cobertura cruzada entre pacotes (usa -coverpkg=./...)
coverage-all:
	$(GO) test $(PKGS) -coverpkg=$(PKGS) -coverprofile=$(COVER_OUT)
	$(GO) tool cover -func=$(COVER_OUT)

## clean: remove artefatos de cobertura
clean:
ifeq ($(OS),Windows_NT)
	- if exist "$(COVER_OUT)" del /Q "$(COVER_OUT)"
	- if exist "$(COVER_HTML)" del /Q "$(COVER_HTML)"
else
	- $(RM) "$(COVER_OUT)" "$(COVER_HTML)"
endif
