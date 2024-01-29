CREATE TABLE IF NOT EXISTS public.users (
    id varchar NOT NULL,
    name varchar NULL,
    surname varchar NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
INSERT INTO public.users (id,"name",surname) VALUES
    ('PeterGonzalesisfna','Peter','Gonzales');
