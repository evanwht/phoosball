CREATE TABLE event_types (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name CHAR(64) NOT NULL
);

INSERT INTO event_types (id, name) VALUES
(1, 'bread'),
(2, 'reverse bread'),
(3, 'toast'),
(4, 'reverse toast'),
(5, 'five man bread'),
(6, 'own goal'),
(7, 'full court bread'),
(8, 'five man goal'),
(9, 'double tap'),
(10, 'front line goal'),
(11, 'defense goal'),
(13, 'goalie goal');