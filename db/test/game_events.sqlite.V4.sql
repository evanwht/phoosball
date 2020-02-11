CREATE TABLE `game_events` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `game` INTEGER NOT NULL,
  `order_num` INTEGER NOT NULL,
  `type` INTEGER NOT NULL,
  `player_by` INTEGER NOT NULL,
  `player_on` INTEGER NOT NULL,
  `rating` INTEGER,
  FOREIGN KEY(game) REFERENCES games(id),
  FOREIGN KEY(type) REFERENCES event_types(id),
  FOREIGN KEY(player_by) REFERENCES players(id),
  FOREIGN KEY(player_on) REFERENCES players(id)
);
