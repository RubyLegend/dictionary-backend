create database Dictionary;
use Dictionary;

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
total int not null default 0,
Primary key (dictionaryID),
Foreign key(userID) references Users(userID)
on delete cascade
);

Create table History (
historyID int not null Auto_Increment,
userID int not null,
wordID int not null,
isCorrect boolean not null,
createdAt DateTime not null,
Primary key (historyID),
Foreign key(userID) references Users(userID)
on delete cascade
);

Create table Words (
wordID int not null Auto_Increment,
name varchar(45) not null,
createdAt DateTime not null,
isLearned boolean not null default false,
Primary key (wordID)
);

Create table DictionariesWords (
dictionaryID int not null Auto_Increment,
wordID int not null,
Foreign key(dictionaryID) references Dictionaries(dictionaryID)
on delete cascade,
Foreign key(wordID) references Words(wordID)
on delete cascade
);

Create table Translation (
translationID int not null Auto_Increment,
wordID int not null,
name varchar(45) not null,
language varchar(45) not null,
primary key(translationID),
Foreign key(wordID) references Words(wordID)
on delete cascade
);

create trigger dictionary_words_sum_insert after insert on DictionariesWords
for each row
update Dictionaries set total = (
    select count(*) 
    from DictionariesWords
    where DictionariesWords.dictionaryID = Dictionaries.dictionaryID
);
create trigger dictionary_words_sum_update after update on DictionariesWords
for each row
update Dictionaries set total = (
    select count(*) 
    from DictionariesWords
    where DictionariesWords.dictionaryID = Dictionaries.dictionaryID
);
create trigger dictionary_words_sum_delete after delete on DictionariesWords
for each row
update Dictionaries set total = (
    select count(*) 
    from DictionariesWords
    where DictionariesWords.dictionaryID = Dictionaries.dictionaryID
);