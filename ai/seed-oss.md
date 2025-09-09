# Seed-OSS

![logo](logo)

Seed-OSS is a series of open-source large language models developed by ByteDance's Seed Team, designed for powerful long-context, reasoning, agent and general capabilities, and versatile developer-friendly features. Although trained with only 12T tokens, Seed-OSS achieves excellent performance on several popular open benchmarks.
Powered by Unsloth's GGUF conversion.

## Intended uses

- **Conversational AI**: Engaging in dialogue with users, providing informative and contextually relevant responses.
- **Reasoning tasks**: Excelling in logical reasoning and problem-solving scenarios.
- **Multi-agent frameworks**: Facilitating interactions between multiple AI agents for complex tasks.

## Characteristics

| Attribute        | Details    |
|------------------|------------|
| **Provider**     | ByteDance  |
| **Architecture** | seed_oss   |
| **Cutoff date**  | July 2024  |
| **Tool calling** | ✅          |
| **License**      | Apache 2.0 |

## Available model variants

| Model Variant               | Parameters | Quantization   | Context Window | VRAM      | Size   |
|---------------------------- |----------- |--------------- |--------------- |---------- |------- |
| {name}:{params]_{quant]     | {param}    | {quant}        | {token}        | {size}GB¹ | {size} | 

¹: VRAM estimates based on model characteristics.

## Use this AI model with Docker Model Runner

```bash
docker model run ai/seed-oss
```

## Links
- https://seed.bytedance.com/en/
- https://huggingface.co/unsloth/Seed-OSS-36B-Instruct-GGUF
