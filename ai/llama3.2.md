# Llama 3.2 Instruct

![logo](https://github.com/docker/model-cards/raw/refs/heads/main/logos/meta-280x184-overview@2x.svg)

Llama 3.2 introduced lightweight 1B and 3B models at bfloat16 (BF16) precision, later adding quantized versions. The quantized models are significantly faster, with a much lower memory footprint and reduced power consumption, while maintaining nearly the same accuracy as their BF16 counterparts. 

## Intended uses

Llama 3.2 instruct models are designed for:

- **AI assistance on edge devices**, Running chatbots and virtual assistants with minimal latency on low-power * hardware.
-  **Code assistance** , Writing, debugging, and optimizing code on mobile or embedded systems.
- **Content generation** ,Drafting emails, summaries, and creative content on lightweight devices.
- **Low-power AI for smart gadgets**, Enhancing voice assistants on wearables and IoT devices.
- **Edge-based data processing**, Summarizing and analyzing data locally for security and efficiency.

## Characteristics

| Attribute             | Details       |
|---------------------- |-------------- |
| **Provider**          | Meta          |
| **Architecture**      | Llama         |
| **Cutoff date**       | December 2023 |
| **Languages**         | English, German, French, Italian, Portuguese, Hindi, Spanish, and Thai |
| **Tool calling**      | ✅            |
| **Input modalities**  | Text          |
| **Output modalities** | Text, Code    |
| **License**           | [Llama 3.2 Community License](https://github.com/meta-llama/llama-models/blob/main/models/llama3_2/LICENSE) |

## Available model variants

| Model Variant | Parameters | Quantization | Context window | VRAM | Size |
|---------------|------------|--------------|----------------|------|-------|
| `ai/llama3.2:1B-F16` | 1.24B | F16 | - | - | 2.48 GB |
| `ai/llama3.2:1B-Q8_0` | 1.24B | Q8_0 | - | - | 1.32 GB |
| `ai/llama3.2:3B-F16` | 3.21B | F16 | - | - | 6.43 GB |
| `ai/llama3.2:latest`<br><br>`ai/llama3.2:3B-Q4_K_M` | 3.21B | IQ2_XXS/Q4_K_M | - | - | 2.02 GB |

¹: VRAM estimation.

> `:latest` → `3B-Q4_K_M`

## Use this AI model with Docker Model Runner

First, pull the model:

```bash
docker model pull ai/llama3.2
```

Then run the model:

```bash
docker model run ai/llama3.2
```

For more information on Docker Model Runner, [explore the documentation](https://docs.docker.com/desktop/features/model-runner/).

## Benchmark performance

| Capability            | Benchmark                | Llama 3.2 1B |
|----------------------|---------------------------|--------------|
| General              | MMLU                      | 49.3         |
| Re-writing           | Open-rewrite eval         | 41.6         |
| Summarization        | TLDR9+ (test)             | 16.8         |
| Instruct. following  | IFEval                    | 59.5         |
| Math                 | GSM8K (CoT)               | 44.4         |
|                      | MATH (CoT)                | 30.6         |
| Reasoning            | ARC-C                     | 59.4         |
|                      | GPQA                      | 27.2         |
|                      | Hellaswag                 | 41.2         |
| Tool Use             | BFCL V2                   | 25.7         |
|                      | Nexus                     | 13.5         |
| Long Context         | InfiniteBench/En.QA       | 20.3         |
|                      | InfiniteBench/En.MC       | 38.0         |
|                      | NIH/Multi-needle          | 75.0         |
| Multilingual         | MGSM (CoT)                | 24.5         |

## Links

- [Llama](https://www.llama.com/)
