# Practice: Bank Account Struct (แบบฝึกหัด: โครงสร้างบัญชีธนาคาร)

แบบฝึกหัดนี้ใช้ทดลองดึงความรู้เรื่อง **Struct**, **Pointer Receiver**, และ **Encapsulation (Package Scope)** มารวมกันเพื่อสร้างระบบบัญชีธนาคาร (Bank Account) แบบจำลองที่มีความปลอดภัย

---

## 🎯 โจทย์การทดลอง
1. สร้างโฟลเดอร์แพ็กเกจใหม่ชื่อ `bank`
2. สร้าง Struct ชื่อ `Account` ไว้เก็บข้อมูลบัญชีธนาคาร ประกอบด้วยฟิลด์ลับ (ตัวพิมพ์เล็ก):
   - `id` (ไอดีบัญชี) -> `string`
   - `name` (ชื่อเจ้าของบัญชี) -> `string`
   - `balance` (ยอดเงินคงเหลือ) -> `float64`
3. เขียนเมธอด (Receiver Functions) ให้กับ `Account` ดังนี้:
   - **`Deposit(amount float64)`**: ใช้สำหรับฝากเงิน (ต้องห้ามฝากยอดติดลบ หรือ 0)
   - **`Withdraw(amount float64) error`**: ใช้สำหรับถอนเงิน (ต้องเช็กว่ายอดเงินพอไหม และห้ามถอนยอดติดลบ ถ้าเงินไม่พอให้ส่ง error กลับไป)
   - **`GetBalance() float64`**: ใช้สำหรับอ่านค่าเงินคงเหลือปัจจุบัน (เนื่องจากฟิลด์ `balance` เป็นตัวเล็กภายนอกจะเข้าถึงตรงๆ ไม่ได้)
4. เขียน **Constructor Function** สำหรับสร้างออบเจกต์บัญชีใหม่ขึ้นมา
5. เขียนไฟล์ `main.go` นำเข้าแพ็กเกจ `bank` มาจำลองขั้นตอนการสร้างบัญชี ฝากเงิน ถอนเงิน และตรวจเช็กความปลอดภัยของระบบ

---

## 📂 โครงสร้างโฟลเดอร์โปรเจกต์
```text
go-basic/
├── go.mod
├── main.go
└── bank/
    └── account.go
```

---

## ขั้นตอนที่ 1: เขียนโมเดลบัญชีธนาคาร (`bank/account.go`)

สร้างไฟล์ `account.go` ภายใต้โฟลเดอร์ `bank` แล้วเขียนโค้ดเพื่อควบคุมธุรกรรมการเงินตามนี้:

**bank/account.go**
```go
package bank

import (
	"errors"
	"fmt"
)

// Account Struct
// ตัวแปรข้างในขึ้นต้นด้วยตัวพิมพ์เล็กทั้งหมดเพื่อเป็น Private ป้องกันการแก้ไขยอดเงินจากภายนอกตรงๆ!
type Account struct {
	id      string
	name    string
	balance float64
}

// 1. NewAccount: ฟังก์ชันสร้างบัญชีใหม่ (Constructor)
// คืนค่ากลับไปเป็น pointer (*Account) เพื่อให้ภายนอกนำไปใช้งานต่อได้ประหยัดหน่วยความจำ
func NewAccount(id string, name string, initialBalance float64) (*Account, error) {
	if initialBalance < 0 {
		return nil, errors.New("ยอดเงินเริ่มต้นต้องไม่ติดลบ")
	}
	return &Account{
		id:      id,
		name:    name,
		balance: initialBalance,
	}, nil
}

// 2. Deposit: เมธอดฝากเงิน (Pointer Receiver)
func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("จำนวนเงินที่ฝากต้องมากกว่า 0")
	}
	a.balance += amount
	fmt.Printf("ฝากเงินสำเร็จ: +%.2f บาท\n", amount)
	return nil
}

// 3. Withdraw: เมธอดถอนเงินพร้อมเช็กยอดคงเหลือ
func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("จำนวนเงินที่ถอนต้องมากกว่า 0")
	}
	if a.balance < amount {
		return errors.New("ยอดเงินในบัญชีไม่เพียงพอสำหรับการถอน")
	}
	a.balance -= amount
	fmt.Printf("ถอนเงินสำเร็จ: -%.2f บาท\n", amount)
	return nil
}

// 4. GetBalance: Getter ดึงยอดเงินปัจจุบัน
func (a *Account) GetBalance() float64 {
	return a.balance
}

// Getter ดึงชื่อเจ้าของบัญชี
func (a *Account) GetName() string {
	return a.name
}
```

