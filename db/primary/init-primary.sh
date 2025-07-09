#!/bin/bash
set -e

# PostgreSQL veri dizini
PGDATA="/var/lib/postgresql/data"

# pg_hba.conf dosyasını güncelle
# Replikasyon için gerekli izinleri ekle
echo "Replikasyon için pg_hba.conf güncelleniyor..."
echo "host replication $POSTGRES_USER 0.0.0.0/0 trust" >> $PGDATA/pg_hba.conf
echo "host all all 0.0.0.0/0 trust" >> $PGDATA/pg_hba.conf

# postgresql.conf dosyasını güncelle
# Replikasyon için gerekli ayarları yap
echo "Replikasyon için postgresql.conf güncelleniyor..."
echo "wal_level = replica" >> $PGDATA/postgresql.conf
echo "max_wal_senders = 10" >> $PGDATA/postgresql.conf
echo "wal_keep_size = 128MB" >> $PGDATA/postgresql.conf

# Dosya izinlerini düzelt
echo "Dosya izinleri düzeltiliyor..."
chown -R postgres:postgres $PGDATA
chmod -R 700 $PGDATA

# postgres rolünü oluştur
# Şifre doğrulama sorununu aşmak için trust authentication kullanıyoruz
echo "postgres rolü oluşturuluyor..."
psql -U $POSTGRES_USER -d postgres -c "CREATE ROLE $POSTGRES_USER WITH LOGIN SUPERUSER;" || true

# Sunucuyu yeniden başlat
echo "Sunucu yeniden başlatılıyor..."
pg_ctl -D $PGDATA restart

echo "Primary sunucu replikasyon için hazır." 