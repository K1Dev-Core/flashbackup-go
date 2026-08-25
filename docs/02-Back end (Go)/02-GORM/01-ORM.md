# Object Relational Mapping (ORM) คืออะไร?

มาทำความเข้าใจแนวคิดของ ORM กันก่อน!
**ORM (Object Relational Mapping)** คือ เทคนิคที่ใช้เชื่อมช่องว่างระหว่างสองฝั่งที่ใช้ภาษาต่างกัน:
1. **ฝั่งโปรแกรมมิ่ง (Object/Struct ในภาษา Go)** ที่ทำงานด้วย Struct, Fields และ Methods
2. **ฝั่งฐานข้อมูล (Relational Database)** ที่ทำงานด้วย Tables, Rows และ Columns

คิดง่ายๆ ว่า ORM ทำหน้าที่เป็น **"ล่ามแปลภาษา"** คอยแปลง Struct ในภาษา Go ให้กลายเป็น Table ในฐานข้อมูล (หรือในทางกลับกันคือ ดึงข้อมูลจาก Table มาแปลงกลับเป็น Struct ใน Go ให้นำมาใช้งานได้ง่าย) โดยไม่จำเป็นต้องเขียนคำสั่ง SQL ดิบๆ (Raw SQL) เองตลอดเวลา

- **Mapping Concept:**
  - **Entity (หรือ Model ใน Go Struct)** $\rightarrow$ **Table (ตารางในฐานข้อมูล)**
  - **Field (ตัวแปรใน Struct)** $\rightarrow$ **Column (คอลัมน์ในตาราง)**
  - **Struct Instance (ตัวแปร Object ที่สร้างจาก Struct)** $\rightarrow$ **Row (แถวข้อมูลในตาราง)**

![[Pasted image 20250306100701.png]]
*ภาพแสดงแนวคิดการแมปข้อมูล (Mapping) ระหว่างออบเจกต์โครงสร้างข้อมูล (Struct) ในภาษา Go กับตารางข้อมูล (Table) ของ Relational Database*

---

# GORM (Go ORM)

สำหรับภาษา Go ตัว ORM ที่ได้รับความนิยมใช้งานมากที่สุดตัวหนึ่งคือ **GORM**

