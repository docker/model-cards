# Gemma 3

<img src="https://github.com/jalonsogo/model-cards/blob/04897016e65199efa2ca12c637f734d22121dbd3/logos/google-gemma3.png" width="250" />

Gemma is a versatile AI model family designed for tasks like question answering, summarization, and reasoning. With open weights and responsible commercial use, it supports image-text input, a 128K token context, and over 140 languages.

## Characteristics

| Attribute             | Details         |
|---------------------- |---------------- |
| **Provider**          | Google DeepMind |
| **Architecture**      | Gemma3          |
| **Cutoff Date**       | -               |
| **Languages**         | 140 languages   |
| **Input Modalities**  | Text, Image     |
| **Output Modalities** | Text, Code      |
| **License**           | [Gemma Terms](https://ai.google.dev/gemma/terms) |

## Available Model Variants

| Model Variant       | Parameters | Quantization   | Context Window | VRAM    | Size   | Download |
|-------------------- |----------- |--------------- |--------------- |-------- |------- |--------- |
| `gemma3:4b-fp16`    | 4B         | fp16           | 128K tokens    |  6.4GB¹ | -      | Link     |
| `gemma3:4b-q4_k_m`  | 4B         | Q4 K M         | 128K tokens    |  3.4GB¹ | -      | Link     |
¹: VRAM extracted from Gemma documentation ([link](https://ai.google.dev/gemma/docs/core#128k-context))

## Intended Uses

Gemma 3 4B model can be used for:

- **Text Generation:** Create poems, scripts, code, marketing copy, and email drafts.  
- **Chatbots and Conversational AI:** Enable virtual assistants and customer service bots.  
- **Text Summarization:** Produce concise summaries of reports and research papers.  
- **Image Data Extraction:** Interpret and summarize visual data for text-based communication.  
- **Language Learning Tools:** Aid in grammar correction and interactive writing practice.  
- **Knowledge Exploration:** Assist researchers by generating summaries and answering questions.  

## How to Run This AI Model

You can pull the model using:
```
docker model pull ai/gemma3
```

To run the model:
```
docker model run ai/gemma3
```

## Benchmark Performance

| Category       | Benchmark          | Value  |
|---------------|--------------------|--------|
| General       | MMLU               | 59.6   |
|               | GSM8K              | 38.4   |
|               | ARC-Challenge      | 56.2   |
|               | BIG-Bench Hard     | 50.9   |
|               | DROP               | 60.1   |
| STEM & Code   | MATH               | 24.2   |
|               | MBPP               | 46.0   |
|               | HumanEval          | 36.0   |
| Multilingual  | MGSM               | 34.7   |
|               | Global-MMLU-Lite   | 57.0   |
|               | XQuAD (all)        | 68.0   |
| Multimodal    | VQAv2              | 63.9   |
|               | TextVQA            | 58.9   |
|               | DocVQA             | 72.8   |



## Links
- [Gemma 3 Model Overview](https://ai.google.dev/gemma/docs/core)
- [Gemma 3 Technical Report](https://storage.googleapis.com/deepmind-media/gemma/Gemma3Report.pdf)
