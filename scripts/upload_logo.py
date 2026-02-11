#!/usr/bin/env python3
"""
Upload repository logo to Docker Hub.
Accepts an explicit logo file path (resolved by the logo-resolver agent)
or falls back to auto-detection from the logos/ directory.
"""
import mimetypes
import os
import sys
import uuid
import requests
from pathlib import Path
from typing import Optional


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


def resolve_logo_path(logo_match_file: Optional[str]) -> Optional[Path]:
    """Resolve the logo file path from the agent's output file.

    Args:
        logo_match_file: Path to the logo-match.txt file written by the agent

    Returns:
        Path to the logo file, or None if no match / agent said 'none'
    """
    if not logo_match_file:
        return None

    match_path = Path(logo_match_file)
    if not match_path.exists():
        print(f"⚠️  Logo match file not found: {match_path}")
        return None

    filename = match_path.read_text(encoding='utf-8').strip()

    if not filename or filename.lower() == 'none':
        print(f"ℹ️  Logo resolver returned 'none' — no matching logo")
        return None

    # The agent writes just the filename; prepend logos/ directory
    logo_path = Path("logos") / filename
    if not logo_path.exists():
        # Try the filename as-is (in case agent wrote a full relative path)
        logo_path = Path(filename)
        if not logo_path.exists():
            print(f"⚠️  Logo file not found: {filename}")
            return None

    print(f"✅ Logo resolved by agent: {logo_path}")
    return logo_path


def main():
    """Main function.

    Usage:
        python3 upload_logo.py <namespace> <repository> [logo_match_file]

    If logo_match_file is provided, reads the logo filename from that file
    (as written by the logo-resolver agent). Otherwise skips.
    """
    if len(sys.argv) < 3:
        print("Usage: python3 upload_logo.py <namespace> <repository> [logo_match_file]", file=sys.stderr)
        return 1

    namespace = sys.argv[1]
    repository = sys.argv[2]
    logo_match_file = sys.argv[3] if len(sys.argv) > 3 else None

    # Get token from environment
    token = os.getenv('DOCKER_HUB_TOKEN')
    if not token:
        print("❌ DOCKER_HUB_TOKEN environment variable must be set", file=sys.stderr)
        return 1

    # Resolve logo path from agent output
    logo_path = resolve_logo_path(logo_match_file)

    if not logo_path:
        print(f"⚠️  No logo resolved for {namespace}/{repository} — skipping logo upload")
        return 0  # Exit successfully (skip, not fail)

    print(f"📤 Uploading logo for {namespace}/{repository}: {logo_path}")

    success = upload_logo(namespace, repository, logo_path, token)
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())
