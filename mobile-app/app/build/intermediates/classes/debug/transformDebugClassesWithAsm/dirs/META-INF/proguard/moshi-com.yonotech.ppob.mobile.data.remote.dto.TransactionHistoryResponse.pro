-keepnames class com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse
-if class com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse
-keep class com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponseJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
-if class com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse
-keepnames class kotlin.jvm.internal.DefaultConstructorMarker
-keepclassmembers class com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse {
    public synthetic <init>(java.lang.String,java.lang.String,java.lang.String,double,double,java.lang.String,java.lang.String,java.lang.String,int,kotlin.jvm.internal.DefaultConstructorMarker);
}
