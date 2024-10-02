CREATE TABLE IF NOT EXISTS `location` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(255) NOT NULL,
  `longitude` double(9,6) DEFAULT NULL,
  `latitude` double(9,6) DEFAULT NULL,
  `fk_Activityid` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `given` (`fk_Activityid`),
  CONSTRAINT `given` FOREIGN KEY (`fk_Activityid`) REFERENCES `activity` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
