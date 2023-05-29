BEGIN;

alter table Dictionaries drop column total;

drop trigger dictionary_words_sum_insert;
drop trigger dictionary_words_sum_update;
drop trigger dictionary_words_sum_delete;

COMMIT;