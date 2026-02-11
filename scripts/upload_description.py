#!/usr/bin/env python3
"""
Upload model card description to Docker Hub.
Reads a markdown file as the full_description and a short description file
for the repository's short description field.
"""
import os
import sys
import requests
from pathlib import Path


def upload_description(namespace: str, repository: str, short_description: str, full_description: str, token: str) -> bool:
    """Upload description to Docker Hub.

    Args:
        namespace: Docker Hub namespace (e.g., 'ai', 'aistaging')
        repository: Docker Hub repository name (e.g., 'qwen3')
        short_description: Short description (max 100 chars)
        full_description: Full markdown content for the repository page
        token: Docker Hub bearer token

    Returns:
        True if upload succeeded, False otherwise
    """
    url = f"https://hub.docker.com/v2/namespaces/{namespace}/repositories/{repository}"
    payload = {
        'description': short_description,
        'full_description': full_description
    }

    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    }

    try:
        print(f"🌐 Making API request to Docker Hub...")
        print(f"📦 Repository: {namespace}/{repository}")
        response = requests.patch(url, json=payload, headers=headers, timeout=30)

        if 200 <= response.status_code < 300:
            print(f"✅ Successfully uploaded description for {namespace}/{repository}")
            print(f"📝 Short description: {short_description}")
            print(f"📄 Full description: {len(full_description)} chars")
            print(f"🔗 https://hub.docker.com/r/{namespace}/{repository}")
            return True
        else:
            error_msg = response.text[:200] + "..." if len(response.text) > 200 else response.text
            print(f"❌ Failed to upload description (HTTP {response.status_code})")
            print(f"Response: {error_msg}")
            return False

    except requests.RequestException as e:
        print(f"❌ Request failed: {e}")
        return False


def main():
    """Main function.

    Usage:
        python3 upload_description.py <namespace> <repository> <model_card_path> <short_description_path>
    """
    if len(sys.argv) != 5:
        print("Usage: python3 upload_description.py <namespace> <repository> <model_card_path> <short_description_path>", file=sys.stderr)
        return 1

    namespace = sys.argv[1]
    repository = sys.argv[2]
    model_card_path = Path(sys.argv[3])
    short_desc_path = Path(sys.argv[4])

    # Validate files exist
    if not model_card_path.exists():
        print(f"❌ Model card file not found: {model_card_path}", file=sys.stderr)
        return 1

    if not short_desc_path.exists():
        print(f"❌ Short description file not found: {short_desc_path}", file=sys.stderr)
        return 1

    # Read model card (full description)
    full_description = model_card_path.read_text(encoding='utf-8')
    if not full_description.strip():
        print(f"❌ Model card file is empty: {model_card_path}", file=sys.stderr)
        return 1

    # Read short description
    short_description = short_desc_path.read_text(encoding='utf-8').strip()
    if not short_description:
        print(f"❌ Short description file is empty: {short_desc_path}", file=sys.stderr)
        return 1

    # Truncate short description to 100 chars (Docker Hub limit)
    if len(short_description) > 100:
        print(f"⚠️  Short description truncated from {len(short_description)} to 100 chars")
        short_description = short_description[:100]

    # Get token from environment
    token = os.getenv('DOCKER_HUB_TOKEN')
    if not token:
        print("❌ DOCKER_HUB_TOKEN environment variable must be set", file=sys.stderr)
        return 1

    print(f"📤 Uploading description for {namespace}/{repository}")

    success = upload_description(namespace, repository, short_description, full_description, token)
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())
