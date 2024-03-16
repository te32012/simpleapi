CREATE TABLE IF NOT EXISTS users 
(
    id         SERIAL PRIMARY KEY,
    username   TEXT NOT NULL,
    hpassword  TEXT NOT NULL,
    permission smallint NOT NULL
);

CREATE TABLE IF NOT EXISTS films 
(
    id SERIAL PRIMARY KEY,
    nameOfFilm TEXT,
    about TEXT,
    releaseDate timestamp,
    rating smallint
);

CREATE TABLE IF NOT EXISTS actors 
(
    id SERIAL PRIMARY KEY,
    nameActor TEXT,
    sex TEXT,
    dataofbirthday timestamp
);

CREATE TABLE IF NOT EXISTS actors_films 
(
    id_films integer references films(id),
    id_actors integer references actors(id)
);

CREATE TABLE IF NOT EXISTS fulltextsearch 
(
    id_films integer references films(id),
    keyworld tsvector
);

CREATE INDEX IF NOT EXISTS indexFullTextSearch ON fulltextsearch
  USING gin(keyworld);

CREATE OR REPLACE FUNCTION public.make_tsvector(
	id_f integer)
    RETURNS tsvector
    LANGUAGE 'plpgsql'
    COST 100
    IMMUTABLE PARALLEL UNSAFE
AS $BODY$
   	DECLARE nameactors TEXT[];
	DECLARE namefilms TEXT[];
	DECLARE alltext TEXT[];
	DECLARE id_actors integer[];
	DECLARE i integer;
BEGIN
	namefilms := (select ARRAY[films.nameoffilm] from films where id_f = films.id);
	id_actors := ARRAY(select actors_films.id_actors from actors_films where actors_films.id_films = id_f);
	i := 1;
	RAISE NOTICE 'cmd: %', i;
	RAISE NOTIcE 'cmd1: %', array_length(id_actors, 1);

	<<e>>
	LOOP 
		EXIT e WHEN array_length(id_actors, 1) < i ;
		nameactors := (nameactors || (select ARRAY[actors.nameactor] from actors where actors.id = id_actors[i]));
		i := i + 1;
		RAISE NOTICE 'cmd: %', i;
		RAISE NOTIcE 'cmd1: %', array_length(id_actors, 1);
	END LOOP e;
	alltext := nameactors || namefilms;
  RETURN array_to_tsvector(alltext);
END
$BODY$;

ALTER FUNCTION public.make_tsvector(integer)
    OWNER TO root;

-- INSERT INTO users(username, hpassword, permission) VALUES ('admin', 'admin', 1) RETURNING id;
-- DECLARE i2 integer := INSERT INTO users(username, hpassword, permission) VALUES ('user', 'user', 2) RETURNING id;

-- DECLARE i3 integer := INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('anna', 'female', '03-05-1998') RETURNING id;
-- DECLARE i4 integer := INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('ivan', 'male','02-08-2005') RETURNING id;
-- DECLARE i5 integer := INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('vasya', 'male', '03-05-1989') RETURNING id;
-- DECLARE i6 integer := INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('petya', 'male', '03-05-2003') RETURNING id;

-- DECLARE i7 integer := INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino1', 'good film', '03-05-1959', 2)RETURNING id; 
-- DECLARE i8 integer := INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino3', 'fine film', '03-12-2139', 2) RETURNING id;
-- DECLARE i9 integer := INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino7', 'bad film', '03-03-1998', 2) RETURNING id;
