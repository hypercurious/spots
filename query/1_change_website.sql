UPDATE "MY_TABLE"
SET website = SUBSTRING(
    website FROM '(?:.*://)?(?:www\.)?([^/]*)'
) -- (https://)(www.)domain.com(/index.php)