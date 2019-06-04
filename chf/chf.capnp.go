// Code generated by capnpc-go. DO NOT EDIT.

package chf

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
	Custom_value_Which_uint64Val  Custom_value_Which = 3
	Custom_value_Which_addrVal    Custom_value_Which = 4
	Custom_value_Which_uint16Val  Custom_value_Which = 5
	Custom_value_Which_uint8Val   Custom_value_Which = 6
)

func (w Custom_value_Which) String() string {
	const s = "uint32Valfloat32ValstrValuint64ValaddrValuint16Valuint8Val"
	switch w {
	case Custom_value_Which_uint32Val:
		return s[0:9]
	case Custom_value_Which_float32Val:
		return s[9:19]
	case Custom_value_Which_strVal:
		return s[19:25]
	case Custom_value_Which_uint64Val:
		return s[25:34]
	case Custom_value_Which_addrVal:
		return s[34:41]
	case Custom_value_Which_uint16Val:
		return s[41:50]
	case Custom_value_Which_uint8Val:
		return s[50:58]

	}
	return "Custom_value_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

// Custom_TypeID is the unique identifier for the type Custom.
const Custom_TypeID = 0xed5d37861203d027

func NewCustom(s *capnp.Segment) (Custom, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1})
	return Custom{st}, err
}

func NewRootCustom(s *capnp.Segment) (Custom, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1})
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
	if s.Struct.Uint16(8) != 0 {
		panic("Which() != uint32Val")
	}
	return s.Struct.Uint32(4)
}

func (s Custom_value) SetUint32Val(v uint32) {
	s.Struct.SetUint16(8, 0)
	s.Struct.SetUint32(4, v)
}

func (s Custom_value) Float32Val() float32 {
	if s.Struct.Uint16(8) != 1 {
		panic("Which() != float32Val")
	}
	return math.Float32frombits(s.Struct.Uint32(4))
}

func (s Custom_value) SetFloat32Val(v float32) {
	s.Struct.SetUint16(8, 1)
	s.Struct.SetUint32(4, math.Float32bits(v))
}

func (s Custom_value) StrVal() (string, error) {
	if s.Struct.Uint16(8) != 2 {
		panic("Which() != strVal")
	}
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Custom_value) HasStrVal() bool {
	if s.Struct.Uint16(8) != 2 {
		return false
	}
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Custom_value) StrValBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Custom_value) SetStrVal(v string) error {
	s.Struct.SetUint16(8, 2)
	return s.Struct.SetText(0, v)
}

func (s Custom_value) Uint64Val() uint64 {
	if s.Struct.Uint16(8) != 3 {
		panic("Which() != uint64Val")
	}
	return s.Struct.Uint64(16)
}

func (s Custom_value) SetUint64Val(v uint64) {
	s.Struct.SetUint16(8, 3)
	s.Struct.SetUint64(16, v)
}

func (s Custom_value) AddrVal() ([]byte, error) {
	if s.Struct.Uint16(8) != 4 {
		panic("Which() != addrVal")
	}
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s Custom_value) HasAddrVal() bool {
	if s.Struct.Uint16(8) != 4 {
		return false
	}
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Custom_value) SetAddrVal(v []byte) error {
	s.Struct.SetUint16(8, 4)
	return s.Struct.SetData(0, v)
}

func (s Custom_value) Uint16Val() uint16 {
	if s.Struct.Uint16(8) != 5 {
		panic("Which() != uint16Val")
	}
	return s.Struct.Uint16(4)
}

func (s Custom_value) SetUint16Val(v uint16) {
	s.Struct.SetUint16(8, 5)
	s.Struct.SetUint16(4, v)
}

func (s Custom_value) Uint8Val() uint8 {
	if s.Struct.Uint16(8) != 6 {
		panic("Which() != uint8Val")
	}
	return s.Struct.Uint8(4)
}

