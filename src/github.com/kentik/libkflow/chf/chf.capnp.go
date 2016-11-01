package chf

// AUTO GENERATED - DO NOT EDIT

import (
	math "math"
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type Custom struct{ capnp.Struct }
type Custom_value Custom
type Custom_value_Which uint16

const (
	Custom_value_Which_uint32Val  Custom_value_Which = 0
	Custom_value_Which_float32Val Custom_value_Which = 1
	Custom_value_Which_strVal     Custom_value_Which = 2
)

func (w Custom_value_Which) String() string {
	const s = "uint32Valfloat32ValstrVal"
	switch w {
	case Custom_value_Which_uint32Val:
		return s[0:9]
	case Custom_value_Which_float32Val:
		return s[9:19]
	case Custom_value_Which_strVal:
		return s[19:25]

	}
	return "Custom_value_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

// Custom_TypeID is the unique identifier for the type Custom.
const Custom_TypeID = 0xed5d37861203d027

func NewCustom(s *capnp.Segment) (Custom, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return Custom{st}, err
}

func NewRootCustom(s *capnp.Segment) (Custom, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return Custom{st}, err
}

func ReadRootCustom(msg *capnp.Message) (Custom, error) {
	root, err := msg.RootPtr()
	return Custom{root.Struct()}, err
}

func (s Custom) String() string {
	str, _ := text.Marshal(0xed5d37861203d027, s.Struct)
	return str
}

func (s Custom) Id() uint32 {
	return s.Struct.Uint32(0)
}

func (s Custom) SetId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Custom) Value() Custom_value { return Custom_value(s) }

func (s Custom_value) Which() Custom_value_Which {
	return Custom_value_Which(s.Struct.Uint16(8))
}
func (s Custom_value) Uint32Val() uint32 {
	return s.Struct.Uint32(4)
}

func (s Custom_value) SetUint32Val(v uint32) {
	s.Struct.SetUint16(8, 0)
	s.Struct.SetUint32(4, v)
}

func (s Custom_value) Float32Val() float32 {
	return math.Float32frombits(s.Struct.Uint32(4))
}

func (s Custom_value) SetFloat32Val(v float32) {
	s.Struct.SetUint16(8, 1)
	s.Struct.SetUint32(4, math.Float32bits(v))
}

func (s Custom_value) StrVal() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Custom_value) HasStrVal() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Custom_value) StrValBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Custom_value) SetStrVal(v string) error {
	s.Struct.SetUint16(8, 2)
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

// Custom_List is a list of Custom.
type Custom_List struct{ capnp.List }

// NewCustom creates a new list of Custom.
func NewCustom_List(s *capnp.Segment, sz int32) (Custom_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	return Custom_List{l}, err
}

func (s Custom_List) At(i int) Custom { return Custom{s.List.Struct(i)} }

func (s Custom_List) Set(i int, v Custom) error { return s.List.SetStruct(i, v.Struct) }

// Custom_Promise is a wrapper for a Custom promised by a client call.
type Custom_Promise struct{ *capnp.Pipeline }

func (p Custom_Promise) Struct() (Custom, error) {
	s, err := p.Pipeline.Struct()
	return Custom{s}, err
}

func (p Custom_Promise) Value() Custom_value_Promise { return Custom_value_Promise{p.Pipeline} }

// Custom_value_Promise is a wrapper for a Custom_value promised by a client call.
type Custom_value_Promise struct{ *capnp.Pipeline }

func (p Custom_value_Promise) Struct() (Custom_value, error) {
	s, err := p.Pipeline.Struct()
	return Custom_value{s}, err
}

type CHF struct{ capnp.Struct }

// CHF_TypeID is the unique identifier for the type CHF.
const CHF_TypeID = 0xa7ab5c68e4bc7b62

func NewCHF(s *capnp.Segment) (CHF, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 224, PointerCount: 10})
	return CHF{st}, err
}

func NewRootCHF(s *capnp.Segment) (CHF, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 224, PointerCount: 10})
	return CHF{st}, err
}

func ReadRootCHF(msg *capnp.Message) (CHF, error) {
	root, err := msg.RootPtr()
	return CHF{root.Struct()}, err
}

func (s CHF) String() string {
	str, _ := text.Marshal(0xa7ab5c68e4bc7b62, s.Struct)
	return str
}

func (s CHF) TimestampNano() int64 {
	return int64(s.Struct.Uint64(0))
}

