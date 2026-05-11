package com.yonotech.ppob.mobile.data.repository;

import com.yonotech.ppob.mobile.data.remote.StaffService;
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
public final class StaffRepository_Factory implements Factory<StaffRepository> {
  private final Provider<StaffService> staffServiceProvider;

  public StaffRepository_Factory(Provider<StaffService> staffServiceProvider) {
    this.staffServiceProvider = staffServiceProvider;
  }

  @Override
  public StaffRepository get() {
    return newInstance(staffServiceProvider.get());
  }

  public static StaffRepository_Factory create(Provider<StaffService> staffServiceProvider) {
    return new StaffRepository_Factory(staffServiceProvider);
  }

  public static StaffRepository newInstance(StaffService staffService) {
    return new StaffRepository(staffService);
  }
}
