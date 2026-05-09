// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'staff.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class StaffAdapter extends TypeAdapter<Staff> {
  @override
  final int typeId = 3;

  @override
  Staff read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Staff(
      id: fields[0] as String,
      mitraId: fields[1] as String,
      name: fields[2] as String,
      phoneNumber: fields[3] as String,
      pinHash: fields[4] as String?,
      marginScheme: fields[5] as String,
      fixedAllowanceAmount: fields[6] as double?,
      marginSharePercentage: fields[7] as double?,
      dailyLimitAmount: fields[8] as double,
      dailyLimitCount: fields[9] as int,
      isActive: fields[10] as bool,
      createdAt: fields[11] as DateTime,
      updatedAt: fields[12] as DateTime,
    );
  }

  @override
  void write(BinaryWriter writer, Staff obj) {
    writer
      ..writeByte(13)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.mitraId)
      ..writeByte(2)
      ..write(obj.name)
      ..writeByte(3)
      ..write(obj.phoneNumber)
      ..writeByte(4)
      ..write(obj.pinHash)
      ..writeByte(5)
      ..write(obj.marginScheme)
      ..writeByte(6)
      ..write(obj.fixedAllowanceAmount)
      ..writeByte(7)
      ..write(obj.marginSharePercentage)
      ..writeByte(8)
      ..write(obj.dailyLimitAmount)
      ..writeByte(9)
      ..write(obj.dailyLimitCount)
      ..writeByte(10)
      ..write(obj.isActive)
      ..writeByte(11)
      ..write(obj.createdAt)
      ..writeByte(12)
      ..write(obj.updatedAt);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is StaffAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
