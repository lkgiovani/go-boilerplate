import "@typespec/http";
import "@typespec/rest";
import "./models.tsp";
import "../common/models.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace GrowthAPI;

@tag("Authentication")
@route("/v1/auth")
interface AuthOperations {
  @doc("Login with email and password")
  @post
  @route("/login")
  @summary("User login")
  login(@body request: LoginRequest): {
    @statusCode statusCode: 200;
    @body body: AuthResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Refresh access token using refresh token")
  @post
  @route("/refresh")
  @summary("Refresh token")
  refresh(@body request: RefreshRequest): {
    @statusCode statusCode: 200;
    @body body: AuthResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Logout current session")
  @post
  @route("/logout")
  @summary("Logout")
  logout(@header Authorization?: string): {
    @statusCode statusCode: 200;
    @body body: LogoutResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Logout from all sessions (requires authentication)")
  @post
  @route("/logout-all")
  @summary("Logout from all sessions")
  logoutAll(@header Authorization?: string): {
    @statusCode statusCode: 200;
    @body body: LogoutResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Sign up a new user")
  @post
  @route("/signup")
  @summary("User signup")
  signup(@body request: SignupRequest): {
    @statusCode statusCode: 201;
    @body body: AuthResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}
