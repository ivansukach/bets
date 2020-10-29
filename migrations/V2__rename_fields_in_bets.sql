ALTER TABLE bets
    ADD COLUMN Number_of_bets INTEGER NOT NULL;
CREATE TABLE blocked_bets(
    Id                  INTEGER PRIMARY KEY,
    Name                VARCHAR(30) NOT NULL,
    Sum_amount          FLOAT       NOT NULL,
    Average_coefficient FLOAT       NOT NULL,
    Number_of_bets      INTEGER     NOT NULL
);