# การจัดการข้อมูลใน GORM (Data Manipulation)

บทนี้จะเป็นการจัดการเรื่องการเพิ่ม ลบ และแก้ไขข้อมูล (Create, Update, Delete) ผ่าน GORM โดย GORM ได้เตรียมฟังก์ชันที่ใช้งานง่ายไว้ให้เรียบร้อยแล้ว

---

## 1. การเพิ่มข้อมูล (Create)

ในการสร้างข้อมูลใหม่ในตาราง จะใช้ฟังก์ชัน `Create()` โดยค่าที่ส่งไปจะต้องเป็น **Address/Pointer** ของ Object ของ Model เสมอ

main.go
```go

	restaurant := models.Restaurant{
		Name:           "ครัวข้าวแกง",
		RestaurantType: "อาหารไทย",
		Location:       "ข้างวัดศรีสวัสดิ์ อ.เมือง จ.มหาสารคาม",
	}
	err = db.Create(&restaurant).Error
	if err != nil {
		log.Fatal("Failed to create restaurant:", err)
		return
	}
	log.Println("Restaurant created successfully.")

	foodItem := models.FoodItem{
		Name:         "ข้าวแกงไข่ดาว",
		FoodType:     "อาหารไทย",
		Price:        100,
		IsAvailable:  true,
		RestaurantID: restaurant.ID,
	}
	result := db.Create(&foodItem)
	if result.Error != nil {
		log.Fatal("Failed to create food item:", err)
		return
	}
	log.Println("Food item created successfully.")
	log.Println("Rows affected:", result.RowsAffected)
```

> [!NOTE]
> **ทำไมต้องส่ง Pointer (`&restaurant`) ไปให้ Create?**
> เพราะว่าระบบฐานข้อมูล (โดยเฉพาะฟิลด์ที่เป็น AUTO_INCREMENT) จะเจนค่า Primary Key คีย์หลักใหม่ขึ้นมาเมื่อข้อมูลบันทึกสำเร็จ GORM จะนำค่านั้นมาใส่กลับคืนในฟิลด์ ID (เช่น `restaurant.ID`) ให้อัตโนมัติ ทำให้สามารถเรียกใช้งาน ID ใหม่นั้นต่อได้ทันทีโดยไม่ต้องส่งคิวรีไปดึงซ้ำ
> และสามารถใช้ `result.RowsAffected` เพื่อตรวจสอบจำนวนแถวที่ได้รับผลกระทบ (ในกรณีนี้คือการสร้างสำเร็จ 1 แถว)

![[Pasted image 20260611072803.png]]
*ภาพแสดงข้อมูลร้านอาหาร "ครัวข้าวแกง" ที่ถูกเพิ่มเข้าไปในตาราง restaurant ของฐานข้อมูล SQLite ผ่านโปรแกรม DB Browser*

![[Pasted image 20260611072820.png]]
*ภาพแสดงข้อมูลเมนูอาหาร "ข้าวแกงไข่ดาว" ที่เพิ่มเข้าไปใหม่พร้อมระบุความสัมพันธ์กับ restaurantId 32*

![[Pasted image 20260611073021.png]]
*ภาพแสดง Log การทำงานบน Terminal ยืนยันการเพิ่มข้อมูลทั้งฝั่งร้านอาหารและฝั่งเมนูอาหารสำเร็จเรียบร้อย*

---

## 2. การลบข้อมูล (Delete)

การสั่งลบข้อมูลด้วย GORM มีโครงสร้างคำสั่งดังนี้:
`Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)`

> [!IMPORTANT]
> **ข้อควรระวัง:** การใช้คำสั่ง Delete ใน GORM จำเป็นต้องระบุคีย์หลัก (Primary Key) เสมอ หรือไม่ก็ต้องใส่เงื่อนไขความปลอดภัยไว้ ไม่อย่างนั้น GORM จะไม่ยอมลบข้อมูลให้เนื่องจากมีระบบป้องกันการลบข้อมูลทั้งหมดในตารางโดยไม่ตั้งใจ (Global Delete Safety)

วิธีการลบข้อมูลสามารถทำได้ 3 วิธีหลักๆ ดังนี้:

### วิธีที่ 1: ลบโดยส่งประเภทโมเดลเปล่าๆ และใส่ Primary Key (ไม่ต้องใช้ตัวแปรแบบ Pointer `&`)

วิธีนี้ทำได้ง่ายโดยส่งอินสแตนซ์โมเดลที่ต้องการระบุตาราง ตามด้วยค่า Primary Key ไปยังอาร์กิวเมนต์ตัวที่สอง

main.go
```go
	err = db.Delete(models.Restaurant{}, 32).Error
	if err != nil {
		log.Fatal("Failed to delete restaurant:", err)
		return
	}
	log.Println("Restaurant deleted successfully.")
```

### วิธีที่ 2: ลบโดยโยนตัวแปร Object ที่มีการกำหนดค่า ID อยู่ภายในแล้ว (ต้องใส่ `&` หน้าออบเจกต์)

