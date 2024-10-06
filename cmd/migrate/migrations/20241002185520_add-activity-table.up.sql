CREATE TABLE IF NOT EXISTS `activity` (
   `id` int(11) NOT NULL AUTO_INCREMENT,
   `name` varchar(255) NOT NULL,
   `description` varchar(255) NOT NULL,
   `basePrice` float NOT NULL,
   `creationDate` datetime NOT NULL DEFAULT current_timestamp(),
   `hidden` tinyint(1) NOT NULL DEFAULT 0,
   `verified` tinyint(1) NOT NULL DEFAULT 0,
   `category` int(11) NOT NULL,
   `averageRating` float NOT NULL DEFAULT 0,
   `fk_Packageid` int(11) NOT NULL,
   PRIMARY KEY (`id`),
   KEY `category` (`category`),
   KEY `fk_Packageid` (`fk_Packageid`),
   CONSTRAINT `activity_ibfk_1` FOREIGN KEY (`category`) REFERENCES `category` (`id_Category`),
   CONSTRAINT `activity_ibfk_2` FOREIGN KEY (`fk_Packageid`) REFERENCES `package` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;