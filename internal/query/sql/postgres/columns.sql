SELECT
	c.column_name,
	c.table_name,
	CASE
		WHEN EXISTS(
			SELECT
				1
			FROM
				INFORMATION_SCHEMA.constraint_column_usage k
			WHERE
				c.table_name = k.table_name
				and k.column_name = c.column_name
		) THEN 1
		ELSE 0
	END as primary_key
FROM
	INFORMATION_SCHEMA.COLUMNS c
	INNER JOIN information_schema.tables t ON t.table_name = c.table_name
WHERE
	t.table_type = 'BASE TABLE'
	AND c.table_schema = 'public'
