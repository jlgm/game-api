CREATE DATABASE game;

\c game

CREATE TABLE player (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    games INT DEFAULT 0,
    score INT DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE friendship (
    player1 UUID NOT NULL,
    player2 UUID NOT NULL,
    PRIMARY KEY (player1, player2)
);

COMMIT;
