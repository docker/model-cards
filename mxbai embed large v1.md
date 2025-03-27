# mxbai-embed-large-v1

![logo](https://github.com/docker/model-cards/raw/refs/heads/main/logos/mixelbread-280x184-overview@2x.svg)

**mxbai-embed-large-v1** is a state-of-the-art English language embedding model developed by Mixedbread AI. It converts text into dense vector representations, capturing the semantic essence of the input. Trained on a vast dataset exceeding 700 million pairs using contrastive training methods and fine-tuned on over 30 million high-quality triplets with the AnglE loss function, this model adapts to a wide range of topics and domains, making it suitable for various real-world applications and Retrieval-Augmented Generation (RAG) use cases.


## Characteristics

| Attribute             | Details          |
|---------------------- |------------------|
| **Provider**          | Mixedbread AI    |
| **Architecture**      | BERT             |
| **Cutoff Date**       | September 2023   |
| **Languages**         | English          |
| **Tool Calling**      | ❌               |
| **Input Modalities**  | Text             |
| **Output Modalities** | Text embeddings  |
| **License**           | Apache 2.0       |

## Available Model Variants

| Model Variant                                                 | Parameters | Quantization   | Context Window | VRAM      | Size   | 
|-------------------------------------------------------------- |----------- |--------------- |--------------- |---------- |------- |
| `ai/mxbai-embed-large:latest` `ai/mxbai-embed-large:335M-F16` | 335M       | F16            | 512 tokens     | 0.8GB¹    | 670MB  | 

¹: VRAM estimates based on model characteristics.

`:latest` → `mxbai-embed-large:335M-F16`

## Intended Uses

mxbai-embed-large-v1 is designed for generating sentence embeddings suitable for various NLP applications.

- **Semantic Search and Information Retrieval:** Specifically designed for RAG, this model enhances search systems by providing relevant document embeddings, improving the accuracy and relevance of search results.
- **Semantic Textual Similarity:** Measures the similarity between sentences, aiding in tasks such as clustering, duplicate detection, and paraphrase identification.
- **Text Classification:** Serves as input features for classifiers in tasks like sentiment analysis, topic categorization, and intent detection.


## Considerations

- **Prompt Usage:** For retrieval tasks, prepend the query with the prompt: "Represent this sentence for searching relevant passages:". This practice helps the model understand the context and improves performance. For other tasks, the text can be used as-is without any additional prompt.
- **Language Limitation:** The model is trained exclusively on English text and is specifically designed for the English language.
- **Sequence Length:** The suggested maximum sequence length is 512 tokens. Longer sequences may be truncated, leading to a loss of information.

## How to Run This AI Model

You can pull the model using:

```
docker model pull mixedbread-ai/mxbai-embed-large-v1
```

## Benchmarks and performance

| Task Category           | mxbai-embed-large-v1 |
|-------------------------|----------------------|
| **Avg (56 datasets)**   | 64.68                |
| **Classification**      | 75.64                |
| **Clustering**          | 46.71                |
| **Pair Classification** | 87.2                 |
| **Reranking**           | 60.11                |
| **Retrieval**           | 54.39                |
| **STS**                 | 85.00                |
| **Summarization**       | 32.71                |

## Links
- [Open Source Strikes Bread - New Fluffy Embedding Model](https://www.mixedbread.com/blog/mxbai-embed-large-v1)
- [Mixelbread Docs:mxbai-embed-large-v1](https://www.mixedbread.com/docs/embeddings/mxbai-embed-large-v1)
