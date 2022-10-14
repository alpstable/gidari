SELECT c.column_name,
       c.table_name,
       CASE
           WHEN EXISTS
                (
                    SELECT 1
                    FROM information_schema.constraint_column_usage k
                    WHERE c.table_name = k.table_name
                          AND k.column_name = c.column_name
                ) THEN
               1
           ELSE
               0
       END AS primary_key,
       pg_relation_size(quote_ident(c.table_name)) AS bytes
FROM information_schema.columns c
    INNER JOIN information_schema.tables t
        ON t.table_name = c.table_name
WHERE t.table_type = 'BASE TABLE'
      AND c.table_schema = 'public'
