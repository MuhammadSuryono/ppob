package com.yonotech.ppob.mobile.data.repository;

import com.yonotech.ppob.mobile.data.remote.TransactionService;
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
public final class TransactionRepository_Factory implements Factory<TransactionRepository> {
  private final Provider<TransactionService> transactionServiceProvider;

  public TransactionRepository_Factory(Provider<TransactionService> transactionServiceProvider) {
    this.transactionServiceProvider = transactionServiceProvider;
  }

  @Override
  public TransactionRepository get() {
    return newInstance(transactionServiceProvider.get());
  }

  public static TransactionRepository_Factory create(
      Provider<TransactionService> transactionServiceProvider) {
    return new TransactionRepository_Factory(transactionServiceProvider);
  }

  public static TransactionRepository newInstance(TransactionService transactionService) {
    return new TransactionRepository(transactionService);
  }
}
