package com.yonotech.ppob.mobile.data.repository;

import com.yonotech.ppob.mobile.data.remote.AuthService;
import dagger.internal.DaggerGenerated;
import dagger.internal.Factory;
import dagger.internal.Provider;
import dagger.internal.QualifierMetadata;
import dagger.internal.ScopeMetadata;
import javax.annotation.processing.Generated;

@ScopeMetadata("javax.inject.Singleton")
@QualifierMetadata
@DaggerGenerated
@Generated(
    value = "dagger.internal.codegen.ComponentProcessor",
    comments = "https://dagger.dev"
)
@SuppressWarnings({
    "unchecked",
    "rawtypes",
    "KotlinInternal",
    "KotlinInternalInJava",
    "cast",
    "deprecation",
    "nullness:initialization.field.uninitialized"
})
public final class AuthRepository_Factory implements Factory<AuthRepository> {
  private final Provider<AuthService> authServiceProvider;

  public AuthRepository_Factory(Provider<AuthService> authServiceProvider) {
    this.authServiceProvider = authServiceProvider;
  }

  @Override
  public AuthRepository get() {
    return newInstance(authServiceProvider.get());
  }

  public static AuthRepository_Factory create(Provider<AuthService> authServiceProvider) {
    return new AuthRepository_Factory(authServiceProvider);
  }

  public static AuthRepository newInstance(AuthService authService) {
    return new AuthRepository(authService);
  }
}
