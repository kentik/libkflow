#include <stdio.h>
#include <stdlib.h>
#include "kflow.h"

int main(int argc, char **argv) {
    int r;
    kflowConfig cfg = {
        .URL = "http://127.0.0.1:8999/chf",
        .API = {
            .email = "test@example.com",
            .token = "token",
            .URL   = "http://127.0.0.1:8999/api/v5",
        },
        .metrics = {
            .interval = 1,
            .URL      = "http://127.0.0.1:8889/metrics",
        },
        .device_id = 1,
        .verbose   = 1,
    };

    if ((r = kflowInit(&cfg)) != 0) {
        printf("error initializing libkflow: %d\n", r);
        exit(1);
    };

    kflowCustom customs[] = {
        { .name = "CUSTOM-STR", .vtype = KFLOWCUSTOMSTR, .value.str = &"foo"[0] },
        { .name = "CUSTOM-U32", .vtype = KFLOWCUSTOMU32, .value.u32 = 42        },
        { .name = "CUSTOM-F32", .vtype = KFLOWCUSTOMF32, .value.f32 = 3.14      },
    };
    uint32_t numCustoms = sizeof(customs) / sizeof(kflowCustom);

    kflow flow = {
        .deviceId    = cfg.device_id,
        .ipv4SrcAddr = 167772161,
        .ipv4DstAddr = 167772162,
        .srcAs       = 1234,
        .inPkts      = 20,
        .inBytes     = 40,
        .srcEthMac   = 1250999896491,
        .dstEthMac   = 226426397786884,
        .customs     = customs,
        .numCustoms  = numCustoms,
    };

    if ((r = kflowSend(&flow)) != 0) {
        printf("error sending flow: %d\n", r);
        exit(1);
    }

    if ((r = kflowStop(10*1000)) != 0) {
        printf("error stopping libkflow: %d\n", r);
        exit(1);
    }

    return 0;
}
