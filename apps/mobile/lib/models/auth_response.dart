/// Auth service response models

class AuthResponse {
  final String accessToken;
  final String refreshToken;
  final String tokenType;
  final int expiresIn;
  final UserResponse user;

  AuthResponse({
    required this.accessToken,
    required this.refreshToken,
    required this.tokenType,
    required this.expiresIn,
    required this.user,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
      tokenType: json['token_type'] as String? ?? 'Bearer',
      expiresIn: json['expires_in'] as int? ?? 3600,
      user: UserResponse.fromJson(json['user'] as Map<String, dynamic>),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'access_token': accessToken,
      'refresh_token': refreshToken,
      'token_type': tokenType,
      'expires_in': expiresIn,
      'user': user.toJson(),
    };
  }
}

class LoginRequest {
  final String username; // email or phone
  final String password;

  LoginRequest({required this.username, required this.password});

  Map<String, dynamic> toJson() {
    return {
      'username': username,
      'password': password,
    };
  }
}

class RegisterRequest {
  final String email;
  final String phone;
  final String name;
  final String password;
  final String pin;
  final String? referralCode;

  RegisterRequest({
    required this.email,
    required this.phone,
    required this.name,
    required this.password,
    required this.pin,
    this.referralCode,
  });

  Map<String, dynamic> toJson() {
    return {
      'email': email,
      'phone': phone,
      'name': name,
      'password': password,
      'pin': pin,
      if (referralCode != null) 'referral_code': referralCode,
    };
  }
}

class VerifyOtpRequest {
  final String phoneNumber;
  final String otp;

  VerifyOtpRequest({required this.phoneNumber, required this.otp});

  Map<String, dynamic> toJson() {
    return {
      'phone_number': phoneNumber,
      'otp': otp,
    };
  }
}

class RefreshTokenRequest {
  final String refreshToken;

  RefreshTokenRequest({required this.refreshToken});

  Map<String, dynamic> toJson() {
    return {
      'refresh_token': refreshToken,
    };
  }
}

class ChangePasswordRequest {
  final String oldPassword;
  final String newPassword;

  ChangePasswordRequest({required this.oldPassword, required this.newPassword});

  Map<String, dynamic> toJson() {
    return {
      'old_password': oldPassword,
      'new_password': newPassword,
    };
  }
}

class ChangePinRequest {
  final String oldPin;
  final String newPin;

  ChangePinRequest({required this.oldPin, required this.newPin});

  Map<String, dynamic> toJson() {
    return {
      'old_pin': oldPin,
      'new_pin': newPin,
    };
  }
}

/// User response from backend
class UserResponse {
  final String id;
  final String email;
  final String phone;
  final String name;
  final String? avatarUrl;
  final String role;
  final int trustScore;
  final String? address;
  final DateTime? dateOfBirth;
  final DateTime createdAt;
  final DateTime updatedAt;

  UserResponse({
    required this.id,
    required this.email,
    required this.phone,
    required this.name,
    this.avatarUrl,
    required this.role,
    this.trustScore = 50,
    this.address,
    this.dateOfBirth,
    required this.createdAt,
    required this.updatedAt,
  });

  factory UserResponse.fromJson(Map<String, dynamic> json) {
    return UserResponse(
      id: json['id'] as String,
      email: json['email'] as String? ?? '',
      phone: json['phone'] as String? ?? '',
      name: json['name'] as String,
      avatarUrl: json['avatar_url'] as String?,
      role: json['role'] as String,
      trustScore: json['trust_score'] as int? ?? 50,
      address: json['address'] as String?,
      dateOfBirth: json['date_of_birth'] != null
          ? DateTime.parse(json['date_of_birth'] as String)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'phone': phone,
      'name': name,
      'avatar_url': avatarUrl,
      'role': role,
      'trust_score': trustScore,
      'address': address,
      'date_of_birth': dateOfBirth?.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}