---

## 💻 ขั้นตอนที่ 2: จำลองธุรกรรมการเงินในไฟล์หลัก (`main.go`)

สร้างไฟล์ `main.go` ในโฟลเดอร์หลัก เพื่อทดสอบการทำงานและเช็กเงื่อนไข Error ต่างๆ:

**main.go**
```go
package main

import (
	"fmt"
	"go-basic/bank"
)

func main() {
	// 1. สร้างบัญชีธนาคารใหม่ ยอดเงินเปิดบัญชี 500 บาท
	myAcc, err := bank.NewAccount("ACC-001", "Potchara", 500.00)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("บัญชี: %s เปิดใช้งานสำเร็จ ยอดเงินเปิดบัญชี: %.2f บาท\n", myAcc.GetName(), myAcc.GetBalance())
	fmt.Println("-------------------------------------------")

	// 2. ทดลองฝากเงิน
	err = myAcc.Deposit(250.50)
	if err != nil {
		fmt.Println("ฝากเงินล้มเหลว:", err)
	}
	fmt.Printf("ยอดเงินคงเหลือปัจจุบัน: %.2f บาท\n", myAcc.GetBalance())
	fmt.Println("-------------------------------------------")

	// 3. ทดลองถอนเงินปกติ
	err = myAcc.Withdraw(150.00)
	if err != nil {
		fmt.Println("ถอนเงินล้มเหลว:", err)
	}
	fmt.Printf("ยอดเงินคงเหลือปัจจุบัน: %.2f บาท\n", myAcc.GetBalance())
	fmt.Println("-------------------------------------------")

	// 4. ทดลองถอนเงินเกินบัญชี (Overdraft) -> ระบบควรแจ้งพังผ่าน error
	fmt.Println("กำลังทดลองถอนเงิน 1,000 บาท...")
	err = myAcc.Withdraw(1000.00)
	if err != nil {
		fmt.Printf("แจ้งเตือนจากธนาคาร: %v\n", err)
	} else {
		fmt.Println("ถอนเงินสำเร็จ!")
	}
	fmt.Printf("ยอดเงินคงเหลือปัจจุบัน: %.2f บาท\n", myAcc.GetBalance())
	fmt.Println("-------------------------------------------")
}
```

---

## 💡 สิ่งที่ได้เรียนรู้เพิ่มจากแบบฝึกหัดนี้:
1. **การสร้าง Constructor (`NewAccount`):** ภาษา Go ไม่มีฟังก์ชันสร้างออบเจกต์สำเร็จรูป จึงมีธรรมเนียมตั้งชื่อฟังก์ชันว่า `New[ชื่อStruct]` เพื่อทำหน้าที่ตั้งค่าเริ่มต้นให้ Struct
2. **การกักบริเวณฟิลด์ด้วยตัวเล็ก:** การใช้ `id`, `name`, `balance` เป็นตัวพิมพ์เล็ก ทำให้ฟิลด์เหล่านี้มีความปลอดภัยสูงมาก คนภายนอกจะไม่สามารถแอบเขียนแก้ไขยอดเงินตรงๆ เช่น `myAcc.balance = 99999` ได้ ต้องทำผ่านเมธอดตรวจสอบความปลอดภัยที่สร้างไว้เท่านั้น
3. **Pointer Receiver สำคัญมาก:** ถ้าเมธอด `Deposit` หรือ `Withdraw` ลืมเครื่องหมาย `*` ใน receiver ยอดเงินในบัญชีในหน้า `main.go` จะไม่มีวันเปลี่ยนแปลงเลย เพราะมันจะไปทำการแก้ไขเงินอยู่ในก้อน Copy ชั่วคราวแทน