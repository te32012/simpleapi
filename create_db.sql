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


with a as (
        INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('anna', 'female', '03-05-1998') RETURNING id
), f as (
    INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino1', 'good film', '03-05-1959', 2) RETURNING id
)
INSERT into actors_films(id_films, id_actors) VALUES ((select id from a), (select id from f));

with a as (
        INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('tanya', 'female', '03-05-1999') RETURNING id
), f as (
    INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino2', 'normal film', '05-05-1969', 5) RETURNING id
)
INSERT into actors_films(id_films, id_actors) VALUES ((select id from a), (select id from f));

with a as (
        INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('ivan', 'male', '06-09-1995') RETURNING id
), f as (
    INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino3', 'bad film', '03-05-1979', 9) RETURNING id
)
INSERT into actors_films(id_films, id_actors) VALUES ((select id from a), (select id from f));

INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('petya', 'male', '03-05-1978') RETURNING id;
INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino7', 'worse film that i watch', '03-05-2000', 7) RETURNING id;
with a as (
        select * from actors where nameActor='petya' and sex='male' and dataofbirthday='03-05-1978'
), f as (
    select * from films where nameOfFilm='kino3' and about='bad film' and releaseDate='03-05-1979' and rating=9
)
INSERT into actors_films(id_films, id_actors) VALUES ((select id from a), (select id from f));
with a as (
        select * from actors where nameActor='petya' and sex='male' and dataofbirthday='03-05-1978'
), f as (
    select * from films where nameOfFilm='kino2' and about='normal film' and releaseDate='05-05-1969' and rating=5
)
INSERT into actors_films(id_films, id_actors) VALUES ((select id from a), (select id from f));
with a as (
        select * from actors where nameActor='tanya' and sex='female' and dataofbirthday='03-05-1999'
), f as (
    select * from films where nameOfFilm='kino1' and about='good film' and releaseDate='03-05-1959' and rating=2
)
INSERT into actors_films(id_films, id_actors) VALUES ((select id from a), (select id from f));


with a as (
        INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('anna', 'female', '03-05-1998') RETURNING id
), f as (
    INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino9', 'good film', '03-05-1959', 2) RETURNING id
)
INSERT into actors_films(id_films, id_actors) VALUES ((select id from a), (select id from f));


-- DECLARE i4 integer := INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('ivan', 'male','02-08-2005') RETURNING id;
-- DECLARE i5 integer := INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('vasya', 'male', '03-05-1989') RETURNING id;
-- DECLARE i6 integer := INSERT INTO actors(nameActor, sex, dataofbirthday) VALUES ('petya', 'male', '03-05-2003') RETURNING id;

-- DECLARE i7 integer := INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino1', 'good film', '03-05-1959', 2)RETURNING id; 
-- DECLARE i8 integer := INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino3', 'fine film', '03-12-2139', 2) RETURNING id;
-- DECLARE i9 integer := INSERT INTO films(nameOfFilm, about, releaseDate, rating) VALUES ('kino7', 'bad film', '03-03-1998', 2) RETURNING id;
