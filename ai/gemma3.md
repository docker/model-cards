# Gemma 3

![logo](https://github.com/docker/model-cards/raw/refs/heads/main/logos/gemma-280x184-overview@2x.svg)

Gemma is a versatile AI model family designed for tasks like question answering, summarization, and reasoning. With open weights and responsible commercial use, it supports image-text input, a 128K token context, and over 140 languages.

## Intended uses

Gemma 3 4B model can be used for:

- **Text generation:** Create poems, scripts, code, marketing copy, and email drafts.  
- **Chatbots and conversational AI:** Enable virtual assistants and customer service bots.  
- **Text summarization:** Produce concise summaries of reports and research papers.  
- **Image data extraction:** Interpret and summarize visual data for text-based communication.  
- **Language learning tools:** Aid in grammar correction and interactive writing practice.  
- **Knowledge exploration:** Assist researchers by generating summaries and answering questions.  

## Characteristics

| Attribute             | Details         |
|---------------------- |---------------- |
| **Provider**          | Google DeepMind |
| **Architecture**      | Gemma3          |
| **Cutoff date**       | -               |
| **Languages**         | 140 languages   |
| **Tool calling**      | ✅              |
| **Input modalities**  | Text, Image     |
| **Output modalities** | Text, Code      |
| **License**           | [Gemma Terms](https://ai.google.dev/gemma/terms) |

## Available model variants

| Model variant | Parameters | Quantization | Context window | VRAM | Size |
|---------------|------------|--------------|----------------|------|-------|
| `ai/gemma3:1B-F16` | 999.89M | F16 | - | - | 2.01 GB |
| `ai/gemma3:1B-Q4_K_M` | 999.89M | IQ2_XXS/Q4_K_M | - | - | 806.06 MB |
| `ai/gemma3:4B-F16` | 3.88B | F16 | - | - | 7.77 GB |
| `ai/gemma3:4B-Q4_0` | 3.88B | Q4_0 | - | - | 2.36 GB |
| `ai/gemma3:latest`<br><br>`ai/gemma3:4B-Q4_K_M` | 3.88B | IQ2_XXS/Q4_K_M | - | - | 2.49 GB |

¹: VRAM estimation.

> `:latest` → `4B-Q4_K_M`

## Use this AI model with Docker Model Runner

First, pull the model:

```bash
docker model pull ai/gemma3
```

Then run the model:

```bash
docker model run ai/gemma3
```

For more information on Docker Model Runner, [explore the documentation](https://docs.docker.com/desktop/features/model-runner/).

## Benchmark performance

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
