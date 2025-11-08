#!/bin/bash

# Script para probar la API de Email Sender

# Colores para la salida
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Iniciando pruebas de la API de Email Sender${NC}"
echo ""

# Verificar que el servidor esté ejecutándose
echo -e "${BLUE}Verificando que el servidor esté ejecutándose en puerto 8080...${NC}"
if ! curl -s http://localhost:8080 > /dev/null 2>&1; then
    echo -e "${RED}❌ El servidor no está ejecutándose. Inicia el servidor con: go run main.go${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Servidor detectado${NC}"
echo ""

# Prueba 1: Email básico
echo -e "${BLUE}📧 Prueba 1: Enviando email básico...${NC}"
curl -X POST http://localhost:8080/send-email \
  -H "Content-Type: application/json" \
  -d '{
    "mail": "Usuario de Prueba",
    "subject": "Email de Prueba",
    "body": "Este es un mensaje de prueba desde la API."
  }' \
  -w "\nStatus: %{http_code}\n" \
  -s

echo ""
echo ""

# Prueba 2: Recomendaciones de productos
echo -e "${BLUE}🛍️  Prueba 2: Enviando recomendaciones de productos...${NC}"
curl -X POST http://localhost:8080/recommendations \
  -H "Content-Type: application/json" \
  -d @example-recommendation-request.json \
  -w "\nStatus: %{http_code}\n" \
  -s

echo ""
echo -e "${GREEN}✅ Pruebas completadas${NC}"
echo -e "${BLUE}💡 Revisa tu bandeja de entrada para ver los emails enviados${NC}"
