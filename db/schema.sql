CREATE TABLE goal ( 
    id   INTEGER NOT NULL,
    name VARCHAR(20) NOT NULL, 
    date VARCHAR(10) NOT NULL, 
    note VARCHAR(50), 
    PRIMARY KEY (id) 
);

CREATE TABLE progress ( 
    id      INTEGER NOT NULL, 
    goal_id INTEGER NOT NULL, 
    value   INTEGER NOT NULL, 
    date    VARCHAR(10) NOT NULL, 
    note    VARCHAR(50), 
    CHECK (value >=0 AND value <=100), 
    PRIMARY KEY (id), 
    FOREIGN KEY (goal_id) REFERENCES goal(id) 
);