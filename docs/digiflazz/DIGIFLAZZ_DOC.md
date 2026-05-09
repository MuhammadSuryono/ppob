# Digiflazz API Reference

Dokumen ini adalah versi yang dirapikan dari catatan integrasi Digiflazz. Isi teknis tetap dipertahankan, tetapi struktur Markdown dibuat konsisten agar lebih mudah dibaca, dicari, dan dijadikan referensi implementasi.

## Daftar Isi

1. Cek Deposit
2. Daftar Harga
3. Deposit
4. Topup Prepaid
5. Cek Tagihan Pascabayar
6. Bayar Tagihan Pascabayar
7. Cek Status
8. Inquiry PLN
9. Response Code
10. Webhooks

## 1. Cek Deposit

Cek deposit memberikan info mengenai sisa deposit yang dimiliki.

### Endpoint

```text
https://api.digiflazz.com/v1/cek-saldo
```

### Request

```json
{
  "cmd": "deposit",
  "username": "username",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

### Parameter Request

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| cmd | Value harus `deposit` | String | Ya |
| username | Username yang diatur pada koneksi API | String | Ya |
| sign | Signature dengan formula `md5(username + apiKey + "depo")` | String | Ya |

### Response

```json
{
  "data": {
    "deposit": 500000000000
  }
}
```

### Parameter Response

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| deposit | Sisa deposit | Float | Ya |

## 2. Daftar Harga

Daftar harga memberikan informasi produk yang sudah disetting.

### Endpoint

```text
https://api.digiflazz.com/v1/price-list
```

> Terdapat limitasi pengecekan daftar harga. Simpan hasil pricelist pada database Anda dan lakukan sinkronisasi berkala.
>
> Query dengan parameter `category`, `brand`, atau `type` tidak real-time. Perbedaan data bisa tertinggal sekitar 10 sampai 15 menit.

### 2.1 Price List Prepaid

#### Request

```json
{
  "cmd": "prepaid",
  "username": "username",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

#### Parameter Request

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| cmd | Command: `prepaid` atau `pasca` | String | Ya |
| username | Username yang diatur pada koneksi API | String | Ya |
| code | Kode produk buyer | String | Tidak |
| category | Kategori produk Digiflazz | String | Tidak |
| brand | Brand produk Digiflazz | String | Tidak |
| type | Tipe produk Digiflazz | String | Tidak |
| sign | Signature dengan formula `md5(username + apiKey + "pricelist")` | String | Ya |

#### Response

```json
{
  "data": [
    {
      "product_name": "Xl 100.000",
      "category": "Pulsa",
      "brand": "XL",
      "type": "Umum",
      "seller_name": "PT. ABC",
      "price": 98000,
      "buyer_sku_code": "X100",
      "buyer_product_status": true,
      "seller_product_status": true,
      "unlimited_stock": true,
      "stock": 0,
      "multi": true,
      "start_cut_off": "23:45",
      "end_cut_off": "00:15",
      "desc": "Pulsa Xl Rp 100.000"
    },
    {
      "product_name": "Telkomsel Pulsa 5.000",
      "category": "Pulsa",
      "brand": "TELKOMSEL",
      "type": "Umum",
      "seller_name": "PT. BCA",
      "price": 5100,
      "buyer_sku_code": "S5",
      "buyer_product_status": true,
      "seller_product_status": false,
      "unlimited_stock": false,
      "stock": 1200,
      "multi": false,
      "start_cut_off": "00:00",
      "end_cut_off": "00:00",
      "desc": "Pulsa Telkomsel Rp 5.000"
    }
  ]
}
```

#### Parameter Response

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| product_name | Nama produk | String | Ya |
| category | Nama kategori | String | Ya |
| brand | Nama brand | String | Ya |
| type | Nama tipe | String | Ya |
| seller_name | Nama seller | String | Ya |
| price | Harga produk dari seller | String | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| buyer_product_status | Status produk buyer | Boolean | Ya |
| seller_product_status | Status produk seller | Boolean | Ya |
| unlimited_stock | Penentu apakah stok terbatas | Boolean | Ya |
| stock | Sisa stok seller | String | Ya |
| multi | Apakah transaksi multi diizinkan | Bool | Ya |
| start_cut_off | Jam mulai cut off (`hh:mm`) | String | Ya |
| end_cut_off | Jam akhir cut off (`hh:mm`) | String | Ya |
| desc | Deskripsi produk | String | Ya |

### 2.2 Price List Pascabayar

#### Request

```json
{
  "cmd": "pasca",
  "username": "username",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

#### Parameter Request

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| cmd | Command: `prepaid` atau `pasca` | String | Ya |
| username | Username yang diatur pada koneksi API | String | Ya |
| code | Kode produk buyer | String | Tidak |
| brand | Brand produk Digiflazz | String | Tidak |
| sign | Signature dengan formula `md5(username + apiKey + "pricelist")` | String | Ya |

#### Response

```json
{
  "data": [
    {
      "product_name": "Pln Postpaid",
      "category": "Pascabayar",
      "brand": "PLN",
      "seller_name": "PT. ABC",
      "admin": 2750,
      "commission": 1800,
      "buyer_sku_code": "pln",
      "buyer_product_status": true,
      "seller_product_status": true,
      "desc": "-"
    },
    {
      "product_name": "aetra",
      "category": "Pascabayar",
      "brand": "PDAM",
      "seller_name": "Mr Ed",
      "admin": 2000,
      "commission": 550,
      "buyer_sku_code": "aetra",
      "buyer_product_status": true,
      "seller_product_status": true,
      "desc": "Provinsi Jakarta"
    }
  ]
}
```

#### Parameter Response

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| product_name | Nama produk | String | Ya |
| category | Nama kategori | String | Ya |
| brand | Nama brand | String | Ya |
| seller_name | Nama seller | String | Ya |
| admin | Biaya admin | Int | Ya |
| commission | Komisi untuk buyer | Int | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| buyer_product_status | Status produk buyer | Boolean | Ya |
| seller_product_status | Status produk seller | Boolean | Ya |
| desc | Deskripsi produk | String | Ya |

## 3. Deposit

Deposit digunakan untuk penarikan tiket deposit.

### Endpoint

```text
https://api.digiflazz.com/v1/deposit
```

### Request

```json
{
  "username": "your-username",
  "amount": 10000000,
  "bank": "BCA",
  "owner_name": "John Doe",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

### Parameter Request

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| username | Username yang diatur pada koneksi API | String | Ya |
| amount | Jumlah deposit yang diinginkan | Int | Ya |
| bank | Bank tujuan. Perorangan: `Flip` / `ShopeePay`. Perusahaan: `BCA` / `MANDIRI` / `BRI` / `BNI` | String | Ya |
| owner_name | Nama pemilik rekening pengirim | String | Ya |
| sign | Signature dengan formula `md5(username + apiKey + "deposit")` | String | Ya |

### Response

```json
{
  "data": {
    "rc": "00",
    "bank": "BCA",
    "payment_method": "Bank Transfer",
    "account_no": "0123 4567 89",
    "notes": "A6R5UPV",
    "amount": 10000001
  }
}
```

### Parameter Response

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| rc | Response code | String | Ya |
| bank | Bank tujuan | String | Ya |
| payment_method | `Bank Transfer` atau `Virtual Account` | String | Ya |
| account_no | Nomor rekening tujuan | String | Ya |
| amount | Jumlah akhir yang harus ditransfer | Int | Ya |
| notes | Berita transfer | String | Ya |

## 4. Topup Prepaid

Seluruh transaksi prepaid diproses secara sinkron. Request akan langsung menerima status `Sukses`, `Gagal`, atau `Pending`.

> Jika respons `Pending`, lakukan pengecekan ulang dengan mengirim topup ulang menggunakan `ref_id` yang sama.

### Endpoint

```text
https://api.digiflazz.com/v1/transaction
```

### Parameter Request

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| username | Username yang diatur pada koneksi API | String | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| customer_no | Nomor pelanggan | String | Ya |
| ref_id | Ref ID unik | String | Ya |
| sign | Signature dengan formula `md5(username + apiKey + ref_id)` | String | Ya |
| testing | Isi `true` untuk development | Boolean | Tidak |
| max_price | Limit harga maksimum | Int | Tidak |
| cb_url | Callback URL | String | Tidak |
| allow_dot | Isi `true` jika `customer_no` dapat berisi titik | Boolean | Tidak |

### Request Example

```json
{
  "username": "username",
  "buyer_sku_code": "xld25",
  "customer_no": "087800001233",
  "ref_id": "some1d",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

### Parameter Response

> Response JSON selalu dibungkus oleh key `data`.

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| ref_id | Ref ID unik Anda | String | Ya |
| customer_no | Nomor pelanggan | String | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| message | Deskripsi status transaksi | String | Ya |
| status | `Sukses`, `Pending`, atau `Gagal` | String | Ya |
| rc | Response code | String | Ya |
| sn | Serial number | String | Tidak |
| buyer_last_saldo | Saldo terakhir setelah transaksi | Float | Tidak |
| price | Harga produk | Integer | Ya |
| tele | Telegram seller | String | Tidak |
| wa | Whatsapp seller | String | Tidak |

### Response Example

```json
{
  "data": {
    "ref_id": "some1d",
    "customer_no": "087800001233",
    "buyer_sku_code": "xld25",
    "message": "Transaksi Pending",
    "status": "Pending",
    "rc": "03",
    "sn": "",
    "buyer_last_saldo": 100000,
    "price": 25000,
    "tele": "@telegram",
    "wa": "081234512345"
  }
}
```

## 5. Cek Tagihan Pascabayar

Inquiry pascabayar juga diproses secara sinkron dan langsung memberikan status `Sukses` atau `Gagal`.

### Endpoint

```text
https://api.digiflazz.com/v1/transaction
```

### Catatan Umum

- Produk dengan format tambahan: `PBB`, `E-Money`, dan `SAMSAT`.
- Response JSON dibungkus oleh key `data`.
- Signature menggunakan formula `md5(username + apiKey + ref_id)`.

### Request Standar

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| commands | Nilai harus `inq-pasca` | String | Ya |
| username | Username yang diatur pada koneksi API | String | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| customer_no | Nomor pelanggan | String | Ya |
| ref_id | Ref ID unik | String | Ya |
| sign | Signature `md5(username + apiKey + ref_id)` | String | Ya |
| testing | Isi `true` untuk development | Boolean | Tidak |

### Request Example

```json
{
  "commands": "inq-pasca",
  "username": "username",
  "buyer_sku_code": "pln",
  "customer_no": "530000000003",
  "ref_id": "some1d",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

### Response Umum

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| ref_id | Ref ID unik | String | Ya |
| customer_no | Nomor pelanggan | String | Ya |
| customer_name | Nama pelanggan | String | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| admin | Total biaya admin | Int | Ya |
| message | Deskripsi status transaksi | String | Ya |
| status | `Sukses` atau `Gagal` | String | Ya |
| rc | Response code | String | Ya |
| periode | Periode tagihan | String | Tidak |
| buyer_last_saldo | Saldo terakhir setelah transaksi | Float | Ya |
| price | Harga yang dipotong dari deposit | Int | Ya |
| selling_price | Harga yang ditagihkan ke client | Int | Ya |
| desc | Objek detail produk | Object | Ya |

### Field `desc` per Produk

#### PLN

| Field | Tipe | Keterangan |
| --- | --- | --- |
| tarif | String | Tarif PLN |
| daya | Int | Daya PLN |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| detail[].periode | String | Periode tagihan |
| detail[].nilai_tagihan | String | Nilai tagihan |
| detail[].admin | String | Admin per tagihan |
| detail[].denda | String | Denda per tagihan |

#### PDAM

| Field | Tipe | Keterangan |
| --- | --- | --- |
| tarif | String | Tarif PDAM |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| alamat | String | Alamat pelanggan |
| jatuh_tempo | String | Tanggal jatuh tempo |
| detail[].periode | String | Periode tagihan |
| detail[].nilai_tagihan | String | Nilai tagihan |
| detail[].denda | String | Denda |
| detail[].meter_awal | String | Meter awal |
| detail[].meter_akhir | String | Meter akhir |
| detail[].biaya_lain | String | Biaya lain |

#### INTERNET

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| detail[].periode | String | Periode tagihan |
| detail[].nilai_tagihan | String | Nilai tagihan |
| detail[].admin | String | Admin per tagihan |

#### BPJS Kesehatan

| Field | Tipe | Keterangan |
| --- | --- | --- |
| jumlah_peserta | String | Jumlah peserta |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| alamat | String | Alamat peserta |
| detail[].periode | String | Banyak periode tagihan |

#### Multifinance

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| item_name | String | Nama benda |
| no_rangka | String | Nomor rangka |
| no_pol | String | Nomor polisi |
| tenor | String | Jumlah angsuran |
| detail[].periode | String | Nomor urut periode |
| detail[].denda | String | Denda |
| detail[].biaya_lain | String | Biaya lain-lain |

#### PBB

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| alamat | String | Alamat |
| tahun_pajak | String | Tahun pajak |
| kelurahan | String | Kelurahan |
| kecamatan | String | Kecamatan |
| kode_kab_kota | String | Kode kabupaten/kota |
| kab_kota | String | Kabupaten/kota |
| luas_tanah | String | Luas tanah |
| luas_gedung | String | Luas gedung |

#### Pajak Daerah Lainnya

Sama seperti PBB, ditambah field `provinsi`.

#### Gas Negara / PERTAGAS

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| alamat | String | Alamat tagihan |
| detail[].periode | String | Periode tagihan |
| detail[].meter_awal | String | Meter awal |
| detail[].meter_akhir | String | Meter akhir |
| detail[].usage | String | Jumlah penggunaan |

#### TV

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| detail[].periode | String | Periode tagihan |
| detail[].nilai_tagihan | String | Nilai tagihan |
| detail[].no_ref | String | Nomor referensi |

#### BPJSTK

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| kode_iuran | String | Kode iuran |
| kode_program | String | Kode program |
| jkk | Int | Jaminan kecelakaan kerja |
| jkm | Int | Jaminan kematian |
| jht | Int | Jaminan hari tua |
| kantor_cabang | String | Kantor cabang |
| tgl_efektif | String | Tanggal efektif |
| tgl_expired | String | Tanggal kedaluwarsa |

#### BPJSTKPU

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| kode_iuran | String | Kode iuran |
| kode_program | String | Kode program |
| jkk | Int | Jaminan kecelakaan kerja |
| jkm | Int | Jaminan kematian |
| jht | Int | Jaminan hari tua |
| jpk | Int | Kehilangan pekerjaan |
| jpn | Int | Jaminan pensiun |
| npp | String | NPP |
| kode_divisi | String | Kode divisi |

#### PLN Nontaglis

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| transaksi | String | Informasi transaksi |
| no_registrasi | String | Nomor registrasi |
| tanggal_registrasi | String | Tanggal registrasi |

#### E-Money

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |

#### SAMSAT

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |
| alamat | String | Alamat |
| nomor_identitas | String | Nomor identitas |
| nomor_rangka | String | Nomor rangka |
| nomor_mesin | String | Nomor mesin |
| nomor_polisi | String | Nomor polisi |
| milik_kenama | String | Milik ke berapa |
| merek_kb | String | Merek kendaraan |
| model_kb | String | Model kendaraan |
| tahun_buatan | String | Tahun pembuatan |
| tgl_akhir_pajak_baru | String | Tanggal akhir pajak baru |
| biaya_pokok_bbn | String | Biaya pokok BBN |
| biaya_pokok_swd | String | Biaya pokok SWD |
| biaya_pokok_pkb | String | Biaya pokok PKB |
| biaya_denda_swd | String | Biaya denda SWD |
| biaya_denda_bbn | String | Biaya denda BBN |
| biaya_denda_pkb | String | Biaya denda PKB |
| biaya_admin_stnk | String | Biaya admin STNK |
| biaya_admin_tnkb | String | Biaya admin TNKB |
| biaya_parkir_pokok | String | Biaya parkir pokok |
| biaya_pajak_progresif | String | Biaya pajak progresif |

#### HP / Lainnya

| Field | Tipe | Keterangan |
| --- | --- | --- |
| lembar_tagihan | Int | Jumlah lembar tagihan |

### Request Khusus

#### PBB

```json
{
  "commands": "inq-pasca",
  "username": "username",
  "buyer_sku_code": "cimahi",
  "customer_no": "329801092375999991",
  "ref_id": "ref-4",
  "year": 2025,
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

Tambahan parameter:

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| year | Tahun pajak tertentu, misalnya `2025` | Int | Tidak |

#### E-Money

```json
{
  "commands": "inq-pasca",
  "username": "username",
  "buyer_sku_code": "emoney",
  "customer_no": "082100000001",
  "ref_id": "some1d",
  "amount": 22500,
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

Tambahan parameter:

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| amount | Denominasi e-money | Int | Ya |

#### SAMSAT

```json
{
  "commands": "inq-pasca",
  "username": "username",
  "buyer_sku_code": "samsat",
  "customer_no": "9658548523568705,0212502110170100",
  "ref_id": "some1d",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

### Sample Response Inquiry PLN

```json
{
  "data": {
    "ref_id": "some1d",
    "customer_no": "530000000001",
    "customer_name": "Nama Pelanggan Pertama",
    "buyer_sku_code": "pln",
    "admin": 2500,
    "message": "Transaksi Sukses",
    "status": "Sukses",
    "rc": "00",
    "periode": "201901",
    "buyer_last_saldo": 100000,
    "price": 10000,
    "selling_price": 11000,
    "desc": {
      "tarif": "R1",
      "daya": 1300,
      "lembar_tagihan": 1,
      "detail": [
        {
          "periode": "201901",
          "nilai_tagihan": "8000",
          "admin": "2500",
          "denda": "500"
        }
      ]
    }
  }
}
```

## 6. Bayar Tagihan Pascabayar

Transaksi pembayaran pascabayar juga diproses sinkron dengan status `Sukses`, `Gagal`, atau `Pending`.

> Pembayaran hanya bisa dilakukan pada tanggal yang sama dengan tanggal inquiry tagihan.
>
> Jika status `Pending`, tunggu webhook atau lakukan cek status.

### Endpoint

```text
https://api.digiflazz.com/v1/transaction
```

### Request

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| commands | Nilai harus `pay-pasca` | String | Ya |
| username | Username yang diatur pada koneksi API | String | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| customer_no | Nomor pelanggan | String | Ya |
| ref_id | Ref ID yang sama dengan saat inquiry | String | Ya |
| sign | Signature `md5(username + apiKey + ref_id)` | String | Ya |
| testing | Isi `true` untuk development | Boolean | Tidak |

### Request Example

```json
{
  "commands": "pay-pasca",
  "username": "username",
  "buyer_sku_code": "pln",
  "customer_no": "530000000003",
  "ref_id": "some1d",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

### Response Umum

Field response sama dengan inquiry, ditambah field berikut:

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| sn | Serial number / reference number | String | Ya |

### Sample Response Bayar PLN

```json
{
  "data": {
    "ref_id": "some1d",
    "customer_no": "530000000001",
    "customer_name": "Nama Pelanggan Pertama",
    "buyer_sku_code": "pln",
    "admin": 2500,
    "message": "Transaksi Sukses",
    "status": "Sukses",
    "rc": "00",
    "periode": "201901",
    "sn": "S1234554321N",
    "buyer_last_saldo": 90000,
    "price": 10000,
    "selling_price": 11000,
    "desc": {
      "tarif": "R1",
      "daya": 1300,
      "lembar_tagihan": 1,
      "detail": [
        {
          "periode": "201901",
          "nilai_tagihan": "8000",
          "admin": "2500",
          "denda": "500",
          "meter_awal": "00080000",
          "meter_akhir": "00090000"
        }
      ]
    }
  }
}
```

### Request Khusus SAMSAT

```json
{
  "commands": "pay-pasca",
  "username": "username",
  "buyer_sku_code": "samsat",
  "customer_no": "9658548523568705,0212502110170100",
  "ref_id": "some1d",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

## 7. Cek Status

### Catatan

- Jangan panggil API berulang untuk transaksi/data yang sama dalam interval kurang dari 1 menit.
- Untuk prepaid, cek status dilakukan dengan topup ulang memakai `ref_id` yang sama.
- Jangan cek status prepaid untuk transaksi yang sudah lebih dari 90 hari karena dapat membuat transaksi baru.
- Untuk postpaid, gunakan command `status-pasca`.
- Cek status postpaid untuk transaksi yang lebih dari 90 hari akan mengembalikan pesan `Data belum ada`.

### Request Status Pascabayar

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| commands | Nilai harus `status-pasca` | String | Ya |
| username | Username yang diatur pada koneksi API | String | Ya |
| buyer_sku_code | Kode produk buyer | String | Ya |
| customer_no | Nomor pelanggan | String | Ya |
| ref_id | Ref ID unik | String | Ya |
| sign | Signature `md5(username + apiKey + ref_id)` | String | Ya |

### Request Example

```json
{
  "commands": "status-pasca",
  "username": "username",
  "buyer_sku_code": "pln",
  "customer_no": "530000000003",
  "ref_id": "some1d",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

## 8. Inquiry PLN

Digunakan untuk validasi ID PLN.

### Endpoint

```text
https://api.digiflazz.com/v1/inquiry-pln
```

### Request

```json
{
  "username": "username",
  "customer_no": "1234554321",
  "sign": "740b00a1b8784e028cc8078edf66d12b"
}
```

### Parameter Request

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| username | Username yang diatur pada koneksi API | String | Ya |
| customer_no | ID PLN customer | String | Ya |
| sign | Signature dengan formula `md5(username + apiKey + customer_no)` | String | Ya |

### Response

> Response JSON dibungkus oleh key `data`.

```json
{
  "data": {
    "message": "Transaksi Sukses",
    "status": "Sukses",
    "rc": "00",
    "customer_no": "1234554321",
    "meter_no": "1234554321",
    "subscriber_id": "523300817840",
    "name": "DAVID",
    "segment_power": "R1 /000001300"
  }
}
```

### Parameter Response

| Parameter | Deskripsi | Tipe Data | Wajib |
| --- | --- | --- | --- |
| message | Deskripsi status transaksi | String | Ya |
| status | `Sukses` atau `Gagal` | String | Ya |
| rc | Response code | String | Ya |
| customer_no | ID PLN customer | String | Ya |
| meter_no | Nomor meter | String | Tidak |
| subscriber_id | Informasi ID customer | String | Tidak |
| name | Nama customer | String | Tidak |
| segment_power | Daya / segmen | String | Tidak |

## 9. Response Code

| RC | Message | Status | Terbentuk Transaksi | Deskripsi |
| --- | --- | --- | --- | --- |
| 00 | Transaksi Sukses | Sukses | Ya | - |
| 01 | Timeout | Gagal | Ya | - |
| 02 | Transaksi Gagal | Gagal | Ya | - |
| 03 | Transaksi Pending | Pending | Ya | - |
| 40 | Payload Error | Gagal | Tidak | Tipe data atau parameter tidak sesuai |
| 41 | Signature tidak valid | Gagal | Tidak | Pastikan formula signature dan mode API sesuai |
| 42 | Gagal memproses API Buyer | Gagal | Tidak | Username belum sesuai |
| 43 | SKU tidak ditemukan atau non-aktif | Gagal | Tidak | - |
| 44 | Saldo tidak cukup | Gagal | Tidak | - |
| 45 | IP tidak dikenali | Gagal | Tidak | Pastikan IP sudah di-whitelist |
| 47 | Transaksi sudah terjadi di buyer lain | Gagal | Tidak | - |
| 49 | Ref ID tidak unik | Gagal | Tidak | - |
| 50 | Transaksi Tidak Ditemukan | Gagal | Ya | - |
| 51 | Nomor Tujuan Diblokir | Gagal | Ya | - |
| 52 | Prefix Tidak Sesuai Dengan Operator | Gagal | Ya | - |
| 53 | Produk Seller Sedang Tidak Tersedia | Gagal | Ya | - |
| 54 | Nomor Tujuan Salah | Gagal | Ya | - |
| 55 | Produk Sedang Gangguan | Gagal | Ya | - |
| 56 | Limit saldo seller | Gagal | Tidak | Deprecated |
| 57 | Jumlah Digit Kurang Atau Lebih | Gagal | Ya | - |
| 58 | Sedang Cut Off | Gagal | Ya | - |
| 59 | Tujuan di Luar Wilayah/Cluster | Gagal | Ya | - |
| 60 | Tagihan belum tersedia | Gagal | Ya | - |
| 61 | Belum pernah melakukan deposit | Gagal | Tidak | - |
| 62 | Seller sedang mengalami gangguan | Gagal | Tidak | - |
| 63 | Tidak support transaksi multi | Gagal | Tidak | - |
| 64 | Tarik tiket gagal, coba nominal lain atau hubungi admin | Gagal | Tidak | - |
| 65 | Limit transaksi multi | Gagal | Tidak | Deprecated |
| 66 | Cut Off (Perbaikan Sistem Seller) | Gagal | Tidak | - |
| 67 | Seller belum ter-verfikasi | Gagal | Tidak | - |
| 68 | Stok habis | Gagal | Tidak | - |
| 69 | Harga seller lebih besar dari ketentuan harga buyer | Gagal | Tidak | - |
| 70 | Timeout Dari Biller | Gagal | Ya | - |
| 71 | Produk Sedang Tidak Stabil | Gagal | Ya | - |
| 72 | Lakukan Unreg Paket Dahulu | Gagal | Ya | - |
| 73 | Kwh Melebihi Batas | Gagal | Ya | - |
| 74 | Transaksi Refund | Gagal | Ya | - |
| 80 | Akun Anda telah diblokir oleh Seller | Gagal | Tidak | - |
| 81 | Seller ini telah diblokir oleh Anda | Gagal | Tidak | - |
| 82 | Akun Anda belum ter-verfikasi | Gagal | Tidak | - |
| 83 | Limitasi pengecekan pricelist tercapai | Gagal | Tidak | Pricelist semua produk maksimal 1x per 5 menit |
| 84 | Nominal tidak valid | Gagal | Ya | - |
| 85 | Limitasi transaksi tercapai | Gagal | Ya | Coba lagi 1 menit |
| 86 | Limitasi pengecekan nomor PLN tercapai | Gagal | Ya | - |
| 87 | Transaksi E-money wajib kelipatan Rp 1.000 | Gagal | Tidak | - |
| 88 | Akun Anda tidak dapat melakukan aksi ini | Gagal | Tidak | - |
| 99 | DF Router Issue | Pending | Ya | - |

## 10. Webhooks

Webhooks memungkinkan aplikasi menerima event saat transaksi dibuat atau berubah status.

### Konfigurasi

Webhook diatur pada menu `Atur Koneksi > API > Webhook`.

### Delivery Headers

| Header | Deskripsi |
| --- | --- |
| X-Digiflazz-Event | Jenis event yang memicu pengiriman |
| X-Hub-Signature | HMAC hex dari response body jika webhook memakai secret |
| User-Agent | Penanda jenis transaksi prepaid atau postpaid |

### Event

| Nama | Deskripsi |
| --- | --- |
| create | Dikirim saat transaksi baru dibuat |
| update | Dikirim saat transaksi berubah status |

### User-Agent

| Nama | Deskripsi |
| --- | --- |
| Digiflazz-Hookshot | Webhook transaksi prepaid |
| Digiflazz-Pasca-Hookshot | Webhook transaksi postpaid |

### Contoh Payload Prepaid

```http
POST /webhook HTTP/1.1
Host: localhost:4567
X-Hub-Signature: sha1=7d6f016c23d03b696e76dada91c07f178cc0af4d
User-Agent: Digiflazz-Hookshot
Content-Type: application/json
Content-Length: 445
X-Digiflazz-Event: create

{
  "data": {
    "ref_id": "30467470",
    "customer_no": "081280556115",
    "buyer_sku_code": "ovo100",
    "message": "Sukses",
    "status": "Sukses",
    "rc": "00",
    "buyer_last_saldo": 326719460,
    "sn": "SEPTIAPAR/20190401214753214742",
    "price": 199800,
    "tele": "@telegram",
    "wa": "081234512345"
  }
}
```

### Contoh Payload Postpaid

```http
POST /webhook HTTP/1.1
Host: localhost:4567
X-Hub-Signature: sha1=debdf6dfb3b62dfd3e98cd39e600027080938f52
User-Agent: Digiflazz-Pasca-Hookshot
Content-Type: application/json
Content-Length: 695
X-Digiflazz-Event: update

{
  "data": {
    "ref_id": "1763103975",
    "customer_no": "530000000000",
    "customer_name": "SUBCRIBER NAME",
    "buyer_sku_code": "plnpsaca",
    "admin": 2750,
    "message": "Transaksi Sukses",
    "status": "Sukses",
    "rc": "00",
    "sn": "004212C9245F1BA43A77CEBD5CD5DA39",
    "periode": "201608",
    "buyer_last_saldo": 326719460,
    "price": 300950,
    "selling_price": 302750,
    "desc": {
      "tarif": "R1",
      "daya": 1300,
      "lembar_tagihan": 1800,
      "detail": [
        {
          "periode": "201608",
          "nilai_tagihan": "300000",
          "admin": "2750",
          "denda": "0",
          "meter_awal": "00080000",
          "meter_akhir": "00080000"
        }
      ]
    }
  }
}
```

### Contoh Verifikasi Signature di Laravel

```php
<?php

use Illuminate\Http\Request;

Route::post('/webhook', function (Request $request) {
    $secret = 'somesecretvalue';

    $postData = file_get_contents('php://input');
    $signature = hash_hmac('sha1', $postData, $secret);

    if ($request->header('X-Hub-Signature') === 'sha1=' . $signature) {
        \Log::info(json_decode($request->getContent(), true));
    }
});
```

### Ping Event

Saat webhook diset, Digiflazz dapat mengirim event ping untuk memastikan endpoint bisa dipakai.

#### Payload Ping

| Key | Value |
| --- | --- |
| sed | Random string dari Digiflazz |
| hook_id | ID webhook |
| hook | Detail konfigurasi webhook |

#### Ping Endpoint

```text
https://api.digiflazz.com/v1/report/hooks/[YOUR-WEBHOOK-ID]/pings
```

#### Contoh Ping

```http
POST /v1/report/hooks/11aaabbb/pings HTTP/1.1
Host: localhost:4567
Accept: */*
Content-Length: 0

HTTP/1.1 200 OK
Content-Length: 155
Content-Type: application/json

{
  "sed": "AgXXtVAHp",
  "hook_id": "11aaabbb",
  "hook": {
    "url": "https://awesomesite.com/webhooks",
    "secret": "somesecretkeywords",
    "type": "application/json",
    "status": 1
  }
}
```