func (s Custom_value) SetUint8Val(v uint8) {
	s.Struct.SetUint16(8, 6)
	s.Struct.SetUint8(4, v)
}

func (s Custom) IsDimension() bool {
	return s.Struct.Bit(80)
}

func (s Custom) SetIsDimension(v bool) {
	s.Struct.SetBit(80, v)
}

// Custom_List is a list of Custom.
type Custom_List struct{ capnp.List }

// NewCustom creates a new list of Custom.
func NewCustom_List(s *capnp.Segment, sz int32) (Custom_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1}, sz)
	return Custom_List{l}, err
}

func (s Custom_List) At(i int) Custom { return Custom{s.List.Struct(i)} }

func (s Custom_List) Set(i int, v Custom) error { return s.List.SetStruct(i, v.Struct) }

func (s Custom_List) String() string {
	str, _ := text.MarshalList(0xed5d37861203d027, s.List)
	return str
}

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
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 224, PointerCount: 14})
	return CHF{st}, err
}

func NewRootCHF(s *capnp.Segment) (CHF, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 224, PointerCount: 14})
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
	return s.Struct.SetText(0, v)
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
	return s.Struct.SetText(1, v)
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
	return s.Struct.SetText(2, v)
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
	return s.Struct.SetText(3, v)
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
	return s.Struct.SetText(4, v)
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
	return s.Struct.SetText(5, v)
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
	return s.Struct.SetText(6, v)
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
	return s.Struct.SetData(7, v)
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
	return s.Struct.SetData(8, v)
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

func (s CHF) Ipv6SrcNextHop() ([]byte, error) {
	p, err := s.Struct.Ptr(10)
	return []byte(p.Data()), err
}

func (s CHF) HasIpv6SrcNextHop() bool {
	p, err := s.Struct.Ptr(10)
	return p.IsValid() || err != nil
}

func (s CHF) SetIpv6SrcNextHop(v []byte) error {
	return s.Struct.SetData(10, v)
}

func (s CHF) Ipv6DstNextHop() ([]byte, error) {
	p, err := s.Struct.Ptr(11)
	return []byte(p.Data()), err
}

func (s CHF) HasIpv6DstNextHop() bool {
	p, err := s.Struct.Ptr(11)
	return p.IsValid() || err != nil
}

func (s CHF) SetIpv6DstNextHop(v []byte) error {
	return s.Struct.SetData(11, v)
}

func (s CHF) Ipv6SrcRoutePrefix() ([]byte, error) {
	p, err := s.Struct.Ptr(12)
	return []byte(p.Data()), err
}

func (s CHF) HasIpv6SrcRoutePrefix() bool {
	p, err := s.Struct.Ptr(12)
	return p.IsValid() || err != nil
}

func (s CHF) SetIpv6SrcRoutePrefix(v []byte) error {
	return s.Struct.SetData(12, v)
}

func (s CHF) Ipv6DstRoutePrefix() ([]byte, error) {
	p, err := s.Struct.Ptr(13)
	return []byte(p.Data()), err
}

func (s CHF) HasIpv6DstRoutePrefix() bool {
	p, err := s.Struct.Ptr(13)
	return p.IsValid() || err != nil
}

func (s CHF) SetIpv6DstRoutePrefix(v []byte) error {
	return s.Struct.SetData(13, v)
}

func (s CHF) IsMetric() bool {
	return s.Struct.Bit(1378)
}

func (s CHF) SetIsMetric(v bool) {
	s.Struct.SetBit(1378, v)
}

// CHF_List is a list of CHF.
type CHF_List struct{ capnp.List }

// NewCHF creates a new list of CHF.
func NewCHF_List(s *capnp.Segment, sz int32) (CHF_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 224, PointerCount: 14}, sz)
	return CHF_List{l}, err
}

func (s CHF_List) At(i int) CHF { return CHF{s.List.Struct(i)} }

func (s CHF_List) Set(i int, v CHF) error { return s.List.SetStruct(i, v.Struct) }

