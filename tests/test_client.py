"""Tests for the Jules client."""

import pytest
from unittest.mock import Mock, patch, MagicMock
from jules_agent_sdk import JulesClient
from jules_agent_sdk.exceptions import JulesAuthenticationError, JulesValidationError


class TestJulesClient:
    """Test cases for JulesClient."""

    def test_client_initialization(self):
        """Test client initializes correctly."""
        client = JulesClient(api_key="test-api-key")
        assert client is not None
        assert client.sessions is not None
        assert client.activities is not None
        assert client.sources is not None

    def test_client_requires_api_key(self):
        """Test client raises error without API key."""
        with pytest.raises(ValueError, match="API key is required"):
            JulesClient(api_key="")

    def test_client_context_manager(self):
        """Test client works as context manager."""
        with JulesClient(api_key="test-api-key") as client:
            assert client is not None

    @patch("jules_agent_sdk.base.BaseClient._request")
    def test_sessions_create(self, mock_request):
        """Test session creation."""
        mock_request.return_value = {
            "name": "sessions/test123",
            "id": "test123",
            "prompt": "Fix bug",
            "sourceContext": {"source": "sources/repo1"},
            "state": "QUEUED",
        }

        client = JulesClient(api_key="test-api-key")
        session = client.sessions.create(
            prompt="Fix bug", source="sources/repo1", starting_branch="main"
        )

        assert session.id == "test123"
        assert session.prompt == "Fix bug"
        mock_request.assert_called_once()

    @patch("jules_agent_sdk.base.BaseClient._request")
    def test_sessions_create_with_automation_mode(self, mock_request):
        """Test session creation with automation mode."""
        mock_request.return_value = {
            "name": "sessions/test456",
            "id": "test456",
            "prompt": "Auto PR feature",
            "sourceContext": {"source": "sources/repo2"},
            "state": "QUEUED",
            "automationMode": "AUTO_CREATE_PR",
        }

        client = JulesClient(api_key="test-api-key")
        session = client.sessions.create(
            prompt="Auto PR feature",
            source="sources/repo2",
            automation_mode="AUTO_CREATE_PR",
        )

        assert session.id == "test456"
        assert session.automation_mode == "AUTO_CREATE_PR"

        # Verify that automationMode was included in the request data
        mock_request.assert_called_once()
        _name, _args, kwargs = mock_request.mock_calls[0]
        assert "json" in kwargs
        assert kwargs["json"]["automationMode"] == "AUTO_CREATE_PR"

    @patch("jules_agent_sdk.base.BaseClient._request")
    def test_sessions_get(self, mock_request):
        """Test getting a session."""
        mock_request.return_value = {
            "name": "sessions/test123",
            "id": "test123",
            "prompt": "Fix bug",
            "sourceContext": {"source": "sources/repo1"},
            "state": "IN_PROGRESS",
        }

        client = JulesClient(api_key="test-api-key")
        session = client.sessions.get("test123")

        assert session.id == "test123"
        assert session.state.value == "IN_PROGRESS"

    @patch("jules_agent_sdk.base.BaseClient._request")
    def test_sessions_list(self, mock_request):
        """Test listing sessions."""
        mock_request.return_value = {
            "sessions": [
                {
                    "name": "sessions/test1",
                    "id": "test1",
                    "prompt": "Task 1",
                    "sourceContext": {"source": "sources/repo1"},
                },
                {
                    "name": "sessions/test2",
                    "id": "test2",
                    "prompt": "Task 2",
                    "sourceContext": {"source": "sources/repo2"},
                },
            ],
            "nextPageToken": "next-page",
        }

        client = JulesClient(api_key="test-api-key")
        result = client.sessions.list(page_size=10)

        assert len(result["sessions"]) == 2
        assert result["nextPageToken"] == "next-page"

    @patch("jules_agent_sdk.base.BaseClient._request")
    def test_activities_list(self, mock_request):
        """Test listing activities."""
        mock_request.return_value = {
            "activities": [
                {
                    "name": "sessions/s1/activities/a1",
                    "id": "a1",
                    "description": "Activity 1",
                },
                {
                    "name": "sessions/s1/activities/a2",
                    "id": "a2",
                    "description": "Activity 2",
                },
            ]
        }

        client = JulesClient(api_key="test-api-key")
        result = client.activities.list("s1")

        assert len(result["activities"]) == 2
        assert result["activities"][0].id == "a1"

    @patch("jules_agent_sdk.base.BaseClient._request")
    def test_sources_list(self, mock_request):
        """Test listing sources."""
        mock_request.return_value = {
            "sources": [
                {
                    "name": "sources/src1",
                    "id": "src1",
                    "githubRepo": {"owner": "test", "repo": "repo1"},
                }
            ]
        }

        client = JulesClient(api_key="test-api-key")
        result = client.sources.list()

        assert len(result["sources"]) == 1
        assert result["sources"][0].id == "src1"
        assert result["sources"][0].github_repo.owner == "test"


class TestErrorHandling:
    """Test error handling."""

    @patch("jules_agent_sdk.base.requests.Session.request")
    def test_authentication_error(self, mock_request):
        """Test authentication error handling."""
        mock_response = Mock()
        mock_response.ok = False
        mock_response.status_code = 401
        mock_response.json.return_value = {"error": {"message": "Invalid API key"}}
        mock_request.return_value = mock_response

        client = JulesClient(api_key="invalid-key")

        with pytest.raises(JulesAuthenticationError):
            client.sessions.list()

    @patch("jules_agent_sdk.base.requests.Session.request")
    def test_validation_error(self, mock_request):
        """Test validation error handling."""
        mock_response = Mock()
        mock_response.ok = False
        mock_response.status_code = 400
        mock_response.json.return_value = {"error": {"message": "Invalid request"}}
        mock_request.return_value = mock_response

        client = JulesClient(api_key="test-key")

        with pytest.raises(JulesValidationError):
            client.sessions.create(prompt="", source="")
