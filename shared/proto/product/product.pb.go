package product

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetInquiryProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brand      string `protobuf:"bytes,1,opt,name=brand,proto3" json:"brand,omitempty"`
	CategoryId uint32 `protobuf:"varint,2,opt,name=category_id,json=categoryId,proto3" json:"category_id,omitempty"`
}

func (x *GetInquiryProductRequest) Reset() { *x = GetInquiryProductRequest{} }
func (x *GetInquiryProductRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetInquiryProductRequest) ProtoMessage() {}
func (x *GetInquiryProductRequest) ProtoReflect() protoreflect.Message { return nil }
func (x *GetInquiryProductRequest) GetBrand() string { if x != nil { return x.Brand }; return "" }
func (x *GetInquiryProductRequest) GetCategoryId() uint32 { if x != nil { return x.CategoryId }; return 0 }

type GetProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId uint32 `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	SkuCode   string `protobuf:"bytes,2,opt,name=sku_code,json=skuCode,proto3" json:"sku_code,omitempty"`
}

func (x *GetProductRequest) Reset() { *x = GetProductRequest{} }
func (x *GetProductRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetProductRequest) ProtoMessage() {}
func (x *GetProductRequest) ProtoReflect() protoreflect.Message { return nil }
func (x *GetProductRequest) GetProductId() uint32 { if x != nil { return x.ProductId }; return 0 }
func (x *GetProductRequest) GetSkuCode() string { if x != nil { return x.SkuCode }; return "" }

type GetProductResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SkuCode      string  `protobuf:"bytes,3,opt,name=sku_code,json=skuCode,proto3" json:"sku_code,omitempty"`
	Price        float64 `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	IsActive     bool    `protobuf:"varint,5,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	CategoryName string  `protobuf:"bytes,6,opt,name=category_name,json=categoryName,proto3" json:"category_name,omitempty"`
}

func (x *GetProductResponse) Reset() { *x = GetProductResponse{} }
func (x *GetProductResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetProductResponse) ProtoMessage() {}
func (x *GetProductResponse) ProtoReflect() protoreflect.Message { return nil }
func (x *GetProductResponse) GetId() uint32 { if x != nil { return x.Id }; return 0 }
func (x *GetProductResponse) GetName() string { if x != nil { return x.Name }; return "" }
func (x *GetProductResponse) GetSkuCode() string { if x != nil { return x.SkuCode }; return "" }
func (x *GetProductResponse) GetPrice() float64 { if x != nil { return x.Price }; return 0 }
func (x *GetProductResponse) GetIsActive() bool { if x != nil { return x.IsActive }; return false }

type ValidateProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId     uint32  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	ExpectedPrice float64 `protobuf:"fixed64,2,opt,name=expected_price,json=expectedPrice,proto3" json:"expected_price,omitempty"`
}

func (x *ValidateProductRequest) Reset() { *x = ValidateProductRequest{} }
func (x *ValidateProductRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*ValidateProductRequest) ProtoMessage() {}
func (x *ValidateProductRequest) ProtoReflect() protoreflect.Message { return nil }

type ValidateProductResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsValid      bool    `protobuf:"varint,1,opt,name=is_valid,json=isValid,proto3" json:"is_valid,omitempty"`
	Message      string  `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	CurrentPrice float64 `protobuf:"fixed64,3,opt,name=current_price,json=currentPrice,proto3" json:"current_price,omitempty"`
}

func (x *ValidateProductResponse) Reset() { *x = ValidateProductResponse{} }
func (x *ValidateProductResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*ValidateProductResponse) ProtoMessage() {}
func (x *ValidateProductResponse) ProtoReflect() protoreflect.Message { return nil }

var file_shared_proto_product_product_proto_rawDesc = []byte{}
var file_shared_proto_product_product_proto_goTypes = []interface{}{
	(*GetProductRequest)(nil),       // 0
	(*GetProductResponse)(nil),      // 1
	(*GetInquiryProductRequest)(nil), // 2
	(*ValidateProductRequest)(nil),  // 3
	(*ValidateProductResponse)(nil), // 4
}
