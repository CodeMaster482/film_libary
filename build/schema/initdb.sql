CREATE TABLE IF NOT EXISTS film (
    film_id         BIGSERIAL PRIMARY KEY,
    title           TEXT      CHECK(length(title) <= 150) NOT NULL,
    "description"   TEXT      CHECK(length("description") <= 1000),
    release_date    DATE,
    rating          INT       CHECK(rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS actor (
    actor_id    BIGSERIAL   PRIMARY KEY,
    "name"      TEXT        CHECK(length("name") <= 100) NOT NULL,
    sex         VARCHAR(1)  CHECK(sex IN ('M', 'W', 'N')),
    birth_date  DATE
);

CREATE TABLE film_actor (
    film_id     BIGINT REFERENCES film(film_id)   ON DELETE CASCADE,
    actor_id    BIGINT REFERENCES actor(actor_id) ON DELETE CASCADE,
    PRIMARY KEY (film_id, actor_id)
);