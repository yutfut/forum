create table data (
    "id"        SERIAL PRIMARY KEY,
    "name"      text,
    "var1"      int[],
    "var2"      int[],
    "var3"      int[]
);

insert into data (id, name, var1, var2, var3) values
('1', 'name_1', '{1,2,3}', '{1,2,3}', '{1,2,3}')