#!/bin/bash
set -e

# PostgreSQL'in zaten çalışıp çalışmadığını kontrol et ve durdur (güvenlik önlemi)
pg_ctl -D "$PGDATA" -m fast -w stop || true

# Ana veritabanının (primary) hazır olmasını bekle
until pg_isready -h postgres-primary -p 5432 -U "$POSTGRES_USER"
do
  echo "Ana veritabaninin hazir olmasi bekleniyor..."
  sleep 2
done
echo "Ana veritabani hazir."

# Replica'nın eski veri klasörünü temizle
rm -rf "$PGDATA"/*

# pg_basebackup ile ana veritabanından yeni bir kopya oluştur
echo "Ana veritabanindan kopya (backup) olusturuluyor..."
pg_basebackup -h postgres-primary -p 5432 -U "$POSTGRES_USER" -D "$PGDATA" -Fp -Xs -R

# -R bayrağı, replikasyon için gerekli olan standby.signal dosyasını ve 
# bağlantı ayarlarını (postgresql.auto.conf) otomatik olarak oluşturur.

echo "Replikasyon kurulumu tamamlandi. Standby (replica) sunucu baslatiliyor."

# Bu script'ten sonra, Postgres'in varsayılan Docker giriş script'i sunucuyu standby modunda başlatacaktır. 