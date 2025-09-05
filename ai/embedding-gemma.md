# Embedding Gemma

![logo](https://github.com/docker/model-cards/raw/refs/heads/main/logos/gemma-280x184-overview@2x.svg)

**Embedding Gemma** is a state-of-the-art text embedding model from Google DeepMind, designed to create high-quality vector representations of text. Built on the Gemma architecture, this model converts text into dense vector embeddings that capture semantic meaning, making it ideal for retrieval-augmented generation (RAG), semantic search, and similarity tasks. With open weights and efficient design, Embedding Gemma provides a powerful foundation for embedding-based applications.

## Intended uses

Embedding Gemma is designed for applications requiring high-quality text embeddings:

- **Semantic search and retrieval**: Excellent for building search systems, document retrieval, and RAG applications that need to find semantically relevant content.
- **Text similarity and clustering**: Generate embeddings for measuring text similarity, document clustering, and content deduplication tasks.
- **Classification and downstream tasks**: Use embeddings as input features for various NLP classification tasks and machine learning pipelines.

## Characteristics

| Attribute             | Details                                                      |
|---------------------- |--------------------------------------------------------------|
| **Provider**          | Google DeepMind                                              |
| **Architecture**      | Gemma Embedding                                              |
| **Cutoff date**       | -                                                            |
| **Languages**         | English                                                      |
| **Tool calling**      | ❌                                                           |
| **Input modalities**  | Text                                                         |
| **Output modalities** | Embedding vectors                                            |
| **License**           | [Gemma Terms](https://ai.google.dev/gemma/terms)            |

## Available model variants

| Model variant                                                        | Parameters | Quantization | Context window | VRAM¹    | Size      |
|----------------------------------------------------------------------|------------|--------------|----------------|----------|-----------|
| `ai/embedding-gemma:latest`<br><br>`ai/embedding-gemma:300M-F16`     | 300M       | F16          | 2K tokens      | 0.68 GiB | 571.25 MB |
| `ai/embedding-gemma:300M-F16`                                        | 300M       | F16          | 2K tokens      | 0.68 GiB | 571.25 MB |

¹: VRAM estimated based on model characteristics.

> `latest` → `300M-F16`

## Use this AI model with Docker Model Runner

First, pull the model:

```bash
docker model pull ai/embedding-gemma
```

Then run the model:

```bash
docker model run ai/embedding-gemma
```

To generate embeddings using the API:

```bash
curl --location 'http://localhost:12434/engines/llama.cpp/v1/embeddings' \
--header 'Content-Type: application/json' \
--data '{
    "model": "ai/embedding-gemma",
    "input": "Your text to embed here"
  }'
```

For more information on Docker Model Runner, [explore the documentation](https://docs.docker.com/desktop/features/model-runner/).

## Considerations

- **Context length**: The model supports up to 2K tokens. Longer texts may need to be chunked for optimal performance.
- **Language support**: Primarily trained on English text, performance on other languages may vary.
- **Embedding dimension**: The model produces 768-dimensional embeddings suitable for most downstream tasks.
- **Normalization**: Embeddings are normalized by default, making them suitable for cosine similarity calculations.

## Benchmark performance

| Task Category       | Embedding Gemma |
|---------------------|----------------|
| Retrieval          | 54.87          |
| STS                | 78.53          |
| Classification     | 73.26          |
| Clustering         | 44.72          |
| Pair Classification| 85.94          |
| Reranking          | 59.36          |

## Links

- [Embedding Gemma Model Card](https://huggingface.co/google/embeddinggemma-300m)
- [Gemma Model Family](https://ai.google.dev/gemma/docs)
- [Gemma Terms of Use](https://ai.google.dev/gemma/terms)