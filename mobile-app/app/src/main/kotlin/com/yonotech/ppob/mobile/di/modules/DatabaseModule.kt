package com.yonotech.ppob.mobile.di.modules

import android.content.Context
import androidx.room.Room
import com.yonotech.ppob.mobile.data.local.PpoDatabase
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): PpoDatabase {
        return Room.databaseBuilder(
            context,
            PpoDatabase::class.java,
            "ppob_database"
        ).build()
    }
}
