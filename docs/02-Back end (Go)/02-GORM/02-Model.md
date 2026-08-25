# การจัดการ Model ใน GORM (Struct to Model)

การออกแบบโครงสร้างข้อมูล (Model) ใน GORM เป็นหัวใจสำคัญ เนื่องจากเป็นสะพานเชื่อมระหว่างโค้ด Go และตารางในฐานข้อมูล

- เอกสารอย่างเป็นทางการเรื่องโมเดลและการเขียนแท็ก: [GORM Models & Field Tags](https://gorm.io/docs/models.html#Fields-Tags)

![[Pasted image 20250306104206.png]]
*ภาพอ้างอิงแสดงโครงสร้างการออกแบบ Model และการกำหนด Field Tags ต่างๆ เพื่ออ้างอิงข้อมูลในฐานข้อมูลของ GORM*

---

## GORM Model

models/restaurant.go
```go
package models

const TableNameRestaurant = "restaurant"

// Restaurant mapped from table <restaurant>
type Restaurant struct {
	ID             int32  `gorm:"column:id;primaryKey" json:"id"`
	Name           string `gorm:"column:name;not null" json:"name"`
	RestaurantType string `gorm:"column:restaurantType;not null" json:"restaurantType"`
	Location       string `gorm:"column:location;not null" json:"location"`
}

// TableName Restaurant's table name
func (*Restaurant) TableName() string {
	return TableNameRestaurant
}

```

models/food_item.go
```go
package models

const TableNameFoodItem = "food_item"

// FoodItem mapped from table <food_item>
type FoodItem struct {
	ID           int32      `gorm:"column:id;primaryKey" json:"id"`
	Name         string     `gorm:"column:name;not null" json:"name"`
	FoodType     string     `gorm:"column:foodType;not null" json:"foodType"`
	ImagePath    string     `gorm:"column:imagePath" json:"imagePath"`
	Price        float64    `gorm:"column:price;not null" json:"price"`
	IsAvailable  bool       `gorm:"column:isAvailable;default:1" json:"isAvailable"`
	RestaurantID int32      `gorm:"column:restaurantId;not null" json:"restaurantId"`
}

// TableName FoodItem's table name
func (*FoodItem) TableName() string {
	return TableNameFoodItem
}

```

## ขั้นตอนการทดสอบดึงข้อมูล (Test Query)

เมื่อเตรียมโมเดลเสร็จแล้ว ทดลองดึงข้อมูลจากตารางมาแสดงผลด้วยฟังก์ชัน `db.Find(&model)`

ฟังก์ชัน `Find` จะดึงข้อมูลทุกแถวที่พบในตาราง แล้วนำมาใส่ไว้ใน Slice (อาร์เรย์แบบไดนามิก) ที่เตรียมไว้

**main.go**
```go
func main() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	log.Println("Database connection established successfully.")

	restaurants := []models.Restaurant{}
	err = db.Find(&restaurants).Error
	if err != nil {
		log.Fatal("Failed to retrieve restaurants:", err)
		return
	}
	log.Printf("Retrieved %d restaurants from the database.\n", len(restaurants))
}
```

> [!IMPORTANT]
> สังเกตว่าใน `db.Find(&restaurants)` ต้องส่งเป็นพอยน์เตอร์ (`&restaurants`) เข้าไป เพื่อให้ GORM สามารถนำข้อมูลที่สืบค้นได้จากฐานข้อมูลเขียนทับลงไปในตัวแปรต้นทางได้โดยตรง

เมื่อสั่งรันโปรเจกต์ ถ้าสำเร็จจะได้ผลลัพธ์ดังนี้:

![[Pasted image 20260608203416.png]]
*ภาพแสดง Log ยืนยันการรันคำสั่งสำเร็จ โดยสามารถดึงข้อมูลร้านอาหารออกมาได้ทั้งหมด 30 ร้านจากฐานข้อมูล*

---

# ความสัมพันธ์ระหว่างโมเดล

ในการพัฒนาระบบจริง ข้อมูลจากแต่ละตารางมักจะเชื่อมโยงกันอยู่แล้ว เช่น ข้อมูลเมนูอาหาร (food_item) ต้องอยู่ในร้านอาหาร (restaurant) ใดร้านหนึ่ง
บทนี้เป็นการเรียนรู้การผูกความสัมพันธ์ (Relation) ระหว่าง Model และการเขียน Query แบบต่างๆ ใน GORM

# การสร้างโมเดล FoodItem และกำหนดความสัมพันธ์

จาก ER Diagram แสดงให้เห็นว่า restaurant มีความสัมพันธ์กับ food_item แบบ One-to-Many
โดยที่ตาราง food_item มีฟิลด์ restaurantId เป็น Foreign Key (FK) เชื่อมโยงไปยัง ID ของ restaurant ที่เป็น Primary Key (PK)

![[Pasted image 20260609070049.png]]
*ภาพแสดงความสัมพันธ์ (ER Diagram) ระหว่างตารางต่างๆ ได้แก่ foodie_user, food_review, food_item และ restaurant*

สามารถประกาศความสัมพันธ์ระหว่างโมเดล `FoodItem` กับ `Restaurant` ได้ดังโค้ดโครงสร้างด้านล่างนี้:

models/food_item.go
```go
package models

const TableNameFoodItem = "food_item"

// FoodItem mapped from table <food_item>
type FoodItem struct {
	ID           int32      `gorm:"column:id;primaryKey" json:"id"`
	Name         string     `gorm:"column:name;not null" json:"name"`
	FoodType     string     `gorm:"column:foodType;not null" json:"foodType"`
	ImagePath    string     `gorm:"column:imagePath" json:"imagePath"`
	Price        float64    `gorm:"column:price;not null" json:"price"`
	IsAvailable  bool       `gorm:"column:isAvailable;default:1" json:"isAvailable"`
	RestaurantID int32      `gorm:"column:restaurantId;not null" json:"restaurantId"`
	Restaurant   Restaurant `gorm:"foreignKey:restaurantId;references:Id" json:"restaurant"`
}

// TableName FoodItem's table name
func (*FoodItem) TableName() string {
	return TableNameFoodItem
}

```

### อธิบายการกำหนดความสัมพันธ์ (Relation Mapping):
- สร้างฟิลด์ชื่อ `Restaurant` ซึ่งมีประเภทข้อมูลเป็นโมเดล `Restaurant`
- ใส่ Tag กำกับความสัมพันธ์ `gorm:"foreignKey:restaurantId;references:Id"` หมายความว่า:
  - `foreignKey:restaurantId`: ฟิลด์เชื่อมโยงในฝั่งโมเดลหลัก คือฟิลด์ `restaurantId` (ตัวแปร int ด้านบน)
  - `references:Id`: ฟิลด์ในตาราง Restaurant ที่ใช้สำหรับอ้างอิงคือคีย์หลัก `Id`

---

## การดึงข้อมูลทั่วไป (แบบไม่โหลดความสัมพันธ์)

หากใช้คำสั่งค้นหาธรรมดาโดยไม่สั่งโหลดตารางสัมพันธ์ล่วงหน้า:

main.go
```go
	foodItems := []models.FoodItem{}
	err = db.Find(&foodItems).Error
	if err != nil {
		log.Fatal("Failed to retrieve food items:", err)
		return
	}
	log.Printf("%v", foodItems)
```

ผลลัพธ์ที่ได้จะไม่มีข้อมูลร้านอาหารติดมาด้วย (ฟิลด์ `Restaurant` จะเป็นค่าว่างเริ่มต้น เช่น `{Id:0 Name:}`) เนื่องจาก GORM จะดึงข้อมูลเฉพาะของตาราง `food_item` เท่านั้นเพื่อประหยัดทรัพยากร:

![[Pasted image 20260609071522.png]]
*ภาพแสดงผลลัพธ์การรันคำสั่งโดยไม่มีการโหลดข้อมูลร้านอาหาร สังเกตในผลลัพธ์ว่าฟิลด์ Restaurant จะไม่มีข้อมูลบรรจุอยู่*

---

# การโหลดข้อมูลสัมพันธ์ด้วย Preload (Eager Loading)

เมื่อต้องการโหลดตารางย่อยที่เกี่ยวข้องเชื่อมโยงกันขึ้นมาพร้อมๆ กันด้วย GORM วิธีแรกที่นิยมใช้คือ **Preload**

- **Preload** คือการสืบค้นแยกกัน โดย GORM จะยิง SQL ตัวแรกเพื่อดึงข้อมูลตารางหลักขึ้นมาก่อน จากนั้นจะยิง SQL คิวรีเพิ่มเติมเพื่อค้นหาข้อมูลร้านอาหารของเมนูอาหารนั้นๆ แล้วนำมาประกอบกันในโค้ดให้อัตโนมัติ
- ใช้รูปแบบคำสั่ง: `Preload("ชื่อฟิลด์ความสัมพันธ์ใน Struct")` (ระบุชื่อฟิลด์ใน Go ซึ่งก็คือ `"Restaurant"` ไม่ใช่ชื่อตารางในฐานข้อมูล)

**main.go**
```go
err = db.Preload("Restaurant").Find(&foodItems).Error
```

![[Pasted image 20260609071814.png]]
*ภาพแสดง Log การดึงข้อมูลที่มีการใช้ Preload ซึ่งจะแสดงผลรายละเอียดข้อมูลร้านอาหารที่ถูกดึงเข้ามาพ่วงรวมกับข้อมูลเมนูอาหารหลักอย่างถูกต้อง*

---

# การโหลดข้อมูลสัมพันธ์ด้วย Joins

อีกวิธีหนึ่งในการโหลดข้อมูลย่อยคือการสั่ง **Joins**

- **Joins** สั่งให้ SQL ทำการ `LEFT JOIN` เชื่อมตารางตั้งแต่ฝั่ง Database ทำให้ดึงข้อมูลสัมพันธ์เสร็จสรรพได้ในคิวรีเดียว (Single Query) ต่างจาก Preload ที่แยกกันสืบค้น
- ใช้รูปแบบคำสั่ง: `Joins("Restaurant")` (ระบุชื่อฟิลด์ความสัมพันธ์ใน Struct)

**main.go**
```go
result := db.Joins("Restaurant").Find(&foodItems)
```

![[Pasted image 20260609071814.png]]
*ภาพแสดง Log การดึงข้อมูลโดยใช้ Joins ซึ่งจะรันคำสั่ง LEFT JOIN เพื่อดึงข้อมูลร้านอาหารมาใส่ไว้ในโครงสร้างโมเดลร่วมกันในการรัน SQL ครั้งเดียว*

---

# ตารางเปรียบเทียบ Preload vs Joins

เปรียบเทียบความแตกต่างด้านประสิทธิภาพและคุณสมบัติ:

### Preload (Eager Loading):
- **กลไกการทำงาน:** รันคิวรีแยกกัน (ดึงตารางหลักเสร็จ ค่อยส่งคิวรีไปตารางรองตาม ID ที่ได้)
- **จำนวนคิวรี:** อย่างน้อยที่สุดคือ 2 คิวรี
- **ความซ้ำซ้อนของข้อมูล:** แทบไม่มีข้อมูลซ้ำซ้อนส่งผ่าน Network เนื่องจากข้อมูลแยกชุดกันมา
- **การจัดเรียง/ตัวกรอง:** ทำได้ยากหากต้องการเรียงตารางหลักด้วยเงื่อนไขของตารางย่อย
- **ความสัมพันธ์ที่เหมาะสม:** ใช้ได้กับทุกรูปแบบความสัมพันธ์ (Belongs To, Has One, Has Many, Many to Many)

### Joins:
- **กลไกการทำงาน:** รัน SQL `JOIN` ในฝั่ง Database สรุปได้ข้อมูลจบในคิวรีเดียว
- **จำนวนคิวรี:** 1 คิวรี
- **ความซ้ำซ้อนของข้อมูล:** อาจเกิดความซ้ำซ้อนของข้อมูลในผลลัพธ์ดิบได้ (ถ้าความสัมพันธ์เป็น 1-to-Many) ซึ่ง GORM จะต้องกรองข้อมูลซ้ำนั้นอีกทีในระดับแอปพลิเคชัน
- **การจัดเรียง/ตัวกรอง:** ทำได้ง่ายมาก สามารถใช้ Where หรือ Order กรองข้ามตารางผ่านเงื่อนไขของตารางย่อยได้ทันที
- **ความสัมพันธ์ที่เหมาะสม:** เหมาะกับความสัมพันธ์แบบ 1-to-1 (เช่น Belongs To หรือ Has One) เท่านั้น

> [!TIP]
> - เลือกใช้ **Joins** เมื่อต้องการประสิทธิภาพสูง ต้องการลด Network round-trips หรือต้องการใช้เงื่อนไขคัดกรองจัดเรียงข้อมูลข้ามตาราง
> - เลือกใช้ **Preload** เมื่อมีระบบความสัมพันธ์แบบ 1-to-Many หรือ Many-to-Many และกังวลเรื่องขนาดชุดข้อมูลที่ Join กันแล้วจะบวมเกินไป
