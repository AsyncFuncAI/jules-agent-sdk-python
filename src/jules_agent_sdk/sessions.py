"""Sessions API module."""

import time
from typing import Optional, List, Dict, Any

from jules_agent_sdk.models import Session, SessionState
from jules_agent_sdk.base import BaseClient
from jules_agent_sdk.exceptions import JulesAPIError

# Constants for session polling
DEFAULT_POLL_INTERVAL = 5
DEFAULT_TIMEOUT = 600


class SessionsAPI:
    """API client for managing Jules sessions."""

    def __init__(self, client: BaseClient) -> None:
        """Initialize the Sessions API.

        Args:
            client: Base HTTP client instance
        """
        self.client = client

    def create(
        self,
        prompt: str,
        source: str,
        starting_branch: Optional[str] = None,
        title: Optional[str] = None,
        require_plan_approval: bool = False,
        automation_mode: Optional[str] = None,
    ) -> Session:
        """Create a new session.

        Args:
            prompt: The prompt to start the session with
            source: The source to use (e.g., "sources/abc123")
            starting_branch: Optional starting branch for GitHub repos
            title: Optional session title
            require_plan_approval: If True, plans require explicit approval
            automation_mode: Optional automation mode (e.g., "AUTO_CREATE_PR")

        Returns:
            Created Session object

        Example:
            >>> client = JulesClient(api_key="your-api-key")
            >>> session = client.sessions.create(
            ...     prompt="Fix the login bug",
            ...     source="sources/my-repo-id",
            ...     starting_branch="main"
            ... )
            >>> print(session.id)
        """
        data: Dict[str, Any] = {
            "prompt": prompt,
            "sourceContext": {"source": source},
        }

        if starting_branch:
            data["sourceContext"]["githubRepoContext"] = {"startingBranch": starting_branch}

        if title:
            data["title"] = title

        if require_plan_approval:
            data["requirePlanApproval"] = require_plan_approval

        if automation_mode:
            data["automationMode"] = automation_mode

        response = self.client.post("sessions", json=data)
        return Session.from_dict(response)

    def get(self, session_id: str) -> Session:
        """Get a single session by ID.

        Args:
            session_id: The session ID or full name (e.g., "sessions/abc123" or "abc123")

        Returns:
            Session object

        Example:
            >>> session = client.sessions.get("abc123")
            >>> print(session.state)
        """
        if not session_id.startswith("sessions/"):
            session_id = f"sessions/{session_id}"

        response = self.client.get(session_id)
        return Session.from_dict(response)

    def list(
        self, page_size: Optional[int] = None, page_token: Optional[str] = None
    ) -> Dict[str, Any]:
        """List all sessions.

        Args:
            page_size: Maximum number of sessions to return
            page_token: Token for pagination

        Returns:
            Dictionary with 'sessions' list and optional 'nextPageToken'

        Example:
            >>> result = client.sessions.list(page_size=10)
            >>> for session in result['sessions']:
            ...     print(session.id, session.state)
        """
        params: Dict[str, Any] = {}
        if page_size is not None:
            params["pageSize"] = page_size
        if page_token:
            params["pageToken"] = page_token

        response = self.client.get("sessions", params=params)

        sessions = []
        if response.get("sessions"):
            sessions = [Session.from_dict(s) for s in response["sessions"]]

        return {
            "sessions": sessions,
            "nextPageToken": response.get("nextPageToken"),
        }

    def approve_plan(self, session_id: str) -> None:
        """Approve a plan in a session.

        Args:
            session_id: The session ID or full name

        Example:
            >>> client.sessions.approve_plan("abc123")
        """
        if not session_id.startswith("sessions/"):
            session_id = f"sessions/{session_id}"

        self.client.post(f"{session_id}:approvePlan")

    def send_message(self, session_id: str, prompt: str) -> None:
        """Send a message from the user to a session.

        Args:
            session_id: The session ID or full name
            prompt: The message to send

        Example:
            >>> client.sessions.send_message("abc123", "Please also add unit tests")
        """
        if not session_id.startswith("sessions/"):
            session_id = f"sessions/{session_id}"

        self.client.post(f"{session_id}:sendMessage", json={"prompt": prompt})

    def wait_for_completion(
        self,
        session_id: str,
        poll_interval: int = DEFAULT_POLL_INTERVAL,
        timeout: Optional[int] = DEFAULT_TIMEOUT,
    ) -> Session:
        """Poll a session until it completes or fails.

        Args:
            session_id: The session ID or full name
            poll_interval: Seconds between polling requests (default: 5)
            timeout: Optional timeout in seconds (default: 600)

        Returns:
            Final Session object

        Raises:
            TimeoutError: If timeout is reached
            JulesAPIError: If session fails

        Example:
            >>> session = client.sessions.create(prompt="Fix bug", source="sources/repo")
            >>> final_session = client.sessions.wait_for_completion(session.id)
            >>> print(final_session.state)
        """
        start_time = time.time()
        terminal_states = {
            SessionState.COMPLETED,
            SessionState.FAILED,
        }

        while True:
            session = self.get(session_id)

            if session.state in terminal_states:
                if session.state == SessionState.FAILED:
                    raise JulesAPIError(f"Session failed: {session_id}")
                return session

            if timeout and (time.time() - start_time) > timeout:
                raise TimeoutError(f"Session polling timed out after {timeout} seconds")

            time.sleep(poll_interval)
