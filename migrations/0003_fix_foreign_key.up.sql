BEGIN;

drop table DictionariesWords;

Create table DictionariesWords (
dictionaryID int not null Auto_Increment,
wordID int not null,
Foreign key(dictionaryID) references Dictionaries(dictionaryID),
Foreign key(wordID) references Words(wordID)
);

COMMIT;