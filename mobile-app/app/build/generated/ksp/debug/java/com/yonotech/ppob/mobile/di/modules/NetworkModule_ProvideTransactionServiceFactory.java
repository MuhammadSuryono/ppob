package com.yonotech.ppob.mobile.di.modules;

import com.yonotech.ppob.mobile.data.remote.TransactionService;
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
public final class NetworkModule_ProvideTransactionServiceFactory implements Factory<TransactionService> {
  private final Provider<Retrofit> retrofitProvider;

  public NetworkModule_ProvideTransactionServiceFactory(Provider<Retrofit> retrofitProvider) {
    this.retrofitProvider = retrofitProvider;
  }

  @Override
  public TransactionService get() {
    return provideTransactionService(retrofitProvider.get());
  }

  public static NetworkModule_ProvideTransactionServiceFactory create(
      Provider<Retrofit> retrofitProvider) {
    return new NetworkModule_ProvideTransactionServiceFactory(retrofitProvider);
  }

  public static TransactionService provideTransactionService(Retrofit retrofit) {
    return Preconditions.checkNotNullFromProvides(NetworkModule.INSTANCE.provideTransactionService(retrofit));
  }
}
