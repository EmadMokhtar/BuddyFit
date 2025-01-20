CREATE TABLE IF NOT EXISTS public.yt_videos (
    id SERIAL PRIMARY KEY,
    author character varying,
    title character varying,
    transcript text,
    created_at timestamp without time zone DEFAULT now()
);