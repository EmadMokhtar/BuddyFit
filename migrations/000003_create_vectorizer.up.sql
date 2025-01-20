SELECT ai.create_vectorizer(
      'yt_videos'::regclass,
      destination => 'yt_videos_embeddings',
      embedding => ai.embedding_ollama('nomic-embed-text', 768),
      chunking => ai.chunking_character_text_splitter('transcript', 128, 10),
      formatting => ai.formatting_python_template('title: $title - author: $author - $chunk')
);