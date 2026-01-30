import "@typespec/http";
import "@typespec/rest";
import "./models.tsp";
import "../common/models.tsp";
import "../users/models.tsp";

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
    @body body: LoginResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Refresh access token using refresh token")
  @post
  @route("/refresh")
  @summary("Refresh token")
  refresh(): {
    @statusCode statusCode: 200;
    @body body: RefreshResponse;
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
    @body body: SuccessResponse;
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
    @body body: SuccessResponse;
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
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}

@tag("Password Recovery")
@route("/v1/password-recovery")
interface PasswordRecoveryOperations {
  @doc("Request a password recovery link")
  @post
  @route("/request")
  @summary("Request recovery")
  requestRecovery(@body request: RequestPasswordRecoveryRequest): {
    @statusCode statusCode: 200;
    @body body: PasswordRecoveryResponse;
  };

  @doc("Verify recovery token (POST)")
  @post
  @route("/verify")
  @summary("Verify token (POST)")
  verifyToken(@body request: VerifyPasswordRecoveryRequest): {
    @statusCode statusCode: 200;
    @body body: VerifyPasswordRecoveryResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };

  @doc("Verify recovery token (GET)")
  @get
  @route("/verify")
  @summary("Verify token (GET)")
  verifyTokenByQuery(@query token: string): {
    @statusCode statusCode: 200;
    @body body: VerifyPasswordRecoveryResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };

  @doc("Reset password using token")
  @post
  @route("/reset")
  @summary("Reset password")
  resetPassword(@body request: ResetPasswordRequest): {
    @statusCode statusCode: 200;
    @body body: PasswordRecoveryResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}
@tag("Mobile Authentication")
@route("/v1/auth/mobile")
interface MobileAuthOperations {
  @doc("Authenticate using Google OAuth2 for mobile")
  @post
  @route("/oauth2/google")
  @summary("Google Mobile Auth")
  authenticateWithGoogleMobile(@body request: MobileOAuth2Request): {
    @statusCode statusCode: 200;
    @body body: MobileLoginResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Refresh mobile access token")
  @post
  @route("/refresh")
  @summary("Mobile Refresh Token")
  refreshMobileToken(@body request: MobileRefreshRequest): {
    @statusCode statusCode: 200;
    @body body: MobileRefreshResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };
}
