#!/usr/bin/env python3
"""
JSON-RPC bridge server for Jules Agent VSCode Extension.

This module provides a JSON-RPC interface between the TypeScript extension
and the Python Jules SDK, enabling seamless communication.
"""

import sys
import json
import os
import asyncio
from typing import Dict, Any, Optional, List
from jules_agent_sdk import JulesClient, AsyncJulesClient
from jules_agent_sdk.exceptions import (
    JulesAPIError,
    JulesAuthenticationError,
    JulesNotFoundError,
    JulesValidationError,
    JulesRateLimitError
)


class JulesBridge:
    """Bridge between VSCode extension and Jules SDK."""

    def __init__(self):
        self.client: Optional[JulesClient] = None
        self.async_client: Optional[AsyncJulesClient] = None
        self.api_key: Optional[str] = None
        self.base_url: Optional[str] = None
        self.timeout: int = 30

    def _ensure_client(self) -> JulesClient:
        """Ensure client is initialized."""
        if not self.client:
            if not self.api_key:
                raise ValueError("API key not configured. Call 'initialize' first.")
            self.client = JulesClient(
                api_key=self.api_key,
                base_url=self.base_url,
                timeout=self.timeout
            )
        return self.client

    def initialize(self, api_key: str, base_url: Optional[str] = None,
                   timeout: int = 30) -> Dict[str, Any]:
        """Initialize the Jules client."""
        try:
            self.api_key = api_key
            self.base_url = base_url
            self.timeout = timeout

            # Test the connection
            self.client = JulesClient(
                api_key=api_key,
                base_url=base_url,
                timeout=timeout
            )

            # Try to list sources to validate API key
            _ = self.client.sources.list(page_size=1)

            return {"success": True, "message": "Client initialized successfully"}
        except JulesAuthenticationError as e:
            return {"success": False, "error": "Invalid API key", "details": str(e)}
        except Exception as e:
            return {"success": False, "error": str(e)}

    def create_session(self, prompt: str, source: str,
                       starting_branch: Optional[str] = None,
                       title: Optional[str] = None,
                       require_plan_approval: bool = False) -> Dict[str, Any]:
        """Create a new Jules session."""
        try:
            client = self._ensure_client()
            session = client.sessions.create(
                prompt=prompt,
                source=source,
                starting_branch=starting_branch,
                title=title,
                require_plan_approval=require_plan_approval
            )
            return {
                "success": True,
                "session": {
                    "id": session.id,
                    "name": session.name,
                    "title": session.title,
                    "state": session.state.value,
                    "url": session.url,
                    "prompt": session.prompt,
                    "createTime": session.create_time,
                    "updateTime": session.update_time
                }
            }
        except Exception as e:
            return self._format_error(e)

    def get_session(self, session_id: str) -> Dict[str, Any]:
        """Get session details."""
        try:
            client = self._ensure_client()
            session = client.sessions.get(session_id)
            return {
                "success": True,
                "session": {
                    "id": session.id,
                    "name": session.name,
                    "title": session.title,
                    "state": session.state.value,
                    "url": session.url,
                    "prompt": session.prompt,
                    "createTime": session.create_time,
                    "updateTime": session.update_time,
                    "outputs": [
                        {
                            "pullRequest": {
                                "url": output.pull_request.url,
                                "title": output.pull_request.title,
                                "description": output.pull_request.description
                            } if output.pull_request else None
                        } for output in session.outputs
                    ]
                }
            }
        except Exception as e:
            return self._format_error(e)

    def list_sessions(self, page_size: int = 50) -> Dict[str, Any]:
        """List all sessions."""
        try:
            client = self._ensure_client()
            result = client.sessions.list(page_size=page_size)
            sessions = [
                {
                    "id": s.id,
                    "name": s.name,
                    "title": s.title,
                    "state": s.state.value,
                    "url": s.url,
                    "createTime": s.create_time,
                    "updateTime": s.update_time
                } for s in result["sessions"]
            ]
            return {
                "success": True,
                "sessions": sessions,
                "nextPageToken": result.get("nextPageToken")
            }
        except Exception as e:
            return self._format_error(e)

    def approve_plan(self, session_id: str) -> Dict[str, Any]:
        """Approve a session plan."""
        try:
            client = self._ensure_client()
            client.sessions.approve_plan(session_id)
            return {"success": True, "message": "Plan approved"}
        except Exception as e:
            return self._format_error(e)

    def send_message(self, session_id: str, prompt: str) -> Dict[str, Any]:
        """Send a message to a session."""
        try:
            client = self._ensure_client()
            client.sessions.send_message(session_id, prompt)
            return {"success": True, "message": "Message sent"}
        except Exception as e:
            return self._format_error(e)

    def list_activities(self, session_id: str, page_size: int = 100) -> Dict[str, Any]:
        """List activities for a session."""
        try:
            client = self._ensure_client()
            activities = client.activities.list_all(session_id)
            return {
                "success": True,
                "activities": [
                    {
                        "id": a.id,
                        "name": a.name,
                        "description": a.description,
                        "createTime": a.create_time,
                        "originator": a.originator,
                        "artifacts": [
                            {
                                "changeSet": {
                                    "gitPatch": {
                                        "unidiffPatch": art.change_set.git_patch.unidiff_patch if art.change_set.git_patch else None,
                                        "suggestedCommitMessage": art.change_set.git_patch.suggested_commit_message if art.change_set.git_patch else None
                                    } if art.change_set else None
                                } if art.change_set else None,
                                "bashOutput": {
                                    "command": art.bash_output.command,
                                    "output": art.bash_output.output,
                                    "exitCode": art.bash_output.exit_code
                                } if art.bash_output else None
                            } for art in a.artifacts
                        ],
                        "agentMessaged": a.agent_messaged,
                        "userMessaged": a.user_messaged,
                        "planGenerated": a.plan_generated,
                        "planApproved": a.plan_approved,
                        "progressUpdated": a.progress_updated,
                        "sessionCompleted": a.session_completed,
                        "sessionFailed": a.session_failed
                    } for a in activities
                ]
            }
        except Exception as e:
            return self._format_error(e)

    def list_sources(self) -> Dict[str, Any]:
        """List all available sources."""
        try:
            client = self._ensure_client()
            sources = client.sources.list_all()
            return {
                "success": True,
                "sources": [
                    {
                        "id": s.id,
                        "name": s.name,
                        "githubRepo": {
                            "owner": s.github_repo.owner,
                            "repo": s.github_repo.repo,
                            "isPrivate": s.github_repo.is_private,
                            "defaultBranch": s.github_repo.default_branch.display_name if s.github_repo.default_branch else None
                        } if s.github_repo else None
                    } for s in sources
                ]
            }
        except Exception as e:
            return self._format_error(e)

    def wait_for_completion(self, session_id: str, poll_interval: int = 5,
                            timeout: int = 600) -> Dict[str, Any]:
        """Wait for session completion."""
        try:
            client = self._ensure_client()
            session = client.sessions.wait_for_completion(
                session_id=session_id,
                poll_interval=poll_interval,
                timeout=timeout
            )
            return {
                "success": True,
                "session": {
                    "id": session.id,
                    "state": session.state.value,
                    "url": session.url
                }
            }
        except TimeoutError:
            return {"success": False, "error": "Session timed out"}
        except Exception as e:
            return self._format_error(e)

    def close(self) -> Dict[str, Any]:
        """Close the client connection."""
        try:
            if self.client:
                self.client.close()
                self.client = None
            return {"success": True, "message": "Client closed"}
        except Exception as e:
            return {"success": False, "error": str(e)}

    def _format_error(self, error: Exception) -> Dict[str, Any]:
        """Format error response."""
        error_type = type(error).__name__
        return {
            "success": False,
            "error": str(error),
            "errorType": error_type
        }


