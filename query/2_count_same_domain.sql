SELECT website, COUNT(*) as spots
FROM "MY_TABLE"
GROUP BY website
HAVING spots>1