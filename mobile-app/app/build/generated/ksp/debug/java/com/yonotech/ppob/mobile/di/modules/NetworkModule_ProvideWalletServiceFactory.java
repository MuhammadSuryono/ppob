package com.yonotech.ppob.mobile.di.modules;

import com.yonotech.ppob.mobile.data.remote.WalletService;
import dagger.internal.DaggerGenerated;
import dagger.internal.Factory;
import dagger.internal.Preconditions;
import dagger.internal.Provider;
import dagger.internal.QualifierMetadata;
import dagger.internal.ScopeMetadata;
import javax.annotation.processing.Generated;
import retrofit2.Retrofit;

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
public final class NetworkModule_ProvideWalletServiceFactory implements Factory<WalletService> {
  private final Provider<Retrofit> retrofitProvider;

  public NetworkModule_ProvideWalletServiceFactory(Provider<Retrofit> retrofitProvider) {
    this.retrofitProvider = retrofitProvider;
  }

  @Override
  public WalletService get() {
    return provideWalletService(retrofitProvider.get());
  }

  public static NetworkModule_ProvideWalletServiceFactory create(
      Provider<Retrofit> retrofitProvider) {
    return new NetworkModule_ProvideWalletServiceFactory(retrofitProvider);
  }

  public static WalletService provideWalletService(Retrofit retrofit) {
    return Preconditions.checkNotNullFromProvides(NetworkModule.INSTANCE.provideWalletService(retrofit));
  }
}
