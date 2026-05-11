package com.yonotech.ppob.mobile.viewmodels.staff;

import com.yonotech.ppob.mobile.data.repository.StaffRepository;
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
public final class StaffViewModel_Factory implements Factory<StaffViewModel> {
  private final Provider<StaffRepository> repositoryProvider;

  public StaffViewModel_Factory(Provider<StaffRepository> repositoryProvider) {
    this.repositoryProvider = repositoryProvider;
  }

  @Override
  public StaffViewModel get() {
    return newInstance(repositoryProvider.get());
  }

  public static StaffViewModel_Factory create(Provider<StaffRepository> repositoryProvider) {
    return new StaffViewModel_Factory(repositoryProvider);
  }

  public static StaffViewModel newInstance(StaffRepository repository) {
    return new StaffViewModel(repository);
  }
}
