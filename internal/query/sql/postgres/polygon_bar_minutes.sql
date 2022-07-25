SELECT * FROM bar_minutes t WHERE t.ticker = $1 AND t.adjusted = $2 AND t.t >= $3 AND t.t <= $4
