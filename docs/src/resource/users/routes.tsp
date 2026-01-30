import "@typespec/http";
import "@typespec/rest";
import "./models.tsp";
import "../common/models.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace GrowthAPI;

// Public User Operations
@tag("Users - Public")
@route("/v1/users/public")
interface PublicUserOperations {
  @doc("Sign up a new user (public)")
  @post
  @route("/signup")
  @summary("Public user signup")
  signupUser(@body request: SignupUserRequest): {
    @statusCode statusCode: 201;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };

  @doc("Resend verification email")
  @post
  @route("/resend-verification")
  @summary("Resend verification email")
  resendVerification(@body request: ResendVerificationRequest): {
    @statusCode statusCode: 200;
    @body body: SuccessResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}

// Authenticated User Operations
@tag("Users - Authenticated")
@route("/v1/users")
interface AuthenticatedUserOperations {
  @doc("Get current authenticated user")
  @get
  @route("/me")
  @summary("Get current user")
  getCurrentUser(@header Authorization?: string): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Update current user information")
  @put
  @route("/")
  @summary("Update user")
  updateUser(
    @header Authorization?: string,
    @body request: UpdateUserRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400 | 401;
    @body body: ErrorResponse;
  };

  @doc("Update current user password")
  @patch
  @route("/password")
  @summary("Update password")
  updatePassword(
    @header Authorization?: string,
    @body request: UpdatePasswordRequest
  ): {
    @statusCode statusCode: 200;
  } | {
    @statusCode statusCode: 400 | 401;
    @body body: ErrorResponse;
  };

  @doc("Add image to user profile")
  @patch
  @route("/add-image")
  @summary("Add user image")
  addImage(
    @header Authorization?: string,
    @body request: UploadImageRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UploadResponse;
  } | {
    @statusCode statusCode: 400 | 401;
    @body body: ErrorResponse;
  };
}

@tag("Uploads")
@route("/v1/uploads")
interface UploadOperations {
  @doc("Get presigned URL for image upload")
  @post
  @route("/images")
  @summary("Get upload URL")
  getUploadUrl(
    @header Authorization?: string,
    @body request: UploadImageRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UploadResponse;
  } | {
    @statusCode statusCode: 400 | 401 | 403;
    @body body: ErrorResponse;
  };
}

// Admin User Operations
@tag("Users - Admin")
@route("/v1/users")
interface AdminUserOperations {
  @doc("Create a new user (admin only)")
  @post
  @route("/")
  @summary("Create user (admin)")
  saveUser(
    @header Authorization?: string,
    @body request: CreateUserRequest
  ): {
    @statusCode statusCode: 201;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400 | 401 | 403;
    @body body: ErrorResponse;
  };

  @doc("Get all users (admin only)")
  @get
  @route("/")
  @summary("List all users (admin)")
  findAllUsers(
    @header Authorization?: string,
    @query page?: int32,
    @query limit?: int32
  ): {
    @statusCode statusCode: 200;
    @body body: UsersListResponse;
  } | {
    @statusCode statusCode: 401 | 403;
    @body body: ErrorResponse;
  };

  @doc("Get user by ID (admin only)")
  @get
  @route("/{id}")
  @summary("Get user by ID (admin)")
  findUserByID(
    @header Authorization?: string,
    @path id: string
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Get user by email (admin only)")
  @get
  @route("/email/{email}")
  @summary("Get user by email (admin)")
  findUserByEmail(
    @header Authorization?: string,
    @path email: string
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Update user by ID (admin only)")
  @put
  @route("/{id}")
  @summary("Update user (admin)")
  updateUserAdmin(
    @header Authorization?: string,
    @path id: string,
    @body request: UpdateUserRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400 | 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Delete user by ID (admin only)")
  @delete
  @route("/{id}")
  @summary("Delete user (admin)")
  deleteUserByID(
    @header Authorization?: string,
    @path id: string
  ): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Delete multiple users by IDs (admin only)")
  @delete
  @route("/")
  @summary("Delete users (admin)")
  deleteUsersByIDs(
    @header Authorization?: string,
    @query ids: string
  ): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400 | 401 | 403;
    @body body: ErrorResponse;
  };

  @doc("Toggle user status (admin only)")
  @patch
  @route("/{userId}/status")
  @summary("Toggle user status (admin)")
  toggleUserStatus(
    @header Authorization?: string,
    @path userId: string,
    @query active: boolean
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400 | 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Update user password (admin only)")
  @patch
  @route("/{id}/password")
  @summary("Update user password (admin)")
  updatePasswordAdmin(
    @header Authorization?: string,
    @path id: string,
    @body request: UpdatePasswordRequest
  ): {
    @statusCode statusCode: 200;
  } | {
    @statusCode statusCode: 400 | 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Update user access mode (admin only)")
  @patch
  @route("/{id}/access-mode")
  @summary("Update access mode (admin)")
  updateAccessMode(
    @header Authorization?: string,
    @path id: string,
    @body request: UpdateAccessModeRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400 | 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Update user features (admin only)")
  @patch
  @route("/{id}/features")
  @summary("Update features (admin)")
  updateFeatures(
    @header Authorization?: string,
    @path id: string,
    @body request: UpdateFeaturesRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400 | 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Update user limits (admin only)")
  @patch
  @route("/{id}/limits")
  @summary("Update limits (admin)")
  updateLimits(
    @header Authorization?: string,
    @path id: string,
    @body request: UpdateLimitsRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 400 | 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Grant lifetime pro to user (admin only)")
  @patch
  @route("/{id}/lifetime-pro")
  @summary("Grant lifetime pro (admin)")
  grantLifetimePro(
    @header Authorization?: string,
    @path id: string,
    @body request: GrantLifetimeProRequest
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Ensure user metadata exists (admin only)")
  @post
  @route("/{id}/ensure-metadata")
  @summary("Ensure metadata (admin)")
  ensureMetadata(
    @header Authorization?: string,
    @path id: string
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 401 | 403 | 404;
    @body body: ErrorResponse;
  };

  @doc("Revoke lifetime pro from user (admin only)")
  @delete
  @route("/{id}/lifetime-pro")
  @summary("Revoke lifetime pro (admin)")
  revokeLifetimePro(
    @header Authorization?: string,
    @path id: string
  ): {
    @statusCode statusCode: 200;
    @body body: UserResponse;
  } | {
    @statusCode statusCode: 401 | 403 | 404;
    @body body: ErrorResponse;
  };
}

// Email Verification Operations
@tag("Email Verification")
@route("/v1/email-verification")
interface EmailVerificationOperations {
  @doc("Send verification email to authenticated user")
  @post
  @route("/send")
  @summary("Send verification email")
  sendVerificationEmail(@header Authorization: string): {
    @statusCode statusCode: 200;
    @body body: EmailVerificationResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @doc("Verify email using token in body")
  @post
  @route("/verify")
  @summary("Verify email (POST)")
  verifyEmail(@body request: VerifyEmailRequest): {
    @statusCode statusCode: 200;
    @body body: EmailVerificationResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };

  @doc("Verify email using token in query")
  @get
  @route("/verify")
  @summary("Verify email (GET)")
  verifyEmailByQuery(@query token: string): {
    @statusCode statusCode: 200;
    @body body: string; // Returns HTML
  } | {
    @statusCode statusCode: 400;
    @body body: string; // Returns HTML
  };

  @doc("Resend verification email to specific address")
  @post
  @route("/resend")
  @summary("Resend verification email")
  resendVerificationEmail(@body request: ResendEmailVerificationRequest): {
    @statusCode statusCode: 200;
    @body body: EmailVerificationResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: ErrorResponse;
  };
}
