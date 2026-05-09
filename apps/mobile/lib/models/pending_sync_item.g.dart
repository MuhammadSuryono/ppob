// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'pending_sync_item.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class PendingSyncItemAdapter extends TypeAdapter<PendingSyncItem> {
  @override
  final int typeId = 5;

  @override
  PendingSyncItem read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return PendingSyncItem(
      id: fields[0] as String,
      type: fields[1] as String,
      data: (fields[2] as Map).cast<String, dynamic>(),
      retryCount: fields[3] as int,
      createdAt: fields[4] as DateTime,
      lastAttemptAt: fields[5] as DateTime?,
      errorMessage: fields[6] as String?,
    );
  }

  @override
  void write(BinaryWriter writer, PendingSyncItem obj) {
    writer
      ..writeByte(7)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.type)
      ..writeByte(2)
      ..write(obj.data)
      ..writeByte(3)
      ..write(obj.retryCount)
      ..writeByte(4)
      ..write(obj.createdAt)
      ..writeByte(5)
      ..write(obj.lastAttemptAt)
      ..writeByte(6)
      ..write(obj.errorMessage);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is PendingSyncItemAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
