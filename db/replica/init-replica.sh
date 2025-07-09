#!/bin/bash
set -e

# PostgreSQL veri dizini
PGDATA="/var/lib/postgresql/data"

# Primary sunucu bilgileri
PRIMARY_HOST="crm-postgres-primary"
PRIMARY_PORT="5432"

# Eğer veri dizini zaten başlatılmışsa (PG_VERSION varsa) hiçbir şey yapma
if [ -f "$PGDATA/PG_VERSION" ]; then
  echo "PGDATA zaten başlatılmış, replikasyon kurulumu atlanıyor."
  exit 0
fi

# Replika veri dizinini temizle (önceki veriler varsa)
echo "Veri dizini temizleniyor..."
rm -rf $PGDATA/*

# Primary sunucudan veri kopyala
echo "Primary sunucudan veri kopyalanıyor..."
pg_basebackup -h $PRIMARY_HOST -p $PRIMARY_PORT -U $POSTGRES_USER -Fp -Xs -v -R -D $PGDATA

# Dosya izinlerini düzelt
echo "Dosya izinleri düzeltiliyor..."
chown -R postgres:postgres $PGDATA
chmod -R 700 $PGDATA

# Docker entrypoint'in sunucuyu başlatmasına izin ver
echo "Replikasyon kurulumu tamamlandı."
echo "Docker entrypoint sunucuyu başlatacak..." 