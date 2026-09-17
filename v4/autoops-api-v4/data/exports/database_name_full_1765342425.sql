-- Database: database_name
-- Export Time: 2025-12-10 12:53:45
-- Tables: 2
-- Export Type: full

-- Table structure for students
CREATE TABLE "students" (
 "id" serial(32) NOT NULL,
 "student_id" varchar(20) NOT NULL,
 "name" varchar(50) NOT NULL,
 "gender" bpchar(1),
 "birth_date" date,
 "email" varchar(100),
 "phone" varchar(20),
 "enrollment_date" date DEFAULT 'CURRENT_TIMESTAMP',
 "major" varchar(100),
 "created_at" timestamp DEFAULT CURRENT_TIMESTAMP, 
PRIMARY KEY ("id","student_id","email")
);
DROP INDEX IF EXISTS students_student_id_key;
DROP INDEX IF EXISTS students_email_key;
CREATE unique INDEX "students_student_id_key" ON "students"("student_id");
CREATE unique INDEX "students_email_key" ON "students"("email");

-- Data for table students
-- No data

-- Table structure for users
CREATE TABLE "users" (
 "id" serial(32) NOT NULL,
 "name" varchar(100) NOT NULL,
 "email" varchar(150) NOT NULL, 
PRIMARY KEY ("id","email")
);
DROP INDEX IF EXISTS users_email_key;
CREATE unique INDEX "users_email_key" ON "users"("email");

-- Data for table users
INSERT INTO "users" ("id","name","email") VALUES (1,'张三','zhangsan@123.com');
INSERT INTO "users" ("id","name","email") VALUES (2,'李四','lisi@123.com');
INSERT INTO "users" ("id","name","email") VALUES (5,'王五','xxxx@123.com');

