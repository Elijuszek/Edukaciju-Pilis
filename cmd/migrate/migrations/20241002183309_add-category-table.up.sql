CREATE TABLE IF NOT EXISTS `category` (
  `id_Category` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(11) NOT NULL,  -- Ensure you have the right column type
  PRIMARY KEY (`id_Category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;