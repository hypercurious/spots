CREATE OR REPLACE FUNCTION change_website(website in VARCHAR)
RETURNS VARCHAR
LANGUAGE PLPGSQL
AS $$
BEGIN
    RETURN SUBSTRING(website FROM '(?:.*://)?(?:www\.)?([^/]*)');
END;
$$