package com.yonotech.ppob.mobile.data.repository;

import com.yonotech.ppob.mobile.data.remote.WalletService;
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
public final class WalletRepository_Factory implements Factory<WalletRepository> {
  private final Provider<WalletService> walletServiceProvider;

  public WalletRepository_Factory(Provider<WalletService> walletServiceProvider) {
    this.walletServiceProvider = walletServiceProvider;
  }

  @Override
  public WalletRepository get() {
    return newInstance(walletServiceProvider.get());
  }

  public static WalletRepository_Factory create(Provider<WalletService> walletServiceProvider) {
    return new WalletRepository_Factory(walletServiceProvider);
  }

  public static WalletRepository newInstance(WalletService walletService) {
    return new WalletRepository(walletService);
  }
}
