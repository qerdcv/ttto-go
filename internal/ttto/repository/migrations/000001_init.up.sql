CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(30) NOT NULL UNIQUE,
    password TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS games
(
    id                SERIAL PRIMARY KEY,
    owner_id          INTEGER     NOT NULL REFERENCES users (id),
    opponent_id       INTEGER REFERENCES users (id),
    current_player_id INTEGER REFERENCES users (id),
    step_count        INTEGER     NOT NULL DEFAULT 0,
    winner_id         INTEGER REFERENCES users (id),
    field             json        NOT NULL DEFAULT '[
      [
        "",
        "",
        ""
      ],
      [
        "",
        "",
        ""
      ],
      [
        "",
        "",
        ""
      ]
    ]',
    current_state     VARCHAR(10) NOT NULL DEFAULT 'pending'
);

CREATE OR REPLACE VIEW games_with_usernames AS
SELECT g.*,
       owner.username          as owner_name,
       opponent.username       as opponent_name,
       current_player.username as current_player_name,
       winner.username         as winner_name
FROM games g
         LEFT JOIN users owner on owner.id = g.owner_id
         LEFT JOIN users opponent on opponent.id = g.opponent_id
         LEFT JOIN users current_player on current_player.id = g.current_player_id
         LEFT JOIN users winner on winner.id = g.winner_id;

CREATE TABLE IF NOT EXISTS games_history
(
    id                SERIAL PRIMARY KEY,
    game_id           INTEGER REFERENCES games (id),
    owner_id          INTEGER     NOT NULL REFERENCES users (id),
    opponent_id       INTEGER REFERENCES users (id),
    current_player_id INTEGER REFERENCES users (id),
    step_count        INTEGER     NOT NULL DEFAULT 0,
    winner_id         INTEGER REFERENCES users (id),
    field             json        NOT NULL DEFAULT '[
      [
        "",
        "",
        ""
      ],
      [
        "",
        "",
        ""
      ],
      [
        "",
        "",
        ""
      ]
    ]',
    current_state     VARCHAR(10) NOT NULL DEFAULT 'pending'
);
