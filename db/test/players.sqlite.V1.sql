CREATE TABLE IF NOT EXISTS `players` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` CHAR(64) NOT NULL,
  `display_name` CHAR(16),
  `email` CHAR(128),
  `favorite_shot` CHAR(64)
);

INSERT INTO `players` (`id`, `name`, `display_name`, `email`, `favorite_shot`) VALUES
(1, 'Evan White', 'king', 'evanwht1@gmail.com', 'pull'),
(2, 'Thomas Mckenna', 'number 2', '', 'push pass'),
(3, 'Zach Volz', 'zachonius', '', 'pull spin pass'),
(4, 'Manny Sahagun', 'still learning', '', 'bullshit'),
(5, 'Artur Jaglowski', 'friend', '', 'push-back'),
(6, 'Rich Kraft', 'straighter', '', 'angle');