import "@typespec/http";

namespace GrowthAPI;

// User Models
model User {
  id: int64;
  email: string;
  name: string;
  imgUrl?: string;
  admin: boolean;
  active: boolean;
  source: string;
  metadata: UserMetadata;
  lastAccess?: utcDateTime;
  createdAt: utcDateTime;
  updatedAt: utcDateTime;
}

model UserMetadata {
  accessMode: string;
  planType: string;
  planExpirationDate?: utcDateTime;
  proSource?: string;
  maxResources?: int32;
  maxRequestsPerMonth?: int32;
  maxAccounts: int32;
  maxCategoriesPerAccount: int32;
  maxTransactionsPerMonth: int32;
  canExportData: boolean;
  canUseReports: boolean;
  canUseAdvancedFeatures: boolean;
  canCreateBudgets: boolean;
  canUseGoals: boolean;
  emailVerified: boolean;
  reputationStatus: string;
  suspiciousActivityCount: int32;
  lastSecurityCheck?: utcDateTime;
  lastPermissionCheck?: utcDateTime;
  notes?: string;
  locale: string;
  currency: string;
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
  imgUrl?: string;
}

model UpdatePasswordRequest {
  currentPassword?: string;
  password: string;
}



model SignupUserRequest {
  email: string;
  password: string;
  name: string;
}

model ResendVerificationRequest {
  email: string;
}



model UpdateAccessModeRequest {
  accessMode: string;
}

model UpdateFeaturesRequest {
  canCreateBudgets?: boolean;
  canExportData?: boolean;
  canUseReports?: boolean;
  canUseGoals?: boolean;
}

model UpdateLimitsRequest {
  maxAccounts?: int32;
  maxTransactionsPerMonth?: int32;
  maxCategoriesPerAccount?: int32;
}

model GrantLifetimeProRequest {
  reason: string;
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

// Upload Models
model UploadImageRequest {
  fileName: string;
  contentType: string;
  contentLength: int64;
}

model UploadResponse {
  uploadSignedUrl: string;
  publicUrl: string;
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
