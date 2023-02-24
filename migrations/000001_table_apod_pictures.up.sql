CREATE TABLE pictures
(
    id            UUID PRIMARY KEY,
    date          DATE                        NOT NULL UNIQUE,
    title         VARCHAR(255)                NOT NULL,
    url           VARCHAR(255)                NOT NULL,
    hd_url        VARCHAR(255),
    local_url     VARCHAR(255)                NOT NULL,
    thumbnail_url VARCHAR(255),
    media_type    VARCHAR(16)                 NOT NULL,
    copyright     VARCHAR(255),
    explanation   VARCHAR(2048)               NOT NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX pictures_date ON pictures (date);
