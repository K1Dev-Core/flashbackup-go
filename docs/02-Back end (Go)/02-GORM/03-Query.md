# การสืบค้นข้อมูล (Query) รูปแบบต่างๆ

GORM มีชุดฟังก์ชันสำเร็จรูปที่จะช่วยสร้างเงื่อนไข SQL ได้อย่างมีประสิทธิภาพ:
รายละเอียดเพิ่มเติมดูได้ที่: [GORM Query Documentation](https://gorm.io/docs/query.html)

- `Where()` $\rightarrow$ ใส่เงื่อนไขตามคอลัมน์ในฐานข้อมูล
- `Order()` $\rightarrow$ จัดเรียงลำดับข้อมูล
- `Limit()` $\rightarrow$ จำกัดจำนวนเรคอร์ดที่ต้องการดึง
- `Distinct()` $\rightarrow$ คัดกรองข้อมูลเฉพาะที่ไม่ซ้ำกัน
- `Find()` $\rightarrow$ ดึงข้อมูลหลายเรคอร์ดใส่ Slice
- `Take()` $\rightarrow$ ดึงข้อมูลมาเพียง 1 เรคอร์ดแบบสุ่ม (ไม่จัดเรียงลำดับ)
- `First()` $\rightarrow$ ดึงเรคอร์ดแรกสุด (เรียงตาม PK)
- `Last()` $\rightarrow$ ดึงเรคอร์ดสุดท้ายสุด (เรียงตาม PK)

---

### วิธีการใช้ Where()

การกรองข้อมูลด้วย `Where` สามารถส่งอาร์กิวเมนต์แบบ Parameter binding ด้วยเครื่องหมาย `?` เพื่อความปลอดภัยจากการถูกโจมตีผ่าน SQL Injection

```go
// ค้นหาผู้ใช้คนแรกที่ชื่อ jinzhu
db.Where("name = ?", "jinzhu").First(&user)

// ใช้ Struct ค้นหา (GORM จะตรวจสอบเฉพาะฟิลด์ที่มีค่า ไม่เป็นค่าเริ่มต้น)
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)

// ต่อเงื่อนไข Where หลายชั้นเข้าด้วยกัน
db.Where("name = ?", "jinzhu").Where("age <> ?", "20").First(&user)
```

ตัวอย่างทดสอบดึงข้อมูลเมนูอาหารที่มีคำว่า "บ้าน" อยู่ในชื่อร้านอาหาร:

main.go
```go
	foodItems := []models.FoodItem{}

	err = db.Joins("Restaurant").Where("restaurant.name like ?", "%บ้าน%").Find(&foodItems).Error
	if err != nil {
		log.Fatal("Failed to retrieve food items:", err)
		return
	}
	log.Printf("%v", foodItems)
```

![[Pasted image 20260609125241.png]]
*ภาพแสดง Log ผลลัพธ์จากการสืบค้นข้อมูลด้วย Where ร่วมกับ Joins เพื่อกรองเฉพาะอาหารของร้านที่มีคำว่า "บ้าน" ในชื่อร้าน*

---

### วิธีการใช้ Take()

`Take()` จะดึงข้อมูลขึ้นมาเพียงเรคอร์ดเดียวจากเงื่อนไข โดยใช้คำสั่งแบบสุ่มที่ไม่มีการพยายามจัดเรียงลำดับ (ดึงแถวแรกสุดที่ฐานข้อมูลสแกนเจอขึ้นมาทันที)

main.go
```go
	foodItems := []models.FoodItem{}
	foodItem := models.FoodItem{}

	err = db.Joins("Restaurant").Where("restaurant.name like ?", "%บ้าน%").Take(&foodItems).Error
	if err != nil {
		log.Fatal("Failed to retrieve food items:", err)
		return
	}
	err = db.Joins("Restaurant").Where("restaurant.name like ?", "%บ้าน%").Take(&foodItem).Error
	if err != nil {
		log.Fatal("Failed to retrieve food items:", err)
		return
	}

	log.Printf("%v", foodItems)
	log.Printf("%v", foodItem)
```

![[Pasted image 20260609125550.png]]
*ภาพแสดง Log ผลลัพธ์การรันคำสั่ง Take ซึ่งดึงข้อมูลเมนูอาหารขึ้นมาแสดงเพียง 1 รายการโดยไม่มีการจัดเรียง*

---

### วิธีการใช้ First() และ Last()

สองฟังก์ชันนี้ใช้ดึงข้อมูลขึ้นมา 1 ตัวคล้าย `Take()` แต่จะมีความพิเศษตรงการจัดเรียงข้อมูล:
- **`First()`**: ดึงข้อมูลตัวแรกสุด โดยเรียงลำดับคีย์หลักจากน้อยไปมาก (ORDER BY primary_key ASC)
- **`Last()`**: ดึงข้อมูลตัวสุดท้ายสุด โดยเรียงลำดับคีย์หลักจากมากไปน้อย (ORDER BY primary_key DESC)

> [!WARNING]
> ทั้ง `First()`, `Last()`, และ `Take()` (ในกรณีที่รับข้อมูลด้วยออบเจกต์เดี่ยว เช่น `&restaurant`) หากไม่พบข้อมูลตามเงื่อนไข GORM จะส่งคืน Error พิเศษคือ `gorm.ErrRecordNotFound` ซึ่งต้องทำการดักจับและตรวจสอบความผิดพลาดนี้ในระบบด้วย ต่างกับ `Find()` ที่หาไม่เจอก็จะได้ Slice ว่างและไม่ส่งคืน Error ออกมา

main.go
```go
	err = db.Joins("Restaurant").Where("restaurant.name like ?", "%บ้าน%").First(&foodItem).Error
	if err != nil {
		log.Fatal("Failed to retrieve food items:", err)
		return
	}

	err = db.Joins("Restaurant").Where("restaurant.name like ?", "%บ้าน%").Last(&foodItem).Error
	if err != nil {
		log.Fatal("Failed to retrieve food items:", err)
		return
	}
	log.Printf("%v", foodItem)
```

---

### วิธีการใช้ Order และ Limit

ใช้ร่วมกันเมื่อต้องการจัดเรียงข้อมูลตามเงื่อนไขที่ต้องการ และตัดจำนวนที่ได้มาเพื่อจำกัดขอบเขตแสดงผลเฉพาะส่วนหัว

**main.go**
```go
	foodItems := []models.FoodItem{}

	err = db.Order("food_item.name DESC").Find(&foodItems).Error
	if err != nil {
		log.Fatal("Failed to retrieve food items:", err)
		return
	}

	log.Printf("%v", foodItems)
```

![[Pasted image 20260610220239.png]]
*ภาพแสดง Log การแสดงผลลัพธ์ข้อมูลที่ดึงมาโดยมีการใช้คำสั่ง Order เพื่อจัดเรียงชื่ออาหารจากอักษรหลังสุดย้อนมาหน้าสุด (DESC)*

---

### วิธีการใช้ Distinct()

ใช้สืบค้นข้อมูลเฉพาะที่ไม่มีค่าซ้ำกันในคอลัมน์นั้นๆ

main.go
```go
	type unique struct {
		RestaurantType string
	}

	uniques := []unique{}
	err = db.Model(&models.Restaurant{}).Distinct("restaurantType").Find(&uniques).Error
	if err != nil {
		log.Fatal("Failed to retrieve unique restaurants:", err)
		return
	}
	log.Printf("%v", uniques)

```

เมื่อใช้ `Distinct` โดยตรงกับ Struct ตัวแบบ GORM อาจจะดึงคอลัมน์อื่นๆ มาด้วย แต่ถ้าต้องการได้เฉพาะคอลัมน์เดี่ยวๆ สามารถกำหนด `Model()` แล้วส่งข้อมูลเข้าไปยัง Struct ชั่วคราว (เช่น `unique`) ได้ทันที

![[Pasted image 20260610222726.png]]
*ภาพแสดง Log ผลลัพธ์การรันคำสั่ง Distinct เพื่อดึงประเภทอาหารของร้านอาหารที่ไม่มีค่าซ้ำกันขึ้นมาแสดงผล*
