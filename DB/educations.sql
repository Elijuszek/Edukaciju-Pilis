SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

CREATE TABLE `activity` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `basePrice` float NOT NULL,
  `creationDate` datetime NOT NULL DEFAULT current_timestamp(),
  `hidden` tinyint(1) NOT NULL DEFAULT 0,
  `verified` tinyint(1) NOT NULL DEFAULT 0,
  `category` int(11) NOT NULL,
  `averageRating` float NOT NULL DEFAULT 0,
  `fk_Packageid` int(11) NOT NULL,
  `fk_Themeid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `administrator` (
  `id` int(11) NOT NULL,
  `securityLevel` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `category` (
  `id_Category` int(11) NOT NULL,
  `name` char(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `category` (`id_Category`, `name`) VALUES
(1, 'Education'),
(2, 'Event'),
(3, 'Service'),
(4, 'Other');

CREATE TABLE `entityimage` (
  `id` int(11) NOT NULL,
  `entityType` varchar(255) NOT NULL,
  `entityFk` int(11) NOT NULL,
  `fk_Imageid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `image` (
  `id` int(11) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `filePath` varchar(255) NOT NULL,
  `url` varchar(255) NOT NULL,
  `uploadDate` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `location` (
  `id` int(11) NOT NULL,
  `address` varchar(255) NOT NULL,
  `longitude` double(9,6) DEFAULT NULL,
  `latitude` double(9,6) DEFAULT NULL,
  `fk_Activityid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `organizer` (
  `id` int(11) NOT NULL,
  `description` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `package` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `price` float NOT NULL,
  `fk_Organizerid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `review` (
  `id` int(11) NOT NULL,
  `date` date NOT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `rating` int(11) NOT NULL,
  `fk_Userid` int(11) NOT NULL,
  `fk_Activityid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `subscribers` (
  `id` int(11) NOT NULL,
  `email` varchar(255) NOT NULL,
  `subscriptionDate` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `theme` (
  `id` int(11) NOT NULL,
  `title` varchar(255) NOT NULL,
  `fk_Organizerid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `user` (
  `id` int(11) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `registrationDate` datetime NOT NULL DEFAULT current_timestamp(),
  `lastLoginDate` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE `activity`
  ADD PRIMARY KEY (`id`),
  ADD KEY `category` (`category`),
  ADD KEY `fk_Packageid` (`fk_Packageid`),
  ADD KEY `contains` (`fk_Themeid`);

ALTER TABLE `administrator`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `category`
  ADD PRIMARY KEY (`id_Category`);

ALTER TABLE `entityimage`
  ADD PRIMARY KEY (`id`),
  ADD KEY `mapping` (`fk_Imageid`);

ALTER TABLE `image`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `location`
  ADD PRIMARY KEY (`id`),
  ADD KEY `given` (`fk_Activityid`);

ALTER TABLE `organizer`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `package`
  ADD PRIMARY KEY (`id`),
  ADD KEY `creates` (`fk_Organizerid`);

ALTER TABLE `review`
  ADD PRIMARY KEY (`id`),
  ADD KEY `writes` (`fk_Userid`),
  ADD KEY `fk_Activityid` (`fk_Activityid`);

ALTER TABLE `subscribers`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `theme`
  ADD PRIMARY KEY (`id`),
  ADD KEY `organizes` (`fk_Organizerid`);

ALTER TABLE `user`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `activity`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `category`
  MODIFY `id_Category` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

ALTER TABLE `entityimage`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `image`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `location`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `package`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `review`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `subscribers`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `theme`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `user`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `activity`
  ADD CONSTRAINT `activity_ibfk_1` FOREIGN KEY (`category`) REFERENCES `category` (`id_Category`),
  ADD CONSTRAINT `activity_ibfk_2` FOREIGN KEY (`fk_Packageid`) REFERENCES `package` (`id`),
  ADD CONSTRAINT `contains` FOREIGN KEY (`fk_Themeid`) REFERENCES `theme` (`id`);

ALTER TABLE `administrator`
  ADD CONSTRAINT `administrator_ibfk_1` FOREIGN KEY (`id`) REFERENCES `user` (`id`) ON DELETE CASCADE;

ALTER TABLE `entityimage`
  ADD CONSTRAINT `mapping` FOREIGN KEY (`fk_Imageid`) REFERENCES `image` (`id`);

ALTER TABLE `location`
  ADD CONSTRAINT `given` FOREIGN KEY (`fk_Activityid`) REFERENCES `activity` (`id`);

ALTER TABLE `organizer`
  ADD CONSTRAINT `organizer_ibfk_1` FOREIGN KEY (`id`) REFERENCES `user` (`id`) ON DELETE CASCADE;

ALTER TABLE `package`
  ADD CONSTRAINT `creates` FOREIGN KEY (`fk_Organizerid`) REFERENCES `organizer` (`id`);

ALTER TABLE `review`
  ADD CONSTRAINT `review_ibfk_1` FOREIGN KEY (`fk_Activityid`) REFERENCES `activity` (`id`),
  ADD CONSTRAINT `writes` FOREIGN KEY (`fk_Userid`) REFERENCES `user` (`id`);

ALTER TABLE `theme`
  ADD CONSTRAINT `organizes` FOREIGN KEY (`fk_Organizerid`) REFERENCES `organizer` (`id`);
COMMIT;

