#include <stdio.h>
#include <stdlib.h>
#include "kflow.h"

int main(int argc, char **argv) {
    int r;
    kflowConfig cfg = {
        .URL = "http://chdev:20012/chf",
        .API = {
            .email = "will@kentik.com",
            .token = "81b7262feceecc94eef3ddafbc2c152f",
            .URL   = "http://chdev:8080/api/v5",
        },
        .device_id = 1001,
        .verbose   = 1,
    };

    if ((r = kflowInit(&cfg)) != 0) {
        printf("error initializing libkflow: %d", r);
        exit(1);
    };

    kflow flow = {
        .deviceId    = cfg.device_id,
        .ipv4SrcAddr = 167772161,
        .ipv4DstAddr = 167772162,
        .srcAs       = 1234,
        .inPkts      = 20,
        .inBytes     = 40,
    };

    if ((r = kflowSend(&flow)) != 0) {
        printf("error sending flow: %d", r);
        exit(1);
    }

    if ((r = kflowStop(10*1000)) != 0) {
        printf("error stopping libkflow: %d", r);
        exit(1);
    }

    return 0;
}
