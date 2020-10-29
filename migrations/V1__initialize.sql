CREATE TABLE blocked_users(
    Id INTEGER PRIMARY KEY
);
CREATE TABLE bets(
    Id INTEGER PRIMARY KEY ,
    Name VARCHAR(30),
    Sum_amount FLOAT NOT NULL,
    Average_coefficient FLOAT NOT NULL
);

