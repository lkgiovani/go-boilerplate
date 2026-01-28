import "@typespec/http";

namespace GrowthAPI;

// User Models
model User {
  id: string;
  email: string;
  name: string;
  role: string;
  image_url?: string;
  status: string;
  access_mode?: string;
  features?: Record<string>;
  limits?: Record<string>;
  lifetime_pro?: boolean;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model CreateUserRequest {
  email: string;
  password: string;
  name: string;
  role?: string;
}

model UpdateUserRequest {
  name?: string;
  email?: string;
}

model UpdatePasswordRequest {
  current_password: string;
  new_password: string;
}

model UpdatePasswordAdminRequest {
  new_password: string;
}

model AddImageRequest {
  image_url: string;
}

model SignupUserRequest {
  email: string;
  password: string;
  name: string;
}

model ResendVerificationRequest {
  email: string;
}

model UpdateUserAdminRequest {
  name?: string;
  email?: string;
  role?: string;
  status?: string;
}

model ToggleUserStatusRequest {
  status: string;
}

model UpdateAccessModeRequest {
  access_mode: string;
}

model UpdateFeaturesRequest {
  features: Record<string>;
}

model UpdateLimitsRequest {
  limits: Record<string>;
}

model DeleteUsersByIDsRequest {
  ids: string[];
}

model UserResponse {
  user: User;
}

model UsersListResponse {
  users: User[];
  total: int64;
  page: int32;
  limit: int32;
}

// Email Verification Models
model VerifyEmailRequest {
  token: string;
}

model ResendEmailVerificationRequest {
  email: string;
}

model EmailVerificationResponse {
  success: boolean;
  message: string;
  userId?: int64;
  email?: string;
}
