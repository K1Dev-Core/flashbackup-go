# Flash Drive Backup CLI

โปรแกรมนี้เป็นโปรแกรมสำรองไฟล์จาก Flash Drive ไปเก็บไว้ในโฟลเดอร์ปลายทางบนเครื่องคอมพิวเตอร์

โปรแกรมทำงานผ่านคำสั่งใน Terminal และใช้ GORM เชื่อมต่อกับ SQLite เพื่อเก็บประวัติไฟล์ที่ย้ายแล้ว

## โปรแกรมทำอะไรได้บ้าง

- ตั้งค่าโฟลเดอร์ต้นทางด้วย `/source`
- ตั้งค่าโฟลเดอร์ปลายทางด้วย `/dest`
- ดูไฟล์ใน source และ destination
- ดูประวัติไฟล์จากฐานข้อมูล
- ย้ายไฟล์ทีละไฟล์หรือย้ายทั้งหมด
- ป้องกันการย้ายไฟล์ซ้ำ
- ตรวจสอบไฟล์ที่หายไปจาก destination
- ลบไฟล์จริงและลบประวัติในฐานข้อมูล
- แสดงเวลาที่ใช้ในการย้ายไฟล์ทั้งหมด

## โครงสร้างโปรเจกต์

```text
.
├── go.mod
├── go.sum
├── src/
│   ├── main.go       จุดเริ่มต้นโปรแกรมและการเปิดฐานข้อมูล
│   ├── app.go        รับคำสั่งจากผู้ใช้
│   ├── models.go     Model ของตาราง files และ settings
│   ├── settings.go   จัดการ source และ destination
│   ├── list.go        แสดงรายการไฟล์และตรวจสอบ integrity
│   ├── backup.go      ย้ายไฟล์และบันทึกประวัติ
│   ├── delete.go      ลบไฟล์และประวัติ
│   └── paths.go       ตรวจสอบ path และชื่อไฟล์
└── docs/              โจทย์และเอกสารประกอบการเรียน
```

ใช้ root `go.mod` เพียงไฟล์เดียว ทุกไฟล์ใน `src` อยู่ใน package เดียวกัน จึง build เป็นโปรแกรมเดียวได้

## สิ่งที่ต้องมี

- Go 1.26.4 หรือใหม่กว่า
- macOS, Linux หรือ Windows
- GORM และ SQLite ซึ่งติดตั้งผ่าน `go.mod` ให้แล้ว

## วิธีรันโปรแกรม

รันจาก root ของโปรเจกต์:

```bash
go run ./src -db /tmp/flashbackup.db
```

ถ้าไม่ใส่ `-db` โปรแกรมจะใช้ไฟล์ `crossplatform.db` ในโฟลเดอร์ที่รันคำสั่ง

ถ้าต้องการ build เป็นไฟล์ executable:

```bash
go build -o flashbackup ./src
./flashbackup -db /tmp/flashbackup.db
```

ไฟล์ที่ build ข้ามแพลตฟอร์มไว้จะอยู่ใน `dist/`:

```text
flashbackup-darwin-amd64
flashbackup-darwin-arm64
flashbackup-linux-amd64
flashbackup-linux-arm64
flashbackup-windows-amd64.exe
flashbackup-windows-arm64.exe
```

GitHub Actions จะรัน test, vet และ build ทั้ง 6 target นี้ทุกครั้งที่ push หรือเปิด pull request

## คำสั่งที่ใช้

```text
/help
/source <path>
/dest <path>
/list source
/list dest
/list db
/list db <path>
/move <file1>, <file2>
/move all
/check
/delete dest <file1>, <file2>
/delete dest all
/exit
```

## ทดสอบด้วย mockdata

ไฟล์ `mockdata` เป็นโปรแกรมจากชุดโจทย์ ใช้สร้างไฟล์จำลองสำหรับทดสอบ โดยปกติจะอยู่ใน `docs/midterm` และไม่ได้ถูก upload ขึ้น Git เพราะเป็น binary และไฟล์ทดสอบขนาดใหญ่

