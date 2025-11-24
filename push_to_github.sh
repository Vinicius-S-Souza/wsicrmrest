#!/bin/bash
# Script para fazer push usando Personal Access Token

echo "========================================"
echo "  Push para GitHub usando Token"
echo "========================================"
echo ""
echo "Você precisa de um Personal Access Token do GitHub"
echo "Criar em: https://github.com/settings/tokens"
echo ""
echo "Permissões necessárias: repo (full control)"
echo ""
read -p "Digite seu Personal Access Token: " TOKEN
echo ""

# Configurar remote com token
git remote set-url origin https://${TOKEN}@github.com/Vinicius-S-Souza/wsicrmrest.git

# Fazer push
echo "Fazendo push..."
git push -u origin main

# Remover token da URL por segurança
git remote set-url origin https://github.com/Vinicius-S-Souza/wsicrmrest.git

echo ""
echo "Push concluído!"
