package com.yonotech.ppob.mobile.viewmodels.wallet;

import com.yonotech.ppob.mobile.data.repository.WalletRepository;
import dagger.internal.DaggerGenerated;
import dagger.internal.Factory;
import dagger.internal.Provider;
import dagger.internal.QualifierMetadata;
import dagger.internal.ScopeMetadata;
import javax.annotation.processing.Generated;

@ScopeMetadata
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
public final class WalletViewModel_Factory implements Factory<WalletViewModel> {
  private final Provider<WalletRepository> repositoryProvider;

  public WalletViewModel_Factory(Provider<WalletRepository> repositoryProvider) {
    this.repositoryProvider = repositoryProvider;
  }

  @Override
  public WalletViewModel get() {
    return newInstance(repositoryProvider.get());
  }

  public static WalletViewModel_Factory create(Provider<WalletRepository> repositoryProvider) {
    return new WalletViewModel_Factory(repositoryProvider);
  }

  public static WalletViewModel newInstance(WalletRepository repository) {
    return new WalletViewModel(repository);
  }
}
