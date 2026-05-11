-keepnames class com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
-if class com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
-keep class com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponseJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
-if class com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
-keepnames class kotlin.jvm.internal.DefaultConstructorMarker
-keepclassmembers class com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse {
    public synthetic <init>(java.lang.String,java.lang.String,java.lang.String,java.lang.String,int,kotlin.jvm.internal.DefaultConstructorMarker);
}
