-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Hôte : 127.0.0.1:3306
-- Généré le : ven. 19 avr. 2024 à 09:44
-- Version du serveur : 8.2.0
-- Version de PHP : 8.2.13

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de données : `forumcour`
--

-- --------------------------------------------------------

--
-- Structure de la table `admin`
--

DROP TABLE IF EXISTS `admin`;
CREATE TABLE IF NOT EXISTS `admin` (
  `topic_id` int NOT NULL,
  `message_id` int NOT NULL,
  `log_id` int NOT NULL,
  PRIMARY KEY (`topic_id`,`message_id`,`log_id`),
  KEY `message_id` (`message_id`),
  KEY `log_id` (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;

-- --------------------------------------------------------

--
-- Structure de la table `admin_logs_table`
--

DROP TABLE IF EXISTS `admin_logs_table`;
CREATE TABLE IF NOT EXISTS `admin_logs_table` (
  `log_id` int NOT NULL AUTO_INCREMENT,
  `action_descrbition` text COLLATE latin1_bin NOT NULL,
  `date_performed` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;

-- --------------------------------------------------------

--
-- Structure de la table `follow`
--

DROP TABLE IF EXISTS `follow`;
CREATE TABLE IF NOT EXISTS `follow` (
  `user_id` int NOT NULL,
  `topic_id` int NOT NULL,
  `followed_date` datetime NOT NULL,
  PRIMARY KEY (`user_id`,`topic_id`),
  KEY `topic_id` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;

-- --------------------------------------------------------

--
-- Structure de la table `friendship`
--

DROP TABLE IF EXISTS `friendship`;
CREATE TABLE IF NOT EXISTS `friendship` (
  `sender_id` int NOT NULL,
  `reciver_id` int NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`sender_id`,`reciver_id`),
  KEY `reciver_id` (`reciver_id`)
) ;

-- --------------------------------------------------------

--
-- Structure de la table `have`
--

DROP TABLE IF EXISTS `have`;
CREATE TABLE IF NOT EXISTS `have` (
  `topic_id` int NOT NULL,
  `tag_id` int NOT NULL,
  PRIMARY KEY (`topic_id`,`tag_id`),
  KEY `tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;

-- --------------------------------------------------------

--
-- Structure de la table `images_table`
--

DROP TABLE IF EXISTS `images_table`;
CREATE TABLE IF NOT EXISTS `images_table` (
  `image_id` int NOT NULL AUTO_INCREMENT,
  `image_origin_name` varchar(1000) COLLATE latin1_bin NOT NULL,
  `image_serv_name` varchar(1000) COLLATE latin1_bin NOT NULL,
  `image_link` varchar(1000) COLLATE latin1_bin NOT NULL DEFAULT '/assets/img/default.png',
  `message_id` int DEFAULT NULL,
  `topic_id` int DEFAULT NULL,
  PRIMARY KEY (`image_id`),
  UNIQUE KEY `image_serv_name` (`image_serv_name`),
  UNIQUE KEY `image_link` (`image_link`),
  KEY `message_id` (`message_id`),
  KEY `topic_id` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;

-- --------------------------------------------------------

--
-- Structure de la table `messages_table`
--

DROP TABLE IF EXISTS `messages_table`;
CREATE TABLE IF NOT EXISTS `messages_table` (
  `message_id` int NOT NULL AUTO_INCREMENT,
  `body` varchar(500) COLLATE latin1_bin NOT NULL,
  `date_sent` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `topic_id` int NOT NULL,
  `base_message_id` int DEFAULT NULL,
  `user_id` int NOT NULL,
  PRIMARY KEY (`message_id`),
  KEY `topic_id` (`topic_id`),
  KEY `base_message_id` (`base_message_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;

-- --------------------------------------------------------

--
-- Structure de la table `react_message`
--

DROP TABLE IF EXISTS `react_message`;
CREATE TABLE IF NOT EXISTS `react_message` (
  `user_id` int NOT NULL,
  `message_id` int NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`user_id`,`message_id`),
  KEY `message_id` (`message_id`)
) ;

-- --------------------------------------------------------

--
-- Structure de la table `react_topic`
--

DROP TABLE IF EXISTS `react_topic`;
CREATE TABLE IF NOT EXISTS `react_topic` (
  `user_id` int NOT NULL,
  `topic_id` int NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`user_id`,`topic_id`),
  KEY `topic_id` (`topic_id`)
) ;

-- --------------------------------------------------------

--
-- Structure de la table `tags_table`
--

DROP TABLE IF EXISTS `tags_table`;
CREATE TABLE IF NOT EXISTS `tags_table` (
  `tag_id` int NOT NULL,
  `tag_name` varchar(50) COLLATE latin1_bin NOT NULL,
  PRIMARY KEY (`tag_id`),
  UNIQUE KEY `tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin;

-- --------------------------------------------------------

--
-- Structure de la table `topics_table`
--

DROP TABLE IF EXISTS `topics_table`;
CREATE TABLE IF NOT EXISTS `topics_table` (
  `topic_id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(50) COLLATE latin1_bin NOT NULL,
  `body` varchar(1000) COLLATE latin1_bin NOT NULL,
  `creation_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` int NOT NULL,
  `is_private` tinyint NOT NULL DEFAULT '0',
  `user_id` int NOT NULL,
  PRIMARY KEY (`topic_id`),
  KEY `user_id` (`user_id`)
) ;

-- --------------------------------------------------------

--
-- Structure de la table `users_table`
--

DROP TABLE IF EXISTS `users_table`;
CREATE TABLE IF NOT EXISTS `users_table` (
  `user_id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(26) COLLATE latin1_bin NOT NULL,
  `email` varchar(50) COLLATE latin1_bin NOT NULL,
  `password` varchar(50) COLLATE latin1_bin NOT NULL,
  `registration_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_login_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `biography` varchar(200) COLLATE latin1_bin DEFAULT NULL,
  `isAdmin` tinyint NOT NULL DEFAULT '0',
  `isModerator` tinyint NOT NULL DEFAULT '0',
  `profile_pic` varchar(100) COLLATE latin1_bin NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ;

--
-- Contraintes pour les tables déchargées
--

--
-- Contraintes pour la table `admin`
--
ALTER TABLE `admin`
  ADD CONSTRAINT `admin_ibfk_1` FOREIGN KEY (`topic_id`) REFERENCES `topics_table` (`topic_id`),
  ADD CONSTRAINT `admin_ibfk_2` FOREIGN KEY (`message_id`) REFERENCES `messages_table` (`message_id`),
  ADD CONSTRAINT `admin_ibfk_3` FOREIGN KEY (`log_id`) REFERENCES `admin_logs_table` (`log_id`);

--
-- Contraintes pour la table `follow`
--
ALTER TABLE `follow`
  ADD CONSTRAINT `follow_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users_table` (`user_id`),
  ADD CONSTRAINT `follow_ibfk_2` FOREIGN KEY (`topic_id`) REFERENCES `topics_table` (`topic_id`);

--
-- Contraintes pour la table `friendship`
--
ALTER TABLE `friendship`
  ADD CONSTRAINT `friendship_ibfk_1` FOREIGN KEY (`sender_id`) REFERENCES `users_table` (`user_id`),
  ADD CONSTRAINT `friendship_ibfk_2` FOREIGN KEY (`reciver_id`) REFERENCES `users_table` (`user_id`);

--
-- Contraintes pour la table `have`
--
ALTER TABLE `have`
  ADD CONSTRAINT `have_ibfk_1` FOREIGN KEY (`topic_id`) REFERENCES `topics_table` (`topic_id`),
  ADD CONSTRAINT `have_ibfk_2` FOREIGN KEY (`tag_id`) REFERENCES `tags_table` (`tag_id`);

--
-- Contraintes pour la table `images_table`
--
ALTER TABLE `images_table`
  ADD CONSTRAINT `images_table_ibfk_1` FOREIGN KEY (`message_id`) REFERENCES `messages_table` (`message_id`),
  ADD CONSTRAINT `images_table_ibfk_2` FOREIGN KEY (`topic_id`) REFERENCES `topics_table` (`topic_id`);

--
-- Contraintes pour la table `messages_table`
--
ALTER TABLE `messages_table`
  ADD CONSTRAINT `messages_table_ibfk_1` FOREIGN KEY (`topic_id`) REFERENCES `topics_table` (`topic_id`),
  ADD CONSTRAINT `messages_table_ibfk_2` FOREIGN KEY (`base_message_id`) REFERENCES `messages_table` (`message_id`),
  ADD CONSTRAINT `messages_table_ibfk_3` FOREIGN KEY (`user_id`) REFERENCES `users_table` (`user_id`);

--
-- Contraintes pour la table `react_message`
--
ALTER TABLE `react_message`
  ADD CONSTRAINT `react_message_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users_table` (`user_id`),
  ADD CONSTRAINT `react_message_ibfk_2` FOREIGN KEY (`message_id`) REFERENCES `messages_table` (`message_id`);

--
-- Contraintes pour la table `react_topic`
--
ALTER TABLE `react_topic`
  ADD CONSTRAINT `react_topic_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users_table` (`user_id`),
  ADD CONSTRAINT `react_topic_ibfk_2` FOREIGN KEY (`topic_id`) REFERENCES `topics_table` (`topic_id`);

--
-- Contraintes pour la table `topics_table`
--
ALTER TABLE `topics_table`
  ADD CONSTRAINT `topics_table_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users_table` (`user_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
