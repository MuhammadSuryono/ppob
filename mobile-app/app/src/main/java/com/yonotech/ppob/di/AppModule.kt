package com.yonotech.ppob.di

import android.content.Context
import com.yonotech.ppob.BuildConfig
import com.yonotech.ppob.data.local.database.PpoDatabase
import com.yonotech.ppob.data.remote.api.AuthApiService
import com.yonotech.ppob.data.remote.api.ProductApiService
import com.yonotech.ppob.data.remote.api.TransactionApiService
import com.yonotech.ppob.data.remote.api.UserApiService
import com.yonotech.ppob.data.remote.api.WalletApiService
import com.yonotech.ppob.data.remote.interceptor.AuthInterceptor
import com.yonotech.ppob.data.remote.interceptor.ErrorInterceptor
import com.yonotech.ppob.data.remote.interceptor.RetryInterceptor
import com.yonotech.ppob.data.repository.*
import com.yonotech.ppob.domain.repository.AuthRepository
import com.yonotech.ppob.domain.repository.ProductRepository
import com.yonotech.ppob.domain.repository.TransactionRepository
import com.yonotech.ppob.domain.repository.UserRepository
import com.yonotech.ppob.domain.repository.WalletRepository
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory
import java.util.concurrent.TimeUnit
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object NetworkModule {

    private const val BASE_URL = "https://fedora.sinauplatform.id/api/v1/"
    private const val CONNECT_TIMEOUT = 30L
    private const val READ_TIMEOUT = 30L
    private const val WRITE_TIMEOUT = 30L

    @Provides
    @Singleton
    fun provideOkHttpClient(
        authInterceptor: AuthInterceptor,
        retryInterceptor: RetryInterceptor,
        loggingInterceptor: HttpLoggingInterceptor,
        errorInterceptor: ErrorInterceptor
    ): OkHttpClient {
        return OkHttpClient.Builder()
            .connectTimeout(CONNECT_TIMEOUT, TimeUnit.SECONDS)
            .readTimeout(READ_TIMEOUT, TimeUnit.SECONDS)
            .writeTimeout(WRITE_TIMEOUT, TimeUnit.SECONDS)
            .addInterceptor(authInterceptor)
            .addInterceptor(retryInterceptor)
            .addInterceptor(errorInterceptor)
            .apply {
                if (BuildConfig.DEBUG) {
                    addInterceptor(loggingInterceptor)
                }
            }
            .build()
    }

    @Provides
    @Singleton
    fun provideRetrofit(okHttpClient: OkHttpClient): Retrofit {
        return Retrofit.Builder()
            .baseUrl(BASE_URL)
            .client(okHttpClient)
            .addConverterFactory(MoshiConverterFactory.create())
            .build()
    }

    @Provides
    @Singleton
    fun provideAuthApiService(retrofit: Retrofit): AuthApiService {
        return retrofit.create(AuthApiService::class.java)
    }

    @Provides
    @Singleton
    fun provideUserApiService(retrofit: Retrofit): UserApiService {
        return retrofit.create(UserApiService::class.java)
    }

    @Provides
    @Singleton
    fun provideWalletApiService(retrofit: Retrofit): WalletApiService {
        return retrofit.create(WalletApiService::class.java)
    }

    @Provides
    @Singleton
    fun provideProductApiService(retrofit: Retrofit): ProductApiService {
        return retrofit.create(ProductApiService::class.java)
    }

    @Provides
    @Singleton
    fun provideTransactionApiService(retrofit: Retrofit): TransactionApiService {
        return retrofit.create(TransactionApiService::class.java)
    }
}

@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): PpoDatabase {
        return PpoDatabase.getInstance(context)
    }

    @Provides
    fun provideUserDao(database: PpoDatabase) = database.userDao()

    @Provides
    fun provideWalletDao(database: PpoDatabase) = database.walletDao()

    @Provides
    fun provideTransactionDao(database: PpoDatabase) = database.transactionDao()

    @Provides
    fun provideProductDao(database: PpoDatabase) = database.productDao()

    @Provides
    fun provideCategoryDao(database: PpoDatabase) = database.categoryDao()

    @Provides
    fun providePendingSyncDao(database: PpoDatabase) = database.pendingSyncDao()
}

@Module
@InstallIn(SingletonComponent::class)
object RepositoryModule {

    @Provides
    @Singleton
    fun provideAuthRepository(
        authApiService: AuthApiService,
        userDao: com.yonotech.ppob.data.local.dao.UserDao,
        authPreferences: com.yonotech.ppob.data.local.datastore.AuthPreferencesDataStore
    ): AuthRepository {
        return AuthRepositoryImpl(authApiService, userDao, authPreferences)
    }

    @Provides
    @Singleton
    fun provideUserRepository(
        userApiService: UserApiService,
        userDao: com.yonotech.ppob.data.local.dao.UserDao
    ): UserRepository {
        return UserRepositoryImpl(userApiService, userDao)
    }

    @Provides
    @Singleton
    fun provideWalletRepository(
        walletApiService: WalletApiService,
        walletDao: com.yonotech.ppob.data.local.dao.WalletDao
    ): WalletRepository {
        return WalletRepositoryImpl(walletApiService, walletDao)
    }

    @Provides
    @Singleton
    fun provideProductRepository(
        productApiService: ProductApiService,
        productDao: com.yonotech.ppob.data.local.dao.ProductDao,
        categoryDao: com.yonotech.ppob.data.local.dao.CategoryDao
    ): ProductRepository {
        return ProductRepositoryImpl(productApiService, productDao, categoryDao)
    }

    @Provides
    @Singleton
    fun provideTransactionRepository(
        transactionApiService: TransactionApiService,
        transactionDao: com.yonotech.ppob.data.local.dao.TransactionDao,
        pendingSyncDao: com.yonotech.ppob.data.local.dao.PendingSyncDao
    ): TransactionRepository {
        return TransactionRepositoryImpl(transactionApiService, transactionDao, pendingSyncDao)
    }
}

@Module
@InstallIn(SingletonComponent::class)
object ViewModelModule {
    // ViewModels are provided automatically via @HiltViewModel
}