func (s CHF_List) String() string {
	str, _ := text.MarshalList(0xa7ab5c68e4bc7b62, s.List)
	return str
}

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

func (s PackedCHF_List) String() string {
	str, _ := text.MarshalList(0xb158a6a28e2d29c2, s.List)
	return str
}

// PackedCHF_Promise is a wrapper for a PackedCHF promised by a client call.
type PackedCHF_Promise struct{ *capnp.Pipeline }

func (p PackedCHF_Promise) Struct() (PackedCHF, error) {
	s, err := p.Pipeline.Struct()
	return PackedCHF{s}, err
}

const schema_c75f49ee0059f55d = "x\xdat\xd7{h\x1c\xd7\x15\x06\xf0\xef\xcc\xecjv" +
	"e\xbdV\xf7&~I\x91_\xb1e\xc7\x0f\xbd\xfcR" +
	"\xe4H\xf2\xabv\x90\xdd\x1d\xadc\x9c`S\xafw\xd7" +
	"\xd2$\xab\xdd\xed\xce\xc8\xaf\xa6$-qi\x8b\xdb\x86" +
	"\x92\x82[bp\x82\xd3\xa6\xd4i\x13p \x14\x07L" +
	"I\xa8\x0b.8\xe0\x80\x0b)\xa4\x90\x82\x03I\x9bB" +
	"J\x9b\xc6\xe9\x96\xefJ\xfb\x90\xaa\xf8/\xedo\xce=" +
	"w\xee\xb9\xf7\x8cg\xbav\x86\x87\xac\xee\xf0I\x1bp" +
	";\xc3u\xa5c\xdf\xb8\xfa\xc1\xf8\xe1_\xfd\x02n\xbd" +
	"\xb4\xfd\xeb\xc8?\x1f\xfd\xdb\xde\xaf\xfd>\xdc\xe4\x00\xea" +
	"\x05)\xaa\x97\xc4\x01z_\x90\x93M\x90\xd2\xefV\xaf" +
	"\xfb\xe1\x8b??\xf4\x1ab\xf5R\x09e\x80\xfak\xec" +
	"E\xf5q\x8c\x7f}\x18\x1b\x84\x94V\xdd\xb4[\xbf\xb3" +
	"\xf9\xc8\xc7\xccj\xcf\x0c]\xd8\xfa}\xb5\xb4\x95\x7f\xb5" +
	"\xb7\xfe\x06\xf2\xf7/><{\xf6\xe0\xc5\xcf\xdd\x98\xd8" +
	"\xd5a\x8f\x88#\x8e\x84\xd4\x85\xd6\xdfBz/\xb4\xfe" +
	"H\xb0\xa3\x94\x1a?\xbe>\x95,\xe4\xd0Q\xe8\xdf\xb1" +
	"gw\\$\x11\x11;\x04\x84\x04P\xab\x9d\"\x90\xe8" +
	"tlI\xf49\x96\xc4\xc4\xd2B\xefvz\x80\xc4Z" +
	"\xfa\x16\xbaek\xb1\x00\xb5\xd1\xe9\x07\x12]\xf4\x01\xba" +
	"\x1d\xd2b\x03j\xab\xf1>\xfa\x10=\x14\xd6\x12\x02\xd4" +
	"6g\x14H\x0c\xd0\xf7\xd0\xc3\xb6\x960\xa0v9\xdb" +
	"\x81\xc4\x10}\x84^\x17\xd2R\x07\xa8\xbd&\xcfNz" +
	"\x9c\xee\xd4k\xb3\xfe}&\xcf\x08\xfd\x10=2OK" +
	"\x04P\x8f\x98\xf88\xfd0=\xda\xa0%\x0a\xa8G\x9d" +
	"c@\xe2\x10=M\xafo\xd4R\x0f\xa8\xa4\xf1\xa3\xf4" +
	",}^\x93\x96y\x80\xf2L\xfeqz@oh\xd6" +
	"\xd2\x00\xa8\xaf\x1b/\xd0\x9f\xa47\xb6hi\x04\xd4i" +
	"\xe71 q\x8a\xfe\x0c\xbd)\xa6\xa5\x09P\xdfr\x1e" +
	"\x06\x12O\xd3\xcf\xd1\x9b[\xb54\x03\xea{\xce\x8f\x81" +
	"\xc49\xfayz\x8b\xd2\xd2\x02\xa8\x9f\x98:?K\x7f" +
	"\x9e\x1e\xd3Zb\x80\xfa\xa9Y\xd7s\xf4\x8b\xf4\xd6{" +
	"\xb4\xb4\x02\xea\x82\xf1\xf3\xf4Ktu\xaf\x16\xc5Sg" +
	"\xe6\xbdH\xbfL\xd7\xf3\xb5h@\xfd\xd2Y\x06$." +
	"\xd1_\xa5\xdf\xb3@\xcb=\x80z\xc5\xe4y\x99~\x85" +
	"~\xefB-\xf7\x02\xea5\xb3/\x97\xe9o\xd0\xe7/" +
	"\xd22\x1fP\xaf\x9b\xba]\xa1_\xa3/X\xace\x01" +
	"\xa0\xde4\xf3^\xa5_\xa7/l\xd2\xb2\x10Po\x1b" +
	"\x7f\x8b~\x93\xbe\xa8Y\xcb\"@\xfd\xd1\xe4\xbfN\xbf" +
	"E_\xbcD\xcbb@\xbdc\xce\xe1M\xfa{\x8e%" +
	"\xd2\xa6\xa5\x0dP\x7f2\xd3\xde&\x7f\xc0\xf0v\xd1\xd2" +
	"\x0e\xa8\xbf\x18\x7f\x9f\xfe\x11\xfd\xbe\xa5Z\xeec\x1f\x99" +
	"m\xb9C\xff\x94\xde\xb1LK\x07\xa0\xfean\xe7\x13" +
	"\xfa]\xfa\x12K\xcb\x12@}f\xfc\xdf\xf4P\xc4\x92" +
	"\xd8\xd2V-K\x01%\x91Q`4bK\xa2\x81\xbc" +
	"\xcc\xd6\xb2\x0cP\xd1\xc8\xe3@\"B\xd7\xf4\xe5!-" +
	"\xcb\x01\x15\x8b|\x1bH\xb4\xd0\xdb\xe8+\xc2ZV\xb0" +
	"kM\xfc\x02\xfa\x0a\xfa\xfduZ\xee\x07\xd4R\x13\xbf" +
	"\x84\xbe\x96\xber\xb9\x96\x95\xecF\x13\xdfI\xef\xa3\xaf" +
	"Z\xa9e\x15\xbb\xd1x\x17}\x80\xde\xb9JK'\xbb" +
	"\xce\xf8\x16\xfaN\xfa\xeaN-\xab\x015l|\x88>" +
	"B_\xb3Z\xcb\x1avW\x84\xe5\xd9C?@\x7f`" +
	"\x8d\x96\x07\x00\xe5\x1a\x8f\xd3\x0fG,\xe9^{4\xac" +
	"e-\xdb(\xc2\xe3s\x80\x17\x8e\xf2\xc2\xbadX\xcb" +
	":@\x1da\x81\x12\x87ya\x9c\x99\xd6\xaf\xd5\xb2\x1e" +
	"P\x99\xc8\x19 \x91\xa6\x17\xe8\x1b\xd6i\xd9\x00\xa8\x09" +
	"\xe3Y\xfa)z\xd7z-]\x80\x9a4\x1e\xd0\x9f\xa6" +
	"wo\xd0\xd2\x0d\xa8o\x1a\x7f\x92\xfe]z\xcf+Z" +
	"z\x00u\xd6\xf83\xf4g\xe9\xbd\xbf\xd6\xd2\x0b\xa8\x1f" +
	"\x18?G?O\xef\xeb\xd2\xd2\xc7\xfe2\x95x\x8e~" +
	"\x91\xbe\xb1[\xcbF\xf6\x91\xf1\xe7\xe9/\xd37\xf5h" +
	"\xd9\x04\xa8\x97\"<X\x97\xe8\xaf\xd27\xf7j\xd9\xcc" +
	"~1~\x99\xfe\x06}\x8b\xa3e\x0b\xfb\xc2\xf8\x15\xfa" +
	"5\xfa\xd6\x88\x96\xad\xec\x0b\xe3W\xe9\xd7\xe9\xfd\x8b\xb4" +
	"\xf4\xb3/L\xdd\xde\xa2\xdf\xa4?\xb8X\xcb\x83\xec\x0b" +
	"\xe37\xe8\xb7\xe9\x03Q-\x03\x80z7\xc2>\xbdI" +
	"\xff\x88\xbe\xad^\xcb6\x1et\xb3\xde;\xf4O\xe9\x0f" +
	"\xcd\xd3\xf2\x10\x0f\xba\xf1O\xe8w\xe9\x83\x0dZ\x06y" +
	"\xd0#?\x03\x12w\xe9\x91\xa8%\xb1\xa1F-C\x80" +
	"\x0aG\xe9\x91(Ot\xd4\x92\xee\xe1ca-\xc3<" +
	"\xd2QvF\x0b/\xb4E-)\x05\xdeD\xc6\x0f\x92" +
	"\x13\xe8(\xecO\xe6\xf2\x12\x86%aHG\xda\x0f\x86" +
	"}\x89\xc0\x92\x08d0\xed\x07_\xc9\xe4k\x7f\xeeK" +
	"\xa6\xca?K\xe3\x99d:S\x1c\xc9@re{\xca" +
	"\xcbm?\x1dd|\x89\xc2\x92(d\xd0\xcb\xc5\x9f\x08" +
	"*?K^\xae0\x19\xc4\xf3EHP\xc9\xea\x15\x12" +
	"\xde\x99L%\xabW8\xd1\xb7\xd3\x0f\x86\xe1\xa4\xd3\xc5" +
	"\x19\x9a(\xa6fi\x96\x913\xd3\x95\xb2\x8c\x9be\xf9" +
	"\xc9\xc0\xcc\x0b\xbbX\xc5B1\x1f\xe4S\xf9,\x80\x8a" +
	"\xf9\xc9\x89B6\x93\x8eK2\xf5D&Hxg\xa4" +
	"rc\x1d~1US\x1a\xbf\x98\xaa-\x8d_L\xd5" +
	"\x96&H\x15vg\x93c~Mn'\xc8WG\x9f" +
	"\xc8&s{\xabU\xe3\xcf\xafN\x063V\xbb?s" +
	"*\x80\xb3'_\xa8\xe8D!\xeb\x1f8]\xc8\xd4\xde" +
	"p~20\x05\xa7M\x17\xf9\xa9\xfcd0\xa3\xe8A" +
	"\xaa0\x9a\x09\x8aIt\xe4\xfc\x09\xaf:\x8b_L\xed" +
	"\xce\xe6O\x1e\x80\x93\x1c\xf3\xa5\x01\x964@Ji?" +
	"\x98C\xa7\x0a3\x9a\x84\x1dT\xb7*\x9d9\xe1\xa52" +
	"{\xd3\xb57t\x9cc\xa7W^\x1e\\9lR(" +
	"\x1f4N\xb3}\xac0\xec\xa39\x9e\x0c\xc6kg\xdf" +
	">V\xd8\x91\x97\x89\x89\xc9\x9c\x17\x9c\xae&\xf1\x8b\xa9" +
	"\xb9\x06L\xf1\x97\x0c`\x0d\xf7\xa09_\xa8n\x1cg" +
	"\x98\x8b\xa7\xf6s4\x83\xe61/\x9f\xab\x8d\x9e\x8b\xa7" +
	"\xa2wx\xb0\x83\xd3\xb3bg\xa2s\xcc\x1b\x13\x81%" +
	"R)\xe2p\x1a\xf2x\xc5\xa6\x8f\xfb~\x0c\xf2\xa6j" +
	"v{\xfa\xc4\xff\xff\x05\xbf\x98\x1a\xcdO\x06\x19\x0c\xc6" +
	"\x8b\x99\xe3\xde\xa9\xda\xf9\xe7\xbeP\x1d1\x92\xc9\x8d\x05" +
	"\xe3R\x07K\xeaf\x8c\x98u\xc1/\xa6\x12\x99T>" +
	"\x87\xe6\xf4\xb0?\xa3\x1as\xb1_L\x1d\x18\xf7\x8ai" +
	"8\xb3\x82\xe7P\xafpbS\xb5\xbd\x1baI\xe3\xb4" +
	"V\xdb\xbb\xac~1\xb5+\x18\xdf\x97\x84\xa4*\x879" +
	"\xed\x07\xb3m05\xe9\x07\xf9\x09i\x82\xc4m\x91\x96" +
	"\xea\x1b4\x84X\xce^-e\xed\xb43\x8a?\xeb~" +
	"F\x85\xd5a5m\xef\xd4\xecQ_r\xd1\xdf\x97\x09" +
	"\x8a^\x8a\xe7\xbf\xbc\xc7\xe5wv)\xf4\xc7\xf9`I" +
	";So\xeen\xa8\xfc\xe2\x1ek\\\x03\xb8\x11[\xdc" +
	"\x15\x964O\xf8c~u5\x95\x8f\x93\xe9\xd5\xd4d" +
	"\xdba\x16\x0e0WC%\xd7\xaeE\x80;d\x8b;" +
	"b\x09\xff\x95\xbf.b{{`\xc5\xac\xb8y\x95\x8f" +
	"u\x1f\x03\xdc.[\xdc\x01Kl/]y\xcc\x9dH" +
	"f'3%\xcf\xdf\xe9Mdr>\x1c\x1e\xfb9\x16" +
	"25\xf5\xe0z\x13\xee\xb6\xd9\xa1\xb6RI\xcc\xc7F" +
	"\xec\xf5Q\xc0\xbdb\x8b{\xcd\x92v\xf9/\xd9\x02b" +
	"o>\x06\xb8Wmq\xaf[\xd2h}Q2_\x1a" +
	"\xb1\xb7\xfb\x01\xf7\x9a-\xee\x0dK\xda\xed\xbb%\xcb|" +
	"P\xc4\xfe\xc0\x1c\xd7mqoY\xd2\x18\xfa\xbcd>" +
	"'b\xefl\x07\xdc\x1b\xb6\xb8\xb7-i\x0f\xff\x87\xc1" +
	"\x0e\x10{\x97\xc1\xb7lq\xdf\xb7\xa4\xbd\xee\xb3R\xc8" +
	"|L\xc4\xfe\xfc0\xe0\xbeg\x8b{\xc7\x92\xd2\xa4\x97" +
	"\x0bz{\x0e&!\xd9\xdagU\x92\x08;\x99\x95z" +
	"XR\xcf'yP<\x98\xccV\x9e\"\x1c\xb7\xa9o" +
	"j\\\xf9\x01\x9bL\xa7MLy\xcf\x19\xd3\xbdi*" +
	"\xc6\x81%\xce\xb4m9\x984\xff\xbbLw\xd5\xff\x02" +
	"\x00\x00\xff\xff\x1f{\xdb\xc8"

func init() {
	schemas.Register(schema_c75f49ee0059f55d,
		0xa7ab5c68e4bc7b62,
		0xb158a6a28e2d29c2,
		0xed5d37861203d027,
		0xfba056008585e9fd)
}
