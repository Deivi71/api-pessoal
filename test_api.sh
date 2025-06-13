#!/bin/bash

# Script de teste para API Gestar Bem
echo "=== TESTE DA API GESTAR BEM ==="

# Token de exemplo (substitua pelo seu token válido)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NDk3MDA5MTQsInVzdWFyaW9JZCI6MTV9.OdlNmMC8Fs6GhcBdNYersuqdikHhrcCtGpe9xC5DgXs"

echo ""
echo "1. Testando busca de usuários:"
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:5000/usuarios | jq '.[0:3]'

echo ""
echo "2. Testando seguir usuário válido (ID 1):"
curl -s -X POST -H "Authorization: Bearer $TOKEN" http://localhost:5000/usuarios/1/seguir
echo "Status: OK"

echo ""
echo "3. Testando seguir usuário inexistente (ID 999):"
curl -s -X POST -H "Authorization: Bearer $TOKEN" http://localhost:5000/usuarios/999/seguir
echo ""

echo ""
echo "4. Testando seguir a si mesmo (ID 15):"
curl -s -X POST -H "Authorization: Bearer $TOKEN" http://localhost:5000/usuarios/15/seguir
echo ""

echo ""
echo "5. Testando feed de publicações:"
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:5000/publicacoes/feed | jq '.'

echo ""
echo "=== TESTE CONCLUÍDO ===" 