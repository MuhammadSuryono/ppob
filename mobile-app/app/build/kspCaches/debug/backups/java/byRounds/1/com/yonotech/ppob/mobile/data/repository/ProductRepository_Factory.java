package com.yonotech.ppob.mobile.data.repository;

import com.yonotech.ppob.mobile.data.remote.ProductService;
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
public final class ProductRepository_Factory implements Factory<ProductRepository> {
  private final Provider<ProductService> productServiceProvider;

  public ProductRepository_Factory(Provider<ProductService> productServiceProvider) {
    this.productServiceProvider = productServiceProvider;
  }

  @Override
  public ProductRepository get() {
    return newInstance(productServiceProvider.get());
  }

  public static ProductRepository_Factory create(Provider<ProductService> productServiceProvider) {
    return new ProductRepository_Factory(productServiceProvider);
  }

  public static ProductRepository newInstance(ProductService productService) {
    return new ProductRepository(productService);
  }
}
