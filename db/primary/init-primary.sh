#!/bin/sh
set -e

# PostgreSQL'in varsayılan pg_hba.conf dosyasına replikasyon için gerekli kuralı ekle.
# Bu kural, aynı ağdaki (0.0.0.0/0) 'user' kullanıcısının replikasyon bağlantısı kurmasına izin verir.
echo "host replication $POSTGRES_USER 0.0.0.0/0 scram-sha-256" >> "$PGDATA/pg_hba.conf" 