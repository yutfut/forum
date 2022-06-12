create extension if not exists citext;

drop table if exists "user"         cascade;
drop table if exists "forum"        cascade;
drop table if exists "thread"       cascade;
drop table if exists "post"         cascade;
drop table if exists "vote"         cascade;
drop table if exists "forum_user"   cascade;

drop function if exists thread_vote();
drop function if exists thread_vote_update();
drop function if exists create_post();
drop function if exists create_thread();

drop trigger if exists "vote_insert"    on "vote";
drop trigger if exists "vote_update"    on "vote";
drop trigger if exists "create_post"    on "post";
drop trigger if exists "create_thread"  on "thread";

create unlogged table if not exists "user" (
    "id"        bigserial                   not null primary key,
    "nickname"  citext collate "ucs_basic"  not null unique,
    "fullname"  citext                      not null,
    "about"     text,
    "email"     citext                      not null unique
);

create unlogged table if not exists "forum" (
    "id"        bigserial not null primary key,
    "title"     text      not null,
    "user"      citext    not null,
    "slug"      citext    not null unique,
    "posts"     bigint    default 0,
    "threads"   int       default 0
);

create unlogged table if not exists "thread" (
    "id"        bigserial     not null primary key,
    "title"     text          not null,
    "author"    citext        not null,
    "forum"     citext,
    "message"   text          not null,
    "votes"     int           default 0,
    "slug"      citext,
    "created"   timestamptz   default now()
);

create unlogged table if not exists "post" (
    "id"        bigserial   not null primary key,
    "parent"    bigint      default 0,
    "author"    citext      not null,
    "message"   text        not null,
    "isEdited"  bool        default false,
    "forum"     citext,
    "thread"    int,
    "created"   timestamptz default now(),
    "path"      bigint[]    not null default '{0}'
);

create unlogged table if not exists "vote" (
    "id"        bigserial   not null primary key,
    "user"      bigint      references "user" (id)   not null,
    "thread"    bigint      references "thread" (id) not null,
    "voice"     int,
    constraint checks       unique ("user", "thread")
);

create function thread_vote() returns trigger as $$
begin
    update "thread"
    set "votes"=(votes + new.voice)
    where "id" = new.thread;
    return new;
end;
$$ language plpgsql;

create trigger "vote_insert"
    after insert
    on "vote"
    for each row
execute procedure thread_vote();

create function thread_vote_update() returns trigger as $$
begin
    update "thread"
    set "votes"=(votes + 2*new.voice)
    where "id" = new.thread;
    return new;
end;
$$ language plpgsql;

create trigger "vote_update"
    after update
    on "vote"
    for each row
execute procedure thread_vote_update();

create function create_post() returns trigger as $$
begin
    update "forum"
    set "posts" = posts + 1
    where "slug" = new.forum;
    new.path = (select "path" from "post" where "id" = new.parent LIMIT 1) || new.id;
    insert into "forum_user" ("user", "forum")
    values ((select "id" from "user" where new.author = nickname), (select "id" from "forum" where new.forum = slug));
    return new;
end
$$ language plpgsql;

create trigger "create_post"
    before insert
    on "post"
    for each row
execute procedure create_post();

create function create_thread() returns trigger as $$
begin
    update "forum"
    set "threads" = threads + 1
    where "slug" = new.forum;
    INSERT INTO "forum_user" ("user", "forum")
    values ((select "id" from "user" where new.author = nickname), (select "id" from "forum" where new.forum = slug));
    return new;
end
$$ language plpgsql;

create trigger "create_thread"
    before insert
    on "thread"
    for each row
execute procedure create_thread();

create unlogged table if not exists "forum_user" (
    "id"    bigserial                           not null primary key,
    "user"  bigint      references "user" (id)  not null,
    "forum" bigint      references "forum" (id) not null
);

drop index if exists index_user_by_nickname;
-- drop index if exists index_thread_by_id_slug;
drop index if exists idxex_thread_by_slug;
drop index if exists index_vote_by_user_thread;
drop index if exists index_user_by_email;
drop index if exists index_forum_by_slug;
drop index if exists index_post_by_thread_path;
drop index if exists index_post_by_thread;
drop index if exists index_forum_user_by_forum;
drop index if exists index_thread_by_forum;

create unique index if not exists index_user_by_nickname on "user" ("nickname");
create unique index if not exists index_user_by_email on "user" (email);
create unique index if not exists index_forum_by_slug on forum (slug);
create index if not exists index_thread_by_slug on thread (slug);

-- create index if not exists index_forum_user_by_user_forum on forum_user (user, forum);

create unique index if not exists index_vote_by_user_thread on vote ("user", thread);

-- create unique index if not exists index_post_by_thread_path on post (thread, path);
-- create index if not exists index_post_by_thread on post (thread);
