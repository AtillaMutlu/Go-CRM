# Frontend Dokümantasyonu

## İçindekiler
1. Genel Bakış
2. Dosya Yapısı
3. HTML Sayfaları
4. JavaScript Modülleri
5. CSS ve Responsive Tasarım
6. API Entegrasyonu
7. Kimlik Doğrulama (JWT)
8. Hata Yönetimi
9. Kullanıcı Arayüzü Bileşenleri
10. Rate Limiting ve Cache
11. Kurumsal Kod Standartları

---

## 1. Genel Bakış
Frontend, vanilla HTML, CSS ve JavaScript (ES6 modülleri) ile geliştirilmiştir. Modern, responsive ve güvenli bir arayüz sunar. Tüm API istekleri JWT ile güvence altındadır. Merkezi hata yönetimi ve kullanıcı dostu bildirimler sağlanır.

---

## 2. Dosya Yapısı
public/
├── index.html          # Giriş (login) sayfası
├── dashboard.html      # Ana uygulama arayüzü
└── js/
    ├── api.js         # API istemcisi
    ├── login.js       # Giriş işlemleri
    └── dashboard.js   # Dashboard işlemleri

---

## 3. HTML Sayfaları
- index.html: Kullanıcı girişi ve JWT token alma
- dashboard.html: Müşteri ve iletişim yönetimi, kullanıcı işlemleri

---

## 4. JavaScript Modülleri
- api.js: Tüm API çağrıları, JWT token yönetimi, hata yönetimi
- login.js: Giriş formu, validasyon, token saklama
- dashboard.js: Müşteri ve iletişim işlemleri, modal yönetimi

---

## 5. CSS ve Responsive Tasarım
- Bootstrap ve modern CSS ile responsive arayüz
- Mobil ve masaüstü uyumlu tasarım

---

## 6. API Entegrasyonu
- Tüm istekler api.js üzerinden fetch ile yapılır
- JWT token localStorage’da saklanır ve her istekte Authorization header’ı ile gönderilir
- Rate limiting ve cache mekanizmaları backend ile entegre çalışır

---

## 7. Kimlik Doğrulama (JWT)
- Giriş sonrası alınan JWT token ile tüm korumalı endpoint’lere erişim sağlanır
- Token süresi dolduğunda otomatik logout ve bildirim

---

## 8. Hata Yönetimi
- Tüm hata mesajları kullanıcı dostu şekilde gösterilir
- API ve ağ hataları merkezi olarak yakalanır

---

## 9. Kullanıcı Arayüzü Bileşenleri
- Bootstrap tabanlı kartlar, tablolar, modal formlar
- Kullanıcı işlemleri için sade ve anlaşılır arayüz

---

## 10. Rate Limiting ve Cache
- API rate limiting ve cache mekanizmaları frontend ile uyumlu çalışır
- Sık yapılan işlemlerde kullanıcıya uyarı gösterilir

---

## 11. Kurumsal Kod Standartları
- Kodda ikon, emoji veya süsleme kullanılmaz
- Tüm değişken ve fonksiyon isimleri İngilizce, açıklamalar Türkçe
- Responsive ve erişilebilirlik standartlarına uygunluk
- Her modül için test ve dokümantasyon zorunludur

---

Daha fazla detay için backend ve API dokümantasyonuna bakınız. Her türlü katkı ve öneri için proje kurallarına uyunuz.