-keepnames class com.yonotech.ppob.mobile.data.remote.dto.transaction.InitiateTransactionRequest
-if class com.yonotech.ppob.mobile.data.remote.dto.transaction.InitiateTransactionRequest
-keep class com.yonotech.ppob.mobile.data.remote.dto.transaction.InitiateTransactionRequestJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
