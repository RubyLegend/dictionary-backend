BEGIN;

Create table Users (
userID int not null Auto_Increment,
email varchar(45) not null,
username varchar(45) not null,
password varchar(45) not null,
createdAt DateTime not null,
primary key (userID)
);

Create table Dictionaries (
dictionaryID int not null Auto_Increment,
userID int not null,
name varchar(45) not null,
createdAt DateTime not null,
Primary key (dictionaryID),
Foreign key(userID) references Users(userID)
);

Create table History (
historyID int not null Auto_Increment,
userID int not null,
wordID int not null,
isCorrect boolean not null,
createdAt DateTime not null,
Primary key (historyID),
Foreign key(userID) references Users(userID)
);

Create table Words (
wordID int not null Auto_Increment,
name varchar(45) not null,
createdAt DateTime not null,
Primary key (wordID)
);

Create table DictionariesWords (
dictionaryID int not null Auto_Increment,
wordID int not null,
Primary key(dictionaryID),
Foreign key(dictionaryID) references Dictionaries(dictionaryID),
Foreign key(wordID) references Words(wordID)
);

Create table Translation (
translationID int not null Auto_Increment,
wordID int not null,
name varchar(45) not null,
language varchar(45) not null,
primary key(translationID),
Foreign key(wordID) references Words(wordID)
);

COMMIT;