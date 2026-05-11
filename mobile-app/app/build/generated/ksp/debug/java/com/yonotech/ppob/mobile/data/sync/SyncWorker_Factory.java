package com.yonotech.ppob.mobile.data.sync;

import android.content.Context;
import androidx.work.WorkerParameters;
import com.yonotech.ppob.mobile.data.local.dao.PendingSyncDao;
import com.yonotech.ppob.mobile.data.remote.TransactionService;
import dagger.internal.DaggerGenerated;
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
public final class SyncWorker_Factory {
  private final Provider<PendingSyncDao> pendingSyncDaoProvider;

  private final Provider<TransactionService> transactionServiceProvider;

  public SyncWorker_Factory(Provider<PendingSyncDao> pendingSyncDaoProvider,
      Provider<TransactionService> transactionServiceProvider) {
    this.pendingSyncDaoProvider = pendingSyncDaoProvider;
    this.transactionServiceProvider = transactionServiceProvider;
  }

  public SyncWorker get(Context ctx, WorkerParameters params) {
    return newInstance(ctx, params, pendingSyncDaoProvider.get(), transactionServiceProvider.get());
  }

  public static SyncWorker_Factory create(Provider<PendingSyncDao> pendingSyncDaoProvider,
      Provider<TransactionService> transactionServiceProvider) {
    return new SyncWorker_Factory(pendingSyncDaoProvider, transactionServiceProvider);
  }

  public static SyncWorker newInstance(Context ctx, WorkerParameters params,
      PendingSyncDao pendingSyncDao, TransactionService transactionService) {
    return new SyncWorker(ctx, params, pendingSyncDao, transactionService);
  }
}
