-- Migration: 016_digiflazz_data_migration.sql
-- Description: Data migration for Digiflazz categories and products
-- +goose Up

-- Add brand column
ALTER TABLE products ADD COLUMN IF NOT EXISTS brand VARCHAR(100);

-- Add code column to categories
ALTER TABLE categories ADD COLUMN IF NOT EXISTS code VARCHAR(100) UNIQUE;

-- Insert Categories
INSERT INTO categories (name, code, is_active) VALUES
    ('Aktivasi Perdana', 'aktivasi_perdana', TRUE),
    ('Aktivasi Voucher', 'aktivasi_voucher', TRUE),
    ('Data', 'data', TRUE),
    ('E-Money', 'e-money', TRUE),
    ('Games', 'games', TRUE),
    ('Gas', 'gas', TRUE),
    ('Masa Aktif', 'masa_aktif', TRUE),
    ('PLN', 'pln', TRUE),
    ('Paket SMS & Telpon', 'paket_sms_and_telpon', TRUE),
    ('Pulsa', 'pulsa', TRUE),
    ('TV', 'tv', TRUE),
    ('Voucher', 'voucher', TRUE)
ON CONFLICT (name) DO UPDATE SET code = EXCLUDED.code;

-- Insert Products
INSERT INTO products (code, name, brand, category, category_id, price, price_api, platform_price, provider, product_type, description, stock, status, is_prepaid, is_active) VALUES
    ('ax10', 'Axis 10.000', 'AXIS', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 10893, 10893, 10893, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('ax5', 'Axis 5.000', 'AXIS', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 5848, 5848, 5848, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('axdj1', 'Axis Data Jawa 2.5 GB 5 Hari', 'AXIS', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 12697, 12697, 12697, 'Digiflazz', 'Jawa Bali Nusra', 'AIGO Mini 2.5GB + Lokal Jawa Bali Nusra 5hr', 0, 'active', TRUE, TRUE),
    ('axdss2', 'Axis Data SS 2 GB 3 Hari', 'AXIS', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 9630, 9630, 9630, 'Digiflazz', 'Aigo SS', 'Axis Data SS AIGO Mini Bronet 24Jam 2GB + Kuota di Kota-mu 3hr', 0, 'active', TRUE, TRUE),
    ('axp3g60d', 'Aktivasi Perdana Axis 3 GB 60 Hari (SP5K SP7K)', 'AXIS', 'Aktivasi Perdana', (SELECT id FROM categories WHERE name = 'Aktivasi Perdana'), 13905, 13905, 13905, 'Digiflazz', 'SP5K SP7K', 'Perdana Bronet 3GB + Kuota di Kota-mu 60hr', 0, 'active', TRUE, TRUE),
    ('byu10', 'by.U 10.000', 'by.U', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 10170, 10170, 10170, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('dana20', 'DANA 20.000', 'DANA', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 20135, 20135, 20135, 'Digiflazz', 'Umum', '-', 0, 'active', TRUE, TRUE),
    ('dana50', 'DANA 50.000', 'DANA', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 50150, 50150, 50150, 'Digiflazz', 'Umum', '-', 0, 'active', TRUE, TRUE),
    ('danacek', 'Cek Nama Pengguna DANA', 'DANA', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 10, 10, 10, 'Digiflazz', 'Umum', '-', 0, 'active', TRUE, TRUE),
    ('ff12', 'Free Fire 12 Diamond', 'FREE FIRE', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 1825, 1825, 1825, 'Digiflazz', 'Umum', 'Jumlah diamond sesuai diamond normal, bonus tidak dihitung', 0, 'active', TRUE, TRUE),
    ('ff140', 'Free Fire 140 Diamond', 'FREE FIRE', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 18040, 18040, 18040, 'Digiflazz', 'Umum', 'Jumlah diamond sesuai diamond normal, bonus tidak dihitung', 0, 'active', TRUE, TRUE),
    ('ff355', 'Free Fire 355 Diamond', 'FREE FIRE', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 45100, 45100, 45100, 'Digiflazz', 'Umum', 'Jumlah diamond sesuai diamond normal, bonus tidak dihitung', 0, 'active', TRUE, TRUE),
    ('ff50', 'Free Fire 50 Diamond', 'FREE FIRE', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 6780, 6780, 6780, 'Digiflazz', 'Umum', 'Jumlah diamond sesuai diamond normal, bonus tidak dihitung', 0, 'active', TRUE, TRUE),
    ('ff70', 'Free Fire 70 Diamond', 'FREE FIRE', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 9035, 9035, 9035, 'Digiflazz', 'Umum', 'Jumlah diamond sesuai diamond normal, bonus tidak dihitung', 0, 'active', TRUE, TRUE),
    ('flash1', 'Telkomsel Data Flash 1 GB 30 Hari', 'TELKOMSEL', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 11410, 11410, 11410, 'Digiflazz', 'Flash', '24 jam nasional.', 0, 'active', TRUE, TRUE),
    ('flash2', 'Telkomsel Data Flash 2 GB 30 Hari', 'TELKOMSEL', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 21175, 21175, 21175, 'Digiflazz', 'Flash', '24 jam nasional.', 0, 'active', TRUE, TRUE),
    ('flash3', 'Telkomsel Data Flash 3 GB 30 Hari', 'TELKOMSEL', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 21255, 21255, 21255, 'Digiflazz', 'Flash', '24 jam nasional.', 0, 'active', TRUE, TRUE),
    ('flexs', 'XL Xtra Combo Flex S 28 Hari', 'XL', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 32075, 32075, 32075, 'Digiflazz', 'Xtra Combo Flex', 'Xtra Combo Flex S', 0, 'active', TRUE, TRUE),
    ('go100', 'Go Pay 100.000', 'GO PAY', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 100789, 100789, 100789, 'Digiflazz', 'Customer', 'Masukan no HP', 0, 'active', TRUE, TRUE),
    ('go50', 'Go Pay 50.000', 'GO PAY', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 50975, 50975, 50975, 'Digiflazz', 'Customer', 'Masukan no HP', 0, 'active', TRUE, TRUE),
    ('gopaycek', 'Cek Nama Pengguna Gopay', 'GO PAY', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 10, 10, 10, 'Digiflazz', 'Umum', '-', 0, 'active', TRUE, TRUE),
    ('happy1', 'Tri Data Happy 1.5 GB 1 Hari', 'TRI', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 5355, 5355, 5355, 'Digiflazz', 'Happy', 'Tri Data Happy 1.5 GB / 1 Hari', 0, 'active', TRUE, TRUE),
    ('happy3', 'Tri Data Happy 3 GB 3 Hari', 'TRI', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 11980, 11980, 11980, 'Digiflazz', 'Happy', 'Tri Data Happy 3 GB / 3 Hari', 0, 'active', TRUE, TRUE),
    ('hotrod3g10d', 'Aktivasi Voucher XL XTRA HotRod Special 3 GB 10 Hari', 'XL', 'Aktivasi Voucher', (SELECT id FROM categories WHERE name = 'Aktivasi Voucher'), 18199, 18199, 18199, 'Digiflazz', 'Hotrod Special', 'Aktivasi Voucher Xtra Hotrod Special 3GB,10hr', 0, 'active', TRUE, TRUE),
    ('i10', 'Indosat 10.000', 'INDOSAT', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 12084, 12084, 12084, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('i20', 'Indosat 20.000', 'INDOSAT', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 20671, 20671, 20671, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('i25', 'Indosat 25.000', 'INDOSAT', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 25595, 25595, 25595, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('i30', 'Indosat 30.000', 'INDOSAT', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 30150, 30150, 30150, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('i5', 'Indosat 5.000', 'INDOSAT', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 6630, 6630, 6630, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('i50', 'Indosat 50.000', 'INDOSAT', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 48530, 48530, 48530, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('iactive90', 'Indosat Tambah Masa Aktif Kartu 90 Hari', 'INDOSAT', 'Masa Aktif', (SELECT id FROM categories WHERE name = 'Masa Aktif'), 32210, 32210, 32210, 'Digiflazz', 'Umum', '-', 0, 'active', TRUE, TRUE),
    ('if2', 'Indosat Freedom Internet 2.5 GB 5 Hari', 'INDOSAT', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 12810, 12810, 12810, 'Digiflazz', 'Freedom Internet', 'Freedom Internet 2.5GB/5hr', 0, 'active', TRUE, TRUE),
    ('if3g30d', 'Indosat Freedom Internet 3 GB 28 Hari', 'INDOSAT', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 20750, 20750, 20750, 'Digiflazz', 'Freedom Internet', 'Freedom Internet 3GB/28 Hari', 0, 'active', TRUE, TRUE),
    ('if3g3d', 'Indosat Freedom Internet 3 GB 3 Hari', 'INDOSAT', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 11615, 11615, 11615, 'Digiflazz', 'Freedom Internet', 'Freedom Internet 3GB/3hr', 0, 'active', TRUE, TRUE),
    ('if5g30d', 'Indosat Freedom Internet 5.5 GB 28 Hari', 'INDOSAT', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 30275, 30275, 30275, 'Digiflazz', 'Freedom Internet', 'Freedom Internet 5.5GB/28 Hari', 0, 'active', TRUE, TRUE),
    ('kvision180d', 'K-Vision & GOL Paket CLING (CL06)  180 Hari', 'K-VISION dan GOL', 'TV', (SELECT id FROM categories WHERE name = 'TV'), 95500, 95500, 95500, 'Digiflazz', 'Cling', 'Siaran National Geographic, Nat Geo Wild, My Family, My Cinema, MTV, Rock , Kids TV, dll', 0, 'active', TRUE, TRUE),
    ('kvision30d', 'K-Vision & GOL Paket CLING (CL01)  30 Hari', 'K-VISION dan GOL', 'TV', (SELECT id FROM categories WHERE name = 'TV'), 19160, 19160, 19160, 'Digiflazz', 'Cling', 'Siaran National Geographic, Nat Geo Wild, My Family, My Cinema, MTV, Rock , Kids TV, dll', 0, 'active', TRUE, TRUE),
    ('ml10', 'MOBILELEGEND - 10 Diamond', 'MOBILE LEGENDS', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 2843, 2843, 2843, 'Digiflazz', 'Umum', 'no pelanggan = gabungan antara user_id dan zone_id', 0, 'active', TRUE, TRUE),
    ('ml12', 'MOBILELEGEND - 12 Diamond', 'MOBILE LEGENDS', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 3363, 3363, 3363, 'Digiflazz', 'Umum', 'no pelanggan = gabungan antara user_id dan zone_id', 0, 'active', TRUE, TRUE),
    ('ml5', 'MOBILELEGEND - 5 Diamond', 'MOBILE LEGENDS', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 1450, 1450, 1450, 'Digiflazz', 'Umum', 'no pelanggan = gabungan antara user_id dan zone_id', 0, 'active', TRUE, TRUE),
    ('mlweek', 'MOBILE LEGENDS Weekly Diamond Pass', 'MOBILE LEGENDS', 'Games', (SELECT id FROM categories WHERE name = 'Games'), 27525, 27525, 27525, 'Digiflazz', 'Membership', '-', 0, 'active', TRUE, TRUE),
    ('ovo100', 'OVO 100.000', 'OVO', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 100595, 100595, 100595, 'Digiflazz', 'Umum', 'OVO 100.000', 0, 'active', TRUE, TRUE),
    ('ovo50', 'OVO 50.000', 'OVO', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 50620, 50620, 50620, 'Digiflazz', 'Umum', 'OVO 50.000', 0, 'active', TRUE, TRUE),
    ('ovocek', 'Cek Nama Pengguna OVO', 'OVO', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 10, 10, 10, 'Digiflazz', 'Umum', '-', 0, 'active', TRUE, TRUE),
    ('pas10', 'Telkomsel Telepon Pas 10.000', 'TELKOMSEL', 'Paket SMS & Telpon', (SELECT id FROM categories WHERE name = 'Paket SMS & Telpon'), 12400, 12400, 12400, 'Digiflazz', 'Telepon Pas', 'Telepon 170 menit sesama, 30 menit semua op 3 Hari (sesuai zona)', 0, 'active', TRUE, TRUE),
    ('pas20', 'Telkomsel Telepon Pas 20.000', 'TELKOMSEL', 'Paket SMS & Telpon', (SELECT id FROM categories WHERE name = 'Paket SMS & Telpon'), 21230, 21230, 21230, 'Digiflazz', 'Telepon Pas', 'Telepon 300 menit semua op 7 Hari (sesuai zona)', 0, 'active', TRUE, TRUE),
    ('pas50', 'Telkomsel Telepon Pas 50.000', 'TELKOMSEL', 'Paket SMS & Telpon', (SELECT id FROM categories WHERE name = 'Paket SMS & Telpon'), 21625, 21625, 21625, 'Digiflazz', 'Telepon Pas', 'Telepon 1000 menit sesama, 100 menit semua op (30 Hari) (manfaat sesuai zona)', 0, 'active', TRUE, TRUE),
    ('pertagas20', 'Pertagas 20.000', 'Pertamina Gas', 'Gas', (SELECT id FROM categories WHERE name = 'Gas'), 21935, 21935, 21935, 'Digiflazz', 'Umum', 'Pertagas 20.000', 0, 'active', TRUE, TRUE),
    ('pln100', 'PLN 100.000', 'PLN', 'PLN', (SELECT id FROM categories WHERE name = 'PLN'), 101744, 101744, 101744, 'Digiflazz', 'Umum', 'masukkan nomor meter/id pelanggan', 0, 'active', TRUE, TRUE),
    ('pln1000', 'PLN 1.000.000', 'PLN', 'PLN', (SELECT id FROM categories WHERE name = 'PLN'), 1001985, 1001985, 1001985, 'Digiflazz', 'Umum', 'masukkan nomor meter/id pelanggan', 0, 'active', TRUE, TRUE),
    ('pln20', 'PLN 20.000', 'PLN', 'PLN', (SELECT id FROM categories WHERE name = 'PLN'), 21880, 21880, 21880, 'Digiflazz', 'Umum', 'masukkan nomor meter/id pelanggan', 0, 'active', TRUE, TRUE),
    ('pln50', 'PLN 50.000', 'PLN', 'PLN', (SELECT id FROM categories WHERE name = 'PLN'), 51800, 51800, 51800, 'Digiflazz', 'Umum', 'masukkan nomor meter/id pelanggan', 0, 'active', TRUE, TRUE),
    ('s10', 'Telkomsel 10.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 10200, 10200, 10200, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('s100', 'Telkomsel 100.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 98525, 98525, 98525, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('s15', 'Telkomsel 15.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 14850, 14850, 14850, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('s20', 'Telkomsel 20.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 19960, 19960, 19960, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('s25', 'Telkomsel 25.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 24619, 24619, 24619, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('s30', 'Telkomsel 30.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 29730, 29730, 29730, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('s5', 'Telkomsel 5.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 5155, 5155, 5155, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('s50', 'Telkomsel 50.000', 'TELKOMSEL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 49310, 49310, 49310, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('shopee100', 'SHOPEE PAY 100.000', 'SHOPEE PAY', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 100250, 100250, 100250, 'Digiflazz', 'Umum', 'SHOPEE PAY 100.000', 0, 'active', TRUE, TRUE),
    ('shopee50', 'SHOPEE PAY 50.000', 'SHOPEE PAY', 'E-Money', (SELECT id FROM categories WHERE name = 'E-Money'), 51000, 51000, 51000, 'Digiflazz', 'Umum', 'SHOPEE PAY 50.000', 0, 'active', TRUE, TRUE),
    ('sm10', 'Smartfren 10.000', 'SMARTFREN', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 10005, 10005, 10005, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('smdu1', 'Smartfren Data Unlimited Harian 1 GB Berlaku 7 Hari', 'SMARTFREN', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 22010, 22010, 22010, 'Digiflazz', 'Unlimited', 'Batas pemakaian wajar 1GB/hari, Unlimited 24 Jam, Gratis Nelpon ke sesama smartfren, Berlaku 7 hari', 0, 'active', TRUE, TRUE),
    ('smdu2', 'Smartfren Data Unlimited Harian 2 GB Berlaku 28 Hari', 'SMARTFREN', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 90525, 90525, 90525, 'Digiflazz', 'Unlimited', 'Batas pemakaian wajar 2GB/hari, Unlimited 24 Jam, Gratis Nelpon ke sesama smartfren, Berlaku 28 hari', 0, 'active', TRUE, TRUE),
    ('t10', 'Three 10.000', 'TRI', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 10125, 10125, 10125, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('t20', 'Three 20.000', 'TRI', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 19622, 19622, 19622, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('t5', 'Three 5.000', 'TRI', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 5505, 5505, 5505, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('tacthappys', 'Aktivasi Perdana Tri Happy S+ 30 Hari', 'TRI', 'Aktivasi Perdana', (SELECT id FROM categories WHERE name = 'Aktivasi Perdana'), 20500, 20500, 20500, 'Digiflazz', 'Happy', 'Aktivasi Perdana Tri Happy S+ 30 Hari', 0, 'active', TRUE, TRUE),
    ('tactive4m', 'Tri Tambah Masa Aktif Kartu  4 Bulan', 'TRI', 'Masa Aktif', (SELECT id FROM categories WHERE name = 'Masa Aktif'), 3655, 3655, 3655, 'Digiflazz', 'Umum', 'Hanya menambah masa aktif kartu 4 bulan, bisa akumulasi', 0, 'active', TRUE, TRUE),
    ('vax1', 'Aktivasi Voucher Axis 1 GB 1 Hari', 'AXIS', 'Aktivasi Voucher', (SELECT id FROM categories WHERE name = 'Aktivasi Voucher'), 6320, 6320, 6320, 'Digiflazz', 'Aigo', 'AIGO Mini Bronet 24Jam 1GB + Kuota di Kota-mu 1hr', 0, 'active', TRUE, TRUE),
    ('vax2', 'Aktivasi Voucher Axis 3 GB 3 Hari', 'AXIS', 'Aktivasi Voucher', (SELECT id FROM categories WHERE name = 'Aktivasi Voucher'), 9350, 9350, 9350, 'Digiflazz', 'Aigo', 'AIGO Mini 3GB + Kuota di Kota-mu 3hr', 0, 'active', TRUE, TRUE),
    ('vflexs', 'Aktivasi Voucher XL Xtra Combo Flex S 28 Hari', 'XL', 'Aktivasi Voucher', (SELECT id FROM categories WHERE name = 'Aktivasi Voucher'), 31908, 31908, 31908, 'Digiflazz', 'Xtra Combo Flex', 'Xtra Combo Flex S', 0, 'active', TRUE, TRUE),
    ('vs2g5d', 'Voucher Telkomsel 2.5 GB 5 Hari (Jawa Barat)', 'TELKOMSEL', 'Voucher', (SELECT id FROM categories WHERE name = 'Voucher'), 12985, 12985, 12985, 'Digiflazz', 'Jawa Barat', 'Voucher paket Internet 1 GB & 1.5 GB Lokal Internet Jawa berlaku selama 5 Hari', 0, 'active', TRUE, TRUE),
    ('x10', 'Xl 10.000', 'XL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 10830, 10830, 10830, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('x5', 'Xl 5.000', 'XL', 'Pulsa', (SELECT id FROM categories WHERE name = 'Pulsa'), 5900, 5900, 5900, 'Digiflazz', 'Umum', 'Reguler', 0, 'active', TRUE, TRUE),
    ('yellow1', 'Indosat Yellow 1 GB 1 Hari', 'INDOSAT', 'Data', (SELECT id FROM categories WHERE name = 'Data'), 5755, 5755, 5755, 'Digiflazz', 'Yellow', 'Online Gaspol 1GB 1 Hari', 0, 'active', TRUE, TRUE)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    brand = EXCLUDED.brand,
    category = EXCLUDED.category,
    category_id = EXCLUDED.category_id,
    price = EXCLUDED.price,
    price_api = EXCLUDED.price_api,
    platform_price = EXCLUDED.platform_price,
    status = EXCLUDED.status,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE products DROP COLUMN IF EXISTS brand;
DELETE FROM products WHERE provider = 'Digiflazz';
