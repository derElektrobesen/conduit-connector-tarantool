DROP DATABASE IF EXISTS mPOP;
CREATE DATABASE mPOP;

USE mPOP;

CREATE TABLE `user` (
  `ID` int unsigned NOT NULL AUTO_INCREMENT,
  `CAGId` bigint NOT NULL DEFAULT '0',
  `Domain` int unsigned NOT NULL DEFAULT '0',
  `Username` varchar(255) NOT NULL DEFAULT '',
  `Flags` int unsigned NOT NULL DEFAULT '0',
  `MailboxLimit` int unsigned DEFAULT NULL,
  `Forward` varchar(128) DEFAULT NULL,
  `RealName` varchar(64) DEFAULT NULL,
  `Comment` varchar(128) DEFAULT NULL,
  `Storage` varchar(128) NOT NULL DEFAULT '',
  `Target` varchar(28) NOT NULL DEFAULT '',
  `PwdCode` int unsigned NOT NULL DEFAULT '0',
  `AttachStorage` varchar(128) NOT NULL DEFAULT '',
  `Flags2` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `User` (`Username`,`Domain`),
  KEY `Domain` (`Domain`,`Storage`),
  KEY `DomainAttach` (`Domain`,`AttachStorage`)
) ENGINE=InnoDB DEFAULT CHARSET=cp1251;