func (s CHF) SetTimestampNano(v int64) {
	s.Struct.SetUint64(0, uint64(v))
}

func (s CHF) DstAs() uint32 {
	return s.Struct.Uint32(8)
}

func (s CHF) SetDstAs(v uint32) {
	s.Struct.SetUint32(8, v)
}

func (s CHF) DstGeo() uint32 {
	return s.Struct.Uint32(12)
}

func (s CHF) SetDstGeo(v uint32) {
	s.Struct.SetUint32(12, v)
}

func (s CHF) DstMac() uint32 {
	return s.Struct.Uint32(16)
}

func (s CHF) SetDstMac(v uint32) {
	s.Struct.SetUint32(16, v)
}

func (s CHF) HeaderLen() uint32 {
	return s.Struct.Uint32(20)
}

func (s CHF) SetHeaderLen(v uint32) {
	s.Struct.SetUint32(20, v)
}

func (s CHF) InBytes() uint64 {
	return s.Struct.Uint64(24)
}

func (s CHF) SetInBytes(v uint64) {
	s.Struct.SetUint64(24, v)
}

func (s CHF) InPkts() uint64 {
	return s.Struct.Uint64(32)
}

func (s CHF) SetInPkts(v uint64) {
	s.Struct.SetUint64(32, v)
}

func (s CHF) InputPort() uint32 {
	return s.Struct.Uint32(40)
}

func (s CHF) SetInputPort(v uint32) {
	s.Struct.SetUint32(40, v)
}

func (s CHF) IpSize() uint32 {
	return s.Struct.Uint32(44)
}

func (s CHF) SetIpSize(v uint32) {
	s.Struct.SetUint32(44, v)
}

func (s CHF) Ipv4DstAddr() uint32 {
	return s.Struct.Uint32(48)
}

func (s CHF) SetIpv4DstAddr(v uint32) {
	s.Struct.SetUint32(48, v)
}

func (s CHF) Ipv4SrcAddr() uint32 {
	return s.Struct.Uint32(52)
}

func (s CHF) SetIpv4SrcAddr(v uint32) {
	s.Struct.SetUint32(52, v)
}

func (s CHF) L4DstPort() uint32 {
	return s.Struct.Uint32(56)
}

func (s CHF) SetL4DstPort(v uint32) {
	s.Struct.SetUint32(56, v)
}

func (s CHF) L4SrcPort() uint32 {
	return s.Struct.Uint32(60)
}

func (s CHF) SetL4SrcPort(v uint32) {
	s.Struct.SetUint32(60, v)
}

func (s CHF) OutputPort() uint32 {
	return s.Struct.Uint32(64)
}

func (s CHF) SetOutputPort(v uint32) {
	s.Struct.SetUint32(64, v)
}

func (s CHF) Protocol() uint32 {
	return s.Struct.Uint32(68)
}

func (s CHF) SetProtocol(v uint32) {
	s.Struct.SetUint32(68, v)
}

func (s CHF) SampledPacketSize() uint32 {
	return s.Struct.Uint32(72)
}

func (s CHF) SetSampledPacketSize(v uint32) {
	s.Struct.SetUint32(72, v)
}

func (s CHF) SrcAs() uint32 {
	return s.Struct.Uint32(76)
}

func (s CHF) SetSrcAs(v uint32) {
	s.Struct.SetUint32(76, v)
}

func (s CHF) SrcGeo() uint32 {
	return s.Struct.Uint32(80)
}

func (s CHF) SetSrcGeo(v uint32) {
	s.Struct.SetUint32(80, v)
}

func (s CHF) SrcMac() uint32 {
	return s.Struct.Uint32(84)
}

func (s CHF) SetSrcMac(v uint32) {
	s.Struct.SetUint32(84, v)
}

func (s CHF) TcpFlags() uint32 {
	return s.Struct.Uint32(88)
}

func (s CHF) SetTcpFlags(v uint32) {
	s.Struct.SetUint32(88, v)
}

func (s CHF) Tos() uint32 {
	return s.Struct.Uint32(92)
}

func (s CHF) SetTos(v uint32) {
	s.Struct.SetUint32(92, v)
}

func (s CHF) VlanIn() uint32 {
	return s.Struct.Uint32(96)
}

func (s CHF) SetVlanIn(v uint32) {
	s.Struct.SetUint32(96, v)
}

