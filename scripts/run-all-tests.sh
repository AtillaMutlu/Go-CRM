#!/bin/bash

# Renkli output için
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 GO MIKROSERVIS TEST SÜİTİ BAŞLATILIYOR...${NC}\n"

# Test sonuçları
INFRASTRUCTURE_SUCCESS=false
UNIT_SUCCESS=false
INTEGRATION_SUCCESS=false
E2E_SUCCESS=false

# 1. Infrastructure Test
echo -e "${YELLOW}📋 1. INFRASTRUCTURE TEST${NC}"
echo "----------------------------------------"
if bash scripts/test-infrastructure.sh; then
    echo -e "${GREEN}✅ Infrastructure test BAŞARILI${NC}\n"
    INFRASTRUCTURE_SUCCESS=true
else
    echo -e "${RED}❌ Infrastructure test BAŞARISIZ${NC}\n"
fi

# 2. Unit Tests
echo -e "${YELLOW}📋 2. UNIT TESTS${NC}"
echo "----------------------------------------"
if go test -v ./tests/unit/...; then
    echo -e "${GREEN}✅ Unit tests BAŞARILI${NC}\n"
    UNIT_SUCCESS=true
else
    echo -e "${RED}❌ Unit tests BAŞARISIZ${NC}\n"
fi

# 3. Integration Tests (sadece infrastructure başarılıysa)
if [ "$INFRASTRUCTURE_SUCCESS" = true ]; then
    echo -e "${YELLOW}📋 3. INTEGRATION TESTS${NC}"
    echo "----------------------------------------"
    # Test database için environment variable set et
    export TEST_DB_URL="postgres://user:pass@localhost:5432/users?sslmode=disable"
    
    if go test -v ./tests/integration/...; then
        echo -e "${GREEN}✅ Integration tests BAŞARILI${NC}\n"
        INTEGRATION_SUCCESS=true
    else
        echo -e "${RED}❌ Integration tests BAŞARISIZ${NC}\n"
    fi
else
    echo -e "${YELLOW}⚠️  Integration tests ATLANDI (Infrastructure başarısız)${NC}\n"
fi

# 4. E2E Tests (sadece önceki testler başarılıysa)
if [ "$INFRASTRUCTURE_SUCCESS" = true ] && [ "$INTEGRATION_SUCCESS" = true ]; then
    echo -e "${YELLOW}📋 4. E2E TESTS${NC}"
    echo "----------------------------------------"
    echo "E2E testleri için API servisini başlatıyoruz..."
    
    # API servisini background'da başlat
    go run cmd/api/main.go &
    API_PID=$!
    sleep 5
    
    # Gateway servisini background'da başlat (eğer gerekirse)
    # go run cmd/gateway/main.go &
    # GATEWAY_PID=$!
    # sleep 3
    
    if go test -v ./tests/e2e/...; then
        echo -e "${GREEN}✅ E2E tests BAŞARILI${NC}\n"
        E2E_SUCCESS=true
    else
        echo -e "${RED}❌ E2E tests BAŞARISIZ${NC}\n"
    fi
    
    # API servisini durdur
    kill $API_PID 2>/dev/null
    # kill $GATEWAY_PID 2>/dev/null
    
else
    echo -e "${YELLOW}⚠️  E2E tests ATLANDI (Önceki testler başarısız)${NC}\n"
fi

# 5. Performance Tests (opsiyonel)
echo -e "${YELLOW}📋 5. PERFORMANCE TESTS (Opsiyonel)${NC}"
echo "----------------------------------------"
if [ "$E2E_SUCCESS" = true ]; then
    echo "Performance testleri için API servisini yeniden başlatıyoruz..."
    go run cmd/api/main.go &
    API_PID=$!
    sleep 3
    
    # Basic performance test
    echo "10 concurrent request testi..."
    for i in {1..10}; do
        curl -s -o /dev/null -w "%{http_code}" "http://localhost:8085/api/login" \
            -H "Content-Type: application/json" \
            -d '{"email":"demo@example.com","password":"demo123"}' &
    done
    wait
    
    kill $API_PID 2>/dev/null
    echo -e "${GREEN}✅ Performance test tamamlandı${NC}\n"
else
    echo -e "${YELLOW}⚠️  Performance tests ATLANDI${NC}\n"
fi

# 6. Benchmark Tests
echo -e "${YELLOW}📋 6. BENCHMARK TESTS${NC}"
echo "----------------------------------------"
if go test -bench=. -benchmem ./tests/unit/...; then
    echo -e "${GREEN}✅ Benchmark tests tamamlandı${NC}\n"
else
    echo -e "${RED}❌ Benchmark tests başarısız${NC}\n"
fi

# Test Özeti
echo -e "${BLUE}📊 TEST ÖZETİ${NC}"
echo "========================================"
echo -e "Infrastructure Tests: $([ "$INFRASTRUCTURE_SUCCESS" = true ] && echo -e "${GREEN}BAŞARILI${NC}" || echo -e "${RED}BAŞARISIZ${NC}")"
echo -e "Unit Tests:          $([ "$UNIT_SUCCESS" = true ] && echo -e "${GREEN}BAŞARILI${NC}" || echo -e "${RED}BAŞARISIZ${NC}")"
echo -e "Integration Tests:   $([ "$INTEGRATION_SUCCESS" = true ] && echo -e "${GREEN}BAŞARILI${NC}" || echo -e "${RED}BAŞARISIZ${NC}")"
echo -e "E2E Tests:           $([ "$E2E_SUCCESS" = true ] && echo -e "${GREEN}BAŞARILI${NC}" || echo -e "${RED}BAŞARISIZ${NC}")"

# Exit code
if [ "$INFRASTRUCTURE_SUCCESS" = true ] && [ "$UNIT_SUCCESS" = true ] && [ "$INTEGRATION_SUCCESS" = true ] && [ "$E2E_SUCCESS" = true ]; then
    echo -e "\n${GREEN}🎉 TÜM TESTLER BAŞARILI!${NC}"
    exit 0
else
    echo -e "\n${RED}💥 BAZI TESTLER BAŞARISIZ!${NC}"
    exit 1
fi 