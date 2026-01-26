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

model RefreshRequest {
  refresh_token?: string;
}

model AuthResponse {
  access_token: string;
  refresh_token: string;
  user: UserInfo;
}

model UserInfo {
  id: string;
  email: string;
  name: string;
  role: string;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model LogoutResponse {
  message: string;
}