func (s CHF) VlanOut() uint32 {
	return s.Struct.Uint32(100)
}

func (s CHF) SetVlanOut(v uint32) {
	s.Struct.SetUint32(100, v)
}

func (s CHF) Ipv4NextHop() uint32 {
	return s.Struct.Uint32(104)
}

func (s CHF) SetIpv4NextHop(v uint32) {
	s.Struct.SetUint32(104, v)
}

func (s CHF) MplsType() uint32 {
	return s.Struct.Uint32(108)
}

func (s CHF) SetMplsType(v uint32) {
	s.Struct.SetUint32(108, v)
}

func (s CHF) OutBytes() uint64 {
	return s.Struct.Uint64(112)
}

func (s CHF) SetOutBytes(v uint64) {
	s.Struct.SetUint64(112, v)
}

func (s CHF) OutPkts() uint64 {
	return s.Struct.Uint64(120)
}

func (s CHF) SetOutPkts(v uint64) {
	s.Struct.SetUint64(120, v)
}

func (s CHF) TcpRetransmit() uint32 {
	return s.Struct.Uint32(128)
}

func (s CHF) SetTcpRetransmit(v uint32) {
	s.Struct.SetUint32(128, v)
}

func (s CHF) SrcFlowTags() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s CHF) HasSrcFlowTags() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s CHF) SrcFlowTagsBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s CHF) SetSrcFlowTags(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s CHF) DstFlowTags() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s CHF) HasDstFlowTags() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s CHF) DstFlowTagsBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s CHF) SetDstFlowTags(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(1, t.List.ToPtr())
}

func (s CHF) SampleRate() uint32 {
	return s.Struct.Uint32(132)
}

func (s CHF) SetSampleRate(v uint32) {
	s.Struct.SetUint32(132, v)
}

func (s CHF) DeviceId() uint32 {
	return s.Struct.Uint32(136)
}

func (s CHF) SetDeviceId(v uint32) {
	s.Struct.SetUint32(136, v)
}

func (s CHF) FlowTags() (string, error) {
	p, err := s.Struct.Ptr(2)
	return p.Text(), err
}

func (s CHF) HasFlowTags() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s CHF) FlowTagsBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(2)
	return p.TextBytes(), err
}

func (s CHF) SetFlowTags(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(2, t.List.ToPtr())
}

func (s CHF) Timestamp() int64 {
	return int64(s.Struct.Uint64(144))
}

func (s CHF) SetTimestamp(v int64) {
	s.Struct.SetUint64(144, uint64(v))
}

func (s CHF) DstBgpAsPath() (string, error) {
	p, err := s.Struct.Ptr(3)
	return p.Text(), err
}

func (s CHF) HasDstBgpAsPath() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s CHF) DstBgpAsPathBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(3)
	return p.TextBytes(), err
}

func (s CHF) SetDstBgpAsPath(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(3, t.List.ToPtr())
}

func (s CHF) DstBgpCommunity() (string, error) {
	p, err := s.Struct.Ptr(4)
	return p.Text(), err
}

func (s CHF) HasDstBgpCommunity() bool {
	p, err := s.Struct.Ptr(4)
	return p.IsValid() || err != nil
}

func (s CHF) DstBgpCommunityBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(4)
	return p.TextBytes(), err
}

func (s CHF) SetDstBgpCommunity(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(4, t.List.ToPtr())
}

func (s CHF) SrcBgpAsPath() (string, error) {
	p, err := s.Struct.Ptr(5)
	return p.Text(), err
}

func (s CHF) HasSrcBgpAsPath() bool {
	p, err := s.Struct.Ptr(5)
	return p.IsValid() || err != nil
}

func (s CHF) SrcBgpAsPathBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(5)
	return p.TextBytes(), err
}

func (s CHF) SetSrcBgpAsPath(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(5, t.List.ToPtr())
}

func (s CHF) SrcBgpCommunity() (string, error) {
	p, err := s.Struct.Ptr(6)
	return p.Text(), err
}

func (s CHF) HasSrcBgpCommunity() bool {
	p, err := s.Struct.Ptr(6)
	return p.IsValid() || err != nil
}

func (s CHF) SrcBgpCommunityBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(6)
	return p.TextBytes(), err
}

func (s CHF) SetSrcBgpCommunity(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(6, t.List.ToPtr())
}

func (s CHF) SrcNextHopAs() uint32 {
	return s.Struct.Uint32(140)
}

