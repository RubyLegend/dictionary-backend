BEGIN;

alter table Dictionaries drop foreign key Dictionaries_ibfk_1;
alter table Dictionaries add constraint Dictionaries_ibfk_1 foreign key(userID) references Users(userID) on delete cascade;

alter table History drop foreign key History_ibfk_1;
alter table History add constraint History_ibfk_1 foreign key(userID) references Users(userID) on delete cascade;

alter table DictionariesWords drop foreign key DictionariesWords_ibfk_1;
alter table DictionariesWords add constraint DictionariesWords_ibfk_1 foreign key(dictionaryID) references Dictionaries(dictionaryID) on delete cascade;

alter table DictionariesWords drop foreign key DictionariesWords_ibfk_2;
alter table DictionariesWords add constraint DictionariesWords_ibfk_2 foreign key(wordID) references Words(wordID) on delete cascade;

alter table Translation drop foreign key Translation_ibfk_1;
alter table Translation add constraint Translation_ibfk_1 foreign key(wordID) references Words(wordID) on delete cascade;

COMMIT;