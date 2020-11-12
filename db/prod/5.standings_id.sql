USE `phoosball`;
CREATE OR REPLACE VIEW `overall_standings` AS (
        select `w`.`id` AS `id`,
            `w`.`name` AS `name`,
            coalesce(`w`.`wins`, 0) AS `wins`,
            coalesce(`l`.`losses`, 0) AS `losses`
        from (
                `phoosball`.`wins` `w`
                left join `phoosball`.`losses` `l` on(`w`.`id` = `l`.`id`)
            )
    )
union
(
    select `l`.`id` AS `id`,
        `l`.`name` AS `name`,
        coalesce(`w`.`wins`, 0) AS `wins`,
        coalesce(`l`.`losses`, 0) AS `losses`
    from (
            `phoosball`.`losses` `l`
            left join `phoosball`.`wins` `w` on(`l`.`id` = `w`.`id`)
        )
)
order by `wins` desc,
    `losses`;