func (s CHF) SetSrcNextHopAs(v uint32) {
	s.Struct.SetUint32(140, v)
}

func (s CHF) DstNextHopAs() uint32 {
	return s.Struct.Uint32(152)
}

func (s CHF) SetDstNextHopAs(v uint32) {
	s.Struct.SetUint32(152, v)
}

func (s CHF) SrcGeoRegion() uint32 {
	return s.Struct.Uint32(156)
}

func (s CHF) SetSrcGeoRegion(v uint32) {
	s.Struct.SetUint32(156, v)
}

func (s CHF) DstGeoRegion() uint32 {
	return s.Struct.Uint32(160)
}

func (s CHF) SetDstGeoRegion(v uint32) {
	s.Struct.SetUint32(160, v)
}

func (s CHF) SrcGeoCity() uint32 {
	return s.Struct.Uint32(164)
}

func (s CHF) SetSrcGeoCity(v uint32) {
	s.Struct.SetUint32(164, v)
}

func (s CHF) DstGeoCity() uint32 {
	return s.Struct.Uint32(168)
}

func (s CHF) SetDstGeoCity(v uint32) {
	s.Struct.SetUint32(168, v)
}

func (s CHF) Big() bool {
	return s.Struct.Bit(1376)
}

func (s CHF) SetBig(v bool) {
	s.Struct.SetBit(1376, v)
}

func (s CHF) SampleAdj() bool {
	return s.Struct.Bit(1377)
}

func (s CHF) SetSampleAdj(v bool) {
	s.Struct.SetBit(1377, v)
}

func (s CHF) Ipv4DstNextHop() uint32 {
	return s.Struct.Uint32(176)
}

func (s CHF) SetIpv4DstNextHop(v uint32) {
	s.Struct.SetUint32(176, v)
}

func (s CHF) Ipv4SrcNextHop() uint32 {
	return s.Struct.Uint32(180)
}

func (s CHF) SetIpv4SrcNextHop(v uint32) {
	s.Struct.SetUint32(180, v)
}

func (s CHF) SrcRoutePrefix() uint32 {
	return s.Struct.Uint32(184)
}

func (s CHF) SetSrcRoutePrefix(v uint32) {
	s.Struct.SetUint32(184, v)
}

func (s CHF) DstRoutePrefix() uint32 {
	return s.Struct.Uint32(188)
}

func (s CHF) SetDstRoutePrefix(v uint32) {
	s.Struct.SetUint32(188, v)
}

func (s CHF) SrcRouteLength() uint8 {
	return s.Struct.Uint8(173)
}

func (s CHF) SetSrcRouteLength(v uint8) {
	s.Struct.SetUint8(173, v)
}

func (s CHF) DstRouteLength() uint8 {
	return s.Struct.Uint8(174)
}

func (s CHF) SetDstRouteLength(v uint8) {
	s.Struct.SetUint8(174, v)
}

func (s CHF) SrcSecondAsn() uint32 {
	return s.Struct.Uint32(192)
}

func (s CHF) SetSrcSecondAsn(v uint32) {
	s.Struct.SetUint32(192, v)
}

func (s CHF) DstSecondAsn() uint32 {
	return s.Struct.Uint32(196)
}

func (s CHF) SetDstSecondAsn(v uint32) {
	s.Struct.SetUint32(196, v)
}

func (s CHF) SrcThirdAsn() uint32 {
	return s.Struct.Uint32(200)
}

func (s CHF) SetSrcThirdAsn(v uint32) {
	s.Struct.SetUint32(200, v)
}

func (s CHF) DstThirdAsn() uint32 {
	return s.Struct.Uint32(204)
}

func (s CHF) SetDstThirdAsn(v uint32) {
	s.Struct.SetUint32(204, v)
}

func (s CHF) Ipv6DstAddr() ([]byte, error) {
	p, err := s.Struct.Ptr(7)
	return []byte(p.Data()), err
}

func (s CHF) HasIpv6DstAddr() bool {
	p, err := s.Struct.Ptr(7)
	return p.IsValid() || err != nil
}

func (s CHF) SetIpv6DstAddr(v []byte) error {
	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(7, d.List.ToPtr())
}

func (s CHF) Ipv6SrcAddr() ([]byte, error) {
	p, err := s.Struct.Ptr(8)
	return []byte(p.Data()), err
}