INSERT INTO `user` (`CAGId`, `Domain`, `Username`, `Flags`, `MailboxLimit`, `Forward`, `RealName`, `Comment`, `Storage`, `Target`, `PwdCode`, `AttachStorage`, `Flags2`) VALUES 
(1, 1, 'user1@domain.com', 0, 100, NULL, 'Имя1', 'Комментарий1', '/var/mail/user1', 'target1', 123, '/var/attach/user1', 0),
(2, 2, 'user2@domain.com', 0, 100, NULL, 'Имя2', 'Комментарий2', '/var/mail/user2', 'target2', 123, '/var/attach/user2', 0),
(3, 3, 'user3@domain.com', 0, 100, NULL, 'Имя3', 'Комментарий3', '/var/mail/user3', 'target3', 123, '/var/attach/user3', 0),
(4, 4, 'user4@domain.com', 0, 100, NULL, 'Имя4', 'Комментарий4', '/var/mail/user4', 'target4', 123, '/var/attach/user4', 0),
(5, 5, 'user5@domain.com', 0, 100, NULL, 'Имя5', 'Комментарий5', '/var/mail/user5', 'target5', 123, '/var/attach/user5', 0),
(6, 6, 'user6@domain.com', 0, 100, NULL, 'Имя6', 'Комментарий6', '/var/mail/user6', 'target6', 123, '/var/attach/user6', 0),
(7, 7, 'user7@domain.com', 0, 100, NULL, 'Имя7', 'Комментарий7', '/var/mail/user7', 'target7', 123, '/var/attach/user7', 0),
(8, 8, 'user8@domain.com', 0, 100, NULL, 'Имя8', 'Комментарий8', '/var/mail/user8', 'target8', 123, '/var/attach/user8', 0),
(9, 9, 'user9@domain.com', 0, 100, NULL, 'Имя9', 'Комментарий9', '/var/mail/user9', 'target9', 123, '/var/attach/user9', 0),
(10, 10, 'user10@domain.com', 0, 100, NULL, 'Имя10', 'Комментарий10', '/var/mail/user10', 'target10', 123, '/var/attach/user10', 0),
(11, 11, 'user11@domain.com', 0, 100, NULL, 'Имя11', 'Комментарий11', '/var/mail/user11', 'target11', 123, '/var/attach/user11', 0),
(12, 12, 'user12@domain.com', 0, 100, NULL, 'Имя12', 'Комментарий12', '/var/mail/user12', 'target12', 123, '/var/attach/user12', 0),
(13, 13, 'user13@domain.com', 0, 100, NULL, 'Имя13', 'Комментарий13', '/var/mail/user13', 'target13', 123, '/var/attach/user13', 0),
(14, 14, 'user14@domain.com', 0, 100, NULL, 'Имя14', 'Комментарий14', '/var/mail/user14', 'target14', 123, '/var/attach/user14', 0),
(15, 15, 'user15@domain.com', 0, 100, NULL, 'Имя15', 'Комментарий15', '/var/mail/user15', 'target15', 123, '/var/attach/user15', 0),
(16, 16, 'user16@domain.com', 0, 100, NULL, 'Имя16', 'Комментарий16', '/var/mail/user16', 'target16', 123, '/var/attach/user16', 0),
(17, 17, 'user17@domain.com', 0, 100, NULL, 'Имя17', 'Комментарий17', '/var/mail/user17', 'target17', 123, '/var/attach/user17', 0),
(18, 18, 'user18@domain.com', 0, 100, NULL, 'Имя18', 'Комментарий18', '/var/mail/user18', 'target18', 123, '/var/attach/user18', 0),
(19, 19, 'user19@domain.com', 0, 100, NULL, 'Имя19', 'Комментарий19', '/var/mail/user19', 'target19', 123, '/var/attach/user19', 0),
(20, 20, 'user20@domain.com', 0, 100, NULL, 'Имя20', 'Комментарий20', '/var/mail/user20', 'target20', 123, '/var/attach/user20', 0),
(21, 21, 'user21@domain.com', 0, 100, NULL, 'Имя21', 'Комментарий21', '/var/mail/user21', 'target21', 123, '/var/attach/user21', 0),
(22, 22, 'user22@domain.com', 0, 100, NULL, 'Имя22', 'Комментарий22', '/var/mail/user22', 'target22', 123, '/var/attach/user22', 0),
(23, 23, 'user23@domain.com', 0, 100, NULL, 'Имя23', 'Комментарий23', '/var/mail/user23', 'target23', 123, '/var/attach/user23', 0),
(24, 24, 'user24@domain.com', 0, 100, NULL, 'Имя24', 'Комментарий24', '/var/mail/user24', 'target24', 123, '/var/attach/user24', 0),
(25, 25, 'user25@domain.com', 0, 100, NULL, 'Имя25', 'Комментарий25', '/var/mail/user25', 'target25', 123, '/var/attach/user25', 0),
(26, 26, 'user26@domain.com', 0, 100, NULL, 'Имя26', 'Комментарий26', '/var/mail/user26', 'target26', 123, '/var/attach/user26', 0),
(27, 27, 'user27@domain.com', 0, 100, NULL, 'Имя27', 'Комментарий27', '/var/mail/user27', 'target27', 123, '/var/attach/user27', 0),
(28, 28, 'user28@domain.com', 0, 100, NULL, 'Имя28', 'Комментарий28', '/var/mail/user28', 'target28', 123, '/var/attach/user28', 0),
(29, 29, 'user29@domain.com', 0, 100, NULL, 'Имя29', 'Комментарий29', '/var/mail/user29', 'target29', 123, '/var/attach/user29', 0),
(30, 30, 'user30@domain.com', 0, 100, NULL, 'Имя30', 'Комментарий30', '/var/mail/user30', 'target30', 123, '/var/attach/user30', 0),
(31, 31, 'user31@domain.com', 0, 100, NULL, 'Имя31', 'Комментарий31', '/var/mail/user31', 'target31', 123, '/var/attach/user31', 0),
(32, 32, 'user32@domain.com', 0, 100, NULL, 'Имя32', 'Комментарий32', '/var/mail/user32', 'target32', 123, '/var/attach/user32', 0),
(33, 33, 'user33@domain.com', 0, 100, NULL, 'Имя33', 'Комментарий33', '/var/mail/user33', 'target33', 123, '/var/attach/user33', 0),
(34, 34, 'user34@domain.com', 0, 100, NULL, 'Имя34', 'Комментарий34', '/var/mail/user34', 'target34', 123, '/var/attach/user34', 0),
(35, 35, 'user35@domain.com', 0, 100, NULL, 'Имя35', 'Комментарий35', '/var/mail/user35', 'target35', 123, '/var/attach/user35', 0),
(36, 36, 'user36@domain.com', 0, 100, NULL, 'Имя36', 'Комментарий36', '/var/mail/user36', 'target36', 123, '/var/attach/user36', 0),
(37, 37, 'user37@domain.com', 0, 100, NULL, 'Имя37', 'Комментарий37', '/var/mail/user37', 'target37', 123, '/var/attach/user37', 0),
(38, 38, 'user38@domain.com', 0, 100, NULL, 'Имя38', 'Комментарий38', '/var/mail/user38', 'target38', 123, '/var/attach/user38', 0),
(39, 39, 'user39@domain.com', 0, 100, NULL, 'Имя39', 'Комментарий39', '/var/mail/user39', 'target39', 123, '/var/attach/user39', 0),
(40, 40, 'user40@domain.com', 0, 100, NULL, 'Имя40', 'Комментарий40', '/var/mail/user40', 'target40', 123, '/var/attach/user40', 0),
(41, 41, 'user41@domain.com', 0, 100, NULL, 'Имя41', 'Комментарий41', '/var/mail/user41', 'target41', 123, '/var/attach/user41', 0),
(42, 42, 'user42@domain.com', 0, 100, NULL, 'Имя42', 'Комментарий42', '/var/mail/user42', 'target42', 123, '/var/attach/user42', 0),
(43, 43, 'user43@domain.com', 0, 100, NULL, 'Имя43', 'Комментарий43', '/var/mail/user43', 'target43', 123, '/var/attach/user43', 0),
(44, 44, 'user44@domain.com', 0, 100, NULL, 'Имя44', 'Комментарий44', '/var/mail/user44', 'target44', 123, '/var/attach/user44', 0),
(45, 45, 'user45@domain.com', 0, 100, NULL, 'Имя45', 'Комментарий45', '/var/mail/user45', 'target45', 123, '/var/attach/user45', 0),
(46, 46, 'user46@domain.com', 0, 100, NULL, 'Имя46', 'Комментарий46', '/var/mail/user46', 'target46', 123, '/var/attach/user46', 0),
(47, 47, 'user47@domain.com', 0, 100, NULL, 'Имя47', 'Комментарий47', '/var/mail/user47', 'target47', 123, '/var/attach/user47', 0),
(48, 48, 'user48@domain.com', 0, 100, NULL, 'Имя48', 'Комментарий48', '/var/mail/user48', 'target48', 123, '/var/attach/user48', 0),
(49, 49, 'user49@domain.com', 0, 100, NULL, 'Имя49', 'Комментарий49', '/var/mail/user49', 'target49', 123, '/var/attach/user49', 0),
(50, 50, 'user50@domain.com', 0, 100, NULL, 'Имя50', 'Комментарий50', '/var/mail/user50', 'target50', 123, '/var/attach/user50', 0),
(51, 51, 'user51@domain.com', 0, 100, NULL, 'Имя51', 'Комментарий51', '/var/mail/user51', 'target51', 123, '/var/attach/user51', 0),
(52, 52, 'user52@domain.com', 0, 100, NULL, 'Имя52', 'Комментарий52', '/var/mail/user52', 'target52', 123, '/var/attach/user52', 0),
(53, 53, 'user53@domain.com', 0, 100, NULL, 'Имя53', 'Комментарий53', '/var/mail/user53', 'target53', 123, '/var/attach/user53', 0),
(54, 54, 'user54@domain.com', 0, 100, NULL, 'Имя54', 'Комментарий54', '/var/mail/user54', 'target54', 123, '/var/attach/user54', 0),
(55, 55, 'user55@domain.com', 0, 100, NULL, 'Имя55', 'Комментарий55', '/var/mail/user55', 'target55', 123, '/var/attach/user55', 0),
(56, 56, 'user56@domain.com', 0, 100, NULL, 'Имя56', 'Комментарий56', '/var/mail/user56', 'target56', 123, '/var/attach/user56', 0),
(57, 57, 'user57@domain.com', 0, 100, NULL, 'Имя57', 'Комментарий57', '/var/mail/user57', 'target57', 123, '/var/attach/user57', 0),
(58, 58, 'user58@domain.com', 0, 100, NULL, 'Имя58', 'Комментарий58', '/var/mail/user58', 'target58', 123, '/var/attach/user58', 0),
(59, 59, 'user59@domain.com', 0, 100, NULL, 'Имя59', 'Комментарий59', '/var/mail/user59', 'target59', 123, '/var/attach/user59', 0),
(60, 60, 'user60@domain.com', 0, 100, NULL, 'Имя60', 'Комментарий60', '/var/mail/user60', 'target60', 123, '/var/attach/user60', 0),
(61, 61, 'user61@domain.com', 0, 100, NULL, 'Имя61', 'Комментарий61', '/var/mail/user61', 'target61', 123, '/var/attach/user61', 0),
(62, 62, 'user62@domain.com', 0, 100, NULL, 'Имя62', 'Комментарий62', '/var/mail/user62', 'target62', 123, '/var/attach/user62', 0),
(63, 63, 'user63@domain.com', 0, 100, NULL, 'Имя63', 'Комментарий63', '/var/mail/user63', 'target63', 123, '/var/attach/user63', 0),
(64, 64, 'user64@domain.com', 0, 100, NULL, 'Имя64', 'Комментарий64', '/var/mail/user64', 'target64', 123, '/var/attach/user64', 0),
(65, 65, 'user65@domain.com', 0, 100, NULL, 'Имя65', 'Комментарий65', '/var/mail/user65', 'target65', 123, '/var/attach/user65', 0),
(66, 66, 'user66@domain.com', 0, 100, NULL, 'Имя66', 'Комментарий66', '/var/mail/user66', 'target66', 123, '/var/attach/user66', 0),
(67, 67, 'user67@domain.com', 0, 100, NULL, 'Имя67', 'Комментарий67', '/var/mail/user67', 'target67', 123, '/var/attach/user67', 0),
(68, 68, 'user68@domain.com', 0, 100, NULL, 'Имя68', 'Комментарий68', '/var/mail/user68', 'target68', 123, '/var/attach/user68', 0),
(69, 69, 'user69@domain.com', 0, 100, NULL, 'Имя69', 'Комментарий69', '/var/mail/user69', 'target69', 123, '/var/attach/user69', 0),
(70, 70, 'user70@domain.com', 0, 100, NULL, 'Имя70', 'Комментарий70', '/var/mail/user70', 'target70', 123, '/var/attach/user70', 0),
(71, 71, 'user71@domain.com', 0, 100, NULL, 'Имя71', 'Комментарий71', '/var/mail/user71', 'target71', 123, '/var/attach/user71', 0),
(72, 72, 'user72@domain.com', 0, 100, NULL, 'Имя72', 'Комментарий72', '/var/mail/user72', 'target72', 123, '/var/attach/user72', 0),
(73, 73, 'user73@domain.com', 0, 100, NULL, 'Имя73', 'Комментарий73', '/var/mail/user73', 'target73', 123, '/var/attach/user73', 0),
(74, 74, 'user74@domain.com', 0, 100, NULL, 'Имя74', 'Комментарий74', '/var/mail/user74', 'target74', 123, '/var/attach/user74', 0),
(75, 75, 'user75@domain.com', 0, 100, NULL, 'Имя75', 'Комментарий75', '/var/mail/user75', 'target75', 123, '/var/attach/user75', 0),
(76, 76, 'user76@domain.com', 0, 100, NULL, 'Имя76', 'Комментарий76', '/var/mail/user76', 'target76', 123, '/var/attach/user76', 0),
(77, 77, 'user77@domain.com', 0, 100, NULL, 'Имя77', 'Комментарий77', '/var/mail/user77', 'target77', 123, '/var/attach/user77', 0),
(78, 78, 'user78@domain.com', 0, 100, NULL, 'Имя78', 'Комментарий78', '/var/mail/user78', 'target78', 123, '/var/attach/user78', 0),
(79, 79, 'user79@domain.com', 0, 100, NULL, 'Имя79', 'Комментарий79', '/var/mail/user79', 'target79', 123, '/var/attach/user79', 0),
(80, 80, 'user80@domain.com', 0, 100, NULL, 'Имя80', 'Комментарий80', '/var/mail/user80', 'target80', 123, '/var/attach/user80', 0),
(81, 81, 'user81@domain.com', 0, 100, NULL, 'Имя81', 'Комментарий81', '/var/mail/user81', 'target81', 123, '/var/attach/user81', 0),
(82, 82, 'user82@domain.com', 0, 100, NULL, 'Имя82', 'Комментарий82', '/var/mail/user82', 'target82', 123, '/var/attach/user82', 0),
(83, 83, 'user83@domain.com', 0, 100, NULL, 'Имя83', 'Комментарий83', '/var/mail/user83', 'target83', 123, '/var/attach/user83', 0),
(84, 84, 'user84@domain.com', 0, 100, NULL, 'Имя84', 'Комментарий84', '/var/mail/user84', 'target84', 123, '/var/attach/user84', 0),
(85, 85, 'user85@domain.com', 0, 100, NULL, 'Имя85', 'Комментарий85', '/var/mail/user85', 'target85', 123, '/var/attach/user85', 0),
(86, 86, 'user86@domain.com', 0, 100, NULL, 'Имя86', 'Комментарий86', '/var/mail/user86', 'target86', 123, '/var/attach/user86', 0),
(87, 87, 'user87@domain.com', 0, 100, NULL, 'Имя87', 'Комментарий87', '/var/mail/user87', 'target87', 123, '/var/attach/user87', 0),
(88, 88, 'user88@domain.com', 0, 100, NULL, 'Имя88', 'Комментарий88', '/var/mail/user88', 'target88', 123, '/var/attach/user88', 0),
(89, 89, 'user89@domain.com', 0, 100, NULL, 'Имя89', 'Комментарий89', '/var/mail/user89', 'target89', 123, '/var/attach/user89', 0),
(90, 90, 'user90@domain.com', 0, 100, NULL, 'Имя90', 'Комментарий90', '/var/mail/user90', 'target90', 123, '/var/attach/user90', 0),
(91, 91, 'user91@domain.com', 0, 100, NULL, 'Имя91', 'Комментарий91', '/var/mail/user91', 'target91', 123, '/var/attach/user91', 0),
(92, 92, 'user92@domain.com', 0, 100, NULL, 'Имя92', 'Комментарий92', '/var/mail/user92', 'target92', 123, '/var/attach/user92', 0),
(93, 93, 'user93@domain.com', 0, 100, NULL, 'Имя93', 'Комментарий93', '/var/mail/user93', 'target93', 123, '/var/attach/user93', 0),
(94, 94, 'user94@domain.com', 0, 100, NULL, 'Имя94', 'Комментарий94', '/var/mail/user94', 'target94', 123, '/var/attach/user94', 0),
(95, 95, 'user95@domain.com', 0, 100, NULL, 'Имя95', 'Комментарий95', '/var/mail/user95', 'target95', 123, '/var/attach/user95', 0),
(96, 96, 'user96@domain.com', 0, 100, NULL, 'Имя96', 'Комментарий96', '/var/mail/user96', 'target96', 123, '/var/attach/user96', 0),
(97, 97, 'user97@domain.com', 0, 100, NULL, 'Имя97', 'Комментарий97', '/var/mail/user97', 'target97', 123, '/var/attach/user97', 0),
(98, 98, 'user98@domain.com', 0, 100, NULL, 'Имя98', 'Комментарий98', '/var/mail/user98', 'target98', 123, '/var/attach/user98', 0),
(99, 99, 'user99@domain.com', 0, 100, NULL, 'Имя99', 'Комментарий99', '/var/mail/user99', 'target99', 123, '/var/attach/user99', 0),
(100, 100, 'user100@domain.com', 0, 100, NULL, 'Имя100', 'Комментарий100', '/var/mail/user100', 'target100', 123, '/var/attach/user100', 0);

CREATE TABLE `security_image` (
  `id` int NOT NULL AUTO_INCREMENT,
  `word` varchar(255) DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  `ip` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `tm` (`time`) USING BTREE,
  KEY `__idx_word` (`word`)
) ENGINE=MEMORY DEFAULT CHARSET=cp1251;

INSERT INTO `security_image` (word, time, ip) VALUES
('word1', '2023-01-01 00:00:00', 123456789),
('word2', '2023-01-02 00:00:00', 987654321),
('word3', '2023-01-03 00:00:00', 111111111),
('word4', '2023-01-04 00:00:00', 222222222),
('word5', '2023-01-05 00:00:00', 333333333),
('word6', '2023-01-06 00:00:00', 444444444),
('word7', '2023-01-07 00:00:00', 555555555),
('word8', '2023-01-08 00:00:00', 666666666),
('word9', '2023-01-09 00:00:00', 777777777),
('word10', '2023-01-10 00:00:00', 888888888);
