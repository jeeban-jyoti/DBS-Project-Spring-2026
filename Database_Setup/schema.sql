-- MySQL dump 10.13  Distrib 9.5.0, for macos26.2 (arm64)
--
-- Host: localhost    Database: gyanpustak
-- ------------------------------------------------------
-- Server version	9.5.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '41e9ff14-f065-11f0-bad4-3351000b39fb:1-547';

--
-- Table structure for table `administrator`
--

DROP TABLE IF EXISTS `administrator`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `administrator` (
  `employee_id` int NOT NULL,
  PRIMARY KEY (`employee_id`),
  CONSTRAINT `administrator_ibfk_1` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`employee_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `administrator`
--

LOCK TABLES `administrator` WRITE;
/*!40000 ALTER TABLE `administrator` DISABLE KEYS */;
INSERT INTO `administrator` VALUES (4),(5);
/*!40000 ALTER TABLE `administrator` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `author`
--

DROP TABLE IF EXISTS `author`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `author` (
  `author_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`author_id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `author`
--

LOCK TABLES `author` WRITE;
/*!40000 ALTER TABLE `author` DISABLE KEYS */;
INSERT INTO `author` VALUES (9,'Author A'),(10,'Author B'),(3,'CLRS'),(1,'Korth'),(4,'Morris Mano'),(2,'Silberschatz');
/*!40000 ALTER TABLE `author` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `book`
--

DROP TABLE IF EXISTS `book`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `book` (
  `book_id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `isbn` varchar(20) DEFAULT NULL,
  `publisher` varchar(255) DEFAULT NULL,
  `publication_date` date DEFAULT NULL,
  `edition` varchar(50) DEFAULT NULL,
  `language` varchar(50) DEFAULT NULL,
  `format` enum('hardcover','softcover','electronic') DEFAULT NULL,
  `type` enum('new','used') DEFAULT NULL,
  `purchase_option` enum('rent','buy') DEFAULT NULL,
  `price` decimal(10,2) DEFAULT NULL,
  `quantity` int NOT NULL DEFAULT '0',
  `category_id` int DEFAULT NULL,
  PRIMARY KEY (`book_id`),
  UNIQUE KEY `isbn` (`isbn`),
  KEY `category_id` (`category_id`),
  CONSTRAINT `book_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `category` (`category_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `book`
--

LOCK TABLES `book` WRITE;
/*!40000 ALTER TABLE `book` DISABLE KEYS */;
INSERT INTO `book` VALUES (1,'Database System Concepts','ISBN001','Pearson','2020-01-01','7th','English','hardcover','new','buy',599.99,9,1),(2,'Introduction to Algorithms','ISBN002','MIT Press','2018-01-01','3rd','English','softcover','new','buy',799.99,6,1),(3,'Digital Design','ISBN003','Pearson','2016-01-01','5th','English','hardcover','new','rent',499.99,16,2),(5,'Updated Book Title','123-456-789','New Publisher','2023-01-01','2nd','English','hardcover','new','buy',59.99,10,1);
/*!40000 ALTER TABLE `book` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `book_author`
--

DROP TABLE IF EXISTS `book_author`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `book_author` (
  `book_id` int NOT NULL,
  `author_id` int NOT NULL,
  PRIMARY KEY (`book_id`,`author_id`),
  KEY `author_id` (`author_id`),
  CONSTRAINT `book_author_ibfk_1` FOREIGN KEY (`book_id`) REFERENCES `book` (`book_id`) ON DELETE CASCADE,
  CONSTRAINT `book_author_ibfk_2` FOREIGN KEY (`author_id`) REFERENCES `author` (`author_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `book_author`
--

LOCK TABLES `book_author` WRITE;
/*!40000 ALTER TABLE `book_author` DISABLE KEYS */;
INSERT INTO `book_author` VALUES (1,1),(1,2),(2,3),(3,4),(5,9),(5,10);
/*!40000 ALTER TABLE `book_author` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `book_keyword`
--

DROP TABLE IF EXISTS `book_keyword`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `book_keyword` (
  `book_id` int NOT NULL,
  `keyword_id` int NOT NULL,
  PRIMARY KEY (`book_id`,`keyword_id`),
  KEY `keyword_id` (`keyword_id`),
  CONSTRAINT `book_keyword_ibfk_1` FOREIGN KEY (`book_id`) REFERENCES `book` (`book_id`) ON DELETE CASCADE,
  CONSTRAINT `book_keyword_ibfk_2` FOREIGN KEY (`keyword_id`) REFERENCES `keyword` (`keyword_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `book_keyword`
--

LOCK TABLES `book_keyword` WRITE;
/*!40000 ALTER TABLE `book_keyword` DISABLE KEYS */;
INSERT INTO `book_keyword` VALUES (5,9),(5,10),(5,11);
/*!40000 ALTER TABLE `book_keyword` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `book_subcategory`
--

DROP TABLE IF EXISTS `book_subcategory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `book_subcategory` (
  `book_id` int NOT NULL,
  `subcategory_id` int NOT NULL,
  PRIMARY KEY (`book_id`,`subcategory_id`),
  KEY `subcategory_id` (`subcategory_id`),
  CONSTRAINT `book_subcategory_ibfk_1` FOREIGN KEY (`book_id`) REFERENCES `book` (`book_id`) ON DELETE CASCADE,
  CONSTRAINT `book_subcategory_ibfk_2` FOREIGN KEY (`subcategory_id`) REFERENCES `subcategory` (`subcategory_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `book_subcategory`
--

LOCK TABLES `book_subcategory` WRITE;
/*!40000 ALTER TABLE `book_subcategory` DISABLE KEYS */;
INSERT INTO `book_subcategory` VALUES (1,1),(2,2),(5,2),(3,3),(5,10);
/*!40000 ALTER TABLE `book_subcategory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cart`
--

DROP TABLE IF EXISTS `cart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cart` (
  `cart_id` int NOT NULL AUTO_INCREMENT,
  `student_id` int DEFAULT NULL,
  `created_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime DEFAULT NULL,
  PRIMARY KEY (`cart_id`),
  UNIQUE KEY `student_id` (`student_id`),
  CONSTRAINT `cart_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `student` (`student_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cart`
--

LOCK TABLES `cart` WRITE;
/*!40000 ALTER TABLE `cart` DISABLE KEYS */;
INSERT INTO `cart` VALUES (1,1,'2026-04-09 23:48:15',NULL);
/*!40000 ALTER TABLE `cart` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cart_item`
--

DROP TABLE IF EXISTS `cart_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cart_item` (
  `cart_id` int NOT NULL,
  `book_id` int NOT NULL,
  `quantity` int NOT NULL,
  PRIMARY KEY (`cart_id`,`book_id`),
  KEY `book_id` (`book_id`),
  CONSTRAINT `cart_item_ibfk_1` FOREIGN KEY (`cart_id`) REFERENCES `cart` (`cart_id`) ON DELETE CASCADE,
  CONSTRAINT `cart_item_ibfk_2` FOREIGN KEY (`book_id`) REFERENCES `book` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cart_item`
--

LOCK TABLES `cart_item` WRITE;
/*!40000 ALTER TABLE `cart_item` DISABLE KEYS */;
/*!40000 ALTER TABLE `cart_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `category`
--

DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `category` (
  `category_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`category_id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `category`
--

LOCK TABLES `category` WRITE;
/*!40000 ALTER TABLE `category` DISABLE KEYS */;
INSERT INTO `category` VALUES (1,'Computer Science'),(2,'Electronics');
/*!40000 ALTER TABLE `category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course`
--

DROP TABLE IF EXISTS `course`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course` (
  `course_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `university_id` int DEFAULT NULL,
  `year` int DEFAULT NULL,
  PRIMARY KEY (`course_id`),
  KEY `course_ibfk_1` (`university_id`),
  CONSTRAINT `course_ibfk_1` FOREIGN KEY (`university_id`) REFERENCES `university` (`university_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course`
--

LOCK TABLES `course` WRITE;
/*!40000 ALTER TABLE `course` DISABLE KEYS */;
/*!40000 ALTER TABLE `course` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_department`
--

DROP TABLE IF EXISTS `course_department`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_department` (
  `course_id` int NOT NULL,
  `department_id` int NOT NULL,
  PRIMARY KEY (`course_id`,`department_id`),
  UNIQUE KEY `course_id` (`course_id`,`department_id`),
  KEY `course_department_ibfk_2` (`department_id`),
  CONSTRAINT `course_department_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `course` (`course_id`) ON DELETE CASCADE,
  CONSTRAINT `course_department_ibfk_2` FOREIGN KEY (`department_id`) REFERENCES `department` (`department_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_cd_course` FOREIGN KEY (`course_id`) REFERENCES `course` (`course_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_department`
--

LOCK TABLES `course_department` WRITE;
/*!40000 ALTER TABLE `course_department` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_department` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `customer_support`
--

DROP TABLE IF EXISTS `customer_support`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `customer_support` (
  `employee_id` int NOT NULL,
  PRIMARY KEY (`employee_id`),
  CONSTRAINT `customer_support_ibfk_1` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`employee_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customer_support`
--

LOCK TABLES `customer_support` WRITE;
/*!40000 ALTER TABLE `customer_support` DISABLE KEYS */;
INSERT INTO `customer_support` VALUES (3);
/*!40000 ALTER TABLE `customer_support` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `department`
--

DROP TABLE IF EXISTS `department`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `department` (
  `department_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `university_id` int DEFAULT NULL,
  PRIMARY KEY (`department_id`),
  UNIQUE KEY `name` (`name`,`university_id`),
  KEY `department_ibfk_1` (`university_id`),
  CONSTRAINT `department_ibfk_1` FOREIGN KEY (`university_id`) REFERENCES `university` (`university_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `department`
--

LOCK TABLES `department` WRITE;
/*!40000 ALTER TABLE `department` DISABLE KEYS */;
INSERT INTO `department` VALUES (5,'ce',1),(3,'cse',1),(4,'me',1);
/*!40000 ALTER TABLE `department` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `employee`
--

DROP TABLE IF EXISTS `employee`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `employee` (
  `employee_id` int NOT NULL,
  `gender` enum('male','female','other') DEFAULT NULL,
  `salary` decimal(10,2) DEFAULT NULL,
  `aadhaar_number` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`employee_id`),
  UNIQUE KEY `aadhaar_number` (`aadhaar_number`),
  CONSTRAINT `employee_ibfk_1` FOREIGN KEY (`employee_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `employee`
--

LOCK TABLES `employee` WRITE;
/*!40000 ALTER TABLE `employee` DISABLE KEYS */;
INSERT INTO `employee` VALUES (3,'male',40000.00,'AADHAAR123'),(4,'female',80000.00,'AADHAAR456'),(5,'male',120000.00,'AADHAAR789');
/*!40000 ALTER TABLE `employee` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `instructor`
--

DROP TABLE IF EXISTS `instructor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `instructor` (
  `instructor_id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) DEFAULT NULL,
  `university_id` int DEFAULT NULL,
  `department_id` int DEFAULT NULL,
  PRIMARY KEY (`instructor_id`),
  KEY `fk_dept` (`department_id`),
  KEY `instructor_ibfk_1` (`university_id`),
  CONSTRAINT `fk_dept` FOREIGN KEY (`department_id`) REFERENCES `department` (`department_id`) ON DELETE CASCADE,
  CONSTRAINT `instructor_ibfk_1` FOREIGN KEY (`university_id`) REFERENCES `university` (`university_id`) ON DELETE CASCADE,
  CONSTRAINT `instructor_ibfk_2` FOREIGN KEY (`department_id`) REFERENCES `department` (`department_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `instructor`
--

LOCK TABLES `instructor` WRITE;
/*!40000 ALTER TABLE `instructor` DISABLE KEYS */;
/*!40000 ALTER TABLE `instructor` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `keyword`
--

DROP TABLE IF EXISTS `keyword`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `keyword` (
  `keyword_id` int NOT NULL AUTO_INCREMENT,
  `keyword` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`keyword_id`),
  UNIQUE KEY `keyword` (`keyword`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `keyword`
--

LOCK TABLES `keyword` WRITE;
/*!40000 ALTER TABLE `keyword` DISABLE KEYS */;
INSERT INTO `keyword` VALUES (9,'Coding'),(11,'Go'),(10,'Programming');
/*!40000 ALTER TABLE `keyword` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order`
--

DROP TABLE IF EXISTS `order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order` (
  `order_id` int NOT NULL AUTO_INCREMENT,
  `student_id` int DEFAULT NULL,
  `created_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `fulfilled_date` datetime DEFAULT NULL,
  `shipping_type` enum('standard','2-day','1-day') DEFAULT NULL,
  `card_number` varchar(20) DEFAULT NULL,
  `card_expiry` varchar(10) DEFAULT NULL,
  `card_holder_name` varchar(100) DEFAULT NULL,
  `card_type` varchar(50) DEFAULT NULL,
  `status` enum('new','processed','awaiting_shipping','shipped','canceled','returned') DEFAULT NULL,
  PRIMARY KEY (`order_id`),
  KEY `student_id` (`student_id`),
  CONSTRAINT `order_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `student` (`student_id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order`
--

LOCK TABLES `order` WRITE;
/*!40000 ALTER TABLE `order` DISABLE KEYS */;
INSERT INTO `order` VALUES (4,1,'2026-04-10 17:55:41',NULL,NULL,NULL,NULL,NULL,NULL,'shipped'),(6,1,'2026-04-10 19:23:05','2026-04-13 01:46:40',NULL,NULL,NULL,NULL,NULL,'processed'),(7,1,'2026-04-12 19:25:36',NULL,NULL,NULL,NULL,NULL,NULL,'canceled'),(8,1,'2026-04-12 20:28:11',NULL,NULL,NULL,NULL,NULL,NULL,'canceled'),(9,1,'2026-04-12 22:50:54',NULL,NULL,NULL,NULL,NULL,NULL,'canceled'),(10,1,'2026-04-13 02:08:08',NULL,NULL,NULL,NULL,NULL,NULL,'canceled'),(11,1,'2026-04-13 23:48:22',NULL,NULL,NULL,NULL,NULL,NULL,'canceled'),(12,1,'2026-04-15 03:05:58',NULL,NULL,NULL,NULL,NULL,NULL,'shipped'),(13,1,'2026-04-16 21:17:22',NULL,NULL,NULL,NULL,NULL,NULL,'awaiting_shipping'),(14,1,'2026-04-16 21:17:44',NULL,NULL,NULL,NULL,NULL,NULL,'canceled');
/*!40000 ALTER TABLE `order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_item`
--

DROP TABLE IF EXISTS `order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_item` (
  `order_id` int NOT NULL,
  `book_id` int NOT NULL,
  `quantity` int NOT NULL,
  `purchase_type` enum('rent','buy') DEFAULT NULL,
  PRIMARY KEY (`order_id`,`book_id`),
  KEY `book_id` (`book_id`),
  CONSTRAINT `order_item_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `order` (`order_id`) ON DELETE CASCADE,
  CONSTRAINT `order_item_ibfk_2` FOREIGN KEY (`book_id`) REFERENCES `book` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_item`
--

LOCK TABLES `order_item` WRITE;
/*!40000 ALTER TABLE `order_item` DISABLE KEYS */;
INSERT INTO `order_item` VALUES (4,1,3,'buy'),(6,2,1,'rent'),(7,1,1,'buy'),(8,1,1,'buy'),(8,2,1,'buy'),(8,3,2,'buy'),(9,2,1,'buy'),(10,2,1,'buy'),(10,3,1,'buy'),(11,1,1,'rent'),(11,5,1,'rent'),(12,1,1,'buy'),(12,3,2,'buy'),(13,3,1,'rent'),(14,1,1,'buy'),(14,2,1,'buy'),(14,3,2,'rent'),(14,5,1,'buy');
/*!40000 ALTER TABLE `order_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `request_logs`
--

DROP TABLE IF EXISTS `request_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `request_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `url` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `user_email` varchar(255) NOT NULL,
  `user_role` varchar(50) NOT NULL,
  `body` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=167 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `request_logs`
--

LOCK TABLES `request_logs` WRITE;
/*!40000 ALTER TABLE `request_logs` DISABLE KEYS */;
INSERT INTO `request_logs` VALUES (73,'/api/v1/fetchBooks','2026-04-16 15:47:18','rahul@student.com','student',''),(74,'/api/v1/fetchBooks','2026-04-16 15:47:18','rahul@student.com','student',''),(75,'/api/v1/cart','2026-04-16 15:47:21','rahul@student.com','student',''),(76,'/api/v1/cart','2026-04-16 15:47:21','rahul@student.com','student',''),(77,'/api/v1/placeOrder','2026-04-16 15:47:22','rahul@student.com','student','{}'),(78,'/api/v1/fetchBooks','2026-04-16 15:47:25','rahul@student.com','student',''),(79,'/api/v1/fetchBooks','2026-04-16 15:47:25','rahul@student.com','student',''),(80,'/api/v1/addToCart','2026-04-16 15:47:27','rahul@student.com','student','{\"bookid\":1,\"quantity\":1}'),(81,'/api/v1/addToCart','2026-04-16 15:47:29','rahul@student.com','student','{\"bookid\":2,\"quantity\":1}'),(82,'/api/v1/addToCart','2026-04-16 15:47:30','rahul@student.com','student','{\"bookid\":3,\"quantity\":1}'),(83,'/api/v1/addToCart','2026-04-16 15:47:32','rahul@student.com','student','{\"bookid\":5,\"quantity\":1}'),(84,'/api/v1/addToCart','2026-04-16 15:47:34','rahul@student.com','student','{\"bookid\":6,\"quantity\":1}'),(85,'/api/v1/cart','2026-04-16 15:47:36','rahul@student.com','student',''),(86,'/api/v1/cart','2026-04-16 15:47:36','rahul@student.com','student',''),(87,'/api/v1/fetchBooks','2026-04-16 15:47:39','rahul@student.com','student',''),(88,'/api/v1/fetchBooks','2026-04-16 15:47:39','rahul@student.com','student',''),(89,'/api/v1/addToCart','2026-04-16 15:47:40','rahul@student.com','student','{\"bookid\":3,\"quantity\":1}'),(90,'/api/v1/cart','2026-04-16 15:47:42','rahul@student.com','student',''),(91,'/api/v1/cart','2026-04-16 15:47:42','rahul@student.com','student',''),(92,'/api/v1/placeOrder','2026-04-16 15:47:44','rahul@student.com','student','{}'),(93,'/api/v1/showMyOrders','2026-04-16 15:48:02','rahul@student.com','student',''),(94,'/api/v1/showMyOrders','2026-04-16 15:48:02','rahul@student.com','student',''),(95,'/api/v1/showAllOrders','2026-04-16 15:49:12','ravi@support.com','support',''),(96,'/api/v1/showAllOrders','2026-04-16 15:49:12','ravi@support.com','support',''),(97,'/api/v1/changeOrderStatus?id=14&status=processed','2026-04-16 15:50:28','ravi@support.com','support',''),(98,'/api/v1/showAllOrders','2026-04-16 15:50:28','ravi@support.com','support',''),(99,'/api/v1/changeOrderStatus?id=13&status=awaiting_shipping','2026-04-16 15:50:30','ravi@support.com','support',''),(100,'/api/v1/showAllOrders','2026-04-16 15:50:30','ravi@support.com','support',''),(101,'/api/v1/changeOrderStatus?id=12&status=shipped','2026-04-16 15:50:32','ravi@support.com','support',''),(102,'/api/v1/showAllOrders','2026-04-16 15:50:32','ravi@support.com','support',''),(103,'/api/v1/fetchBooks','2026-04-16 15:51:10','rahul@student.com','student',''),(104,'/api/v1/fetchBooks','2026-04-16 15:51:10','rahul@student.com','student',''),(105,'/api/v1/showMyOrders','2026-04-16 15:51:13','rahul@student.com','student',''),(106,'/api/v1/showMyOrders','2026-04-16 15:51:13','rahul@student.com','student',''),(107,'/api/v1/fetchBooks','2026-04-16 15:52:21','rahul@student.com','student',''),(108,'/api/v1/fetchBooks','2026-04-16 15:52:21','rahul@student.com','student',''),(109,'/api/v1/showAllOrders','2026-04-16 15:52:34','ravi@support.com','support',''),(110,'/api/v1/showAllOrders','2026-04-16 15:52:34','ravi@support.com','support',''),(111,'/api/v1/changeOrderStatus?id=14&status=new','2026-04-16 15:52:36','ravi@support.com','support',''),(112,'/api/v1/showAllOrders','2026-04-16 15:52:36','ravi@support.com','support',''),(113,'/api/v1/fetchBooks','2026-04-16 15:52:44','rahul@student.com','student',''),(114,'/api/v1/fetchBooks','2026-04-16 15:52:44','rahul@student.com','student',''),(115,'/api/v1/showMyOrders','2026-04-16 15:52:47','rahul@student.com','student',''),(116,'/api/v1/showMyOrders','2026-04-16 15:52:47','rahul@student.com','student',''),(117,'/api/v1/cancelOrder?id=14','2026-04-16 15:52:59','rahul@student.com','student','{}'),(118,'/api/v1/fetchBooks','2026-04-16 15:53:04','rahul@student.com','student',''),(119,'/api/v1/fetchBooks','2026-04-16 15:53:04','rahul@student.com','student',''),(120,'/api/v1/fetchBooks','2026-04-16 15:54:54','neha@admin.com','admin',''),(121,'/api/v1/fetchBooks','2026-04-16 15:54:54','neha@admin.com','admin',''),(122,'/api/v1/fetchUniversities','2026-04-16 15:54:56','neha@admin.com','admin',''),(123,'/api/v1/fetchUniversities','2026-04-16 15:54:56','neha@admin.com','admin',''),(124,'/api/v1/fetchDepartments','2026-04-16 15:54:56','neha@admin.com','admin',''),(125,'/api/v1/fetchUniversities','2026-04-16 15:54:56','neha@admin.com','admin',''),(126,'/api/v1/fetchDepartments','2026-04-16 15:54:56','neha@admin.com','admin',''),(127,'/api/v1/fetchCourses','2026-04-16 15:54:57','neha@admin.com','admin',''),(128,'/api/v1/fetchUniversities','2026-04-16 15:54:57','neha@admin.com','admin',''),(129,'/api/v1/fetchCourses','2026-04-16 15:54:57','neha@admin.com','admin',''),(130,'/api/v1/fetchBooks','2026-04-16 15:54:57','neha@admin.com','admin',''),(131,'/api/v1/fetchSemesters','2026-04-16 15:54:57','neha@admin.com','admin',''),(132,'/api/v1/fetchUniversities','2026-04-16 15:54:57','neha@admin.com','admin',''),(133,'/api/v1/fetchInstructor','2026-04-16 15:54:57','neha@admin.com','admin',''),(134,'/api/v1/fetchSemesters','2026-04-16 15:54:57','neha@admin.com','admin',''),(135,'/api/v1/fetchBooks','2026-04-16 15:54:58','neha@admin.com','admin',''),(136,'/api/v1/fetchBooks','2026-04-16 15:54:58','neha@admin.com','admin',''),(137,'/api/v1/removeBook?id=1','2026-04-16 15:55:02','neha@admin.com','admin',''),(138,'/api/v1/removeBook?id=1','2026-04-16 15:55:02','neha@admin.com','admin',''),(139,'/api/v1/removeBook?id=1','2026-04-16 15:55:02','neha@admin.com','admin',''),(140,'/api/v1/removeBook?id=1','2026-04-16 15:55:02','neha@admin.com','admin',''),(141,'/api/v1/removeBook?id=1','2026-04-16 15:55:02','neha@admin.com','admin',''),(142,'/api/v1/fetchBooks','2026-04-16 15:55:44','neha@admin.com','admin',''),(143,'/api/v1/fetchBooks','2026-04-16 15:55:44','neha@admin.com','admin',''),(144,'/api/v1/removeBook?id=1','2026-04-16 15:55:47','neha@admin.com','admin',''),(145,'/api/v1/removeBook?id=1','2026-04-16 15:55:47','neha@admin.com','admin',''),(146,'/api/v1/fetchBooks','2026-04-16 16:01:36','neha@admin.com','admin',''),(147,'/api/v1/fetchBooks','2026-04-16 16:01:36','neha@admin.com','admin',''),(148,'/api/v1/removeBook?id=6','2026-04-16 16:01:43','neha@admin.com','admin',''),(149,'/api/v1/removeBook?id=1','2026-04-16 16:02:51','neha@admin.com','admin',''),(150,'/api/v1/fetchBooks','2026-04-16 16:05:11','neha@admin.com','admin',''),(151,'/api/v1/fetchBooks','2026-04-16 16:05:11','neha@admin.com','admin',''),(152,'/api/v1/removeBook?id=6','2026-04-16 16:05:15','neha@admin.com','admin',''),(153,'/api/v1/removeBook?id=1','2026-04-16 16:05:26','neha@admin.com','admin',''),(154,'/api/v1/removeBook?id=2','2026-04-16 16:05:35','neha@admin.com','admin',''),(155,'/api/v1/fetchBooks','2026-04-16 16:08:59','neha@admin.com','admin',''),(156,'/api/v1/fetchBooks','2026-04-16 16:08:59','neha@admin.com','admin',''),(157,'/api/v1/removeBook?id=6','2026-04-16 16:09:02','neha@admin.com','admin',''),(158,'/api/v1/fetchBooks','2026-04-16 16:11:04','neha@admin.com','admin',''),(159,'/api/v1/fetchBooks','2026-04-16 16:11:04','neha@admin.com','admin',''),(160,'/api/v1/removeBook?id=6','2026-04-16 16:11:08','neha@admin.com','admin',''),(161,'/api/v1/fetchAdmins','2026-04-16 16:26:58','amit@superadmin.com','superadmin',''),(162,'/api/v1/fetchAdmins','2026-04-16 16:26:58','amit@superadmin.com','superadmin',''),(163,'/api/v1/fetchAdmins','2026-04-16 16:27:08','amit@superadmin.com','superadmin',''),(164,'/api/v1/fetchAdmins','2026-04-16 16:27:08','amit@superadmin.com','superadmin',''),(165,'/api/v1/fetchAdmins','2026-04-16 16:27:13','amit@superadmin.com','superadmin',''),(166,'/api/v1/fetchAdmins','2026-04-16 16:27:13','amit@superadmin.com','superadmin','');
/*!40000 ALTER TABLE `request_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `review`
--

DROP TABLE IF EXISTS `review`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `review` (
  `review_id` int NOT NULL AUTO_INCREMENT,
  `student_id` int DEFAULT NULL,
  `book_id` int DEFAULT NULL,
  `rating` int DEFAULT NULL,
  `review_text` text,
  `review_date` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`review_id`),
  UNIQUE KEY `unique_review` (`student_id`,`book_id`),
  KEY `book_id` (`book_id`),
  CONSTRAINT `review_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `student` (`student_id`),
  CONSTRAINT `review_ibfk_2` FOREIGN KEY (`book_id`) REFERENCES `book` (`book_id`),
  CONSTRAINT `review_chk_1` CHECK ((`rating` between 1 and 5))
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `review`
--

LOCK TABLES `review` WRITE;
/*!40000 ALTER TABLE `review` DISABLE KEYS */;
INSERT INTO `review` VALUES (1,1,1,5,'Excellent DB book','2026-04-09 23:35:51'),(2,2,1,4,'Very useful','2026-04-09 23:35:51'),(3,1,2,5,'Algorithm bible','2026-04-09 23:35:51');
/*!40000 ALTER TABLE `review` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `semester`
--

DROP TABLE IF EXISTS `semester`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `semester` (
  `sem_id` int NOT NULL AUTO_INCREMENT,
  `year` int NOT NULL,
  `season` varchar(10) NOT NULL,
  `course_id` int NOT NULL,
  `instructor_id` int NOT NULL,
  `university_id` int NOT NULL,
  PRIMARY KEY (`sem_id`),
  UNIQUE KEY `year` (`year`,`season`,`course_id`,`instructor_id`),
  KEY `fk_sem_course` (`course_id`),
  KEY `fk_sem_instructor` (`instructor_id`),
  CONSTRAINT `fk_sem_course` FOREIGN KEY (`course_id`) REFERENCES `course` (`course_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sem_instructor` FOREIGN KEY (`instructor_id`) REFERENCES `instructor` (`instructor_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `semester`
--

LOCK TABLES `semester` WRITE;
/*!40000 ALTER TABLE `semester` DISABLE KEYS */;
/*!40000 ALTER TABLE `semester` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `semester_book`
--

DROP TABLE IF EXISTS `semester_book`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `semester_book` (
  `sem_id` int NOT NULL,
  `book_id` int NOT NULL,
  PRIMARY KEY (`sem_id`,`book_id`),
  KEY `fk_sb_book` (`book_id`),
  CONSTRAINT `fk_sb_book` FOREIGN KEY (`book_id`) REFERENCES `book` (`book_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sb_sem` FOREIGN KEY (`sem_id`) REFERENCES `semester` (`sem_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `semester_book`
--

LOCK TABLES `semester_book` WRITE;
/*!40000 ALTER TABLE `semester_book` DISABLE KEYS */;
/*!40000 ALTER TABLE `semester_book` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `student`
--

DROP TABLE IF EXISTS `student`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student` (
  `student_id` int NOT NULL,
  `date_of_birth` date DEFAULT NULL,
  `university_id` int DEFAULT NULL,
  `major` varchar(100) DEFAULT NULL,
  `status` enum('graduate','undergraduate') DEFAULT NULL,
  `year_of_study` int DEFAULT NULL,
  PRIMARY KEY (`student_id`),
  KEY `fk_student_university` (`university_id`),
  CONSTRAINT `fk_student_university` FOREIGN KEY (`university_id`) REFERENCES `university` (`university_id`),
  CONSTRAINT `student_ibfk_2` FOREIGN KEY (`student_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `student`
--

LOCK TABLES `student` WRITE;
/*!40000 ALTER TABLE `student` DISABLE KEYS */;
INSERT INTO `student` VALUES (1,'2002-05-10',1,'Computer Science','undergraduate',3),(2,'2001-08-20',1,'Electronics','graduate',1),(9,'2003-05-12',1,'Computer Science','undergraduate',2),(10,'2002-08-21',1,'Mechanical Engineering','undergraduate',3),(11,'2003-09-10',1,'Information Tefchnology','undergraduate',1),(12,'2000-03-15',1,'Civil Engineering','graduate',2);
/*!40000 ALTER TABLE `student` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subcategory`
--

DROP TABLE IF EXISTS `subcategory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `subcategory` (
  `subcategory_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `category_id` int DEFAULT NULL,
  PRIMARY KEY (`subcategory_id`),
  KEY `category_id` (`category_id`),
  CONSTRAINT `subcategory_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `category` (`category_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subcategory`
--

LOCK TABLES `subcategory` WRITE;
/*!40000 ALTER TABLE `subcategory` DISABLE KEYS */;
INSERT INTO `subcategory` VALUES (1,'DBMS',1),(2,'Algorithms',1),(3,'Digital Electronics',2),(10,'Data Structures',1);
/*!40000 ALTER TABLE `subcategory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `super_admin`
--

DROP TABLE IF EXISTS `super_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `super_admin` (
  `employee_id` int NOT NULL,
  PRIMARY KEY (`employee_id`),
  CONSTRAINT `super_admin_ibfk_1` FOREIGN KEY (`employee_id`) REFERENCES `administrator` (`employee_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `super_admin`
--

LOCK TABLES `super_admin` WRITE;
/*!40000 ALTER TABLE `super_admin` DISABLE KEYS */;
INSERT INTO `super_admin` VALUES (5);
/*!40000 ALTER TABLE `super_admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ticket`
--

DROP TABLE IF EXISTS `ticket`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ticket` (
  `ticket_id` int NOT NULL AUTO_INCREMENT,
  `created_by_user_id` int DEFAULT NULL,
  `category` enum('user_profile','products','cart','orders','other') DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `description` text,
  `solution_description` text,
  `created_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `completion_date` datetime DEFAULT NULL,
  `status` enum('new','assigned','in-process','completed') DEFAULT NULL,
  `assigned_admin_id` int DEFAULT NULL,
  PRIMARY KEY (`ticket_id`),
  KEY `created_by_user_id` (`created_by_user_id`),
  KEY `assigned_admin_id` (`assigned_admin_id`),
  CONSTRAINT `ticket_ibfk_1` FOREIGN KEY (`created_by_user_id`) REFERENCES `user` (`user_id`),
  CONSTRAINT `ticket_ibfk_2` FOREIGN KEY (`assigned_admin_id`) REFERENCES `administrator` (`employee_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ticket`
--

LOCK TABLES `ticket` WRITE;
/*!40000 ALTER TABLE `ticket` DISABLE KEYS */;
INSERT INTO `ticket` VALUES (3,3,'cart','cart failed','haha testing','jygj\n','2026-04-11 16:07:42','2026-04-15 04:35:00','completed',4),(4,3,'cart','cart failed','haha testing','hihihii','2026-04-11 19:07:44','2026-04-15 04:34:43','completed',4),(5,1,'cart','dasfads','dsfsa',NULL,'2026-04-12 20:11:20',NULL,'assigned',5),(6,1,'other','hihihaha','hahahihi',NULL,'2026-04-12 20:13:03',NULL,'assigned',4);
/*!40000 ALTER TABLE `ticket` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ticket_status_history`
--

DROP TABLE IF EXISTS `ticket_status_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ticket_status_history` (
  `history_id` int NOT NULL AUTO_INCREMENT,
  `ticket_id` int DEFAULT NULL,
  `changed_by` int DEFAULT NULL,
  `old_status` varchar(50) DEFAULT NULL,
  `new_status` varchar(50) DEFAULT NULL,
  `change_date` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`history_id`),
  KEY `ticket_id` (`ticket_id`),
  KEY `changed_by` (`changed_by`),
  CONSTRAINT `ticket_status_history_ibfk_1` FOREIGN KEY (`ticket_id`) REFERENCES `ticket` (`ticket_id`),
  CONSTRAINT `ticket_status_history_ibfk_2` FOREIGN KEY (`changed_by`) REFERENCES `user` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ticket_status_history`
--

LOCK TABLES `ticket_status_history` WRITE;
/*!40000 ALTER TABLE `ticket_status_history` DISABLE KEYS */;
INSERT INTO `ticket_status_history` VALUES (1,4,3,'new','assigned','2026-04-11 19:10:04'),(2,4,4,'assigned','in-process','2026-04-11 19:15:30'),(3,3,3,'new','assigned','2026-04-13 01:31:49'),(4,5,3,'new','assigned','2026-04-13 01:51:14'),(5,4,4,'in-process','completed','2026-04-15 10:04:43'),(6,3,4,'assigned','in-process','2026-04-15 10:04:52'),(7,3,4,'in-process','completed','2026-04-15 10:04:59'),(8,6,3,'new','assigned','2026-04-15 10:07:14');
/*!40000 ALTER TABLE `ticket_status_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `university`
--

DROP TABLE IF EXISTS `university`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `university` (
  `university_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `address` text,
  `rep_first_name` varchar(100) DEFAULT NULL,
  `rep_last_name` varchar(100) DEFAULT NULL,
  `rep_email` varchar(255) DEFAULT NULL,
  `rep_phone` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`university_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `university`
--

LOCK TABLES `university` WRITE;
/*!40000 ALTER TABLE `university` DISABLE KEYS */;
INSERT INTO `university` VALUES (1,'IIT Bhubaneswar','Odisha',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `university` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `user_id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(100) NOT NULL,
  `last_name` varchar(100) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `address` text,
  `phone` varchar(20) DEFAULT NULL,
  `password_hash` varchar(255) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'Rahul','Sharma','rahul@student.com','Odisha','1111111111','pass123'),(2,'Ananya','Patel','ananya@student.com','Delhi','2222222222','pass123'),(3,'Ravi','Kumar','ravi@support.com','Mumbai','3333333333','pass123'),(4,'Neha','Singh','neha@admin.com','Bangalore','4444444444','p'),(5,'Amit','Verma','amit@superadmin.com','Hyderabad','5555555555','pass123'),(9,'Aman','Sharma','aman1@student.com','Delhi','9000000001','iils3xvu1R'),(10,'Riya','Patel','riya2@student.com','Mumbai','9000000002','kJFVpuDnZ8'),(11,'haha','hihi','googoogaga@student.com','Lucknow','9000000005','j1R1sA3oap'),(12,'Sneha','Reddy','sneha4@student.com','Hyderabad','9000000004','eF2uiM0XuH');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-04-20 21:28:32
