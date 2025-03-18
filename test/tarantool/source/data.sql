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

CREATE TABLE `security_image` (
  `id` int NOT NULL AUTO_INCREMENT,
  `word` varchar(255) DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  `ip` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `tm` (`time`) USING BTREE,
  KEY `__idx_word` (`word`)
) ENGINE=MEMORY DEFAULT CHARSET=cp1251;
