#!/usr/bin/env python3
"""
Upload repository logo to Docker Hub.
Auto-detects the logo file from the logos/ directory based on the repository name,
then uploads it to the Docker Hub media service.
"""
import mimetypes
import os
import re
import sys
import uuid
import requests
from pathlib import Path
from typing import Optional


def find_logo(logos_dir: Path, repository: str, logo_prefix: Optional[str] = None) -> Optional[Path]:
    """Auto-detect a logo file from the logos directory.

    Search strategy:
    1. If logo_prefix is provided, search for logos/{logo_prefix}-*
    2. Try the full repository name: logos/{repository}-*
    3. Strip trailing version numbers: e.g., qwen3 → qwen, llama3.1 → llama
    4. Among matches, prefer SVG over PNG, prefer larger sizes

    Args:
        logos_dir: Path to the logos directory
        repository: Repository name (e.g., 'qwen3', 'llama3.1')
        logo_prefix: Optional override prefix for logo search

    Returns:
        Path to the best matching logo file, or None if not found
    """
    if not logos_dir.exists():
        print(f"⚠️  Logos directory not found: {logos_dir}")
        return None

    # Build list of prefixes to try (in priority order)
    prefixes = []
    if logo_prefix:
        prefixes.append(logo_prefix)
    prefixes.append(repository)

    # Strip trailing version numbers: qwen3 → qwen, llama3.1 → llama, deepseek3.2 → deepseek
    stripped = re.sub(r'[\d]+[.\d]*$', '', repository).rstrip('-')
    if stripped and stripped != repository:
        prefixes.append(stripped)

    # Also try with hyphens removed from version parts: granite-4.0-h-micro → granite
    base = repository.split('-')[0]
    if base and base not in prefixes:
        prefixes.append(base)

    print(f"🔍 Searching for logo with prefixes: {prefixes}")

    for prefix in prefixes:
        candidates = list(logos_dir.glob(f"{prefix}-*"))
        if not candidates:
            continue

        print(f"📁 Found {len(candidates)} candidates for prefix '{prefix}'")

        # Score candidates: prefer SVG, prefer larger sizes
        def score(path: Path) -> tuple:
            name = path.name.lower()
            is_svg = name.endswith('.svg')
            # Extract size hint from filename
            size = 0
            if '280x' in name:
                size = 280
            elif '120x' in name:
                size = 120
            elif '32x' in name:
                size = 32
            # Prefer non-retina for upload (simpler), but retina is fine too
            return (is_svg, size, path.name)

        candidates.sort(key=score, reverse=True)
        best = candidates[0]
        print(f"✅ Best match: {best.name}")
        return best

    print(f"⚠️  No logo found for repository '{repository}' in {logos_dir}")
    return None


def upload_logo(namespace: str, repository: str, image_path: Path, token: str) -> bool:
    """Upload logo to the Docker Hub media service.

    Args:
        namespace: Docker Hub namespace
        repository: Docker Hub repository name
        image_path: Path to the logo file
        token: Docker Hub bearer token

    Returns:
        True if upload succeeded, False otherwise
    """
    repo = f"{namespace}%2F{repository}"
    url = f"https://hub.docker.com/api/media/repos_logo/v1/{repo}/media"

    content_type = mimetypes.guess_type(str(image_path))[0]
    if not content_type:
        content_type = "application/octet-stream"

    # Generate boundary for multipart form
    boundary = f"----formdata-python-{uuid.uuid4().hex}"

    # Read file data
    with open(image_path, 'rb') as f:
        file_data = f.read()

    # Manually construct multipart body
    body_parts = []
    body_parts.append(f'--{boundary}')
    body_parts.append('Content-Disposition: form-data; name="file"; filename="file"')
    body_parts.append(f'Content-Type: {content_type}')
    body_parts.append('')

    body_text = '\r\n'.join(body_parts) + '\r\n'
    body = body_text.encode('utf-8') + file_data + b'\r\n'

    # Add type field
    body += f'--{boundary}\r\n'.encode('utf-8')
    body += b'Content-Disposition: form-data; name="type"\r\n'
    body += b'\r\n'
    body += b'logo\r\n'

    # Add dark field
    body += f'--{boundary}\r\n'.encode('utf-8')
    body += b'Content-Disposition: form-data; name="dark"\r\n'
    body += b'\r\n'
    body += b'false\r\n'

    # End boundary
    body += f'--{boundary}--\r\n'.encode('utf-8')

    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': f'multipart/form-data; boundary={boundary}'
    }

    print(f"🌐 Uploading logo to Docker Hub media service...")
    print(f"📦 Repository: {namespace}/{repository}")
    print(f"📁 File: {image_path.name} ({content_type})")
    print(f"📏 Body length: {len(body)} bytes")

    try:
        response = requests.post(url, data=body, headers=headers, timeout=30)

        print(f"📥 Response: {response.status_code} {response.reason}")
        if response.text:
            print(f"📄 Response body: {response.text[:200]}")

        if response.status_code in [200, 201]:
            print(f"✅ Successfully uploaded logo for {namespace}/{repository}")
            return True
        else:
            print(f"❌ Logo upload failed: {response.status_code}")
            return False

    except Exception as e:
        print(f"❌ Request failed: {e}")
        return False


def main():
    """Main function.

    Usage:
        python3 upload_logo.py <namespace> <repository> [logo_prefix]
    """
    if len(sys.argv) < 3:
        print("Usage: python3 upload_logo.py <namespace> <repository> [logo_prefix]", file=sys.stderr)
        return 1

    namespace = sys.argv[1]
    repository = sys.argv[2]
    logo_prefix = sys.argv[3] if len(sys.argv) > 3 else None

    # Get token from environment
    token = os.getenv('DOCKER_HUB_TOKEN')
    if not token:
        print("❌ DOCKER_HUB_TOKEN environment variable must be set", file=sys.stderr)
        return 1

    # Auto-detect logo
    logos_dir = Path("logos")
    logo_path = find_logo(logos_dir, repository, logo_prefix)

    if not logo_path:
        print(f"⚠️  No logo found for {namespace}/{repository} — skipping logo upload")
        return 0  # Exit successfully (skip, not fail)

    print(f"📤 Uploading logo for {namespace}/{repository}: {logo_path}")

    success = upload_logo(namespace, repository, logo_path, token)
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())
