-- phpMyAdmin SQL Dump
-- version 4.5.4.1deb2ubuntu2.1
-- http://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: Nov 05, 2019 at 03:00 PM
-- Server version: 5.7.27-0ubuntu0.16.04.1
-- PHP Version: 7.0.33-8+ubuntu16.04.1+deb.sury.org+1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `evidentor`
--

-- --------------------------------------------------------

--
-- Table structure for table `activities`
--

CREATE TABLE `activities` (
  `id` int(11) NOT NULL,
  `url` varchar(255) NOT NULL,
  `method` varchar(6) NOT NULL,
  `url_regex` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `casbin_rule`
--

CREATE TABLE `casbin_rule` (
  `p_type` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `casbin_rule`
--

INSERT INTO `casbin_rule` (`p_type`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES
('p', 'admin', '/*', '*', NULL, NULL, NULL),
('p', 'member', '/logout', '*', NULL, NULL, NULL),
('p', 'member', '/member/*', '*', NULL, NULL, NULL),
('p', 'anonymous', '/users', '*', NULL, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `item`
--

CREATE TABLE `item` (
  `id` int(11) NOT NULL,
  `data` text,
  `type` varchar(32) NOT NULL,
  `status` enum('active','inactive') NOT NULL DEFAULT 'active',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `uuid` varchar(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `item`
--

INSERT INTO `item` (`id`, `data`, `type`, `status`, `created_at`, `updated_at`, `uuid`) VALUES
(1, NULL, 'child', 'active', '2019-11-05 00:00:00', '2019-11-05 00:00:00', ''),
(2, NULL, 'child', 'active', '2019-11-04 09:20:39', '2019-11-04 13:41:50', ''),
(3, NULL, 'parent', 'active', '2019-11-05 00:00:00', '2019-11-05 00:00:00', '');

-- --------------------------------------------------------

--
-- Table structure for table `item_access`
--

CREATE TABLE `item_access` (
  `id` int(11) NOT NULL,
  `item_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `item_property`
--

CREATE TABLE `item_property` (
  `id` int(11) NOT NULL,
  `item_id` int(11) NOT NULL,
  `first_name` varchar(35) NOT NULL,
  `last_name` varchar(35) NOT NULL,
  `jmbg` varchar(15) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `item_type`
--

CREATE TABLE `item_type` (
  `type` varchar(32) NOT NULL,
  `name` varchar(255) NOT NULL,
  `data` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `item_type`
--

INSERT INTO `item_type` (`type`, `name`, `data`) VALUES
('child', 'Child', ''),
('parent', 'Parent', ''),
('teacher', 'Teacher', ''),
('cleaner', 'Cleaner', '');

-- --------------------------------------------------------

--
-- Table structure for table `role`
--

CREATE TABLE `role` (
  `id` int(10) NOT NULL,
  `name` varchar(30) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `role`
--

INSERT INTO `role` (`id`, `name`, `created_at`, `updated_at`) VALUES
(1, 'admin', '2019-09-20 00:00:00', '2019-09-20 00:00:00'),
(2, 'member', '2019-09-20 00:00:00', '2019-09-20 00:00:00');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(10) UNSIGNED NOT NULL,
  `first_name` varchar(30) DEFAULT NULL,
  `last_name` varchar(40) DEFAULT NULL,
  `gender` enum('male','female','other','unknown') NOT NULL DEFAULT 'unknown',
  `birthday` date DEFAULT NULL,
  `email` varchar(150) NOT NULL DEFAULT '',
  `address` varchar(255) DEFAULT NULL,
  `optin` enum('yes','no','unknown') NOT NULL DEFAULT 'unknown',
  `password` char(60) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `first_name`, `last_name`, `gender`, `birthday`, `email`, `address`, `optin`, `password`) VALUES
(1, 'Boris', 'Kundacina', 'male', NULL, 'boris@horisen.com', NULL, 'yes', ''),
(2, 'Iva', 'Stankovic', 'female', NULL, 'iva@horisen.com', 'Viz. Bul. 78', 'yes', '$2a$04$LRUCcaVxnHyZHnMezk1T.OpDOMjgwnv7GfUDnefCT/2.tIbhlhVCO'),
(3, 'Mile', 'Kitic', 'male', NULL, 'mile@gmail.com', NULL, 'unknown', '$2a$04$krYFRHZEK8EGoqJ.gusqwuCLh5FTVq5Nxg4.Op7y1wuEUlMeT90Dq');

-- --------------------------------------------------------

--
-- Table structure for table `user_roles`
--

CREATE TABLE `user_roles` (
  `id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `user_roles`
--

INSERT INTO `user_roles` (`id`, `role_id`, `user_id`) VALUES
(1, 1, 1);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `activities`
--
ALTER TABLE `activities`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `item`
--
ALTER TABLE `item`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `item_access`
--
ALTER TABLE `item_access`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `item_property`
--
ALTER TABLE `item_property`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `role`
--
ALTER TABLE `role`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`);

--
-- Indexes for table `user_roles`
--
ALTER TABLE `user_roles`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `activities`
--
ALTER TABLE `activities`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `item`
--
ALTER TABLE `item`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
--
-- AUTO_INCREMENT for table `item_access`
--
ALTER TABLE `item_access`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `item_property`
--
ALTER TABLE `item_property`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `role`
--
ALTER TABLE `role`
  MODIFY `id` int(10) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
--
-- AUTO_INCREMENT for table `user_roles`
--
ALTER TABLE `user_roles`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
