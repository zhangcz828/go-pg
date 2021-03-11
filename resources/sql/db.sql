CREATE TABLE Hero (
    Name varchar(20) PRIMARY KEY,
    Detail text check (length(Detail) > 8),
    AttackPower int,
    DefensePower int,
    Blood int
);

INSERT INTO Hero VALUES ('张无忌', '武功：九阳神功，乾坤大挪移，圣火令神功', 20, 10, 100);
INSERT INTO Hero VALUES ('韦小宝', '武功： 轻功, 神行百变，大擒拿手', 10, 30, 100);


CREATE TABLE Boss (
    Name varchar(20),
    Detail text check (length(Detail) > 4),
    AttackPower int,
    DefensePower int,
    Blood int,
    Level int UNIQUE
);

INSERT INTO Boss VALUES ('东方不败','武功：葵花宝典', 20, 5, 100, 1);
INSERT INTO Boss VALUES ('玄冥二老','武功：玄冥神掌', 15, 10, 100, 2);

CREATE TABLE Session(
    UID int primary key,
    HeroName varchar(20),
    HeroBlood int,
    BossBlood int,
    CurrentLevel int,
    Score int,
    ArchiveDate timestamp default now()
);

INSERT INTO Session VALUES ('4','张无忌', 100, 100, 1, 0, '2021-03-11T18:25:06.1577213+08:00');




