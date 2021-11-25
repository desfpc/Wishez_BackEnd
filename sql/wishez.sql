-- --------------------------------------------------------
-- Хост:                         127.0.0.1
-- Версия сервера:               8.0.19 - MySQL Community Server - GPL
-- Операционная система:         Win64
-- HeidiSQL Версия:              11.2.0.6213
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Дамп структуры для таблица wishez.chat
DROP TABLE IF EXISTS `chat`;
CREATE TABLE IF NOT EXISTS `chat` (
                                      `id` int NOT NULL AUTO_INCREMENT,
                                      `chat_id` varchar(22) NOT NULL DEFAULT '0',
                                      `author` int NOT NULL DEFAULT '0',
                                      `group` int DEFAULT '0',
                                      `receiver` int DEFAULT '0',
                                      `content` varchar(255) NOT NULL,
                                      `date_add` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      `date_edit` datetime DEFAULT NULL,
                                      PRIMARY KEY (`id`),
                                      KEY `IDX_chat_chat_id` (`chat_id`),
                                      KEY `FK_chat_author` (`author`),
                                      KEY `FK_chat_group` (`group`),
                                      KEY `FK_chat_receiver` (`receiver`),
                                      CONSTRAINT `FK_chat_author` FOREIGN KEY (`author`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                                      CONSTRAINT `FK_chat_group` FOREIGN KEY (`group`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                                      CONSTRAINT `FK_chat_receiver` FOREIGN KEY (`receiver`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.files
DROP TABLE IF EXISTS `files`;
CREATE TABLE IF NOT EXISTS `files` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `name` varchar(255) NOT NULL,
                                       `ext` varchar(8) DEFAULT NULL,
                                       `size` varchar(8) NOT NULL,
                                       `author` int DEFAULT '0',
                                       `kind` enum('chat','wish','avatar') NOT NULL DEFAULT 'wish',
                                       `chat` int DEFAULT '0',
                                       `wish` int DEFAULT '0',
                                       `avatar` int DEFAULT '0',
                                       `date_add` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       PRIMARY KEY (`id`),
                                       KEY `IDX_files_kind` (`kind`) USING BTREE,
                                       KEY `FK_files_users` (`avatar`),
                                       KEY `FK_files_wish` (`wish`),
                                       CONSTRAINT `FK_files_users` FOREIGN KEY (`avatar`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                                       CONSTRAINT `FK_files_wish` FOREIGN KEY (`wish`) REFERENCES `wish` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.group
DROP TABLE IF EXISTS `group`;
CREATE TABLE IF NOT EXISTS `group` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `author` int NOT NULL,
                                       `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
                                       `visible` enum('hidden','normal','public') NOT NULL DEFAULT 'normal',
                                       `open_sum` decimal(20,4) NOT NULL DEFAULT '0.0000',
                                       `closed_sum` decimal(20,4) NOT NULL DEFAULT '0.0000',
                                       `date_add` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       PRIMARY KEY (`id`),
                                       KEY `FK_group_author` (`author`),
                                       CONSTRAINT `FK_group_author` FOREIGN KEY (`author`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.group_users
DROP TABLE IF EXISTS `group_users`;
CREATE TABLE IF NOT EXISTS `group_users` (
                                             `group_id` int NOT NULL,
                                             `user_id` int NOT NULL,
                                             `right` enum('admin','user') NOT NULL DEFAULT 'user',
                                             `date_add` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                             PRIMARY KEY (`group_id`,`user_id`),
                                             KEY `FK_user` (`user_id`),
                                             KEY `FK_group` (`group_id`),
                                             KEY `IDX_right` (`right`) USING BTREE,
                                             CONSTRAINT `FK_group` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                                             CONSTRAINT `FK_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.mounth_budget
DROP TABLE IF EXISTS `mounth_budget`;
CREATE TABLE IF NOT EXISTS `mounth_budget` (
                                               `user_id` int NOT NULL,
                                               `mounth` date NOT NULL,
                                               `budget` decimal(19,4) NOT NULL DEFAULT '0.0000',
                                               `used` decimal(19,4) DEFAULT '0.0000',
                                               `date_add` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                               `date_edit` datetime DEFAULT NULL,
                                               PRIMARY KEY (`user_id`,`mounth`),
                                               KEY `IDX_mounth_budget_mounth` (`mounth`),
                                               CONSTRAINT `FK_mounth_budget_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Бюджет пользователя на месяц, для авто-выбора желания';

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.recommendation
DROP TABLE IF EXISTS `recommendation`;
CREATE TABLE IF NOT EXISTS `recommendation` (
                                                `id` int NOT NULL AUTO_INCREMENT,
                                                `mounth` date NOT NULL,
                                                `group_id` int NOT NULL,
                                                `wish_id` int NOT NULL,
                                                PRIMARY KEY (`id`),
                                                KEY `IDX_recommendation_mounth` (`mounth`),
                                                KEY `FK_recommendation_group` (`group_id`),
                                                KEY `FK_recommendation_wish` (`wish_id`),
                                                CONSTRAINT `FK_recommendation_group` FOREIGN KEY (`group_id`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                                                CONSTRAINT `FK_recommendation_wish` FOREIGN KEY (`wish_id`) REFERENCES `wish` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.users
DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `email` varchar(255) NOT NULL DEFAULT '',
                                       `pass` varchar(255) NOT NULL DEFAULT '',
                                       `fio` varchar(255) DEFAULT '',
                                       `sex` enum('male','female','other') DEFAULT NULL,
                                       `telegram` varchar(255) DEFAULT NULL,
                                       `instagrtam` varchar(255) DEFAULT NULL,
                                       `twitter` varchar(255) DEFAULT NULL,
                                       `facebook` varchar(255) DEFAULT NULL,
                                       `phone` varchar(255) DEFAULT NULL,
                                       `role` enum('user','moderator','administrator') NOT NULL,
                                       `avatar` int DEFAULT NULL,
                                       `google` varchar(255) DEFAULT NULL,
                                       `date_add` datetime DEFAULT NULL,
                                       PRIMARY KEY (`id`),
                                       UNIQUE KEY `UNQ_email` (`email`) USING BTREE,
                                       KEY `FK_users_files` (`avatar`),
                                       CONSTRAINT `FK_users_files` FOREIGN KEY (`avatar`) REFERENCES `files` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.users_friends
DROP TABLE IF EXISTS `users_friends`;
CREATE TABLE IF NOT EXISTS `users_friends` (
                                               `user_id` int NOT NULL,
                                               `friend_id` int NOT NULL,
                                               `approved` enum('Y','N') NOT NULL DEFAULT 'N',
                                               `date_approved` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                               PRIMARY KEY (`user_id`,`friend_id`),
                                               KEY `FK_users_friends_users_1` (`user_id`),
                                               KEY `FK_users_friends_users_2` (`friend_id`),
                                               CONSTRAINT `FK_users_friends_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                                               CONSTRAINT `FK_users_friends_users_2` FOREIGN KEY (`friend_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Экспортируемые данные не выделены.

-- Дамп структуры для таблица wishez.wish
DROP TABLE IF EXISTS `wish`;
CREATE TABLE IF NOT EXISTS `wish` (
                                      `id` int NOT NULL AUTO_INCREMENT,
                                      `author` int NOT NULL,
                                      `name` varchar(255) NOT NULL,
                                      `content` text,
                                      `cost` decimal(19,4) NOT NULL,
                                      `currency` enum('₽','$','€') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '₽',
                                      `group` int NOT NULL,
                                      `status` enum('open','closed') DEFAULT NULL,
                                      `date_add` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      `date_edit` datetime DEFAULT NULL,
                                      `priority` int NOT NULL DEFAULT '4',
                                      PRIMARY KEY (`id`),
                                      KEY `IDX_wish_status` (`status`) USING BTREE,
                                      KEY `FK_wish_group` (`group`),
                                      KEY `FK_wish_users` (`author`),
                                      CONSTRAINT `FK_wish_group` FOREIGN KEY (`group`) REFERENCES `group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
                                      CONSTRAINT `FK_wish_users` FOREIGN KEY (`author`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Желание';

-- Экспортируемые данные не выделены.

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;