- **เว็บไซต์หลักอย่างเป็นทางการ:** [https://gorm.io](https://gorm.io)
- **เอกสารแพ็กเกจบน pkg.go.dev:** [https://pkg.go.dev/gorm.io/gorm](https://pkg.go.dev/gorm.io/gorm)

![[Pasted image 20250306101134.png]]
*ภาพแสดงเว็บไซต์อย่างเป็นทางการของ GORM (gorm.io) ซึ่งรวบรวมเอกสารแนะนำและวิธีการใช้งานระบบ ORM ในภาษา Go*

## การติดตั้ง GORM (Get Library)


## สร้าง Project

```sh
mkdir backend
cd backend
code .

go mod init backend
```

![[Pasted image 20260608190321.png]]
*ภาพแสดงขั้นตอนการสร้างโฟลเดอร์โปรเจกต์ใหม่และสั่งรัน go mod init เพื่อกำหนดขอบเขตของโมดูลใน VSCode Terminal*

ก่อนที่จะเริ่มเขียนโค้ด ต้องติดตั้ง GORM Library เข้ามาในโปรเจกต์ รันคำสั่งนี้ใน Terminal:

```sh
go get gorm.io/gorm
```

### ตัวขับเคลื่อนฐานข้อมูล (Dialector / Database Drivers)

เนื่องจาก GORM เป็นตัวกลางคอยจัดการเรื่อง Logic แต่ฐานข้อมูลแต่ละตัว (เช่น MySQL, PostgreSQL, SQLite, SQL Server) มีวิธีการเรียกใช้หรือไวยากรณ์ที่แตกต่างกันเล็กน้อย GORM จึงต้องการตัวช่วยเฉพาะทางที่เรียกว่า **Dialector** หรือ **Driver** เพื่อให้สามารถแปลงคำสั่งเป็น SQL ที่ฐานข้อมูลนั้นๆ เข้าใจได้

สามารถเข้าไปดูรายชื่อฐานข้อมูลที่รองรับได้ที่: [Connecting to the Database Docs](https://gorm.io/docs/connecting_to_the_database.html)

สำหรับบทเรียนนี้จะใช้ **SQLite** ดังนั้นต้องดาวน์โหลด Driver สำหรับ SQLite มาด้วย:

- **SQLite Driver Package:** https://pkg.go.dev/github.com/glebarez/sqlite

```sh
go get github.com/glebarez/sqlite
```

---

# SQLite

https://sqlitebrowser.org/


```sql
DROP TABLE IF EXISTS "food_review";
DROP TABLE IF EXISTS "food_item";
DROP TABLE IF EXISTS "restaurant";
DROP TABLE IF EXISTS "foodie_user";

CREATE TABLE "foodie_user" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255),
    "name" VARCHAR(100) NOT NULL
);

CREATE TABLE "restaurant" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" VARCHAR(255) NOT NULL,
    "restaurantType" VARCHAR(100) NOT NULL,
    "location" TEXT NOT NULL
);

CREATE TABLE "food_item" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" VARCHAR(255) NOT NULL,
    "foodType" VARCHAR(100) NOT NULL,
    "imagePath" VARCHAR(512),
    "price" DECIMAL(10, 2) NOT NULL,
    "isAvailable" BOOLEAN DEFAULT 1,
    "restaurantId" INTEGER NOT NULL,
    FOREIGN KEY ("restaurantId") REFERENCES "restaurant" ("id") ON DELETE CASCADE
);

CREATE TABLE "food_review" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "rating" INTEGER NOT NULL,
    "comment" TEXT,
    "imagePath" VARCHAR(512),
    "reviewedAt" DATETIME DEFAULT CURRENT_TIMESTAMP,
    "userId" INTEGER NOT NULL,
    "foodItemId" INTEGER NOT NULL,
    FOREIGN KEY ("userId") REFERENCES "foodie_user" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("foodItemId") REFERENCES "food_item" ("id") ON DELETE CASCADE
);


INSERT INTO "foodie_user" ("email", "name") VALUES
('narin.p@example.com', 'นรินทร์ พรมมา'),
('siriporn.k@example.com', 'ศิริพร คำดี'),
('thanawat.s@example.com', 'ธนวัฒน์ สายทอง'),
('patchara.m@example.com', 'พัชรา มณีวงศ์'),
('kittipong.r@example.com', 'กิตติพงศ์ รัตนา'),
('oranuch.t@example.com', 'อรอนงค์ ทองสุข'),
('preecha.j@example.com', 'ปรีชา ใจดี'),
('chutima.n@example.com', 'ชุติมา นาคทอง'),
('sompong.w@example.com', 'สมพงษ์ วัฒนะ'),
('wanida.l@example.com', 'วนิดา ลีลาวัฒน์'),
('apichai.y@example.com', 'อภิชัย ยิ้มแย้ม'),
('kamonchanok.b@example.com', 'กมลชนก บุญมาก'),
('rattanaporn.c@example.com', 'รัตนาภรณ์ ชัยดี'),
('phuwadon.p@example.com', 'ภูวดล พิทักษ์'),
('yarinda.s@example.com', 'ญารินดา สุขใจ'),
('ekkarat.d@example.com', 'เอกรัฐ ดวงดี'),
('kanchana.h@example.com', 'กาญจนา หอมหวล'),
('teerapat.v@example.com', 'ธีรภัทร วงศ์ศรี'),
('napatsorn.a@example.com', 'ณภัสสร อ่อนน้อม'),
('somchai.p@example.com', 'สมชาย ประเสริฐ'),
('chonnikan.k@example.com', 'ชนนิกานต์ คงคา'),
('worawut.f@example.com', 'วรวุฒิ ฟุ้งเฟื่อง'),
('pimchanok.r@example.com', 'พิมพ์ชนก รุ่งเรือง'),
('natthawat.s@example.com', 'ณัฐวัฒน์ ศรีบุญ'),
('supansa.t@example.com', 'สุพรรษา เทพดี'),
('kritsada.l@example.com', 'กฤษดา ลาภเพิ่ม'),
('arisa.m@example.com', 'อริสา มีชัย'),
('boonyarit.n@example.com', 'บุญฤทธิ์ นาคำ'),
('thidarat.w@example.com', 'ธิดารัตน์ วงษ์คำ'),
('pakorn.g@example.com', 'ภากร เกรียงไกร');

INSERT INTO "restaurant" ("name", "restaurantType", "location") VALUES
('ครัวบ้านสวน', 'อาหารไทย', '89/3 ถนนลาดพร้าว แขวงจอมพล เขตจตุจักร กรุงเทพมหานคร'),
('ก๋วยเตี๋ยวเรือท่าช้าง', 'อาหารเส้น', '15 ซอยท่าช้าง 2 อำเภอเมืองนนทบุรี จังหวัดนนทบุรี'),
('ข้าวแกงปักษ์ใต้แม่อร', 'อาหารไทย', '221 ถนนพัฒนาการ แขวงสวนหลวง เขตสวนหลวง กรุงเทพมหานคร'),
('ส้มตำแซ่บนัว', 'อาหารไทย', '44/8 ถนนรัชดาภิเษก แขวงดินแดง เขตดินแดง กรุงเทพมหานคร'),
('ชาบูนายพล', 'ปิ้งย่าง', '102 หมู่ 5 ตำบลบางรักพัฒนา อำเภอบางบัวทอง จังหวัดนนทบุรี'),
('ปิ้งย่างริมคลอง', 'ปิ้งย่าง', '76 ถนนกาญจนาภิเษก แขวงบางแค เขตบางแค กรุงเทพมหานคร'),
('บ้านกะเพราพริกแห้ง', 'อาหารไทย', '58/1 ซอยรามคำแหง 24 แขวงหัวหมาก เขตบางกะปิ กรุงเทพมหานคร'),
('ข้าวมันไก่เฮียตี๋', 'อาหารไทย', '12 ถนนประชาชื่น แขวงบางซื่อ เขตบางซื่อ กรุงเทพมหานคร'),
('บะหมี่เกี๊ยวเจริญผล', 'อาหารเส้น', '399 ถนนสุขุมวิท แขวงคลองตันเหนือ เขตวัฒนา กรุงเทพมหานคร'),
('ร้านขนมจีนยายแสง', 'อาหารเส้น', '9 หมู่ 2 ตำบลบางกร่าง อำเภอเมืองนนทบุรี จังหวัดนนทบุรี'),
('ซีฟู้ดบางปูฟ้าใส', 'อาหารไทย', '128 ถนนสุขุมวิท ตำบลท้ายบ้าน อำเภอเมืองสมุทรปราการ จังหวัดสมุทรปราการ'),
('ครัวฮาลาลริมทาง', 'อาหารนานาชาติ', '67 ถนนนวมินทร์ แขวงคลองกุ่ม เขตบึงกุ่ม กรุงเทพมหานคร'),
('วีแกนกรีนเดย์', 'คาเฟ่', '31 ซอยสุขุมวิท 49 แขวงคลองตันเหนือ เขตวัฒนา กรุงเทพมหานคร'),
('สเต๊กเฮาส์ลุงป้อม', 'อาหารนานาชาติ', '22/5 ถนนพระราม 2 แขวงแสมดำ เขตบางขุนเทียน กรุงเทพมหานคร'),
('คาเฟ่ขนมหวานบ้านหอม', 'คาเฟ่', '45 ถนนราชพฤกษ์ แขวงบางระมาด เขตตลิ่งชัน กรุงเทพมหานคร'),
('ร้านผัดไทยโบราณ', 'อาหารเส้น', '88 ซอยเพชรเกษม 48 แขวงบางด้วน เขตภาษีเจริญ กรุงเทพมหานคร'),
('ครัวปลาเผาแม่กลอง', 'อาหารไทย', '14/2 ถนนพระราม 5 ตำบลบางไผ่ อำเภอเมืองนนทบุรี จังหวัดนนทบุรี'),
('กุ้งอบวุ้นเส้นคลองสาน', 'อาหารนานาชาติ', '73 ถนนเจริญนคร แขวงคลองสาน เขตคลองสาน กรุงเทพมหานคร'),
('ข้าวขาหมูเจริญชัย', 'อาหารไทย', '55/7 ถนนพระราม 3 แขวงบางโพงพาง เขตยานนาวา กรุงเทพมหานคร'),
('โจ๊กหมูทรงเครื่อง', 'อาหารไทย', '19 ซอยอารีย์สัมพันธ์ แขวงสามเสนใน เขตพญาไท กรุงเทพมหานคร'),
('เป็ดย่างฮ่องเต้', 'อาหารนานาชาติ', '41 ถนนโชคชัย 4 แขวงลาดพร้าว เขตลาดพร้าว กรุงเทพมหานคร'),
('ครัวญี่ปุ่นซากุระไทย', 'อาหารนานาชาติ', '110 ถนนสุขุมวิท 63 แขวงคลองตันเหนือ เขตวัฒนา กรุงเทพมหานคร'),
('พิซซ่าเตาถ่านบางเขน', 'อาหารนานาชาติ', '24 ถนนพหลโยธิน แขวงอนุสาวรีย์ เขตบางเขน กรุงเทพมหานคร'),
('เบอร์เกอร์บ้านเพื่อน', 'อาหารนานาชาติ', '6 ซอยลาดปลาเค้า 12 แขวงจรเข้บัว เขตลาดพร้าว กรุงเทพมหานคร'),
('หมูกระทะคุณย่า', 'ปิ้งย่าง', '77 ถนนรัตนาธิเบศร์ ตำบลเสาธงหิน อำเภอบางใหญ่ จังหวัดนนทบุรี'),
('ครัวอินเดียมาซาลา', 'อาหารนานาชาติ', '39 ซอยสุขุมวิท 11 แขวงคลองเตยเหนือ เขตวัฒนา กรุงเทพมหานคร'),
('ข้าวหน้าเนื้อญี่ปุ่นย่านอโศก', 'อาหารนานาชาติ', '21 ถนนสุขุมวิท 21 แขวงคลองเตยเหนือ เขตวัฒนา กรุงเทพมหานคร'),
('ร้านติ่มซำยามเช้า', 'อาหารนานาชาติ', '13 ถนนประชาราษฎร์สาย 1 แขวงบางซื่อ เขตบางซื่อ กรุงเทพมหานคร'),
('เฝอเวียดนามต้นตำรับ', 'อาหารเส้น', '92 ซอยเอกมัย 12 แขวงคลองตันเหนือ เขตวัฒนา กรุงเทพมหานคร'),
('บ้านโรตีชาชัก', 'คาเฟ่', '5 ถนนเพชรบุรี แขวงมักกะสัน เขตราชเทวี กรุงเทพมหานคร');

INSERT INTO "food_item" ("name", "foodType", "imagePath", "price", "isAvailable", "restaurantId") VALUES
('ผัดกะเพราหมูสับไข่ดาว', 'ข้าว', '/images/food.png', 65.00, 1, 1),
('ก๋วยเตี๋ยวเรือน้ำตกเนื้อ', 'เส้น', '/images/food.png', 75.00, 1, 2),
('แกงเหลืองปลากะพง', 'กับข้าว', '/images/food.png', 120.00, 1, 3),
('ส้มตำไทยไข่เค็ม', 'กับข้าว', '/images/food.png', 70.00, 1, 4),
('ชุดชาบูหมูรวม', 'ชุดเซต', '/images/food.png', 299.00, 1, 5),
('ชุดปิ้งย่างหมูคุโรบูตะ', 'ชุดเซต', '/images/food.png', 359.00, 1, 6),
('ข้าวกะเพราเนื้อพริกแห้ง', 'ข้าว', '/images/food.png', 89.00, 1, 7),
('ข้าวมันไก่ต้มพิเศษ', 'ข้าว', '/images/food.png', 60.00, 1, 8),
('บะหมี่เกี๊ยวหมูแดง', 'เส้น', '/images/food.png', 65.00, 1, 9),
('ขนมจีนแกงเขียวหวานไก่', 'เส้น', '/images/food.png', 55.00, 1, 10),
('ปลากะพงนึ่งมะนาว', 'กับข้าว', '/images/food.png', 320.00, 1, 11),
('ข้าวหมกไก่', 'ข้าว', '/images/food.png', 85.00, 1, 12),
('ข้าวไรซ์เบอร์รีผัดเต้าหู้', 'ข้าว', '/images/food.png', 95.00, 1, 13),
('สเต๊กหมูพริกไทยดำ', 'กับข้าว', '/images/food.png', 189.00, 1, 14),
('ฮันนี่โทสต์ผลไม้รวม', 'ของหวาน', '/images/food.png', 149.00, 1, 15),
('ผัดไทยกุ้งสด', 'เส้น', '/images/food.png', 90.00, 1, 16),
('ปลาทับทิมเผาเกลือ', 'กับข้าว', '/images/food.png', 280.00, 1, 17),
('กุ้งอบวุ้นเส้นหม้อดิน', 'เส้น', '/images/food.png', 210.00, 1, 18),
('ข้าวขาหมูคากิ', 'ข้าว', '/images/food.png', 75.00, 1, 19),
('โจ๊กหมูไข่ออนเซ็น', 'ข้าว', '/images/food.png', 65.00, 1, 20),
('ข้าวหน้าเป็ดย่าง', 'ข้าว', '/images/food.png', 95.00, 1, 21),
('แซลมอนดงบุริ', 'ข้าว', '/images/food.png', 229.00, 1, 22),
('พิซซ่ามาร์เกริต้า', 'กับข้าว', '/images/food.png', 259.00, 1, 23),
('ชีสเบอร์เกอร์ดับเบิล', 'กับข้าว', '/images/food.png', 199.00, 1, 24),
('ชุดหมูกระทะรวมมิตร', 'ชุดเซต', '/images/food.png', 329.00, 1, 25),
('ไก่ทิกก้ามาซาลา', 'กับข้าว', '/images/food.png', 245.00, 1, 26),
('ข้าวหน้าเนื้อวากิวซอสหวาน', 'ข้าว', '/images/food.png', 269.00, 1, 27),
('ติ่มซำรวมเข่ง', 'ชุดเซต', '/images/food.png', 169.00, 1, 28),
('เฝอเนื้อสไลซ์', 'เส้น', '/images/food.png', 139.00, 1, 29),
('โรตีแกงเขียวหวานไก่', 'กับข้าว', '/images/food.png', 99.00, 1, 30);

INSERT INTO "food_review" ("rating", "comment", "imagePath", "userId", "foodItemId") VALUES
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 1, 1),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 2, 2),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 3, 3),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 4, 4),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 5, 5),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 6, 6),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 7, 7),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 8, 8),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 9, 9),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 10, 10),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 11, 11),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 12, 12),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 13, 13),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 14, 14),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 15, 15),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 16, 16),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 17, 17),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 18, 18),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 19, 19),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 20, 20),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 21, 21),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 22, 22),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 23, 23),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 24, 24),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 25, 25),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 26, 26),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 27, 27),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 28, 1),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 29, 2),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 30, 3),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 1, 4),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 2, 5),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 3, 6),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 4, 7),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 5, 8),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 6, 9),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 7, 10),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 8, 11),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 9, 12),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 10, 13),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 11, 14),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 12, 15),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 13, 16),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 14, 17),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 15, 18),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 16, 19),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 17, 20),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 18, 21),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 19, 22),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 20, 23),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 21, 24),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 22, 25),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 23, 26),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 24, 27),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 25, 1),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 26, 2),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 27, 3),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 28, 4),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 29, 5),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 30, 6),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 1, 7),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 2, 8),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 3, 9),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 4, 10),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 5, 11),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 6, 12),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 7, 13),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 8, 14),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 9, 15),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 10, 16),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 11, 17),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 12, 18),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 13, 19),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 14, 20),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 15, 21),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 16, 22),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 17, 23),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 18, 24),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 19, 25),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 20, 26),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 21, 27),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 22, 1),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 23, 2),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 24, 3),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 25, 4),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 26, 5),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 27, 6),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 28, 7),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 29, 8),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 30, 9),
(5, 'วัตถุดิบสด รสชาติกลมกล่อม', '/images/food.png', 1, 10),
(4, 'ปริมาณเหมาะสม ราคาสมเหตุสมผล', '/images/food.png', 2, 11),
(4, 'รสชาติกำลังดี ทานง่าย', '/images/food.png', 3, 12),
(3, 'เผ็ดกำลังพอดี กลิ่นหอมชัด', '/images/food.png', 4, 13),
(5, 'เสิร์ฟไว อาหารยังร้อน', '/images/food.png', 5, 14),
(4, 'หน้าตาน่าทาน รสชาติดี', '/images/food.png', 6, 15),
(5, 'ซอสเข้มข้น เข้ากับเมนู', '/images/food.png', 7, 16),
(3, 'เนื้อสัมผัสดี เคี้ยวง่าย', '/images/food.png', 8, 17),
(4, 'โดยรวมถูกใจ มีสั่งซ้ำแน่นอน', '/images/food.png', 9, 18),
(5, 'รสชาติคงที่ คุณภาพดี', '/images/food.png', 10, 19);
```


![[foodie.db]]
*ภาพอ้างอิงไฟล์ foodie.db ซึ่งเป็นตัวไฟล์ฐานข้อมูล SQLite สำเร็จรูปที่ได้จากการสร้างตารางและกรอกข้อมูลเริ่มต้น*


![[Pasted image 20260608193945.png]]
*ภาพแสดงโครงสร้างไฟล์ของโปรเจกต์ backend บน VSCode Explorer ซึ่งประกอบด้วยไฟล์ foodie.db, go.mod, go.sum และ main.go*

## รูปแบบการเขียน DSN (Data Source Name) สำหรับ SQLite

**DSN** หรือ String ที่ใช้กำหนดค่าเชื่อมต่อฐานข้อมูล มีโครงสร้างแบบนี้:

```text
dsn := "foodie.db"
```

---

# เริ่มต้นสร้างโปรเจกต์ใหม่ (New Project)

โปรเจกต์ตัวอย่างนี้จะใช้ชื่อว่า `backend`

ตัวอย่างไฟล์แรกสำหรับการเริ่มต้นเขียนโค้ดเพื่อทดสอบการเชื่อมต่อ:

**main.go**
```go
package main

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("foodie.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
		return nil, err
	}
	return db, nil
}

func main() {
	_, err := ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	log.Println("Database connection established successfully.")
}

```

> [!NOTE]
> ในเวอร์ชันจริงของ GORM คำสั่ง `gorm.Open` มักจะต้องส่งพารามิเตอร์ตัวที่สองเป็น Config เข้าไปด้วย เช่น `gorm.Open(dialactor, &gorm.Config{})` เพื่อปรับตั้งค่าความสามารถต่างๆ ของ GORM

เมื่อรันโค้ดด้านบนนี้แล้ว ถ้าทุกอย่างถูกต้องและฐานข้อมูลปลายทางเปิดให้เข้าใช้งานได้ตามปกติ ก็จะมีข้อความขึ้นแจ้งเตือนแบบนี้เลย!

![[Pasted image 20260608194501.png]]
*ภาพแสดงการทำงานของโปรแกรมในช่อง Output/Terminal ยืนยันผลการเชื่อมต่อฐานข้อมูล foodie.db สำเร็จลุล่วง*

---
.env
```env
DATABASE=foodie.db
```

อ่าน .env

```sh
go get github.com/joho/godotenv
```

main.go
```go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return nil, err
	}
	dbHost := os.Getenv("DATABASE")
	db, err := gorm.Open(sqlite.Open(dbHost), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
		return nil, err
	}
	return db, nil
}

func main() {
	_, err := ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	log.Println("Database connection established successfully.")
}

```