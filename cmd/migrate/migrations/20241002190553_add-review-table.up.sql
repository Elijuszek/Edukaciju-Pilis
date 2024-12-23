CREATE TABLE IF NOT EXISTS `review` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `comment` varchar(255) DEFAULT NULL,
  `rating` int(11) NOT NULL,
  `fk_Userid` int(11) NOT NULL,
  `fk_Activityid` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `writes` (`fk_Userid`),
  KEY `fk_Activityid` (`fk_Activityid`),
  CONSTRAINT `review_ibfk_1` FOREIGN KEY (`fk_Activityid`) REFERENCES `activity` (`id`) ON DELETE CASCADE,
  CONSTRAINT `writes` FOREIGN KEY (`fk_Userid`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;