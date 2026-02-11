#!/usr/bin/env python3
"""
Docker Hub token authentication script.
Authenticates with Docker Hub using HUB_USER and HUB_PAT environment variables.
"""
import json
import os
import sys
import requests
from typing import Optional


def get_base_url(stage: bool = False) -> str:
    """Get the Docker Hub API base URL"""
    if stage:
        return "https://hub-stage.docker.com"
    return "https://hub.docker.com"


def get_token(stage: bool = False, max_retries: int = 3, retry_delay: float = 2.0) -> str:
    """
    Get Docker Hub authentication token.

    Args:
        stage: Whether to use staging environment
        max_retries: Maximum number of retry attempts for transient errors
        retry_delay: Delay between retries in seconds

    Returns:
        Authentication token string

    Raises:
        Exception: If authentication fails or environment variables are missing
    """
    import time

    username = os.getenv("HUB_USER")
    if not username:
        raise Exception("HUB_USER is not set")

    password = os.getenv("HUB_PAT")
    if not password:
        raise Exception("HUB_PAT is not set")

    last_error = None
    for attempt in range(max_retries):
        try:
            return _get_token_v2(username, password, stage)
        except Exception as e:
            error_str = str(e)
            last_error = e

            is_transient = any(code in error_str for code in ['500', '502', '503', '504'])

            if is_transient and attempt < max_retries - 1:
                wait_time = retry_delay * (2 ** attempt)
                print(f"v2 login failed with transient error ({e}), retrying in {wait_time}s... (attempt {attempt + 1}/{max_retries})", file=sys.stderr)
                time.sleep(wait_time)
            else:
                break

    raise Exception(f"Docker Hub login failed after {max_retries} attempts: {last_error}")


def _get_token_v2(username: str, password: str, stage: bool = False) -> str:
    """Try v2 users/login endpoint"""
    payload = {
        "username": username,
        "password": password
    }

    url = f"{get_base_url(stage)}/v2/users/login/"
    headers = {"Content-Type": "application/json", "User-Agent": "curl/8.0"}

    try:
        response = requests.post(url, json=payload, headers=headers, timeout=30)

        if response.status_code != 200:
            try:
                error_data = response.json()
                error_msg = error_data.get('detail') or error_data.get('message') or response.text
            except Exception:
                error_msg = response.text or response.reason

            raise Exception(f"failed to login: {response.status_code} - {error_msg}")

        login_response = response.json()
        token = login_response.get("token")

        if not token:
            raise Exception(f"no token found in login response. Response keys: {list(login_response.keys())}")

        return token

    except requests.RequestException as e:
        raise Exception(f"request failed: {str(e)}")
    except json.JSONDecodeError as e:
        raise Exception(f"failed to parse login response: {str(e)}")


def main():
    """Main function for command line usage"""
    stage = "--stage" in sys.argv

    try:
        token = get_token(stage)
        print(token)
        return 0
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
