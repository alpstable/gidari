SELECT * FROM candle_minutes t WHERE t.product_id = $1 AND t.unix >= $2 AND t.unix <= $3
