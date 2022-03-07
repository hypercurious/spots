SELECT NAME AS "Spots with duplicate domain" FROM "MY_TABLE"
WHERE website IN (
    SELECT website FROM (
        SELECT website, COUNT(*) AS spots FROM "MY_TABLE" GROUP BY website
    ) AS foo WHERE spots>1
);