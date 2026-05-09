// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'wallet.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class WalletAdapter extends TypeAdapter<Wallet> {
  @override
  final int typeId = 4;

  @override
  Wallet read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Wallet(
      id: fields[0] as String,
      userId: fields[1] as String,
      role: fields[2] as String,
      ownerName: fields[3] as String?,
      availableBalance: fields[4] as double,
      heldBalance: fields[5] as double,
      dailySpentAmount: fields[6] as double,
      dailyTransactionCount: fields[7] as int,
      date: fields[8] as DateTime,
      updatedAt: fields[9] as DateTime,
    );
  }

  @override
  void write(BinaryWriter writer, Wallet obj) {
    writer
      ..writeByte(10)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.userId)
      ..writeByte(2)
      ..write(obj.role)
      ..writeByte(3)
      ..write(obj.ownerName)
      ..writeByte(4)
      ..write(obj.availableBalance)
      ..writeByte(5)
      ..write(obj.heldBalance)
      ..writeByte(6)
      ..write(obj.dailySpentAmount)
      ..writeByte(7)
      ..write(obj.dailyTransactionCount)
      ..writeByte(8)
      ..write(obj.date)
      ..writeByte(9)
      ..write(obj.updatedAt);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is WalletAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
