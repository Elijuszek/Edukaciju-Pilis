CREATE TABLE IF NOT EXISTS `theme` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `fk_Organizerid` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `organizes` (`fk_Organizerid`),
  CONSTRAINT `organizes` FOREIGN KEY (`fk_Organizerid`) REFERENCES `organizer` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