def handle_request(bridge: JulesBridge, request: Dict[str, Any]) -> Dict[str, Any]:
    """Handle a single JSON-RPC request."""
    method = request.get("method")
    params = request.get("params", {})
    request_id = request.get("id")

    # Map methods to bridge functions
    method_map = {
        "initialize": lambda: bridge.initialize(**params),
        "createSession": lambda: bridge.create_session(**params),
        "getSession": lambda: bridge.get_session(**params),
        "listSessions": lambda: bridge.list_sessions(**params),
        "approvePlan": lambda: bridge.approve_plan(**params),
        "sendMessage": lambda: bridge.send_message(**params),
        "listActivities": lambda: bridge.list_activities(**params),
        "listSources": lambda: bridge.list_sources(**params),
        "waitForCompletion": lambda: bridge.wait_for_completion(**params),
        "close": lambda: bridge.close()
    }

    if method not in method_map:
        return {
            "jsonrpc": "2.0",
            "id": request_id,
            "error": {
                "code": -32601,
                "message": f"Method not found: {method}"
            }
        }

    try:
        result = method_map[method]()
        return {
            "jsonrpc": "2.0",
            "id": request_id,
            "result": result
        }
    except Exception as e:
        return {
            "jsonrpc": "2.0",
            "id": request_id,
            "error": {
                "code": -32603,
                "message": str(e)
            }
        }


def main():
    """Main entry point for the bridge server."""
    bridge = JulesBridge()

    try:
        # Read from stdin, write to stdout
        for line in sys.stdin:
            line = line.strip()
            if not line:
                continue

            try:
                request = json.loads(line)
                response = handle_request(bridge, request)
                print(json.dumps(response), flush=True)
            except json.JSONDecodeError as e:
                error_response = {
                    "jsonrpc": "2.0",
                    "error": {
                        "code": -32700,
                        "message": f"Parse error: {str(e)}"
                    }
                }
                print(json.dumps(error_response), flush=True)
    finally:
        bridge.close()


if __name__ == "__main__":
    main()
