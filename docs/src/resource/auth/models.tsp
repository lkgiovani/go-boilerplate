import "@typespec/http";

namespace GrowthAPI;

// Auth Models
model LoginRequest {
  email: string;
  password: string;
}

model SignupRequest {
  email: string;
  password: string;
  name: string;
}

model LoginResponse {
  userId: int64;
  email: string;
  accessToken: string;
  expiresIn: int64;
}

model RefreshResponse {
  accessToken: string;
  expiresIn: int64;
}

model UserInfo {
  id: int64;
  email: string;
  name: string;
  admin: boolean;
}

model LogoutResponse {
  message: string;
}

// Password Recovery Models
model RequestPasswordRecoveryRequest {
  email: string;
}

model VerifyPasswordRecoveryRequest {
  token: string;
}

model ResetPasswordRequest {
  token: string;
  password: string;
}

model PasswordRecoveryResponse {
  message: string;
}

model VerifyPasswordRecoveryResponse {
  valid: boolean;
}
// Mobile Auth Models
model MobileOAuth2Request {
  idToken: string;
  deviceId?: string;
}

model MobileLoginResponse {
  userId: int64;
  email: string;
  name: string;
  accessToken: string;
  refreshToken: string;
  expiresIn: int64;
  isNewUser: boolean;
}

model MobileRefreshRequest {
  refreshToken: string;
}

model MobileRefreshResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: int64;
}
