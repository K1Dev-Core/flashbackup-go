![[Pasted image 20260606213003.png]]

### คำอธิบายแผนภาพสถาปัตยกรรม Wails + Nuxt 4 (Architecture Diagram)

แผนภาพนี้แสดงโครงสร้างสถาปัตยกรรมของแอปพลิเคชันเดสก์ท็อปรองรับหลายระบบปฏิบัติการ (Cross-Platform Desktop Application) โดยใช้ Wails ทำหน้าที่เป็นหลังบ้าน ร่วมกับ Nuxt 4 ทำหน้าที่เป็นหน้าบ้าน โดยแบ่งออกเป็นส่วนหลักๆ ดังนี้:

1. **Wails (Go Backend - ฝั่งซ้าย):**
   - **Go Runtime (Business Logic):** ส่วนการทำงานหลักของโปรแกรม เขียนด้วยภาษา Go จัดการลอจิกต่างๆ ของแอปพลิเคชัน
   - **API Handlers (Events/Bindings):** ตัวรับส่งฟังก์ชันและการเชื่อมต่อข้อมูลระหว่างหน้าบ้านและหลังบ้าน
   - **OS Integration (Files, Native APIs, Shell):** ส่วนเข้าถึงระบบปฏิบัติการโดยตรง เช่น การจัดการไฟล์, การเรียกใช้งาน API ของระบบปฏิบัติการ และคำสั่ง Shell
   - **Application Data (Database/Storage):** ส่วนติดต่อฐานข้อมูล SQLite และการเก็บข้อมูลถาวรภายในเครื่อง

2. **Wails Binding & IPC (สะพานเชื่อมต่อตรงกลาง):**
   - **Bidirectional Communication (Go <-> Nuxt):** ช่องทางสื่อสารแบบสองทิศทางข้ามกระบวนการทำงาน
   - **Go Functions Exposed as JS Promises:** การนำฟังก์ชันหลังบ้านใน Go มาแปลงเป็น Promise ใน JavaScript ให้ฝั่งหน้าบ้านเรียกใช้งานได้ทันที
   - **Native Window:** หน้าต่างแสดงผลดั้งเดิมของแต่ละระบบปฏิบัติการ (เช่น WebKit ใน macOS, WebView2 ใน Windows, GTK WebKit ใน Linux) โดยไม่ต้องพึ่งเบราว์เซอร์แยกภายนอก

3. **Nuxt 4 (Vue/Vite Frontend - ฝั่งขวา):**
   - **Nuxt 4 (Client-Side Rendering):** จัดการวาดหน้าตาแอปและนำเสนอผลลัพธ์ที่ฝั่งผู้ใช้งาน
   - **Vue 3 Components (UI/UX):** ส่วนประกอบหน้าจอและหน้าตาโปรแกรมเพื่อการโต้ตอบที่ลื่นไหล
   - **Nuxt Pages/Routing:** ระบบนำเส้นทางข้ามหน้าเว็บและการจัดการเพจ
   - **State Management (Pinia/Stores):** ส่วนควบคุมข้อมูลและสถานะของหน้าบ้านให้สอดคล้องกันทั้งแอป
   - **Assets (CSS/Images):** ไฟล์ตกแต่งและสไตล์ชีทต่างๆ เช่น CSS, ไฟล์รูปภาพ

4. **OS Desktop Environment & Distribution (ด้านล่าง):**
   - แอปสามารถรันได้บนทุกระบบปฏิบัติการหลัก ทั้ง Windows, macOS และ Linux
   - ขั้นตอนการคอมไพล์สุดท้าย (Build Artifacts) จะได้ออกมาเป็นไฟล์รัน Wails Binary และไฟล์เว็บ Nuxt Static Files รวมกันเป็นโปรแกรมเดียว
   - สามารถแปลงผลลัพธ์เป็นตัวติดตั้ง (Installers) หรือแพ็กเกจแอปพลิเคชันสำเร็จรูป (App Packages) สำหรับแจกจ่ายใช้งานจริงได้ทันที