#!/bin/bash
set -e

# Script de build para release com incremento automático de versão

VERSION_FILE="VERSION"
BUILD_DIR="bin"

# Ler versão atual
if [ ! -f "$VERSION_FILE" ]; then
    echo "erro: arquivo VERSION não encontrado" >&2
    exit 1
fi

CURRENT_VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')

# Validar formato
if ! echo "$CURRENT_VERSION" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'; then
    echo "erro: versão inválida no arquivo VERSION: $CURRENT_VERSION" >&2
    exit 1
fi

# Incrementar versão PATCH
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$PATCH"

echo "Incrementando versão: $CURRENT_VERSION -> $NEW_VERSION"

# Atualizar arquivo VERSION
echo "$NEW_VERSION" > "$VERSION_FILE"

# Criar diretório de build
mkdir -p "$BUILD_DIR"

# Build para plataformas alvo
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$PLATFORM"
    OUTPUT="$BUILD_DIR/specs-$GOOS-$GOARCH"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi
    
    echo "Building $PLATFORM -> $OUTPUT"
    
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w -X main.version=$NEW_VERSION" \
        -trimpath \
        -o "$OUTPUT" \
        ./cmd/specs
done

echo "Build concluído. Versão: $NEW_VERSION"

# Criar tag Git
TAG="v$NEW_VERSION"
echo "Criando tag Git: $TAG"

# Verificar se tag já existe
if git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "erro: tag $TAG já existe" >&2
    exit 1
fi

# Commit da atualização de versão
git add "$VERSION_FILE"
git commit -m "chore: bump version to $NEW_VERSION" || true

# Criar tag
git tag -a "$TAG" -m "Release $TAG"

echo "Tag $TAG criada localmente"
echo "Execute 'git push origin $TAG' para enviar a tag e acionar o CI/CD"

