DO $$
DECLARE
    embedding_table varchar;
    vectorizer_id integer;
BEGIN
    SELECT v.target_table INTO embedding_table FROM ai.vectorizer v WHERE v.source_table = 'yt_videos' LIMIT 1;
    SELECT v.id INTO vectorizer_id FROM ai.vectorizer v WHERE v.source_table = 'yt_videos' LIMIT 1;

    PERFORM ai.drop_vectorizer(vectorizer_id);
    EXECUTE 'DROP TABLE IF EXISTS ' || embedding_table || ' CASCADE';
END $$;