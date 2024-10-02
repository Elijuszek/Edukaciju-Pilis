CREATE TABLE IF NOT EXISTS `entityimage` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `entityType` varchar(255) NOT NULL,
  `entityFk` int(11) NOT NULL,
  `fk_Imageid` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `mapping` (`fk_Imageid`),
  CONSTRAINT `mapping` FOREIGN KEY (`fk_Imageid`) REFERENCES `image` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
