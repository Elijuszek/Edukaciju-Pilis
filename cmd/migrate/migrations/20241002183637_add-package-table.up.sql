CREATE TABLE IF NOT EXISTS `package` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `price` float NOT NULL,
  `fk_Organizerid` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `creates` (`fk_Organizerid`),
  CONSTRAINT `creates` FOREIGN KEY (`fk_Organizerid`) REFERENCES `organizer` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
