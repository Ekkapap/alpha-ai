# MORE REQUIREMENTS FROM USER (Check done = [x], pending = [ ])

[x] Update Sync Logic
ผมอยากให้ การทำงานของ alpha --sync หรือ /alpha-sync ทำงานแบบที่ผมบอกให้คุณทำ
คือให้ agent เป็นคนสรุป context ก่อน แล้วส่งสรุปนั้นเข้าไปให้ sync "args"
ภายใน ให้ sync ทำสั่ง 
1. สร้าง session-[date-time].md จาก args ที่รับเข้ามา
2. update หรือ merge session-summary.md จาก args ที่รับเข้ามา 
3. graphify update 
4. understand update

[x] ข้อ 1-2 คงต้องรับเข้ามาสองรอบไหม หรือส่งแบบ array มาได้เลย เพราะว่า
1. เป็นข้อมูลสรุปของใหม่เฉพาะครั้งล่าสุด
2. เป็นข้อมูลสรุปของเก่าที่ต้องวิเคราห์ก่อนว่า จะอัพเดทข้อมูลเดิมตรงไหนไหม หรือเพิ่มต่อท้าย หรือทำทั้งสองอย่าง หรือตัดส่วนที่ไม่จำเป็นทิ้งไป

