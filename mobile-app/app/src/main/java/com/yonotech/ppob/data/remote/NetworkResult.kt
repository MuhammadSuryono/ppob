package com.yonotech.ppob.data.remote

sealed class NetworkResult<T>(
    val data: T? = null,
    val error: Throwable? = null
) {
    class Success<T>(data: T) : NetworkResult<T>(data)
    class Error<T>(error: Throwable, data: T? = null) : NetworkResult<T>(data, error)
    class Loading<T>(data: T? = null) : NetworkResult<T>(data)
}

// Convenience extension functions
fun <T> NetworkResult<T>.isSuccess(): Boolean = this is NetworkResult.Success && data != null
fun <T> NetworkResult<T>.getOrNull(): T? = (this as? NetworkResult.Success)?.data