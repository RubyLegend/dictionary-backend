BEGIN;

drop table DictionariesWords;

Create table DictionariesWords (
dictionaryID int not null Auto_Increment,
wordID int not null,
primary key(dictionaryID),
Foreign key(dictionaryID) references Dictionaries(dictionaryID),
Foreign key(wordID) references Words(wordID)
);

COMMIT;