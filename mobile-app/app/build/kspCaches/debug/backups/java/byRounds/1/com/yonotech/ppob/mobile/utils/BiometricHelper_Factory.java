package com.yonotech.ppob.mobile.utils;

import dagger.internal.DaggerGenerated;
import dagger.internal.Factory;
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
public final class BiometricHelper_Factory implements Factory<BiometricHelper> {
  @Override
  public BiometricHelper get() {
    return newInstance();
  }

  public static BiometricHelper_Factory create() {
    return InstanceHolder.INSTANCE;
  }

  public static BiometricHelper newInstance() {
    return new BiometricHelper();
  }

  private static final class InstanceHolder {
    static final BiometricHelper_Factory INSTANCE = new BiometricHelper_Factory();
  }
}
