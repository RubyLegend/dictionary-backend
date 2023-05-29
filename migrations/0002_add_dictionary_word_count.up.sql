START TRANSACTION;

alter table Dictionaries add column total int not null;

create trigger dictionary_words_sum_insert after insert on DictionariesWords for each row
update Dictionaries set total = (
    select count(*) 
    from DictionariesWords
    where DictionariesWords.dictionaryID = Dictionaries.dictionaryID
);
create trigger dictionary_words_sum_update after update on DictionariesWords for each row
update Dictionaries set total = (
    select count(*) 
    from DictionariesWords
    where DictionariesWords.dictionaryID = Dictionaries.dictionaryID
);
create trigger dictionary_words_sum_delete after delete on DictionariesWords for each row
update Dictionaries set total = (
    select count(*) 
    from DictionariesWords
    where DictionariesWords.dictionaryID = Dictionaries.dictionaryID
);

COMMIT;