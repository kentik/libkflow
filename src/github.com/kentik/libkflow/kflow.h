#ifndef KFLOW_H
#define KFLOW_H

#include <stdint.h>

// struct kflowConfig defines the flow sending configuration.
typedef struct {
    char *URL;                   // URL of receiving HTTP server
    struct {
        char *email;             // Kentik API email address
        char *token;             // Kentik API access token
        char *URL;               // URL of API HTTP server
    } API;
    struct {
        int interval;            // metrics flush interval (s)
        char *URL;               // URL of metrics server
    } metrics;
    int device_id;               // device ID
    int timeout;                 // flow sending timeout (ms)
    int verbose;                 // logging verbosity level
} kflowConfig;

// struct kflowCustom defines a custom flow field which may
// contain a string, uint32, or float32 value.
typedef struct {
    char *name;                  // field name
    uint64_t id;                 // IGNORE
    int vtype;                   // value type
    union {
        char *str;               // string value
        uint32_t u32;            // uint32 value
        float f32;               // float32 value
    } value;                     // field value
} kflowCustom;

// custom value types:

#define KFLOWCUSTOMSTR 1
#define KFLOWCUSTOMU32 2
#define KFLOWCUSTOMF32 3

// struct kflow defines the flow fields that may be sent to Kentik.
// MAC and IPv4 addresses are represented as bytes packed in network
// byte order, 6 bytes for MAC and 4 for IPv4. IPv6 addresses are
// 16 bytes in network byte order.
typedef struct {
    int64_t timestampNano;       // IGNORE
    uint32_t dstAs;              // destination AS
    uint32_t dstGeo;             // IGNORE
    uint32_t dstMac;             // IGNORE
    uint32_t headerLen;          // IGNORE
    uint64_t inBytes;            // number of bytes in
    uint64_t inPkts;             // number of packets in
    uint32_t inputPort;          // input interface identifier
    uint32_t ipSize;             // IGNORE
    uint32_t ipv4DstAddr;        // IPv4 dst address
    uint32_t ipv4SrcAddr;        // IPv4 src address
    uint32_t l4DstPort;          // layer 4 dst port
    uint32_t l4SrcPort;          // layer 4 src port
    uint32_t outputPort;         // output interface identifier
    uint32_t protocol;           // IP protocol number
    uint32_t sampledPacketSize;  // IGNORE
    uint32_t srcAs;              // source AS
    uint32_t srcGeo;             // IGNORE
    uint32_t srcMac;             // IGNORE
    uint32_t tcpFlags;           // TCP flags
    uint32_t tos;                // IPv4 ToS (DSCP + ECN)
    uint32_t vlanIn;             // input VLAN number
    uint32_t vlanOut;            // output VLAN number
    uint32_t ipv4NextHop;        // IPv4 next-hop address
    uint32_t mplsType;           // IGNORE
    uint64_t outBytes;           // number of bytes out
    uint64_t outPkts;            // number of packets out
    uint32_t tcpRetransmit;      // number of packets retransmitted
    char *srcFlowTags;           // IGNORE
    char *dstFlowTags;           // IGNORE
    uint32_t sampleRate;         // IGNORE
    uint32_t deviceId;           // IGNORE
    char *flowTags;              // IGNORE
    int64_t timestamp;           // IGNORE
    char *dstBgpAsPath;          // IGNORE
    char *dstBgpCommunity;       // IGNORE
    char *srcBgpAsPath;          // IGNORE
    char *srcBgpCommunity;       // IGNORE
    uint32_t srcNextHopAs;       // 1st AS in AS path to src
    uint32_t dstNextHopAs;       // 1st AS in AS path to dst
    uint32_t srcGeoRegion;       // IGNORE
    uint32_t dstGeoRegion;       // IGNORE
    uint32_t srcGeoCity;         // IGNORE
    uint32_t dstGeoCity;         // IGNORE
    uint8_t big;                 // IGNORE
    uint8_t sampleAdj;           // IGNORE
    uint32_t ipv4DstNextHop;     // IPv4 next-hop address for dst IP
    uint32_t ipv4SrcNextHop;     // IPv4 next-hop address for src IP
    uint32_t srcRoutePrefix;     // BGP table prefix for src IP
    uint32_t dstRoutePrefix;     // BGP table prefix for dst IP
    uint8_t srcRouteLength;      // BGP prefix length for src IP
    uint8_t dstRouteLength;      // BGP prefix length for dst IP
    uint32_t srcSecondAsn;       // 2nd AS in AS path to src
    uint32_t dstSecondAsn;       // 2nd AS in AS path to dst
    uint32_t srcThirdAsn;        // 3rd AS in AS path to src
    uint32_t dstThirdAsn;        // 3rd AS in AS path to dst
    uint8_t *ipv6DstAddr;        // IPv6 dst address
    uint8_t *ipv6SrcAddr;        // IPv6 src address
    uint64_t srcEthMac;          // src Ethernet MAC address
    uint64_t dstEthMac;          // dst Ethernet MAC address

    kflowCustom *customs;        // custom field array
    uint32_t numCustoms;         // custom field count
} kflow;

// kflowInit initializes the library and must be called prior
// to any other functions. It returns 0 on success or an error
// code on failure.
int kflowInit(kflowConfig *);

// kflowSend asynchronously dispatches a kflow record to the
// server. All fields of the record are copied and may be
// released after the function returns. It returns 0 on
// success or an error code on failure.
int kflowSend(kflow *);

// kflowStop stops the asynchronous flow sending process and
// releases all resources, waiting up to the supplied timeout
// in milliseconds. It returns 0 on success or an error code
// indicating timeout or failure.
int kflowStop(int);

// kflow error codes:

#define EKFLOWCONFIG   1         // configuration invalid
#define EKFLOWNOINIT   2         // kflowInit(...) not called
#define EKFLOWNOMEM    3         // out of memory
#define EKFLOWTIMEOUT  4         // request timed out
#define EKFLOWSEND     5         // flow could not be sent
#define EKFLOWNOCUSTOM 6         // custom field does not exist

#endif // KFLOW_H
