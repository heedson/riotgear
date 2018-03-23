CREATE TABLE IF NOT EXISTS player_lp (
    player_id integer NOT NULL,
    match_id integer NOT NULL,
    lp_gain integer NOT NULL,
    total_lp integer NOT NULL,

    PRIMARY KEY(player_id, match_id)
);
