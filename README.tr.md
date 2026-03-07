# prayertime-cli

`prayertime-cli`, terminal ve ajan kullanımına uygun, açık kaynaklı bir namaz vakti aracıdır.

## Amaçlar

- Günlük namaz vakitlerini ve geri sayımları kararlı bir CLI sözleşmesiyle sunmak.
- Katı JSON desteği ve öngörülebilir çıkış kodlarıyla yapay zeka ajanları ve kabuk betikleri için uygun olmak.
- İlk aşamada Open-Meteo konum çözümleme ve AlAdhan `method=13` (Diyanet) ile stateless bir MVP sunmak.

## Planlanan Komut Yüzeyi

```text
prayertime-cli locations search --query <metin> [--country-code TR] [--limit 5] [--json]
prayertime-cli times get (--query <metin> | --lat <float> --lon <float>) [--country-code TR] [--date YYYY-MM-DD|today] [--json] [--field <alan>] [--quiet]
prayertime-cli times countdown (--query <metin> | --lat <float> --lon <float>) --target next-prayer|fajr|sunrise|dhuhr|asr|maghrib|isha|imsak|iftar [--at RFC3339] [--json] [--quiet]
prayertime-cli version
```

## İlkeler

- Komutlar ve bayraklar İngilizce ve kanoniktir.
- Türkçe destek yalnızca namaz adı ve alan seçici takma adlarıyla sınırlıdır.
- JSON çıktıları `stdout` üzerinden verilir; tanılayıcı bilgiler `stderr` üzerinde kalır.
- CLI etkileşimli onay istemez.

## Geliştirme

Bu depo Go 1.26 kullanır ve aşamalı, stacked PR iş akışını takip eder.

```bash
go test ./...
go build ./cmd/prayertime-cli
```

## Kurulum

Etiketli sürümler çapraz platform ikili dosyalar olarak yayımlanır. Homebrew Cask ve Scoop otomasyonu hazırdır.

```bash
# Homebrew
brew tap SeeknnDestroy/homebrew-tap
brew install --cask prayertime-cli

# Scoop
scoop bucket add prayertime-cli https://github.com/SeeknnDestroy/scoop-bucket
scoop install prayertime-cli

# Go
go install github.com/SeeknnDestroy/prayertime-cli/cmd/prayertime-cli@latest
```

## Örnekler

```bash
prayertime-cli locations search --query Istanbul --country-code TR --json
prayertime-cli times get --query Istanbul --field iftar --quiet
prayertime-cli times countdown --query Istanbul --target next-prayer --json
```

## Çıkış Kodları

- `0`: başarılı
- `1`: iç hata
- `2`: kullanım hatası
- `3`: bulunamadı veya belirsiz girdi
- `4`: ağ veya upstream zaman aşımı
- `5`: gelecekteki state conflict kodu

## Lisans

MIT
