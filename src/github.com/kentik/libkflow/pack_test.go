package main

import (
	"bytes"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/kentik/libkflow/chf"
	"github.com/stretchr/testify/assert"

	"zombiezen.com/go/capnproto2"
)

func TestPack(t *testing.T) {
	assert := assert.New(t)

	ipv6srcaddr := randbytes(16)
	ipv6dstaddr := randbytes(16)

	f := &Ckflow{
		dstAs:             _Ctype_uint32_t(rand.Int31()),
		dstGeo:            _Ctype_uint32_t(rand.Int31()),
		dstMac:            _Ctype_uint32_t(rand.Int31()),
		headerLen:         _Ctype_uint32_t(rand.Int31()),
		inBytes:           _Ctype_uint64_t(rand.Int63()),
		inPkts:            _Ctype_uint64_t(rand.Int63()),
		inputPort:         _Ctype_uint32_t(rand.Int31()),
		ipSize:            _Ctype_uint32_t(rand.Int31()),
		ipv4DstAddr:       _Ctype_uint32_t(rand.Int31()),
		ipv4SrcAddr:       _Ctype_uint32_t(rand.Int31()),
		l4DstPort:         _Ctype_uint32_t(rand.Int31()),
		l4SrcPort:         _Ctype_uint32_t(rand.Int31()),
		outputPort:        _Ctype_uint32_t(rand.Int31()),
		protocol:          _Ctype_uint32_t(rand.Int31()),
		sampledPacketSize: _Ctype_uint32_t(rand.Int31()),
		srcAs:             _Ctype_uint32_t(rand.Int31()),
		srcGeo:            _Ctype_uint32_t(rand.Int31()),
		srcMac:            _Ctype_uint32_t(rand.Int31()),
		tcpFlags:          _Ctype_uint32_t(rand.Int31()),
		tos:               _Ctype_uint32_t(rand.Int31()),
		vlanIn:            _Ctype_uint32_t(rand.Int31()),
		vlanOut:           _Ctype_uint32_t(rand.Int31()),
		ipv4NextHop:       _Ctype_uint32_t(rand.Int31()),
		mplsType:          _Ctype_uint32_t(rand.Int31()),
		outBytes:          _Ctype_uint64_t(rand.Int63()),
		outPkts:           _Ctype_uint64_t(rand.Int63()),
		tcpRetransmit:     _Ctype_uint32_t(rand.Int31()),
		srcFlowTags:       (*_Ctype_char)(nil),
		dstFlowTags:       (*_Ctype_char)(nil),
		sampleRate:        _Ctype_uint32_t(rand.Int31()),
		deviceId:          _Ctype_uint32_t(rand.Int31()),
		flowTags:          (*_Ctype_char)(nil),
		timestamp:         _Ctype_int64_t(rand.Int31()),
		dstBgpAsPath:      (*_Ctype_char)(nil),
		dstBgpCommunity:   (*_Ctype_char)(nil),
		srcBgpAsPath:      (*_Ctype_char)(nil),
		srcBgpCommunity:   (*_Ctype_char)(nil),
		srcNextHopAs:      _Ctype_uint32_t(rand.Int31()),
		dstNextHopAs:      _Ctype_uint32_t(rand.Int31()),
		srcGeoRegion:      _Ctype_uint32_t(rand.Int31()),
		dstGeoRegion:      _Ctype_uint32_t(rand.Int31()),
		srcGeoCity:        _Ctype_uint32_t(rand.Int31()),
		dstGeoCity:        _Ctype_uint32_t(rand.Int31()),
		big:               _Ctype_uint8_t(rand.Int31()),
		sampleAdj:         _Ctype_uint8_t(rand.Int31()),
		ipv4DstNextHop:    _Ctype_uint32_t(rand.Int31()),
		ipv4SrcNextHop:    _Ctype_uint32_t(rand.Int31()),
		srcRoutePrefix:    _Ctype_uint32_t(rand.Int31()),
		dstRoutePrefix:    _Ctype_uint32_t(rand.Int31()),
		srcRouteLength:    _Ctype_uint8_t(byte(rand.Int31n(256))),
		dstRouteLength:    _Ctype_uint8_t(byte(rand.Int31n(256))),
		srcSecondAsn:      _Ctype_uint32_t(rand.Int31()),
		dstSecondAsn:      _Ctype_uint32_t(rand.Int31()),
		srcThirdAsn:       _Ctype_uint32_t(rand.Int31()),
		dstThirdAsn:       _Ctype_uint32_t(rand.Int31()),
		ipv6DstAddr:       (*_Ctype_uint8_t)(&ipv6srcaddr[0]),
		ipv6SrcAddr:       (*_Ctype_uint8_t)(&ipv6dstaddr[0]),
		srcEthMac:         _Ctype_uint64_t(rand.Int63()),
		dstEthMac:         _Ctype_uint64_t(rand.Int63()),
	}

	msgs, err := roundtrip(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(1, msgs.Len())

	kflow := msgs.At(0)

	noerr := func(v interface{}, e error) interface{} {
		assert.NoError(e)
		return v
	}

	assert.EqualValues(f.dstAs, kflow.DstAs())
	assert.EqualValues(f.dstGeo, kflow.DstGeo())
	assert.EqualValues(f.dstMac, kflow.DstMac())
	assert.EqualValues(f.headerLen, kflow.HeaderLen())
	assert.EqualValues(f.inBytes, kflow.InBytes())
	assert.EqualValues(f.inPkts, kflow.InPkts())
	assert.EqualValues(f.inputPort, kflow.InputPort())
	assert.EqualValues(f.ipSize, kflow.IpSize())
	assert.EqualValues(f.ipv4DstAddr, kflow.Ipv4DstAddr())
	assert.EqualValues(f.ipv4SrcAddr, kflow.Ipv4SrcAddr())
	assert.EqualValues(f.l4DstPort, kflow.L4DstPort())
	assert.EqualValues(f.l4SrcPort, kflow.L4SrcPort())
	assert.EqualValues(f.outputPort, kflow.OutputPort())
	assert.EqualValues(f.protocol, kflow.Protocol())
	assert.EqualValues(f.sampledPacketSize, kflow.SampledPacketSize())
	assert.EqualValues(f.srcAs, kflow.SrcAs())
	assert.EqualValues(f.srcGeo, kflow.SrcGeo())
	assert.EqualValues(f.srcMac, kflow.SrcMac())
	assert.EqualValues(f.tcpFlags, kflow.TcpFlags())
	assert.EqualValues(f.tos, kflow.Tos())
	assert.EqualValues(f.vlanIn, kflow.VlanIn())
	assert.EqualValues(f.vlanOut, kflow.VlanOut())
	assert.EqualValues(f.ipv4NextHop, kflow.Ipv4NextHop())
	assert.EqualValues(f.mplsType, kflow.MplsType())
	assert.EqualValues(f.outBytes, kflow.OutBytes())
	assert.EqualValues(f.outPkts, kflow.OutPkts())
	assert.EqualValues(f.tcpRetransmit, kflow.TcpRetransmit())
	assertEqualValues(t, f.srcFlowTags, noerr(kflow.SrcFlowTags()))
	assertEqualValues(t, f.dstFlowTags, noerr(kflow.DstFlowTags()))
	assert.EqualValues(f.sampleRate, kflow.SampleRate())
	assert.EqualValues(f.deviceId, kflow.DeviceId())
	assertEqualValues(t, f.flowTags, noerr(kflow.FlowTags()))
	assert.EqualValues(f.timestamp, kflow.Timestamp())
	assertEqualValues(t, f.dstBgpAsPath, noerr(kflow.DstBgpAsPath()))
	assertEqualValues(t, f.dstBgpCommunity, noerr(kflow.DstBgpCommunity()))
	assertEqualValues(t, f.srcBgpAsPath, noerr(kflow.SrcBgpAsPath()))
	assertEqualValues(t, f.srcBgpCommunity, noerr(kflow.SrcBgpCommunity()))
	assert.EqualValues(f.srcNextHopAs, kflow.SrcNextHopAs())
	assert.EqualValues(f.dstNextHopAs, kflow.DstNextHopAs())
	assert.EqualValues(f.srcGeoRegion, kflow.SrcGeoRegion())
	assert.EqualValues(f.dstGeoRegion, kflow.DstGeoRegion())
	assert.EqualValues(f.srcGeoCity, kflow.SrcGeoCity())
	assert.EqualValues(f.dstGeoCity, kflow.DstGeoCity())
	assert.EqualValues(f.big == 1, kflow.Big())
	assert.EqualValues(f.sampleAdj == 1, kflow.SampleAdj())
	assert.EqualValues(f.ipv4DstNextHop, kflow.Ipv4DstNextHop())
	assert.EqualValues(f.ipv4SrcNextHop, kflow.Ipv4SrcNextHop())
	assert.EqualValues(f.srcRoutePrefix, kflow.SrcRoutePrefix())
	assert.EqualValues(f.dstRoutePrefix, kflow.DstRoutePrefix())
	assert.EqualValues(f.srcRouteLength, kflow.SrcRouteLength())
	assert.EqualValues(f.dstRouteLength, kflow.DstRouteLength())
	assert.EqualValues(f.srcSecondAsn, kflow.SrcSecondAsn())
	assert.EqualValues(f.dstSecondAsn, kflow.DstSecondAsn())
	assert.EqualValues(f.srcThirdAsn, kflow.SrcThirdAsn())
	assert.EqualValues(f.dstThirdAsn, kflow.DstThirdAsn())
	assertEqualValues(t, f.ipv6DstAddr, noerr(kflow.Ipv6DstAddr()))
	assertEqualValues(t, f.ipv6SrcAddr, noerr(kflow.Ipv6SrcAddr()))
	assert.EqualValues(f.srcEthMac, kflow.SrcEthMac())
	assert.EqualValues(f.dstEthMac, kflow.DstEthMac())

	runtime.KeepAlive(ipv6srcaddr)
	runtime.KeepAlive(ipv6dstaddr)
}

func TestPackCustoms(t *testing.T) {
	assert := assert.New(t)

	customs := []Custom{
		{ID: 1, Value: string("foo")},
		{ID: 2, Value: uint32(42)},
		{ID: 3, Value: float32(3.14)},
	}

	ckcust := make([]_Ctype_kflowCustom, len(customs))
	for i, c := range customs {
		ckcust[i] = c.ToC()
	}

	flow := &Ckflow{
		numCustoms: _Ctype_uint32_t(len(ckcust)),
	}
	*(**_Ctype_kflowCustom)(unsafe.Pointer(&flow.customs)) = &ckcust[0]

	msgs, err := roundtrip(flow)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(1, msgs.Len())

	list, _ := msgs.At(0).Custom()
	for i := 0; i < list.Len(); i++ {
		c := list.At(i)

		assert.EqualValues(customs[i].ID, c.Id())

		switch expected := customs[i].Value.(type) {
		case string:
			val, _ := c.Value().StrVal()
			assert.Equal(chf.Custom_value_Which_strVal, c.Value().Which())
			assert.Equal(expected, val)
		case uint32:
			val := c.Value().Uint32Val()
			assert.Equal(chf.Custom_value_Which_uint32Val, c.Value().Which())
			assert.Equal(expected, val)
		case float32:
			val := c.Value().Float32Val()
			assert.Equal(chf.Custom_value_Which_float32Val, c.Value().Which())
			assert.Equal(expected, val)
		default:
			t.Fatal("unsupported custom column type", reflect.TypeOf(expected))
		}
	}
}

func pack(flows ...*Ckflow) (*capnp.Message, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	root, err := chf.NewRootPackedCHF(seg)
	if err != nil {
		return nil, err
	}

	msgs, err := root.NewMsgs(int32(len(flows)))
	if err != nil {
		return nil, err
	}

	for i, cflow := range flows {
		kflow, err := Pack(seg, cflow)
		if err != nil {
			return nil, err
		}

		msgs.Set(i, kflow)

	}

	root.SetMsgs(msgs)

	return msg, nil
}

func roundtrip(flows ...*Ckflow) (msgs chf.CHF_List, err error) {
	buf := &bytes.Buffer{}

	msg, err := pack(flows...)
	if err != nil {
		return
	}

	err = capnp.NewPackedEncoder(buf).Encode(msg)
	if err != nil {
		return
	}

	msg, err = capnp.NewPackedDecoder(buf).Decode()
	if err != nil {
		return
	}

	root, err := chf.ReadRootPackedCHF(msg)
	if err != nil {
		return
	}

	msgs, err = root.Msgs()
	if err != nil {
		return
	}

	return
}

func assertEqualValues(t *testing.T, expected interface{}, actual interface{}) bool {
	switch v := expected.(type) {
	case *_Ctype_uint8_t:
		if v == nil {
			expected = []byte(nil)
			break
		}
		n := len(actual.([]byte))
		h := reflect.SliceHeader{Data: (uintptr)(unsafe.Pointer(v)), Len: n, Cap: n}
		expected = *(*[]byte)(unsafe.Pointer(&h))
	case *_Ctype_char:
		if v == nil {
			expected = ""
			break
		}
		n := len(actual.(string))
		h := reflect.StringHeader{Data: (uintptr)(unsafe.Pointer(v)), Len: n}
		expected = *(*string)(unsafe.Pointer(&h))
	}

	return assert.EqualValues(t, expected, actual)
}

func randbytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Int31n(256))
	}
	return b
}

type Custom struct {
	ID    uint32
	Value interface{}
}

func (c *Custom) ToC() _Ctype_kflowCustom {
	kc := _Ctype_kflowCustom{
		id: _Ctype_uint64_t(c.ID),
	}

	p := unsafe.Pointer(&kc.value[0])
	switch v := c.Value.(type) {
	case string:
		kc.vtype = 1
		array := make([]_Ctype_char, len(v)+1)
		for i := range v {
			array[i] = _Ctype_char(v[i])
		}
		*(**_Ctype_char)(p) = &array[0]
	case uint32:
		kc.vtype = 2
		*(*_Ctype_uint32_t)(p) = _Ctype_uint32_t(v)
	case float32:
		kc.vtype = 3
		*(*_Ctype_float)(p) = _Ctype_float(v)
	}

	return kc
}
