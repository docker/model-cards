# Llama 3.2 Instruct

![meta](https://github.com/user-attachments/assets/5fc304e6-44de-40b8-af95-64f85d8ac3c0)

Llama 3.2 introduced lightweight 1B and 3B models at bfloat16 (BF16) precision, later adding quantized versions. These quantized models are significantly faster, with a much lower memory footprint and reduced power consumption, while maintaining nearly the same accuracy as their BF16 counterparts. 

## Characteristics

| Attribute             | Details       |
|---------------------- |-------------- |
| **Provider**          | Meta          |
| **Architecture**      | Llama         |
| **Cutoff Date**       | December 2023 |
| **Languages**         | English, German, French, Italian, Portuguese, Hindi, Spanish, and Thai |
| **Input Modalities**  | Text          |
| **Output Modalities** | Text, Code    |
| **License**           | [Llama 3.2 Community License](https://github.com/meta-llama/llama-models/blob/main/models/llama3_2/LICENSE) |

## Available Model Variants

| Model Variant               | Parameters | Quantization   | Context Window | VRAM    | Size   | Download |
|---------------------------- |----------- |--------------- |--------------- |-------- |------- |--------- |
| `llama3.2-1b-instruct:fp16`  | 1B         | fp16           | 16K tokens     |  7.5GB¹ | -      | Link     |
| `llama3.2-1b-instruct:q8_0` | 1B         | q8             | 16K tokens     |  3.5GB¹ | -      | Link     |
¹: VRAM estimated based on model characteristics.

## Intended Uses

Llama 3.2 instruct models are designed for:

- **AI Assistance on Edge Devices**, Running chatbots and virtual assistants with minimal latency on low-power * hardware.
-  **Code Assistance** , Writing, debugging, and optimizing code on mobile or embedded systems.
- **Content Generation** ,Drafting emails, summaries, and creative content on lightweight devices.
- **Low-Power AI for Smart Gadgets**, Enhancing voice assistants on wearables and IoT devices.
- **Edge-Based Data Processing**, Summarizing and analyzing data locally for security and efficiency.

## How to Run This AI Model

You can pull the model using:
```
docker model pull llama3.2-1b-instruct
```

To run the model:
```
docker model run llama3.2-1b-instruct
```

## Benchmark Performance

| Capability            | Benchmark                | Llama 3.2 1B      |
|----------------------|---------------------------|-------------------|
| General              | MMLU                      | 49.3              |
| Re-writing           | Open-rewrite eval         | 41.6              |
| Summarization        | TLDR9+ (test)             | 16.8              |
| Instruct. following  | IFEval                    | 59.5              |
| Math                 | GSM8K (CoT)               | 44.4              |
|                      | MATH (CoT)                | 30.6              |
| Reasoning            | ARC-C                     | 59.4              |
|                      | GPQA                      | 27.2              |
|                      | Hellaswag                 | 41.2              |
| Tool Use             | BFCL V2                   | 25.7              |
|                      | Nexus                     | 13.5              |
| Long Context         | InfiniteBench/En.QA       | 20.3              |
|                      | InfiniteBench/En.MC       | 38.0              |
|                      | NIH/Multi-needle          | 75.0              |
| Multilingual         | MGSM (CoT)                | 24.5              |




## Links
- [Llama](https://www.llama.com/)
