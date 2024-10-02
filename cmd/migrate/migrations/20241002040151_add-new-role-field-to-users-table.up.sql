ALTER TABLE users 
ADD COLUMN role ENUM('admin', 'user') DEFAULT 'user' NOT NULL;