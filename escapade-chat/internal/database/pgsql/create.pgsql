CREATE TABLE Player (
    id SERIAL PRIMARY KEY,
    name varchar(30) NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(30) NOT NULL,
    photo_title varchar(50),
    --FirstSeen   timestamp without time zone NOT NULL,
	--LastSeen    timestamp without time zone NOT NULL,
    best_score  int default 0 CHECK (best_score > -1),
    best_time   int default 0 CHECK (best_time > -1),
    GamesTotal  int default 0 CHECK (GamesTotal > -1),
	SingleTotal int default 0 CHECK (SingleTotal > -1),
	OnlineTotal int default 0 CHECK (OnlineTotal > -1),
	SingleWin   int default 0 CHECK (SingleWin > -1),
	OnlineWin   int default 0 CHECK (OnlineWin > -1),
	MinsFound   int default 0 CHECK (MinsFound > -1)
    
);

CREATE Table Session (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    session_code varchar(30) NOT NULL,
    expiration timestamp without time zone NOT NULL
);

ALTER TABLE Session
ADD CONSTRAINT session_player
   FOREIGN KEY (player_id)
   REFERENCES Player(id)
   ON DELETE CASCADE;

CREATE Table Game (
    id SERIAL PRIMARY KEY,
    player_id   int NOT NULL,
    FieldWidth  int CHECK (FieldWidth > -1),
    FieldHeight int CHECK (FieldHeight > -1),
    MinsTotal   int CHECK (MinsTotal > -1),
    MinsFound   int CHECK (MinsFound > -1),
    Finished bool NOT NULL,
    Exploded bool NOT NULL,
    Date timestamp without time zone NOT NULL,
    FOREIGN KEY (player_id) REFERENCES Player (id)
);


INSERT INTO Player(name, password, email, best_score, best_time) VALUES
    ('tiger', 'Bananas', 'tinan@mail.ru', 1000, 10),
    ('panda', 'apple', 'today@mail.ru', 2323, 20),
    ('catmate', 'juice', 'allday@mail.ru', 10000, 5),
    ('hotdog', 'where', 'three@mail.ru', 88, 1000);

    /*
    id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name varchar(30) NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(30) NOT NULL,
    photo_id int,
    best_score int,
    FOREIGN KEY (photo_id) REFERENCES Photo (id)
    */

INSERT INTO Game(player_id, FieldWidth, FieldHeight,
MinsTotal, MinsFound, Finished, Exploded, Date) VALUES
    (1, 50, 50, 100, 20, true, true, date '2001-09-28'),
    (1, 50, 50, 80, 30, false, false, date '2018-09-27'),
    (1, 50, 50, 70, 70, true, false, date '2018-09-26'),
    (1, 50, 50, 60, 30, true, true, date '2018-09-23'),
    (1, 50, 50, 50, 50, true, false, date '2018-09-24'),
    (1, 50, 50, 40, 30, true, false, date '2018-09-25'),
    (2, 25, 25, 80, 30, false, false, date '2018-08-27'),
    (2, 25, 25, 70, 70, true, false, date '2018-08-26'),
    (2, 25, 25, 60, 30, true, true, date '2018-08-23'),
    (3, 10, 10, 10, 10, true, false, date '2018-10-26'),
    (3, 10, 10, 20, 19, true, true, date '2018-10-23'),
    (3, 10, 10, 30, 30, true, false, date '2018-10-24'),
    (3, 10, 10, 40, 5, true, false, date '2018-10-25');

    /*
CREATE Table Game (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    FieldWidth int CHECK (FieldWidth > -1),
    FieldHeight int CHECK (FieldHeight > -1),
    MinsTotal int CHECK (MinsTotal > -1),
    MinsFound int CHECK (MinsFound > -1),
    Finished bool NOT NULL,
    Exploded bool NOT NULL,
    Date timestamp without time zone NOT NULL,
    FOREIGN KEY (player_id) REFERENCES Player (id)
);
    */

--GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO escapade;