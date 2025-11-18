CREATE DATABASE IF NOT EXISTS merchant_admin 
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_general_ci;

USE merchant_admin;

CREATE TABLE IF NOT EXISTS business (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    name VARCHAR(255) NOT NULL COMMENT '商家名称',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱（唯一）',
    address VARCHAR(255) NOT NULL COMMENT '地址',
    type VARCHAR(100) NOT NULL COMMENT '商家类型',
    contact VARCHAR(255) NOT NULL COMMENT '联系方式',
    rating FLOAT DEFAULT 0 COMMENT '评分（0-5）',
    latitude DOUBLE NULL COMMENT '纬度（可空）',
    longitude DOUBLE NULL COMMENT '经度（可空）',
    otherInfo TEXT NULL COMMENT '其他信息（可空）',
    imageBase64 LONGTEXT NULL COMMENT 'base64图片（可空）',
    description TEXT NULL COMMENT '描述（可空）',
    status VARCHAR(50) DEFAULT 'active' COMMENT '状态（active, inactive, suspended）',
    phone VARCHAR(20) NULL COMMENT '电话号码',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_business_type (type),
    INDEX idx_business_status (status),
    INDEX idx_business_rating (rating),
    INDEX idx_business_email (email),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='商家信息表';

-- 插入示例数据
INSERT INTO business 
(name, email, address, type, contact, rating, latitude, longitude, otherInfo, imageBase64, description, status, phone)
VALUES
(
    'Star Coffee',
    'contact@starcoffee.com',
    'Tokyo, Shibuya 1-2-3',
    'restaurant',
    '03-1234-5678',
    4.5,
    35.6595,
    139.7005,
    'Open 24 hours',
    NULL,
    'A cozy café located in Shibuya.',
    'active',
    '03-1234-5678'
),
(
    'Tech Repair Shop',
    'support@techrepair.jp',
    'Osaka, Namba 4-5-6',
    'service',
    '06-9876-5432',
    4.0,
    34.6666,
    135.5000,
    'Specialized in phone & laptop repairs',
    NULL,
    'Fast and affordable technical repair service.',
    'active',
    '06-9876-5432'
),
(
    'Fresh Market',
    'info@freshmarket.jp',
    'Nagoya, Sakae 7-8-9',
    'retail',
    '052-9988-7766',
    3.8,
    35.1709,
    136.8815,
    'Organic food available',
    NULL,
    'Local supermarket offering fresh products daily.',
    'active',
    '052-9988-7766'
),
(
    'Game Center',
    'fun@gamecenter.jp',
    'Fukuoka, Tenjin 2-3-4',
    'entertainment',
    '092-1122-3344',
    4.2,
    33.5904,
    130.4017,
    'Latest arcade games and VR experiences',
    NULL,
    'Modern entertainment facility with various gaming options.',
    'active',
    '092-1122-3344'
),
(
    'Book Store Plus',
    'books@storeplus.jp',
    'Kyoto, Kawaramachi 5-6-7',
    'retail',
    '075-5544-3322',
    4.7,
    35.0116,
    135.7681,
    'Wide selection of books and magazines',
    NULL,
    'Traditional bookstore with modern reading space.',
    'active',
    '075-5544-3322'
);