func (s CHF) HasIpv6SrcAddr() bool {
	p, err := s.Struct.Ptr(8)
	return p.IsValid() || err != nil
}

func (s CHF) SetIpv6SrcAddr(v []byte) error {
	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(8, d.List.ToPtr())
}

func (s CHF) SrcEthMac() uint64 {
	return s.Struct.Uint64(208)
}

func (s CHF) SetSrcEthMac(v uint64) {
	s.Struct.SetUint64(208, v)
}

func (s CHF) DstEthMac() uint64 {
	return s.Struct.Uint64(216)
}

func (s CHF) SetDstEthMac(v uint64) {
	s.Struct.SetUint64(216, v)
}

func (s CHF) Custom() (Custom_List, error) {
	p, err := s.Struct.Ptr(9)
	return Custom_List{List: p.List()}, err
}

func (s CHF) HasCustom() bool {
	p, err := s.Struct.Ptr(9)
	return p.IsValid() || err != nil
}

func (s CHF) SetCustom(v Custom_List) error {
	return s.Struct.SetPtr(9, v.List.ToPtr())
}

// NewCustom sets the custom field to a newly
// allocated Custom_List, preferring placement in s's segment.
func (s CHF) NewCustom(n int32) (Custom_List, error) {
	l, err := NewCustom_List(s.Struct.Segment(), n)
	if err != nil {
		return Custom_List{}, err
	}
	err = s.Struct.SetPtr(9, l.List.ToPtr())
	return l, err
}

// CHF_List is a list of CHF.
type CHF_List struct{ capnp.List }

// NewCHF creates a new list of CHF.
func NewCHF_List(s *capnp.Segment, sz int32) (CHF_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 224, PointerCount: 10}, sz)
	return CHF_List{l}, err
}

func (s CHF_List) At(i int) CHF { return CHF{s.List.Struct(i)} }

func (s CHF_List) Set(i int, v CHF) error { return s.List.SetStruct(i, v.Struct) }

// CHF_Promise is a wrapper for a CHF promised by a client call.
type CHF_Promise struct{ *capnp.Pipeline }

func (p CHF_Promise) Struct() (CHF, error) {
	s, err := p.Pipeline.Struct()
	return CHF{s}, err
}

type PackedCHF struct{ capnp.Struct }

// PackedCHF_TypeID is the unique identifier for the type PackedCHF.
const PackedCHF_TypeID = 0xb158a6a28e2d29c2

func NewPackedCHF(s *capnp.Segment) (PackedCHF, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return PackedCHF{st}, err
}

func NewRootPackedCHF(s *capnp.Segment) (PackedCHF, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return PackedCHF{st}, err
}

func ReadRootPackedCHF(msg *capnp.Message) (PackedCHF, error) {
	root, err := msg.RootPtr()
	return PackedCHF{root.Struct()}, err
}

func (s PackedCHF) String() string {
	str, _ := text.Marshal(0xb158a6a28e2d29c2, s.Struct)
	return str
}

func (s PackedCHF) Msgs() (CHF_List, error) {
	p, err := s.Struct.Ptr(0)
	return CHF_List{List: p.List()}, err
}

