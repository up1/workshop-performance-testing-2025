-- Insert 10,000 users
INSERT INTO users (name)
SELECT CONCAT('User_', seq)
FROM (
    SELECT @rownum := @rownum + 1 AS seq
    FROM information_schema.columns, (SELECT @rownum := 0) r
    LIMIT 10000
) AS temp;