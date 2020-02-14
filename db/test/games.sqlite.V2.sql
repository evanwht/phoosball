PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS `games` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `game_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `input_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `team_1_p1` INT NOT NULL,
  `team_1_p2` INT NOT NULL,
  `team_2_p1` INT NOT NULL,
  `team_2_p2` INT NOT NULL,
  `team_1_half` INT NOT NULL,
  `team_1_final` INT NOT NULL,
  `team_2_half` INT NOT NULL,
  `team_2_final` INT NOT NULL,
  `input_by` INT,
  FOREIGN KEY(team_1_p1) REFERENCES players(id),
  FOREIGN KEY(team_1_p2) REFERENCES players(id),
  FOREIGN KEY(team_2_p1) REFERENCES players(id),
  FOREIGN KEY(team_2_p2) REFERENCES players(id),
  FOREIGN KEY(input_by) REFERENCES players(id)
);

INSERT INTO games (id, game_date, input_date, team_1_p1, team_1_p2, team_2_p1, team_2_p2, team_1_half, team_1_final, team_2_half, team_2_final, input_by) VALUES
(8, '2020-01-21 08:32:14', '2020-01-21 14:32:14', 2, 1, 5, 3, 5, 12, 3, 10, NULL),
(12, '2020-01-22 12:22:28', '2020-01-22 18:22:28', 4, 2, 6, 1, 5, 10, 2, 5, NULL),
(14, '2020-01-22 15:17:48', '2020-01-22 21:17:48', 2, 1, 4, 3, 5, 10, 3, 8, NULL),
(15, '2020-01-22 15:26:43', '2020-01-22 21:26:43', 1, 2, 6, 4, 5, 10, 3, 7, NULL),
(17, '2020-01-21 12:15:07', '2020-01-21 18:15:07', 4, 1, 3, 2, 3, 10, 5, 8, NULL),
(18, '2020-01-21 12:29:05', '2020-01-21 18:29:05', 2, 3, 6, 1, 5, 15, 4, 14, NULL),
(19, '2020-01-22 12:10:50', '2020-01-22 18:10:50', 6, 1, 5, 2, 5, 10, 3, 8, NULL);