ถ้ามีตัวแปรออบเจกต์ที่เตรียมข้อมูลไว้อยู่แล้ว หรือดึงมาจากฐานข้อมูลแล้วมีคีย์หลักอยู่ข้างใน สามารถโยน Address ของออบเจกต์นั้นไปสั่งลบได้เลย

main.go
```go

	restaurant := models.Restaurant{
		ID: 33,
	}
	err = db.Delete(&restaurant).Error
	if err != nil {
		log.Fatal("Failed to delete restaurant:", err)
		return
	}
	log.Println("Restaurant deleted successfully.")
```

### วิธีที่ 3: ลบข้อมูลแบบระบุเงื่อนไขตามฟิลด์ต่างๆ (Delete with condition)

ถ้าต้องการลบแถวข้อมูลตามคีย์เวิร์ด หรือเงื่อนไขที่ต้องการ สามารถโยน SQL query template และพารามิเตอร์ต่อท้ายได้เลย

main.go
```go
	err = db.Delete(models.Restaurant{}, "id = ?", 33).Error
	if err != nil {
		log.Fatal("Failed to delete restaurant:", err)
		return
	}
	log.Println("Restaurant deleted successfully.")
```

---

## 3. การแก้ไขข้อมูล (Update)

เมื่อต้องการอัปเดตข้อมูลเดิมที่มีอยู่ในระบบ สามารถทำได้สองแบบหลักๆ ตามลักษณะการใช้งาน:

### แบบที่ 1: อัปเดตข้อมูลทุกคอลัมน์ของเรคอร์ด (ใช้ฟังก์ชัน Save)

ฟังก์ชัน `Save()` จะทำการบันทึกความเปลี่ยนแปลงทุกคอลัมน์ของออบเจกต์กลับไปยังตาราง โดยอาศัยคีย์หลัก (Primary Key) ในตัวแปรเพื่อชี้ว่ากำลังแก้ไขแถวไหนอยู่

ขั้นตอนการเขียนโค้ด:
1. ดึงเรคอร์ดตัวเดิมจากฐานข้อมูลขึ้นมาก่อน
2. แก้ไขฟิลด์ของออบเจกต์ที่ต้องการเปลี่ยน
3. ส่งออบเจกต์ตัวนั้นผ่าน `db.Save(&object)`

![[Pasted image 20260611073623.png]]
*ภาพแสดงข้อมูลแถวที่ 30 ร้าน "บ้านโรตีชาชัก" ในฐานข้อมูล SQLite ก่อนทำการอัปเดตข้อมูลผ่านคำสั่ง Save*

main.go
```go
	restaurant := models.Restaurant{}
	db.Find(&restaurant, "id = ?", 30)
	log.Println("Restaurant found:", restaurant)
	restaurant.Name = "บ้านโรตีปักษ์ใต้"
	err = db.Save(&restaurant).Error
	if err != nil {
		log.Fatal("Failed to update restaurant:", err)
		return
	}
	log.Println("Restaurant updated successfully.")
```

![[Pasted image 20260611073645.png]]
*ภาพแสดงข้อมูลร้านอาหารแถวที่ 30 ที่ถูกอัปเดตชื่อร้านใหม่เป็น "บ้านโรตีปักษ์ใต้" สำเร็จเรียบร้อย*

### แบบที่ 2: อัปเดตเฉพาะคอลัมน์ที่กำหนด (ใช้ฟังก์ชัน Where และ Update)

ถ้าไม่อยากส่งค่าแก้ไขไปทุกคอลัมน์ (เช่น เปลี่ยนแค่บางฟิลด์ข้อมูล) สามารถเจาะจงเฉพาะตัวคอลัมน์นั้นได้โดยจับคู่ฟังก์ชัน `Where()` ร่วมกับ `Update()` หรือ `Updates()` เพื่อระบุเป้าหมายอย่างปลอดภัย

main.go
```go
	restaurant := models.Restaurant{}
	err = db.Model(&restaurant).Where("id = ?", 30).Update("name", "บ้านโรตีปักษ์ใต้ 2").Error
	if err != nil {
		log.Fatal("Failed to update restaurant:", err)
		return
	}
	log.Println("Restaurant updated successfully.")
```

![[Pasted image 20260611073949.png]]
*ภาพแสดงผลลัพธ์การอัปเดตเฉพาะฟิลด์ชื่อร้านผ่านคำสั่ง Update เป็น "บ้านโรตีปักษ์ใต้ 2"*

> [!TIP]
> - หากต้องการแก้ไขแค่คอลัมน์เดียว: ใช้ `.Update("ชื่อคอลัมน์", ค่าใหม่)`
> - หากต้องการแก้ไขพร้อมกันหลายคอลัมน์: ใช้ `.Updates(map[string]interface{}{"name": "บ้านโรตีปักษ์ใต้ 2", "location": "ม.มหาสารคาม"})` หรือส่งผ่าน Struct `.Updates(model.Restaurant{name: "บ้านโรตีปักษ์ใต้ 2"})` ได้เช่นกัน