func (s PackedCHF) HasMsgs() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s PackedCHF) SetMsgs(v CHF_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMsgs sets the msgs field to a newly
// allocated CHF_List, preferring placement in s's segment.
func (s PackedCHF) NewMsgs(n int32) (CHF_List, error) {
	l, err := NewCHF_List(s.Struct.Segment(), n)
	if err != nil {
		return CHF_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// PackedCHF_List is a list of PackedCHF.
type PackedCHF_List struct{ capnp.List }

// NewPackedCHF creates a new list of PackedCHF.
func NewPackedCHF_List(s *capnp.Segment, sz int32) (PackedCHF_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return PackedCHF_List{l}, err
}

func (s PackedCHF_List) At(i int) PackedCHF { return PackedCHF{s.List.Struct(i)} }

func (s PackedCHF_List) Set(i int, v PackedCHF) error { return s.List.SetStruct(i, v.Struct) }

// PackedCHF_Promise is a wrapper for a PackedCHF promised by a client call.
type PackedCHF_Promise struct{ *capnp.Pipeline }

func (p PackedCHF_Promise) Struct() (PackedCHF, error) {
	s, err := p.Pipeline.Struct()
	return PackedCHF{s}, err
}

const schema_c75f49ee0059f55d = "x\xdat\xd6}h\x1c\xc7\x19\x06\xf0y\xefN\xbbg" +
	"Y\xb2\xee|#[\xb6\xe5H\xb6e[r$Y_" +
	"\xf1\x87\"Pd\xc9\xae]d\xf7N\xb2\x8d[l\xe2" +
	"\xf5\xddZ\xda\xe4\xbe\xb8[\xf9#mI)u)\xa5" +
	"-\xa5\xb4\xd0?Rp\xd2\xb8ui\xd2:\xe0BZ" +
	"\\p\x8bKSH\xc1\x85\x14\x12H!\x05\x17\x12H" +
	"\xda\x14\x1c\x9a\xd2\xba\xd7\xe7\x19I{'Y1\xac\xd1" +
	"\xfev\xf6\x9d\xd9w\xe6\x9d\x9b\xbe\x1fF\x9e\x08\xf5\xd7" +
	"]\x0c+\x95\xea\xac\xb3*\xe7>\x7f\xeb\xde\xec\xe9\x9f" +
	"\xfeX\xa5\xea\xa5\xf5_g>\xfa\xec\xdf\x8f<\xf9\xfb" +
	"\xbaz[\xa9\xc4\xf3RJ\\\x13\xfc5\xf8\xbc<\xd9" +
	"\xa8\xa4\xf2\xdb\xae\x9eo\xbd\xf0\xa3S\xaf\xa8x\xbd\x04" +
	"M\xd9 \xf1\x83\xd8\x0b\x89k1\xf3Rl\x14-w" +
	"\xde\x0d\xaf\xfd\xea\xde3\x1f0jhi\xd3?\xc5\xbe" +
	"\x9ex\xcb4\xfds\xec\xa2\x92\x7f<x\xef\xca\x95\x93" +
	"W\xff\x93\x8aK\xa8\xfa\xda\x09\xb1%,\x91Ds\xfc" +
	"WJ\xf0\xff\xcf\xd5x%={\xbe7\xed\x14\xf3\xaa" +
	"\xad8<~\xf8PR$u?\x1cQ*\"\x1c\x80" +
	"URj\xfa9+,\xd3\xd7\xad\x90 \x98\x16\xfa5" +
	"k\x00~\x95\xfe\x12=\x14\xd6\x12\x82\xff\xc4\x1a\x86\xbf" +
	"H\xbfA\x0fG\xb4 %\x89\x97\x8d_\xa7\xdf\xa4G" +
	"\xea\xb4\xa0\x8f\xc4+\xd6\x14\xfc\x06\xfd\x16\xbd\x0eq\xea" +
	"\xe0\xbf\xb4\x0e\xc0o\xd2o\xd3-\xc4\xb1\xe0\xbf6q" +
	"^\xa5\xdf\xa1\xdb\xf5\xda|\xfdoL\x9c\xdb\xf4\xd7\xe9" +
	"\xd1\xd5Z\xa2\xf0?\x98\xf6w\xe8w\xe9\xab\x1a\xb4\xac" +
	"\x82\xff\xd1:\x07\x7f\x9d\xfe&\xbd\xbeQK=sg" +
	"\xfc\x0d\xfa;\xf4\xd5k\xb4\xac\x86\xff\xc5\xc4\x7f\x9b\xfe" +
	".\xbd\xa1IK\x03\xfco\xc6\xef\xd1?\xa47\xc6\xb4" +
	"4\xc2?\xb0>\x07\x7f\x9f\xfe1}M\\\xcb\x1a\xf8" +
	"G\xd6\xa7\xe1\xf7\xe1S6\xb8i\xad\x96&\xf0\x03\xeb" +
	";JM\xd9h\xdd@\x8e%\xb4\xc4\xc0\xablf9" +
	"B\x8f\xd1\xe3ZK\x1c\xdeh\xf3\xab\xa2tM_\xdb" +
	"\xace-<n\xbc\x81\xdeBO\xac\xd3\x92\x807\xdb" +
	"\xecU\xd3\xdb\xe9z\xbd\x16\x0d\xdflo\x85\xb7\xd0;" +
	"\xe8\xcd-Z\x9a\xe1[L\x9cVz'}\xdd\x06-" +
	"\xeb\xe0\xdbm\xceJ;\xbd\x9b\xbe~\xa3\x96\xf5\xf0." +
	"\x9bY\xeb\xa4\x0f\xd1[6ii\x81\xf7\x9b~\xfb\xe8" +
	"#\xf4\x0d\xc8\xe6\x06\xf8~\xe3\xfb\xe8\x13\xf4\x8d\xc8\xe6" +
	"F\xf8\x98\x89?B?L\xdf\xd4\xaee\x13\xfc\xa0\xcd" +
	"U8AO\xc2\xa5UK+\xf8\xa8\xe9v\x92|\x8a" +
	"\xcd7\x8b\x96\xcd\xf0\x13\xc6\x8f\xd3\xcf\xd2\x1f\xd9\xa2\xe5" +
	"\x11\xf8\x19\x9b\x93r\x9a>Ko\xdb\xaa\xa5\x0d\xee\x9a" +
	"\xe1d\xe8Ez;\x16y;<g<K\xbfD\xdf" +
	"\x82\xd9\xda\x02\x9f\xb39\xe9>\xfdK\xf4\xadX\xb4H" +
	"d\xe2\x8b\xf6S\xf0/\xd0\xbfF\xdf\x86E\xbb\x0d~" +
	"\xc5\xfe2\xfc+\xf4o\xd3;\xb0\xf8;\xe0\xdf4\xed" +
	"\xbfA\xff>}\xbb\xa5e;\xfc{\xa6\xfdw\xe9W" +
	"\xe9;\xb6i\xd9\xc1b4\xed\x9f\xa3_\xa7\xef\xdc\xa1" +
	"e'\x8b\xd1\xf8\x8b\xf4\x1b\xf4\xce\x9dZ:Yt\xc6" +
	"_\xa2\xbfJ\xef\xea\xd4\xd2\x05\xff\x85\xf1\x9b\xf4\xdb\xf4" +
	"]]Zv\xb1\xb8L~n\xd1_\xa3?\xbaK\xcb" +
	"\xa3\xf0\xdf\x19\xbfC\xbf\x0b\xef\xef>\x8b\x0f\xe8f\x15" +
	"\x99\xf5\xf3\x1a\x1f\xbc\xc1\x07=\x0e\x1e\xf4p32\x19" +
	"\xba\xcb\x07o3Ro\xb7\x96^\xf8[\xf63\xf07" +
	"\xe9\xf7\xe8\xbb{\xb4\xec\x86\xff\xd5\xf8;\xf4\xf7\xe9}" +
	"\xbdZ\xfa\xe0\xef\x19\x7f\x97~\x9f\xde\xbf[K?\xfc" +
	"\x9f\xc6?\xa4\xff\x97>\xf0\xb2\x16\x94J\xe2\xdf\xc6?" +
	"\xa6G\xa2\xf0\xc1\x9fi\x19\x84K\x14>\x15e}\x91" +
	"\x87\xfa\xb4\x0c\xb1\xbe\xa2LD\x94\xae\xe9\x8f\xf5ky" +
	"\x8cud<Fo\xa5\xef\x19\xd0\xb2\x07\xbe!\xca\x85" +
	"\xd5B\xef\xa0\xef\x1d\xd4\xb2\x97\xf5b\xbc\x9d\xdeM\xdf" +
	"gk\xd9\xc7\xba0\xdeI\x1f\xa2\xef\x8fj\xd9\xcf\xba" +
	"0\xdeG\x1f\xa1\x0f\xa3\x8e\x86Y\x17Q\xa6m\x1f}" +
	"\x82\xfe8\xea\xe8q\xd6\x85\xf1'\xe8\x93\xf4\x91UZ" +
	"F\xe0G\xa2\xac\xd3\x09\xfaYx\xc5\xf7rn\xd9w" +
	"r\xd8\xca\x8f9\xf9\x02v\xd2\x10.i\xcb\x94\xfd\xb1" +
	"2\xf6\xc3\x10.\x19\xc5\xdd\xa7\xdcB\xed\xedQ'\xbd" +
	"x[\x99u\x9d\x8c[\x9at\x95\xe4\x17\xedY/\x7f" +
	"\xe0\xb2\xef\x96\xb1s\x86p\xc9\xa8\x97O>\xed\x07\xb7" +
	"\x15/_\x9c\xf3\x93\x85\x92\x12?\x88\xea\x15\xa7\xbdg" +
	"\xdc \xaaW\xbc04\x81Q(;\x93)-\xd1\xe9" +
	"Rz\x99f\xd9ri8\x18\xda-\xb3\xc2\x9co\xfa" +
	"U\xe1R\x15\x8b\xa5\x82_H\x17\xb2J\xa9\xc0\xcaN" +
	"\xae\x98u3Iq\xd2O\xbb>\xc6%\xc1\xc0\xda\xca" +
	"\xe8\xbe\x9a\x1a\xdc\xd5\xa6\x06\xb7\xb5\xa9\xf1\xd3\xc5CY" +
	"g\xa6\\\x13\xdb\xf6\x0b\xd5\xb7/d\x9d\xfc\x91j\xd6" +
	"x\xfb\x999\x7f\xc9\xd7\x1es/\xf9\xca>\\(\x06" +
	"\x8a\xa1\x95\x8f_.\xba\xb5\x03\xc6\x97\x99\x84\xd3\x16\x92" +
	"\xfc,lI\xd21\x98)\xd7/9\xaa-_\xcey" +
	"\xd5^0\xe6C\xd9\xc2\xc5\xe3\xca\xc6H\xf1\x1b\x15\xc2" +
	"%\x15L\xf2\x0a:\x9f\x98)G\x85\xfd\xeaTe\xdc" +
	"\x0b^\xda=\x92\xa9\x1d\xd0y\xbe\xbb\xf0\xe5\x8b/\x07" +
	"\x8bM\x8a\x8b\x0b\x8d\xdd\x1c\x98)\x8e\x95US\xd2\xf1" +
	"gk{\x07\x8f\x17$\x97\x9b\xcb{\xfe\xe5j\x10\x8c" +
	"v\xa5\x17\xe6\xf9\x13^`\x0e\x0f\xab\xa6B\xb1:q" +
	"\xeca%\x9e\x9f\xcf)W5\xcdx\x85|m\xeb\x95" +
	"x\xbe\xf5\xb8\x87t\\^\xd6v)\xda\xe7\xbc\x19\x1c" +
	"\x83\xf0K\x14$q,\xa3\xe4\xa9\xc0\x16\x96\xfb15" +
	"\xcaA\xd5\xcc\xf6\xc2\x8a\x7f\xf8\x01\xfa\x9e\xc2\x14\xbbj" +
	"4Yr\xcf{\x97j\xfb_\xf9A\xf5\x8dI7?" +
	"\x83\xe4Yx`-yc\xd9\x03\xbc1\xed\xa6\x0by" +
	"\xd5\x94\x19+/\xc9\xc6J\x8c\xd6\xc7g\xbdRF\xd9" +
	"\xcb\x1a\xaf\xa0\xf8\xac=\xd5\xf2n\x846.h\xb5\xbc" +
	"\x17\x15q\x0f\xfa\xb3G\x1d%\xe9`1#\xear\x1b" +
	"M\xcf\x95\xfdB\x0e\xa7&I\x86Eb\xd5C\xac\x12" +
	"bpn\x95\xe2p\x92\x85\x9d\xb1\x17N\xaf\x91\xc5\xd3" +
	"k\xbc\x11\xbfj)\xec\x8e\xa9\x8e\x904\xe5\xcaX\xf9" +
	"A\xb4\xe0|\xfep\xb4q\xd3\xb1R\x8c\x15\x0dbu" +
	"\xe1\x80\x92\xea@\xac>L1\xfe-\x1e\xb0\xe3=\x03" +
	"*\x14\xf62\xc1\x8er\xc1\xc9\xce\xb9\x0f\xc7\x1b\xed5" +
	"\x0fR\x0d\xe1Hk\xa5\"\xe6\x18\x1d?\x88}=5" +
	"\x81\xa0\xc9\x90l\x96\xff\x91q\x8a\x8e\x1f\xc5\xefmj" +
	"\x12|\x0a)\x0b=\xa8\x983t\xfc\x046\xfb\x14F" +
	"\x9f:\x8d\xad~\xce\xcb\xfb\x83\x03'\x91\xb1lm\x95" +
	":D\x15v\xb28\xcd\x86pa\x0f\xf3K'q\xbb" +
	"P?\xff\x0f\x00\x00\xff\xff\xa3\xc2r\xe6"

func init() {
	schemas.Register(schema_c75f49ee0059f55d,
		0xa7ab5c68e4bc7b62,
		0xb158a6a28e2d29c2,
		0xed5d37861203d027,
		0xfba056008585e9fd)
}
