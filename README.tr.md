# prayertime-cli

`prayertime-cli`, namaz vakitleri ve geri sayımlar için stateless bir CLI aracıdır. Ajanlar, kabuk betikleri ve doğrudan terminal kullanımı için tasarlanmıştır.

## MVP 1

- Open-Meteo ile konum arama
- AlAdhan `method=13` (Diyanet) ile günlük vakit alma
- Bir sonraki namaza veya belirli bir vakte geri sayım yapma
- `--json` ile `stdout` üzerinde JSON ya da `--quiet` ile tek değer döndürme

## Girdi Modeli

- MVP 1 içinde kayıtlı varsayılan konum yoktur.
- Her `times` komutu ya `--query <yer>` ya da birlikte `--lat <float>` ve `--lon <float>` ister.
- `--date today`, hedef konumun saat dilimine göre çözülür.

## Yaygın Görevler

Önce aday konumları bul:

```bash
prayertime-cli locations search --query "Springfield" --country-code US --json
```

Bugünün tam namaz vakitlerini al:

```bash
prayertime-cli times get --query Istanbul --json
```

Otomasyon için tek bir alan çıkar:

```bash
prayertime-cli times get --query Ankara --country-code TR --field yatsi --quiet
```

"Ezana kaç dakika kaldı?" gibi genel bir geri sayım sor:

```bash
prayertime-cli times countdown --query Istanbul --target next-prayer --json
```

Belirli bir vakte, örneğin iftara, kalan süreyi sor:

```bash
prayertime-cli times countdown --query Istanbul --target iftar --quiet
```

Yer adı yerine koordinat kullan:

```bash
prayertime-cli times get --lat 41.01384 --lon 28.94966 --date today --json
```

Yazım hatası veya belirsiz şehir sorgusunu düzelt:

```bash
prayertime-cli locations search --query Istnbul --json
```

## Çıktı Modları

- `--json`: yapılandırılmış çıktı `stdout` üzerindedir. Hatalar da JSON olarak `stdout` üzerindedir.
- `--quiet`: tek bir yalın değer döndürür. `times get` için `--field` gerekir; `times countdown --quiet` varsayılan olarak `seconds_remaining` döndürür.
- `--output text|json|value`: aynı çıktı modelinin genel biçimidir. Tüm komutlarda tek bir çıktı anahtarı kullanmak istiyorsan `--output` kullan.
- Varsayılan insan modu: okunabilir çıktı `stdout`, hata ve yönlendirme `stderr`.
- Kesin çıkış kodları gerekiyorsa derlenmiş ikiliyi çalıştır. `go run` sıfır olmayan çıkışları sarar.

## Takma Adlar

- Komutlar ve bayraklar İngilizce ve kanoniktir.
- Namaz hedefleri ve alan seçiciler için Türkçe anlamsal takma adlar desteklenir.
- `iftar` ve `aksam`, `maghrib` olur; `yatsi`, `isha` olur; `ogle`, `dhuhr` olur.

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

## Geliştirme Ve Dokümantasyon

```bash
make verify
make docs
make build
make release-check
```

## Çıkış Kodları

- `0`: başarılı
- `1`: iç hata
- `2`: kullanım hatası
- `3`: bulunamadı veya belirsiz girdi
- `4`: ağ veya upstream zaman aşımı
- `5`: gelecekteki state conflict kodu

## Ek Dokümanlar

- [Agent Workflows](docs/agent-workflows.md)
- [CLI Contract](docs/cli-contract.md)
- [Agent Evaluation](docs/agent-evaluation.md)
- [CLI Reference](docs/cli/prayertime-cli.md)
- [ADR 0002: Data Sources](docs/adr/0002-data-sources.md)

## Lisans

MIT
