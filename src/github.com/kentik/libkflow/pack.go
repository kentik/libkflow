package main

// #include "kflow.h"
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/kentik/libkflow/chf"
	"zombiezen.com/go/capnproto2"
)

// Ckflow is an alias for C.kflow used because tests run by
// `go test` cannot reference the "C" package.
type Ckflow C.kflow

// Pack allocates a new capnp struct in the supplied segment
// and packets the contents of the C flow structure into it.
func Pack(seg *capnp.Segment, cflow *Ckflow) (kflow chf.CHF, err error) {
	kflow, err = chf.NewCHF(seg)
	if err != nil {
		return
	}

	kflow.SetTimestampNano(int64(cflow.timestampNano))
	kflow.SetDstAs(uint32(cflow.dstAs))
	kflow.SetDstGeo(uint32(cflow.dstGeo))
	kflow.SetDstMac(uint32(cflow.dstMac))
	kflow.SetHeaderLen(uint32(cflow.headerLen))
	kflow.SetInBytes(uint64(cflow.inBytes))
	kflow.SetInPkts(uint64(cflow.inPkts))
	kflow.SetInputPort(uint32(cflow.inputPort))
	kflow.SetIpSize(uint32(cflow.ipSize))
	kflow.SetIpv4DstAddr(uint32(cflow.ipv4DstAddr))
	kflow.SetIpv4SrcAddr(uint32(cflow.ipv4SrcAddr))
	kflow.SetL4DstPort(uint32(cflow.l4DstPort))
	kflow.SetL4SrcPort(uint32(cflow.l4SrcPort))
	kflow.SetOutputPort(uint32(cflow.outputPort))
	kflow.SetProtocol(uint32(cflow.protocol))
	kflow.SetSampledPacketSize(uint32(cflow.sampledPacketSize))
	kflow.SetSrcAs(uint32(cflow.srcAs))
	kflow.SetSrcGeo(uint32(cflow.srcGeo))
	kflow.SetSrcMac(uint32(cflow.srcMac))
	kflow.SetTcpFlags(uint32(cflow.tcpFlags))
	kflow.SetTos(uint32(cflow.tos))
	kflow.SetVlanIn(uint32(cflow.vlanIn))
	kflow.SetVlanOut(uint32(cflow.vlanOut))
	kflow.SetIpv4NextHop(uint32(cflow.ipv4NextHop))
	kflow.SetMplsType(uint32(cflow.mplsType))
	kflow.SetOutBytes(uint64(cflow.outBytes))
	kflow.SetOutPkts(uint64(cflow.outPkts))
	kflow.SetTcpRetransmit(uint32(cflow.tcpRetransmit))
	kflow.SetSrcFlowTags(C.GoString(cflow.srcFlowTags))
	kflow.SetDstFlowTags(C.GoString(cflow.dstFlowTags))
	kflow.SetSampleRate(uint32(cflow.sampleRate))
	kflow.SetDeviceId(uint32(cflow.deviceId))
	kflow.SetFlowTags(C.GoString(cflow.flowTags))
	kflow.SetTimestamp(int64(cflow.timestamp))
	kflow.SetDstBgpAsPath(C.GoString(cflow.dstBgpAsPath))
	kflow.SetDstBgpCommunity(C.GoString(cflow.dstBgpCommunity))
	kflow.SetSrcBgpAsPath(C.GoString(cflow.srcBgpAsPath))
	kflow.SetSrcBgpCommunity(C.GoString(cflow.srcBgpCommunity))
	kflow.SetSrcNextHopAs(uint32(cflow.srcNextHopAs))
	kflow.SetDstNextHopAs(uint32(cflow.dstNextHopAs))
	kflow.SetSrcGeoRegion(uint32(cflow.srcGeoRegion))
	kflow.SetDstGeoRegion(uint32(cflow.dstGeoRegion))
	kflow.SetSrcGeoCity(uint32(cflow.srcGeoCity))
	kflow.SetDstGeoCity(uint32(cflow.dstGeoCity))
	kflow.SetBig(cflow.big == 1)
	kflow.SetSampleAdj(cflow.sampleAdj == 1)
	kflow.SetIpv4DstNextHop(uint32(cflow.ipv4DstNextHop))
	kflow.SetIpv4SrcNextHop(uint32(cflow.ipv4SrcNextHop))
	kflow.SetSrcRoutePrefix(uint32(cflow.srcRoutePrefix))
	kflow.SetDstRoutePrefix(uint32(cflow.dstRoutePrefix))
	kflow.SetSrcRouteLength(uint8(cflow.srcRouteLength))
	kflow.SetDstRouteLength(uint8(cflow.dstRouteLength))
	kflow.SetSrcSecondAsn(uint32(cflow.srcSecondAsn))
	kflow.SetDstSecondAsn(uint32(cflow.dstSecondAsn))
	kflow.SetSrcThirdAsn(uint32(cflow.srcThirdAsn))
	kflow.SetDstThirdAsn(uint32(cflow.dstThirdAsn))
	kflow.SetIpv6DstAddr(bts(cflow.ipv6DstAddr, 16))
	kflow.SetIpv6SrcAddr(bts(cflow.ipv6SrcAddr, 16))
	kflow.SetSrcEthMac(uint64(cflow.srcEthMac))
	kflow.SetDstEthMac(uint64(cflow.dstEthMac))

	list, err := PackCustoms(seg, cflow)
	if err != nil {
		return
	}
	kflow.SetCustom(*list)

	return
}

func PackCustoms(seg *capnp.Segment, cflow *Ckflow) (*chf.Custom_List, error) {
	customs := *(*[]C.kflowCustom)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (uintptr)(unsafe.Pointer(cflow.customs)),
		Len:  int(cflow.numCustoms),
		Cap:  int(cflow.numCustoms),
	}))

	list, err := chf.NewCustom_List(seg, int32(len(customs)))
	if err != nil {
		return nil, err
	}

	for i, ccc := range customs {
		c, err := chf.NewCustom(seg)
		if err != nil {
			return nil, err
		}

		p := unsafe.Pointer(&ccc.value[0])
		v := c.Value()

		c.SetId(uint32(ccc.id))
		switch ccc.vtype {
		case C.KFLOWCUSTOMSTR:
			v.SetStrVal(C.GoString(*(**C.char)(p)))
		case C.KFLOWCUSTOMU32:
			v.SetUint32Val(uint32(*(*C.uint32_t)(p)))
		case C.KFLOWCUSTOMF32:
			v.SetFloat32Val(float32(*(*C.float)(p)))
		}

		list.Set(i, c)
	}

	return &list, nil
}

func bts(p *C.uint8_t, len C.int) []byte {
	if p == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), len)
}
