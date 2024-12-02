CREATE OR REPLACE FUNCTION generate_rag_response(query_text TEXT)
RETURNS TEXT AS $$
DECLARE
   context_chunks TEXT;
   response TEXT;
BEGIN
   -- Perform similarity search to find relevant blog posts
   SELECT string_agg(title || ': ' || chunk, ' ') INTO context_chunks
   FROM (
       SELECT title, chunk
       FROM public.yt_videos_embeddings
       ORDER BY embedding <=> ai.ollama_embed('nomic-embed-text', query_text)
       LIMIT 5
   ) AS relevant_yt_videos;

   -- Generate a summary using gpt-4o-mini
   SELECT ai.ollama_chat_complete(
       'llama3.1:latest',
       jsonb_build_array(
           jsonb_build_object('role', 'system', 'content', 'You are a helpful gym personal trainer and professional bodybuilder. Use only the context provided to answer the question. Also mention the titles of the youtube videos you use to answer the question.'),
           jsonb_build_object('role', 'user', 'content', format('Context: %s\n\nUser Question: %s\n\nAssistant:', context_chunks, query_text))
       )
   )->'message'->>'content' INTO response;

   RETURN response;
END;
$$ LANGUAGE plpgsql;