```bash
cd /Users/k1god/Downloads/cp/docs/midterm
chmod +x mockdata
./mockdata
```

โปรแกรมจะสร้าง:

- `data10` จำนวน 10 ไฟล์
- `data1000` จำนวน 1,000 ไฟล์

แนะนำให้เริ่มจาก `data10` ก่อน

## ทดสอบโปรแกรมตัวอย่างอาจารย์

ตัวอย่างและโปรแกรมของเราจะย้ายไฟล์ออกจาก source ดังนั้นต้องใช้ source คนละชุดกัน

```bash
rm -rf /tmp/flashbackup-test
mkdir -p /tmp/flashbackup-test/example/source
mkdir -p /tmp/flashbackup-test/example/dest
cp -R data10/. /tmp/flashbackup-test/example/source/
cp flashbackup /tmp/flashbackup-test/example/
cd /tmp/flashbackup-test/example
./flashbackup
```

ในโปรแกรมตัวอย่าง พิมพ์:

```text
/help
/source /tmp/flashbackup-test/example/source
/dest /tmp/flashbackup-test/example/dest
/list source
/move all
/list dest
/list db
/check
/exit
```

## ทดสอบโปรแกรมของเรา

เตรียม source อีกชุดหนึ่ง:

```bash
mkdir -p /tmp/flashbackup-test/ours/source
mkdir -p /tmp/flashbackup-test/ours/dest
cp -R data10/. /tmp/flashbackup-test/ours/source/
cd /Users/k1god/Downloads/cp
go run ./src -db /tmp/flashbackup-test/ours/flashbackup.db
```

ในโปรแกรม พิมพ์:

```text
/list db
/source /tmp/flashbackup-test/ours/source
/dest /tmp/flashbackup-test/ours/dest
/list source
/list dest
/list db
/move all
/list source
/list dest
/list db
/check
/exit
```

ผลที่ควรได้:

- ก่อนย้าย source มี 10 ไฟล์
- หลังย้าย source ว่างลง
- destination มี 10 ไฟล์
- `/list db` แสดงประวัติ 10 รายการ
- `/check` แจ้งว่า integrity ผ่าน
- มีเวลารวมการย้ายแสดงเป็น milliseconds

## ทดสอบไฟล์ซ้ำ

หลังจากย้ายไฟล์เสร็จแล้ว ให้เลือกชื่อไฟล์จาก `/list dest` แล้ว copy กลับไปที่ source:

```bash
cp /tmp/flashbackup-test/ours/dest/<ชื่อไฟล์> \
   /tmp/flashbackup-test/ours/source/<ชื่อไฟล์>
```

เปิดโปรแกรมด้วย DB เดิม แล้วสั่ง:

```text
/source /tmp/flashbackup-test/ours/source
/dest /tmp/flashbackup-test/ours/dest
/move <ชื่อไฟล์>
```

โปรแกรมต้องแจ้งว่าไฟล์มีอยู่ใน destination แล้ว และต้องไม่เขียนทับไฟล์เดิม

ถ้าลบไฟล์จาก destination แต่ยังไม่ลบประวัติใน DB แล้วลอง `/move` อีกครั้ง โปรแกรมจะตรวจพบชื่อไฟล์ซ้ำจากประวัติใน DB

## ทดสอบไฟล์หายและการลบ

ลบไฟล์หนึ่งไฟล์จาก destination ด้วยคำสั่งภายนอก แล้วเปิดโปรแกรมด้วย DB เดิม:

```text
/check
```

โปรแกรมต้องแสดงชื่อไฟล์ที่หายไปด้วย `Missing:`

ทดสอบลบประวัติและไฟล์:

```text
/delete dest <ชื่อไฟล์>
/list dest
/list db
/delete dest all
/exit
```

## ตรวจสอบโค้ดก่อนส่ง

```bash
go test ./...
go vet ./...
go build -o /tmp/flashbackup ./src
```
