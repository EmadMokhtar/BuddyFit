CREATE OR REPLACE FUNCTION get_related_docs(query_text TEXT)
RETURNS TEXT AS $$
DECLARE
    context_chunks TEXT;
    embedding_model_name TEXT = 'nomic-embed-text';
BEGIN
   -- Perform similarity search to find relevant YouTube videos based on the query text
   SELECT string_agg(title || ': ' || chunk, ' ') INTO context_chunks
   FROM (
       SELECT title, chunk
       FROM public.yt_videos_embeddings
       ORDER BY embedding <=> ai.ollama_embed(embedding_model_name, query_text)
       LIMIT 5
   ) AS relevant_yt_videos;

RETURN context_chunks;
END;
$$ LANGUAGE plpgsql;