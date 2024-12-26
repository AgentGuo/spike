package io.pixelsdb.pixels.spike.handler.impl;

import io.pixelsdb.pixels.spike.handler.RequestHandler;
import io.pixelsdb.pixels.spike.handler.SpikeWorker;

public class Example implements RequestHandler {
    @Override
    public SpikeWorker.CallWorkerFunctionResp execute(SpikeWorker.CallWorkerFunctionReq request) {
        int sleepSeconds = Integer.parseInt(request.getPayload());
        try {
            Thread.sleep(sleepSeconds * 1000L);
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }
        // example：反转有效负载字符串
        return SpikeWorker.CallWorkerFunctionResp.newBuilder()
                .setRequestId(request.getRequestId())
                .setPayload(new StringBuilder(request.getPayload()).reverse().toString())
                .build();
    }
}
