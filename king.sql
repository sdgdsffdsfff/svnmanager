CREATE DATABASE  IF NOT EXISTS `king` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `king`;
-- MySQL dump 10.13  Distrib 5.5.41, for debian-linux-gnu (i686)
--
-- Host: 127.0.0.1    Database: king
-- ------------------------------------------------------
-- Server version	5.5.41-0ubuntu0.14.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `config`
--

DROP TABLE IF EXISTS `config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) DEFAULT NULL,
  `content` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config`
--

LOCK TABLES `config` WRITE;
/*!40000 ALTER TABLE `config` DISABLE KEYS */;
/*!40000 ALTER TABLE `config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group`
--

DROP TABLE IF EXISTS `group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) DEFAULT NULL,
  `desc` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group`
--

LOCK TABLES `group` WRITE;
/*!40000 ALTER TABLE `group` DISABLE KEYS */;
INSERT INTO `group` VALUES (1,'Frontend Group','The CSS transform property lets you modify the coordinate space of the CSS visual formatting model.'),(2,'Backend Group','The transform-origin CSS property lets you modify the origin for transformations of an element.'),(3,'Manage Server','The transform-style CSS property determines if the children of the element are positioned in the 3D-space or are flattened in the plane of the element.');
/*!40000 ALTER TABLE `group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `up_file`
--

DROP TABLE IF EXISTS `up_file`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `up_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `path` varchar(200) DEFAULT NULL,
  `action` tinyint(4) DEFAULT NULL,
  `version` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `up_file`
--

LOCK TABLES `up_file` WRITE;
/*!40000 ALTER TABLE `up_file` DISABLE KEYS */;
INSERT INTO `up_file` VALUES (1,'traversing/pnum.js',1,135),(2,'var/pnum.js',3,135);
/*!40000 ALTER TABLE `up_file` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `version`
--

DROP TABLE IF EXISTS `version`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `version` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `version` int(11) DEFAULT NULL,
  `backup_path` varchar(300) DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  `comment` varchar(500) DEFAULT NULL,
  `list` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `version`
--

LOCK TABLES `version` WRITE;
/*!40000 ALTER TABLE `version` DISABLE KEYS */;
INSERT INTO `version` VALUES (1,133,NULL,'2015-03-26 09:14:09','','{\"traversing/push.js\":1,\"var/push.js\":3}'),(2,134,NULL,'2015-03-26 09:15:26','','{\"traversing/indexOf.js\":3,\"var/indexOf.js\":1}'),(3,135,NULL,'2015-03-26 09:19:10','','{\"traversing/pnum.js\":1,\"var/pnum.js\":3}'),(4,136,NULL,'2015-03-31 03:25:15','','{\"traversing/findFilter.js\":3,\"var/findFilter.js\":1}'),(5,137,NULL,'2015-03-31 03:28:13','','{\"traversing/push.js\":3,\"var/push.js\":1}'),(6,138,NULL,'2015-03-31 03:32:43','','{\"dimensions.js\":3,\"var/dimensions.js\":1}'),(7,139,NULL,'2015-04-01 07:55:40','','{\"core/arr.js\":1,\"traversing/arr.js\":3}'),(8,140,NULL,'2015-04-01 08:02:10','','{\"core/ready.js\":3,\"traversing/ready.js\":1}'),(9,141,NULL,'2015-04-02 03:04:40','','{\"core/class2type.js\":1,\"traversing/class2type.js\":3}'),(10,142,NULL,'2015-04-02 03:05:54','','{\"core/class2type.js\":2}');
/*!40000 ALTER TABLE `version` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `web_server`
--

DROP TABLE IF EXISTS `web_server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `web_server` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ip` varchar(45) DEFAULT NULL,
  `port` varchar(45) DEFAULT NULL,
  `name` varchar(45) DEFAULT NULL,
  `deploy_path` varchar(200) DEFAULT NULL,
  `version` int(11) DEFAULT NULL,
  `config` tinyint(4) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `group` int(2) DEFAULT '0',
  `internal_ip` varchar(45) DEFAULT NULL,
  `backup_path` varchar(200) DEFAULT NULL,
  `un_deploy_list` text,
  PRIMARY KEY (`id`),
  KEY `group_idx` (`group`),
  KEY `fk_web_server_1_idx` (`group`),
  KEY `grp` (`group`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `web_server`
--

LOCK TABLES `web_server` WRITE;
/*!40000 ALTER TABLE `web_server` DISABLE KEYS */;
INSERT INTO `web_server` VALUES (1,'115.29.186.80','4000','','/home/languid/svn/download/Vanille',124,0,2,1,'10.161.169.5','/home/languid/svn/download/backup','{\"core/class2type.js\":2,\"core/ready.js\":3,\"traversing/class2type.js\":3,\"traversing/ready.js\":1}'),(2,'121.199.28.9','4001','','/home/languid/svn/download/Lighting',124,NULL,2,1,'10.132.37.229','','{\"core/class2type.js\":2,\"core/ready.js\":3,\"traversing/class2type.js\":3,\"traversing/ready.js\":1}'),(3,'192.168.1.111','4000','111','/opt/wings',16508,NULL,0,0,'192.168.1.111','','{\"core/class2type.js\":2,\"core/ready.js\":3,\"traversing/class2type.js\":3,\"traversing/ready.js\":1}'),(4,'192.168.1.117','4000','Test','',16407,NULL,0,0,'192.168.1.117','','{\"core/class2type.js\":2,\"core/ready.js\":3,\"traversing/class2type.js\":3,\"traversing/ready.js\":1}'),(5,'192.168.31.152','4000','','',0,NULL,0,0,'192.168.31.152','','{\"core/class2type.js\":2,\"core/ready.js\":3,\"traversing/class2type.js\":3,\"traversing/ready.js\":1}');
/*!40000 ALTER TABLE `web_server` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-04-03 16:39:09
