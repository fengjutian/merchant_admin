CREATE DATABASE IF NOT EXISTS merchant_admin 
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_general_ci;

USE merchant_admin;

CREATE TABLE IF NOT EXISTS business (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '业务名称',
    email VARCHAR(255) NOT NULL COMMENT '联系方式或标识',
    address VARCHAR(255) NOT NULL COMMENT '地址',
    type VARCHAR(100) NOT NULL COMMENT '类型',
    contact VARCHAR(255) NOT NULL COMMENT '通讯方式',
    rating FLOAT DEFAULT 0 COMMENT '评价 0~5',
    latitude DOUBLE NULL COMMENT '纬度',
    longitude DOUBLE NULL COMMENT '经度',
    otherInfo TEXT NULL COMMENT '其他信息',
    imageBase64 LONGTEXT NULL COMMENT '图片 Base64',
    description TEXT NULL COMMENT '描述/简介',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO business 
(name, email, address, type, contact, rating, latitude, longitude, otherInfo, imageBase64, description)
VALUES
(
    'Star Coffee',
    'contact@starcoffee.com',
    'Tokyo, Shibuya 1-2-3',
    'Cafe',
    '03-1234-5678',
    4.5,
    35.6595,
    139.7005,
    'Open 24 hours',
    NULL,
    'A cozy café located in Shibuya.'
),
(
    'Tech Repair Shop',
    'support@techrepair.jp',
    'Osaka, Namba 4-5-6',
    'Repair',
    '06-9876-5432',
    4.0,
    34.6666,
    135.5000,
    'Specialized in phone & laptop repairs',
    NULL,
    'Fast and affordable technical repair service.'
),
(
    'Fresh Market',
    'info@freshmarket.jp',
    'Nagoya, Sakae 7-8-9',
    'Supermarket',
    '052-9988-7766',
    3.8,
    35.1709,
    136.8815,
    'Organic food available',
    NULL,
    'Local supermarket offering fresh products daily.'
);
