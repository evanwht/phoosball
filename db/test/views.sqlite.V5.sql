DROP VIEW IF EXISTS last_games;
DROP VIEW IF EXISTS wins;
DROP VIEW IF EXISTS losses;
DROP VIEW IF EXISTS overall_standings;

CREATE VIEW last_games AS 
SELECT
    `g`.`id` AS `id`,
    `g`.`game_date` AS `game_date`,
    `p1`.`name` AS `team_1_Defense`,
    `p2`.`name` AS `team_1_Offense`,
    `p3`.`name` AS `team_2_Defense`,
    `p4`.`name` AS `team_2_Offense`,
    `g`.`team_1_half` AS `team_1_half`,
    `g`.`team_2_half` AS `team_2_half`,
    `g`.`team_1_final` AS `team_1_final`,
    `g`.`team_2_final` AS `team_2_final`
FROM
    (
        (
            (
                (
                    `games` `g`
                    LEFT JOIN `players` `p1`
                    ON `p1`.`id` = `g`.`team_1_p1`
                )
                LEFT JOIN `players` `p2`
                ON `p2`.`id` = `g`.`team_1_p2`
            )
            LEFT JOIN `players` `p3`
            ON `p3`.`id` = `g`.`team_2_p1`
        )
        LEFT JOIN `players` `p4`
        ON `p4`.`id` = `g`.`team_2_p2`
    )
ORDER BY `g`.`game_date` DESC;

CREATE VIEW wins AS
SELECT
    `p`.`id` AS `id`,
    `p`.`name` AS `name`,
    COUNT(0) AS `wins`
FROM
    `players` `p`
    JOIN
    (
    SELECT
        `g`.`id` AS `id`,
        `g`.`team_1_p1` AS `w1`,
        `g`.`team_1_p2` AS `w2`
    FROM `games` `g`
    ) `w`
    ON `p`.`id` = `w`.`w1` OR `p`.`id` = `w`.`w2`
WHERE `p`.`name` IS NOT NULL
GROUP BY `p`.`id`, `p`.`name`;

CREATE VIEW losses AS
SELECT
    `p`.`id` AS `id`,
    `p`.`name` AS `name`,
    COUNT(0) AS `losses`
FROM
    `players` `p`
    JOIN
    (
    SELECT
        `g`.`id` AS `id`,
        `g`.`team_2_p1` AS `l1`,
        `g`.`team_2_p2` AS `l2`
    FROM `games` `g`
    ) `l`
    ON `p`.`id` = `l`.`l1` OR `p`.`id` = `l`.`l2`
WHERE `p`.`name` IS NOT NULL
GROUP BY `p`.`id`, `p`.`name`;

CREATE VIEW overall_standings AS
SELECT
    `w`.`name` AS `name`,
    COALESCE(`w`.`wins`, 0) AS `wins`,
    COALESCE(`l`.`losses`, 0) AS `losses`
FROM
    `wins` `w`
    LEFT JOIN `losses` `l`
    ON `w`.`id` = `l`.`id`
UNION
SELECT
    `l`.`name` AS `name`,
    COALESCE(`w`.`wins`, 0) AS `wins`,
    COALESCE(`l`.`losses`, 0) AS `losses`
FROM
    `losses` `l`
    LEFT JOIN `wins` `w`
    ON `l`.`id` = `w`.`id`
ORDER BY `wins` DESC, `losses` DESC;